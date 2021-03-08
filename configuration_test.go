package wrappederror

import "testing"

func TestConfigurationSet(t *testing.T) {
	c := newConfiguration()
	c.Set(false, false, 0, false, 0, false, false)

	if c.captureCaller.get().(bool) != false {
		t.Error("Capture caller not set.")
	}

	if c.captureProcess.get().(bool) != false {
		t.Error("Capture process not set.")
	}

	if c.sourceFragmentRadius.get().(int) != 0 {
		t.Error("Source fragment radius not set.")
	}

	if c.ignoreBreakpoints.get().(bool) != false {
		t.Error("Ignore breakpoints not set.")
	}

	if c.nextErrorIndex.get().(int) != 0 {
		t.Error("Next error index not set.")
	}

	if c.trackSimilarErrors.get().(bool) != false {
		t.Error("Track similar errors not set.")
	}

	if c.marshalMinimalJSON.get().(bool) != false {
		t.Error("Marshal minimal JSON not set.")
	}
}

func TestConfigurationIgnoreBreakpoints(t *testing.T) {
	c := newConfiguration()
	if c.IgnoreBreakpoints() != true {
		t.Error("Unexpected ignore breakpoints value.")
	}
}

func TestConfigurationNextErrorIndex(t *testing.T) {
	packageState.reset()
	if packageState.config.NextErrorIndex() != 1 {
		t.Errorf("Expected next error index value 1: %d.\n", packageState.config.NextErrorIndex())
	}

	_ = New(nil, "test")
	if packageState.config.NextErrorIndex() != 2 {
		t.Errorf("Expected next error index value 2: %d.\n", packageState.config.NextErrorIndex())
	}
}
