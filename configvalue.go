package wrappederror

import "sync"

// Stores a configuration value and provides methods to safely access and set
// the value.
type configValue struct {

	// The configuration value's inner value. Do not access this property
	// directly unless you plan on utilizing the value's mutex.
	value interface{}

	// Mutex for safe access to value.
	mutex *sync.RWMutex
}

// Initializers

// newConfigValue creates and returns a new configuration value.
func newConfigValue(cv interface{}) *configValue {
	return &configValue{
		value: cv,
		mutex: new(sync.RWMutex),
	}
}

// Non-exported methods

// set sets a configuration value.
func (v *configValue) set(cv interface{}) {
	v.mutex.Lock()
	v.value = cv
	v.mutex.Unlock()
}

// get gets a configuration value.
func (v *configValue) get() interface{} {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.value
}
