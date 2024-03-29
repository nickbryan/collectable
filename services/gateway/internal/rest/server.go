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

// Server defines a HTTP server for handling Rest requests.
type Server struct {
	router *mux.Router
	logger *zap.Logger
}

// NewServer initialises a new Server with a router.
func NewServer(logger *zap.Logger) *Server {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		newResponder(logger, w).Respond(http.StatusNotFound).WithErrors(errors.New("resource not found"))
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newResponder(logger, w).Respond(http.StatusMethodNotAllowed).WithErrors(fmt.Errorf("method %s is not allowed on this resource", r.Method))
	})

	return &Server{router: router, logger: logger}
}

// Start the server and listen for incoming requests.
func (s *Server) Start(address string) error {
	srv := &http.Server{
		Addr:         address,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      s,
	}

	// Here we start the web server listing for connections and serving responses. When the server
	// closes we may get back an error so we create a channel that allows us to receive that error
	// to be reported on later. If the error is http.ErrServerClosed then we ignore it as we expect
	// that to happen at some point. We start this in a separate go routine to allow us to block on
	// an os.Interrupt signal later.
	errChan := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down the server: %w", err)
	}

	// Here we return the result of our error channel from earlier. If there was no error then we will receive nil
	// otherwise we will receive the specified error.
	return <-errChan
}

// ServeHTTP requests via the internal router.
// This is what allows us to use our Server struct as the http.Server Handler in the Start method.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// RegisterHandlers with the router. This allows a Handler to define their route with the router.
func (s *Server) RegisterHandlers(handlers ...Handler) {
	for _, h := range handlers {
		h.AddRoute(s.router, s.logger)
	}
}
