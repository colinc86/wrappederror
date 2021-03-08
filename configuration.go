package wrappederror

// Configuration types keep track of the package's configuration.
type Configuration struct {
	captureCaller          *configValue
	captureProcess         *configValue
	marshalMinimalJSON     *configValue
	captureSourceFragments *configValue
	sourceFragmentRadius   *configValue
	ignoreBreakpoints      *configValue
	nextErrorIndex         *configValue
	trackSimilarErrors     *configValue
}

// Initializers

// newConfiguration creates and returns a new configuration.
func newConfiguration() *Configuration {
	return &Configuration{
		captureCaller:          newConfigValue(true),
		captureProcess:         newConfigValue(true),
		marshalMinimalJSON:     newConfigValue(true),
		captureSourceFragments: newConfigValue(true),
		sourceFragmentRadius:   newConfigValue(2),
		ignoreBreakpoints:      newConfigValue(true),
		nextErrorIndex:         newConfigValue(1),
		trackSimilarErrors:     newConfigValue(true),
	}
}

// Exported methods

// Set configures the behavior of various aspects of the wrappederror package at
// once.
//
// Configure individual properties using the Set(Property) methods.
func (c *Configuration) Set(
	captureCaller bool,
	captureProcess bool,
	captureSourceFragments bool,
	sourceFragmentRadius int,
	ignoreBreakpoints bool,
	nextErrorIndex int,
	trackSimilarErrors bool,
	marshalMinimalJSON bool,
) {
	c.SetCaptureCaller(captureCaller)
	c.SetCaptureProcess(captureProcess)
	c.SetCaptureSourceFragments(captureSourceFragments)
	c.SetSourceFragmentRadius(sourceFragmentRadius)
	c.SetIgnoreBreakpoints(ignoreBreakpoints)
	c.SetNextErrorIndex(nextErrorIndex)
	c.SetTrackSimilarErrors(trackSimilarErrors)
	c.SetMarshalMinimalJSON(marshalMinimalJSON)
}

// Error interface values

// SetCaptureCaller sets a flag to determine if new errors capture their caller
// information.
func (c *Configuration) SetCaptureCaller(capture bool) {
	c.captureCaller.set(capture)
}

// CaptureCaller gets a boolean that indicates whether or not new errors capture
// their caller information.
func (c *Configuration) CaptureCaller() bool {
	return c.captureCaller.get().(bool)
}

// SetCaptureProcess sets a flag to determine if new errors capture their
// process information.
func (c *Configuration) SetCaptureProcess(capture bool) {
	c.captureProcess.set(capture)
}

// CaptureProcess gets a boolean that indicates whether or not new errors
// capture their process information.
func (c *Configuration) CaptureProcess() bool {
	return c.captureProcess.get().(bool)
}

// SetMarshalMinimalJSON determines how errors are marshaled in to JSON. When
// this value is true, a smaller JSON object is created without size-inflating
// data like stack traces and source fragments.
func (c *Configuration) SetMarshalMinimalJSON(minimal bool) {
	c.marshalMinimalJSON.set(minimal)
}

// MarshalMinimalJSON gets a boolean that indicates whether or not errors will
// be marshaled in to a minimal version of JSON.
func (c *Configuration) MarshalMinimalJSON() bool {
	return c.marshalMinimalJSON.get().(bool)
}

// Caller interface values

// CaptureSourceFragments gets a boolean that indicates whether or not new
// errors capture source fragments.
func (c *Configuration) CaptureSourceFragments() bool {
	return c.captureSourceFragments.get().(bool)
}

// SetCaptureSourceFragments sets a flag to determine whether or not new errors
// capture source fragments.
func (c *Configuration) SetCaptureSourceFragments(capture bool) {
	c.captureSourceFragments.set(capture)
}

// SetSourceFragmentRadius sets the radius of the source fragment obtained from
// source files at the line that the caller was created on.
func (c *Configuration) SetSourceFragmentRadius(radius int) {
	c.sourceFragmentRadius.set(radius)
}

// SourceFragmentRadius gets the radius of source fragments obtained from source
// files.
func (c *Configuration) SourceFragmentRadius() int {
	return c.sourceFragmentRadius.get().(int)
}

// Process interface values

// SetIgnoreBreakpoints tells all calls to `Break` on `Process` types to either
// handle or ignore invocations.
func (c *Configuration) SetIgnoreBreakpoints(ignore bool) {
	c.ignoreBreakpoints.set(ignore)
}

// IgnoreBreakpoints returns whether or not calls to `Break` on `Process` types
// will be ignored. This value defaults to true.
func (c *Configuration) IgnoreBreakpoints() bool {
	return c.ignoreBreakpoints.get().(bool)
}

// Metadata interface values

// SetNextErrorIndex sets the next error index that will be used when creating
// an error.
func (c *Configuration) SetNextErrorIndex(index int) {
	c.nextErrorIndex.set(index)
}

// NextErrorIndex gets the next error index that will be used when creating an
// error.
func (c *Configuration) NextErrorIndex() int {
	return c.nextErrorIndex.get().(int)
}

// SetTrackSimilarErrors enables or prohibits similar error tracking.
func (c *Configuration) SetTrackSimilarErrors(track bool) {
	c.trackSimilarErrors.set(track)
}

// TrackSimilarErrors returns whether or not similar errors are being tracked.
func (c *Configuration) TrackSimilarErrors() bool {
	return c.trackSimilarErrors.get().(bool)
}

// Non-exported methods

// getAndIncrementNextErrorIndex gets the next error index and increments the
// value.
func (c *Configuration) getAndIncrementNextErrorIndex() int {
	c.nextErrorIndex.mutex.Lock()
	v := c.nextErrorIndex.value.(int)
	c.nextErrorIndex.value = c.nextErrorIndex.value.(int) + 1
	c.nextErrorIndex.mutex.Unlock()
	return v
}
