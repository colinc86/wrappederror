package wrappederror

import "sync"

// Stores a configuration value and provides methods to safely access and set
// the value across goroutines.
type safeValue struct {

	// The value's inner value. Do not access this property directly. If you need
	// to modify the value but the set method is insufficient, then use the
	// transform method and perform your operation inside it's parameter function.
	value interface{}

	// Mutex for safe access to value.
	mutex *sync.RWMutex
}

// Initializers

// newSafeValue creates and returns a new thread-safe value.
func newSafeValue(v interface{}) *safeValue {
	return &safeValue{
		value: v,
		mutex: new(sync.RWMutex),
	}
}

// Non-exported methods

// set sets the value.
func (v *safeValue) set(cv interface{}) {
	v.mutex.Lock()
	v.value = cv
	v.mutex.Unlock()
}

// get gets the value.
func (v *safeValue) get() interface{} {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.value
}

// transform takes a transform function that takes the value of the receiver
// and returns a transformed value that replace's the receiver's value.
func (v *safeValue) transform(t func(v interface{}) interface{}) {
	v.mutex.Lock()
	v.value = t(v.value)
	v.mutex.Unlock()
}
