package rest

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name                 string
		panicArg             any
		expectedErrorMessage string
	}{
		{
			name:                 "recovers from string panic",
			panicArg:             "panic from test action",
			expectedErrorMessage: "panic from test action",
		},
		{
			name:                 "recovers from error panic",
			panicArg:             errors.New("error panic from test action"),
			expectedErrorMessage: "error panic from test action",
		},
		{
			name:                 "recovers from unexpected type panic",
			panicArg:             12345,
			expectedErrorMessage: "unknown error",
		},
	}

	for _, test := range tests {
		tc := test

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := Handler{
				Route: func(r *mux.Route) {
					r.Path("/test").Methods(http.MethodGet)
				},
				Action: func(res Responder, req *Request) {
					panic(tc.panicArg)
				},
			}

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			t.Cleanup(cancel)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/test", nil)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("unable to create request: %v", err))
			}

			res := httptest.NewRecorder()

			r := mux.NewRouter()
			core, logs := observer.New(zap.DebugLevel)
			handler.AddRoute(r, zap.New(core))
			r.ServeHTTP(res, req)

			assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)

			expectedLog := "application panicked"
			lgs := logs.FilterMessage(expectedLog)
			assert.Equal(t, lgs.Len(), 1, "expected to find a log with the given message: %s, got: %+v", expectedLog, logs.All())
			assert.Equal(t, zap.ErrorLevel, lgs.All()[0].Level, "expected error level got: %+v", lgs.All()[0].Level)
			assert.Equal(t, tc.expectedErrorMessage, lgs.All()[0].Context[0].Interface.(error).Error())
		})
	}
}

func TestResponse(t *testing.T) {
	t.Run("sets status code on writer when created through Responder", func(t *testing.T) {
		rec := httptest.NewRecorder()

		newResponder(nil, rec).Respond(http.StatusNoContent)

		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("sets content-type to application/json when responding with data", func(t *testing.T) {
		rec := httptest.NewRecorder()

		Response{
			logger: nil,
			writer: rec,
		}.WithData(nil)

		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("writes data to the writer when responding with data", func(t *testing.T) {
		rec := httptest.NewRecorder()

		Response{
			logger: nil,
			writer: rec,
		}.WithData(struct {
			SomeKey string
		}{
			SomeKey: "test-value",
		})

		require.JSONEq(t, `{"SomeKey":"test-value"}`, rec.Body.String())
	})

	t.Run("sets status code to 500 if json encoding fails when responding with data", func(t *testing.T) {
		rec := httptest.NewRecorder()
		core, _ := observer.New(zap.DebugLevel)

		Response{
			logger: zap.New(core),
			writer: rec,
		}.WithData(math.Inf(1)) // Infinity is an invalid value so json.UnsupportedValueError will be thrown.

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("writes an error log if json encoding fails when responding with data", func(t *testing.T) {
		rec := httptest.NewRecorder()
		core, logs := observer.New(zap.DebugLevel)

		Response{
			logger: zap.New(core),
			writer: rec,
		}.WithData(math.Inf(1)) // Infinity is an invalid value so json.UnsupportedValueError will be thrown.

		lgs := logs.FilterLevelExact(zap.ErrorLevel)
		require.Equal(t, 1, lgs.Len())
		require.Equal(t, "unable to encode the response as json", lgs.All()[0].Message)
		require.Equal(t, "json: unsupported value: +Inf", lgs.All()[0].Context[0].Interface.(error).Error())
	})

	t.Run("sets content-type to application/json when responding with errors", func(t *testing.T) {
		rec := httptest.NewRecorder()
		core, _ := observer.New(zap.DebugLevel)

		Response{
			logger: zap.New(core),
			writer: rec,
		}.WithErrors()

		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("writes data to the writer when responding with errors", func(t *testing.T) {
		rec := httptest.NewRecorder()
		core, _ := observer.New(zap.DebugLevel)

		Response{
			logger: zap.New(core),
			writer: rec,
		}.WithErrors(errors.New("my test error"), errors.New("my second test error"))

		require.JSONEq(t, `{"errors":[{"message":"my test error"},{"message":"my second test error"}]}`, rec.Body.String())
	})

	// TODO: test the logs on the error response
}
