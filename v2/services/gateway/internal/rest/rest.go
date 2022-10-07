package rest

import (
	"errors"
	"fmt"
	"net/http"
)

type NoBody struct{} // TODO: test equality

var errResourceNotFound = errors.New("resource not found")

func defaultNotFoundHandler() http.Handler {
	return JSON(func(_ Request[NoBody]) (*Response[NoBody], error) {
		return nil, NewErrorResponse(http.StatusNotFound, errResourceNotFound)
	})
}

type methodNotAllowedError struct {
	method string
}

func (e *methodNotAllowedError) Error() string {
	return fmt.Sprintf("method %s is not allowed on this resource", e.method)
}

func defaultMethodNotAllowedHandler() http.Handler {
	return JSON(func(r Request[NoBody]) (*Response[NoBody], error) {
		return nil, NewErrorResponse(http.StatusMethodNotAllowed, &methodNotAllowedError{method: r.Method})
	})
}
