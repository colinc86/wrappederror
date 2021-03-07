package wrappederror

import (
	"fmt"
	"testing"
)

func TestFormatter(t *testing.T) {
	e1 := New(nil, "error 1")
	e2 := New(e1, "error 2")

	ef := fmt.Sprintf(
		"Error #%s at %s (%s:%s): %s",
		ErrorFormatTokenIndex,
		ErrorFormatTokenTime,
		ErrorFormatTokenFile,
		ErrorFormatTokenLine,
		ErrorFormatTokenChain,
	)

	fmt.Println(e2.Format(ef))
}
