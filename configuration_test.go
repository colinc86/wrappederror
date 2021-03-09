package wrappederror

import "testing"

func TestNewConfiguration(t *testing.T) {
	c := newConfiguration()
	t.Run("Capture caller", func(t *testing.T) { testConfigurationValue(t, c.captureCaller, true) })
	t.Run("Capture process", func(t *testing.T) { testConfigurationValue(t, c.captureProcess, true) })
	t.Run("Capture source fragments", func(t *testing.T) { testConfigurationValue(t, c.captureSourceFragments, true) })
	t.Run("Source fragment radius", func(t *testing.T) { testConfigurationValue(t, c.sourceFragmentRadius, 2) })
	t.Run("Ignore breakpoints", func(t *testing.T) { testConfigurationValue(t, c.ignoreBreakpoints, true) })
	t.Run("Next error index", func(t *testing.T) { testConfigurationValue(t, c.nextErrorIndex, 1) })
	t.Run("Track similar errors", func(t *testing.T) { testConfigurationValue(t, c.trackSimilarErrors, true) })
	t.Run("Marshal minimal JSON", func(t *testing.T) { testConfigurationValue(t, c.marshalMinimalJSON, true) })
}

func TestConfigurationSet(t *testing.T) {
	c := newConfiguration()
	c.Set(false, false, false, 0, false, 0, false, false)
	t.Run("Capture caller", func(t *testing.T) { testConfigurationValue(t, c.captureCaller, false) })
	t.Run("Capture process", func(t *testing.T) { testConfigurationValue(t, c.captureProcess, false) })
	t.Run("Capture source fragments", func(t *testing.T) { testConfigurationValue(t, c.captureSourceFragments, false) })
	t.Run("Source fragment radius", func(t *testing.T) { testConfigurationValue(t, c.sourceFragmentRadius, 0) })
	t.Run("Ignore breakpoints", func(t *testing.T) { testConfigurationValue(t, c.ignoreBreakpoints, false) })
	t.Run("Next error index", func(t *testing.T) { testConfigurationValue(t, c.nextErrorIndex, 0) })
	t.Run("Track similar errors", func(t *testing.T) { testConfigurationValue(t, c.trackSimilarErrors, false) })
	t.Run("Marshal minimal JSON", func(t *testing.T) { testConfigurationValue(t, c.marshalMinimalJSON, false) })
}

func testConfigurationValue(t *testing.T, v *safeValue, ev interface{}) {
	if v.get() != ev {
		t.Errorf("Expected %v but received %v.\n", ev, v.get())
	}
}

func TestConfigurationGetAndIncrementNextErrorIndex(t *testing.T) {
	c := newConfiguration()
	t.Run("Default error index", func(t *testing.T) { testConfigurationNextErrorIndex(t, c, 1) })
	t.Run("Next error index", func(t *testing.T) { testConfigurationNextErrorIndex(t, c, 2) })
	t.Run("Third error index", func(t *testing.T) { testConfigurationNextErrorIndex(t, c, 3) })
	t.Run("Reset error index", func(t *testing.T) {
		c.SetNextErrorIndex(1)
		testConfigurationNextErrorIndex(t, c, 1)
	})
}

func testConfigurationNextErrorIndex(t *testing.T, c *Configuration, i int) {
	if c.getAndIncrementNextErrorIndex() != i {
		t.Errorf("Expected %d but recived %+v.\n", i, c.NextErrorIndex())
	}
}
