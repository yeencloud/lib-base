package errors

import (
	"net/http"
)

// MARK: Not Implemented
type NotImplementedError struct {
}

func (NotImplementedError) Error() string {
	return "not implemented"
}

func (NotImplementedError) RestCode() int {
	return http.StatusNotImplemented
}

// MARK: Module Disabled
type ModuleDisabledError struct {
	Module string
}

func (e *ModuleDisabledError) Error() string {
	return "module " + e.Module + " is disabled"
}

func (e *ModuleDisabledError) HowToFix() string {
	return "set `Use" + e.Module + "` to true in the service options when calling `Run`"
}

// MARK: Module Not Initialized

type ModuleNotInitializedError struct {
	Module string
}

func (e *ModuleNotInitializedError) Error() string {
	return "module " + e.Module + " is not initialized"
}
