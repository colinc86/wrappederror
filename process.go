package wrappederror

import (
	"fmt"
	"runtime"
)

// Values to use when we can't get components of the process.
const (
	processRoutinesNumberUnknown int = -1
	processCPUsNumberUnknown     int = -1
	processCGONumberUnknown      int = 0
)

// Process types contain process information at their time of creation.
type Process struct {

	// The number of go routines.
	Routines int `json:"goroutines"`

	// The number of available logical CPUs.
	CPUs int `json:"cpus"`

	// The number of cgo calls made by the process.
	CGO int `json:"cgos"`

	// Memory statistics about the process.
	Memory *runtime.MemStats `json:"memory,omitempty"`
}

// Initializers

// newProcess creates and returns a new process.
func newProcess() *Process {
	ms := new(runtime.MemStats)
	runtime.ReadMemStats(ms)

	return &Process{
		runtime.NumGoroutine(),
		runtime.NumCPU(),
		int(runtime.NumCgoCall()),
		ms,
	}
}

// Exported methods

// Break executes a breakpoint trap if the configuration's ignore breakpoints
// value is false.
func (p Process) Break() {
	if packageState.config.IgnoreBreakpoints() {
		return
	}

	runtime.Breakpoint()
}

// Stringer interface methods

func (p Process) String() string {
	return fmt.Sprintf(
		"goroutines: %d, cpus: %d, cgos: %d",
		p.Routines,
		p.CPUs,
		p.CGO,
	)
}
