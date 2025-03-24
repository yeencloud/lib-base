package errors

// MARK: UnsupportedDatabaseEngineError
type UnsupportedDatabaseEngineError struct {
	Engine string
}

func (e *UnsupportedDatabaseEngineError) Error() string {
	return "unsupported database engine: " + e.Engine
}
