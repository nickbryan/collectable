package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("responds when route not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		t.Cleanup(cancel)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/not/a/route", nil)
		if err != nil {
			assert.FailNow(t, fmt.Sprintf("unable to create request: %v", err))
		}

		res := httptest.NewRecorder()

		core, logs := observer.New(zap.DebugLevel)
		s := NewServer(zap.New(core))
		s.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.JSONEq(t, `{"errors": [{"message": "resource not found"}]}`, res.Body.String())

		expectedLog := "responding with errors"
		lgs := logs.FilterMessage(expectedLog)
		assert.Equal(t, lgs.Len(), 1, "expected to find a log with the given message: %s, got: %+v", expectedLog, logs.All())
		assert.Equal(t, zap.InfoLevel, lgs.All()[0].Level, "expected info level got: %+v", lgs.All()[0].Level)
		// TODO: test for the errors in the log context on all of these tests
	})

	t.Run("responds when method not allowed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		t.Cleanup(cancel)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/existing/route", nil)
		if err != nil {
			assert.FailNow(t, fmt.Sprintf("unable to create request: %v", err))
		}

		res := httptest.NewRecorder()

		handler := Handler{
			Route: func(r *mux.Route) {
				r.Path("/existing/route").Methods(http.MethodPost)
			},
		}
		core, logs := observer.New(zap.DebugLevel)
		s := NewServer(zap.New(core))
		s.RegisterHandlers(handler)
		s.ServeHTTP(res, req)

		assert.Equal(t, http.StatusMethodNotAllowed, res.Result().StatusCode)
		assert.JSONEq(t, `{"errors": [{"message": "method GET is not allowed on this resource"}]}`, res.Body.String())

		expectedLog := "responding with errors"
		lgs := logs.FilterMessage(expectedLog)
		assert.Equal(t, lgs.Len(), 1, "expected to find a log with the given message: %s, got: %+v", expectedLog, logs.All())
		assert.Equal(t, zap.InfoLevel, lgs.All()[0].Level, "expected info level got: %+v", lgs.All()[0].Level)
	})

	t.Run("registers handlers", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		t.Cleanup(cancel)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/existing/route", nil)
		if err != nil {
			assert.FailNow(t, fmt.Sprintf("unable to create request: %v", err))
		}

		res := httptest.NewRecorder()

		called := false
		core, _ := observer.New(zap.DebugLevel)
		s := NewServer(zap.New(core))
		s.RegisterHandlers(func(called *bool) Handler {
			return Handler{
				Route: func(r *mux.Route) {
					r.Path("/existing/route").Methods(http.MethodGet)
				},
				Action: func(res Responder, req *Request) {
					*called = true
					res.Respond(http.StatusNoContent)
				},
			}
		}(&called))
		s.ServeHTTP(res, req)

		assert.True(t, called, "expected handler action to be called")
		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)
	})
}
