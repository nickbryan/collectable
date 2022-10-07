package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	router *mux.Router
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	router := mux.NewRouter()

	router.NotFoundHandler = defaultNotFoundHandler()
	router.MethodNotAllowedHandler = defaultMethodNotAllowedHandler()

	svr := &Server{router: router, logger: logger}
	svr.UseMiddleware(recoverPanicMiddleware(logger))

	return svr
}

func (s *Server) Start(ctx context.Context, address string) error { // TODO: handler context
	const defaultTimeout = 30 * time.Second // TODO configuration here

	svr := &http.Server{
		Addr:              address,
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultTimeout,
		ReadHeaderTimeout: defaultTimeout,
		Handler:           s,
	}

	// Here we start the web server listing for connections and serving responses. When the server
	// closes we may get back an error so we create a channel that allows us to receive that error
	// to be reported on later. If the error is http.ErrServerClosed then we ignore it as we expect
	// that to happen at some point. We start this in a separate go routine to allow us to block on
	// an os.Interrupt signal later.
	errChan := make(chan error, 1)

	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("listening and serving http requests: %w", err)
		}

		close(errChan)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan // Block awaiting an interrupt signal

	// Once we receive an interrupt signal, we know it is time to shut down the web server. We will
	// allow 30 seconds for graceful shutdown before we force the close. This is handled through
	// our context.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutting down the server should trigger our http.ErrServerClosed that we ignore
	// in the above go routine.
	if err := svr.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down the server: %w", err)
	}

	// Here we return the result of our error channel from earlier. If there was no error then we will receive nil
	// otherwise we will receive the specified error.
	return <-errChan
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type loggerSetter interface {
	setLogger(l *zap.Logger)
}

func bindHandler(l *zap.Logger, r *mux.Router, h Handler) {
	if logSetter, ok := h.Action.(loggerSetter); ok {
		logSetter.setLogger(l)
	}

	r.Path(h.Path).Methods(h.Methods...).Handler(h.Action)
}

func (s *Server) RegisterHandlerGroup(pathPrefix string, handlers ...Handler) {
	subRouter := s.router.PathPrefix(pathPrefix).Subrouter()

	for _, h := range handlers {
		bindHandler(s.logger, subRouter, h)
	}
}

func (s *Server) UseMiddleware(middleware Middleware) {
	s.router.Use(mux.MiddlewareFunc(middleware))
}
