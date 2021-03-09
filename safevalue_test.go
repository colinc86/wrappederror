package wrappederror

import "testing"

func TestNewSafeValue(t *testing.T) {
	sv := newSafeValue(1)
	t.Run("New safe value 0", func(t *testing.T) { testSafeValueInt(t, sv, 1) })
}

func testSafeValueInt(t *testing.T, sv *safeValue, i int) {
	if si, ok := sv.value.(int); !ok {
		t.Error("Expected integer value.")
	} else if si != i {
		t.Errorf("Expected %d but received %d.\n", i, si)
	}
}

func TestSafeValueSet(t *testing.T) {
	sv := newSafeValue(1)
	sv.set(2)
	t.Run("Safe value set 0", func(t *testing.T) {
		testSafeValueInt(t, sv, 2)
	})
	sv.set(true)
	t.Run("Safe value set 1", func(t *testing.T) {
		testSafeValueBool(t, sv, true)
	})
	sv.set(1)
	t.Run("Safe value set 2", func(t *testing.T) {
		testSafeValueInt(t, sv, 1)
	})
}

func testSafeValueBool(t *testing.T, sv *safeValue, b bool) {
	if sb, ok := sv.value.(bool); !ok {
		t.Error("Expected boolean value.")
	} else if sb != b {
		t.Errorf("Expected %t but received %t.\n", b, sb)
	}
}

func TestSafeValueGet(t *testing.T) {
	sv := newSafeValue(1)
	t.Run("Safe value get 0", func(t *testing.T) {
		testSafeValueIntGet(t, sv, 1)
	})
}

func testSafeValueIntGet(t *testing.T, sv *safeValue, i int) {
	if si, ok := sv.get().(int); !ok {
		t.Error("Expected integer value.")
	} else if si != i {
		t.Errorf("Expected %d but received %d.\n", i, si)
	}
}

func TestSaveValueTransform(t *testing.T) {
	sv := newSafeValue(1)
	sv.transform(func(v interface{}) interface{} {
		return 2
	})
	t.Run("Safe value get 0", func(t *testing.T) {
		testSafeValueInt(t, sv, 2)
	})
}
