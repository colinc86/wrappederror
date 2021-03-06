package wrappederror

import (
	"fmt"
	"runtime"
)

// Process types contain process information.
type Process interface {
	fmt.Stringer

	// The number of go routines when the process was created.
	Routines() int

	// The number of available logical CPUs.
	CPUs() int

	// The number of cgo calls made by the process.
	CGO() int

	// The memory statistics when the process was created.
	Memory() *runtime.MemStats

	// Break executes a breakpoint trap at the point that this method is called.
	Break()
}
