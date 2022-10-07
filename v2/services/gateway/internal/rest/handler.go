package rest

import (
	"net/http"
)

type Handler struct {
	Action  http.Handler
	Methods []string
	Path    string
}

type Action[Req, Res any] func(r Request[Req]) (*Response[Res], error)

type Request[T any] struct {
	*http.Request

	Data T

	PathValues map[string]string
}

type Response[T any] struct {
	Header http.Header

	data       T
	statusCode int
}

func newResponse[T any](statusCode int, data T) *Response[T] {
	return &Response[T]{
		Header:     make(http.Header),
		data:       data,
		statusCode: statusCode,
	}
}

func NewResponse[T any](statusCode int, data T) *Response[T] {
	return newResponse(statusCode, data)
}

func NewEmptyResponse[T any]() *Response[NoBody] {
	return newResponse(http.StatusNoContent, NoBody{})
}

type ResponseError struct {
	err    error
	Header http.Header
	status int
}

// TODO: test all of this and clean up code
// TODO: copy remaining functionality from old gateway code
// TODO: move Earthfile and this gateway into their expected places in the code

func (e *ResponseError) Error() string { // TODO: improve this with unwrap etc.
	return e.err.Error()
}

func NewErrorResponse(statusCode int, err error) *ResponseError {
	return &ResponseError{
		Header: make(http.Header),
		status: statusCode,
		err:    err,
	}
}
