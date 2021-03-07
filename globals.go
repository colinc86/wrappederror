package wrappederror

// The package's current state.
//
// Do not set this after launch.
var packageState = newState()

// Exported functions

// Config returns the package's configuration.
func Config() *Configuration {
	return packageState.configuration
}

// ResetState resets the package's state to that at process launch.
func ResetState() {
	packageState.reset()
}
