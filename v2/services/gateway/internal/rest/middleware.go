package rest

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Middleware func(next http.Handler) http.Handler

// ErrUnknown will be logged when the panic recovery has an unknown type.
var ErrUnknown = errors.New("unknown error")

func recoverPanicMiddleware(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					var err error

					switch e := rec.(type) {
					case string:
						err = fmt.Errorf("%w", err)
					case error:
						err = e
					default:
						err = ErrUnknown
					}

					w.WriteHeader(http.StatusInternalServerError)
					logger.Error("application panicked", zap.Error(err))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
