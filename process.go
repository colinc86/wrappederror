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

// Exported functions

// SetIgnoreBreakpoints tells all calls to `Break` on `Process` types to either
// handle or ignore invocations.
func SetIgnoreBreakpoints(ignore bool) {
	ignoreBreakpointsMutex.Lock()
	ignoreBreakpoints = ignore
	ignoreBreakpointsMutex.Unlock()
}

// IgnoreBreakpoints returns whether or not calls to `Break` on `Process` types
// will be ignored. This value defaults to true.
func IgnoreBreakpoints() bool {
	ignoreBreakpointsMutex.RLock()
	defer ignoreBreakpointsMutex.RUnlock()
	return ignoreBreakpoints
}
