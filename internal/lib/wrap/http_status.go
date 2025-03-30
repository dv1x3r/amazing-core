package wrap

import (
	"errors"
	"net/http"
)

type HTTPStatusError struct {
	err    error
	status int
}

func (e HTTPStatusError) Error() string {
	return e.err.Error()
}

func (e HTTPStatusError) Unwrap() error {
	return e.err
}

func (e HTTPStatusError) HTTPStatus() int {
	return e.status
}

func WithHTTPStatus(err error, status int) error {
	return HTTPStatusError{
		err:    err,
		status: status,
	}
}

func HTTPStatus(err error) int {
	var statusErr HTTPStatusError
	if errors.As(err, &statusErr) {
		return statusErr.HTTPStatus()
	}
	return http.StatusInternalServerError
}
