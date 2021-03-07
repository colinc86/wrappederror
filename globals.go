package wrappederror

// Error interface values

// Variables used to configure the behavior of wError instances.
var (
	// Whether or not errors should capture their caller information.
	captureCaller = newConfigValue(true)

	// Whether or not errors should capture their process information.
	captureProcess = newConfigValue(true)
)

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

// Caller interface values

// Variables used to configure the behavior of wCaller instances.
var (
	// The radius around the detected line number when creating a source code
	// fragment.
	sourceFragmentRadius = newConfigValue(2)
)

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

// Process interface values

// Variables used to configure the behavior of wProcess instances.
var (
	// Whether or not process types should ignore breakpoints when their `Break`
	// method is called.
	ignoreBreakpoints = newConfigValue(true)
)

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

// Metadata interface values

// Variables used to configure metadata instances.
var (
	// The next error index that will be used when creating error metadata.
	nextErrorIndex = newConfigValue(1)

	// Whether or not similar errors should be tracked.
	trackSimilarErrors = newConfigValue(true)
)

// SetNextErrorIndex sets the next error index that will be used when creating
// an error.
func SetNextErrorIndex(index int) {
	nextErrorIndex.set(index)
}

// NextErrorIndex gets the next error index that will be used when creating an
// error.
func NextErrorIndex() int {
	return nextErrorIndex.get().(int)
}

// SetTrackSimilarErrors enables or prohibits similar error tracking.
func SetTrackSimilarErrors(track bool) {
	trackSimilarErrors.set(track)
}

// TrackSimilarErrors returns whether or not similar errors are being tracked.
func TrackSimilarErrors() bool {
	return trackSimilarErrors.get().(bool)
}

// Non-exported functions and variables

// getAndIncrementNextErrorIndex gets the next error index and increments the
// value.
func getAndIncrementNextErrorIndex() int {
	nextErrorIndex.mutex.Lock()
	v := nextErrorIndex.value.(int)
	nextErrorIndex.value = nextErrorIndex.value.(int) + 1
	nextErrorIndex.mutex.Unlock()
	return v
}

// The hash map responsible for tracking similar errors.
var errorHashMap = newErrorMap()

// getSimilarErrorCount gets and returns the number of errors in the error hash
// map equal to err.
func getSimilarErrorCount(err error) int {
	if !TrackSimilarErrors() || err == nil {
		return 0
	}

	s := errorHashMap.similarErrors(err)
	errorHashMap.addError(err)
	return s
}
