package wrappederror

import (
	"testing"
)

func TestFindIndexes_1(t *testing.T) {
	f := newFormatter()
	ef := ".  .  .  "
	idx := f.findIndexes(ef, ".")

	if len(idx) != 3 {
		t.Fatalf("Expected 3 elements but received %d.\n", len(idx))
	}
}

func TestFindIndexes_2(t *testing.T) {
	f := newFormatter()
	ef := "asdf"
	idx := f.findIndexes(ef, "asdfj")

	if len(idx) != 0 {
		t.Fatalf("Expected 0 elements but received %d.\n", len(idx))
	}
}

func TestReplaceTokens_1(t *testing.T) {
	f := newFormatter()
	ef := "asdf"
	efc, tokens := f.replaceTokens(ef, nil)

	if efc != ef {
		t.Errorf("Expected %s but received %s.\n", ef, efc)
	}

	if len(tokens) > 0 {
		t.Errorf("Expected 0 tokens but received %d.\n", len(tokens))
	}
}

func TestReplaceTokens_2(t *testing.T) {
	f := newFormatter()
	ef := "${{LIN}}"
	efc, tokens := f.replaceTokens(ef, []int{0})

	if efc != "%d" {
		t.Errorf("Incorrect error format string: %s.\n", efc)
	}

	if len(tokens) != 1 {
		t.Errorf("Expected 1 token but received %d.\n", len(tokens))
	}
}

func TestReplaceTokens_3(t *testing.T) {
	f := newFormatter()
	ef := " ${{LIN}}"
	efc, tokens := f.replaceTokens(ef, []int{0})

	if efc != ef {
		t.Errorf("Incorrect error format string: %s.\n", efc)
	}

	if len(tokens) != 0 {
		t.Errorf("Expected 0 tokens but received %d.\n", len(tokens))
	}
}

func TestReplaceTokens_4(t *testing.T) {
	f := newFormatter()
	ef := "${{LIN}} "
	efc, tokens := f.replaceTokens(ef, []int{0})

	if efc != "%d " {
		t.Errorf("Incorrect error format string: %s.\n", efc)
	}

	if len(tokens) != 1 {
		t.Errorf("Expected 1 token but received %d.\n", len(tokens))
	}
}

func TestReplaceTokens_5(t *testing.T) {
	f := newFormatter()
	ef := "${{LIN}}${{LIN}}"
	efc, tokens := f.replaceTokens(ef, []int{0, 8})

	if efc != "%d%d" {
		t.Errorf("Incorrect error format string: %s.\n", efc)
	}

	if len(tokens) != 2 {
		t.Errorf("Expected 1 token but received %d.\n", len(tokens))
	}
}

func TestReplaceToken_1(t *testing.T) {
	f := newFormatter()
	ef := "${{FIL}}"
	efc, token := f.replaceToken(ef, 0)

	if efc != "%s" {
		t.Errorf("Incorrect error format string: %s\n", efc)
	}

	if token != ErrorFormatTokenFile {
		t.Errorf("Incorrect token: %s\n", token)
	}
}

func TestReplaceToken_2(t *testing.T) {
	f := newFormatter()
	ef := "${{FIL}}  "
	efc, token := f.replaceToken(ef, 1)

	if efc != ef {
		t.Errorf("Incorrect error format string: %s\n", efc)
	}

	if token != errorFormatTokenNone {
		t.Errorf("Incorrect token: %s\n", token)
	}
}

func TestReplaceToken_3(t *testing.T) {
	f := newFormatter()
	ef := "${{FIL}}"
	efc, token := f.replaceToken(ef, 1)

	if efc != ef {
		t.Errorf("Incorrect error format string: %s\n", efc)
	}

	if token != errorFormatTokenNone {
		t.Errorf("Incorrect token: %s\n", token)
	}
}
