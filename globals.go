package wrappederror

// The package's current state.
//
// Do not set this after launch.
var packageState = newState()

// Exported functions

// Config returns the package's configuration.
func Config() *Configuration {
	return packageState.config
}

// ResetState resets the package's state to that at process launch.
func ResetState() {
	packageState.reset()
}

// RegisterErrorSeverity registers the error severity with the package. If the
// severity has already been registered, then a ErrSeverityAlreadyRegistered
// error is returned.
func RegisterErrorSeverity(severity *ErrorSeverity) error {
	return packageState.registerSeverity(severity)
}

// UnregisterErrorSeverity unregisters the error severity from the package. If
// the severity wasn't already registered, then this function does nothing.
func UnregisterErrorSeverity(severity *ErrorSeverity) {
	packageState.unregisterSeverity(severity)
}
