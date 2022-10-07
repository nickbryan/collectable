package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Action func(res Responder, req *Request)

type Handler struct {
	Route  func(r *mux.Route)
	Action Action
}

func (h Handler) AddRoute(r *mux.Router, logger *zap.Logger) {
	h.Route(r.NewRoute().HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		recoverPanicMiddleware(h.Action, logger)(newResponder(logger, w), &Request{Request: req})
	}))
}

// ErrUnknown will be logged when the panic recovery has an unknown type.
var ErrUnknown = errors.New("unknown error")

func recoverPanicMiddleware(next Action, logger *zap.Logger) Action {
	return func(res Responder, req *Request) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error

				switch e := rec.(type) {
				case string:
					err = errors.New(e)
				case error:
					err = e
				default:
					err = ErrUnknown
				}

				res.Respond(http.StatusInternalServerError)
				logger.Error("application panicked", zap.Error(err))
			}
		}()

		next(res, req)
	}
}

// Request wraps the http.Request so that we can add custom methods.
type Request struct {
	*http.Request
}

// Decode de-serialises the JSON body of the request into the passed destination object.
func (r Request) Decode(dest interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dest)
	if err != nil {
		return fmt.Errorf("decoding request as json: %w", err)
	}

	return nil
}

type Responder interface {
	Respond(statusCode int) Response
}

type responder struct {
	logger *zap.Logger
	writer http.ResponseWriter
}

func newResponder(l *zap.Logger, w http.ResponseWriter) Responder {
	return responder{logger: l, writer: w}
}

func (r responder) Respond(statusCode int) Response {
	r.writer.WriteHeader(statusCode)

	return Response(r)
}

type Response struct {
	logger *zap.Logger
	writer http.ResponseWriter
}

func (r Response) WithData(data any) {
	r.writer.Header().Set("Content-Type", "application/json")

	if data != nil {
		if err := json.NewEncoder(r.writer).Encode(data); err != nil {
			r.logger.Error("unable to encode the response as json", zap.Error(err))
			r.writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (r Response) WithErrors(errors ...error) {
	type errorDefinition struct {
		Message string `json:"message"`
	}

	type response struct {
		Errors []errorDefinition `json:"errors"`
	}

	resp := response{}
	for _, err := range errors {
		resp.Errors = append(resp.Errors, errorDefinition{Message: err.Error()})
	}

	r.logger.Info("responding with errors", zap.Errors("errors", errors))

	r.WithData(resp)
}
