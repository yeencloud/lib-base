package errors

import (
	"net/http"
)

type NotImplementedError struct {
}

func (NotImplementedError) Error() string {
	return "not implemented"
}

func (NotImplementedError) RestCode() int {
	return http.StatusNotImplemented
}
