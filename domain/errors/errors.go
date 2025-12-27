package errors

// MARK: Module Disabled
type ModuleDisabledError struct {
	Module string
}

func (e *ModuleDisabledError) Error() string {
	return "module " + e.Module + " is disabled"
}

func (e *ModuleDisabledError) TroubleshootingTip() string {
	return "set `Use" + e.Module + "` to true in the service options when calling `Run`"
}

// MARK: Module Not Initialized

type ModuleNotInitializedError struct {
	Module string
}

func (e *ModuleNotInitializedError) Error() string {
	return "module " + e.Module + " is not initialized"
}
