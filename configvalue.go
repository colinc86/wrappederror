package wrappederror

import "sync"

// Variables used to configure the behavior of wError instances.
var (
	// Whether or not errors should capture their caller information.
	captureCaller = newConfigValue(true)

	// Whether or not errors should capture their process information.
	captureProcess = newConfigValue(true)
)

// Variables used to configure the behavior of wCaller instances.
var (
	// The radius around the detected line number when creating a source code
	// fragment.
	sourceFragmentRadius = newConfigValue(2)
)

// Variables used to configure the behavior of wProcess instances.
var (
	// Whether or not process types should ignore breakpoints when their `Break`
	// method is called.
	ignoreBreakpoints = newConfigValue(true)
)

// Implementation

// Stores a configuration value and provides methods to safely access and set
// the value.
type configValue struct {

	// The configuration value's inner value. Do not access this property
	// directly.
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

// Exported functions

// SetCaptureCaller sets a flag to determine if new errors capture their caller
// information.
func SetCaptureCaller(capture bool) {
	captureCaller.set(capture)
}

// CaptureCaller gets a boolean that indicates whether or not new errors capture
// their caller information.
func CaptureCaller() bool {
	return captureCaller.get().(bool)
}

// SetCaptureProcess sets a flag to determine if new errors capture their
// process information.
func SetCaptureProcess(capture bool) {
	captureProcess.set(capture)
}

// CaptureProcess gets a boolean that indicates whether or not new errors
// capture their process information.
func CaptureProcess() bool {
	return captureProcess.get().(bool)
}

// SetSourceFragmentRadius sets the radius of the source fragment obtained from
// source files at the line that the caller was created on.
func SetSourceFragmentRadius(radius int) {
	sourceFragmentRadius.set(radius)
}

// SourceFragmentRadius gets the radius of source fragments obtained from source
// files.
func SourceFragmentRadius() int {
	return sourceFragmentRadius.get().(int)
}

// SetIgnoreBreakpoints tells all calls to `Break` on `Process` types to either
// handle or ignore invocations.
func SetIgnoreBreakpoints(ignore bool) {
	ignoreBreakpoints.set(ignore)
}

// IgnoreBreakpoints returns whether or not calls to `Break` on `Process` types
// will be ignored. This value defaults to true.
func IgnoreBreakpoints() bool {
	return ignoreBreakpoints.get().(bool)
}
