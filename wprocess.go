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

// A type containing process information.
type wProcess struct {
	NumRoutines int               `json:"goroutines"`
	NumCPUs     int               `json:"cpus"`
	NumCGO      int               `json:"cgos"`
	MemStats    *runtime.MemStats `json:"memory,omitempty"`
}

// Initializers

// newWProcess creates a new process with the specified components.
func newWProcess(
	numRoutines int,
	numCPUs int,
	numCGO int,
	memStats *runtime.MemStats,
) *wProcess {
	return &wProcess{
		NumRoutines: numRoutines,
		NumCPUs:     numCPUs,
		NumCGO:      numCGO,
		MemStats:    memStats,
	}
}

// Methods

// currentProcess gets the current process.
func currentProcess() *wProcess {
	ms := new(runtime.MemStats)
	runtime.ReadMemStats(ms)

	return newWProcess(
		runtime.NumGoroutine(),
		runtime.NumCPU(),
		int(runtime.NumCgoCall()),
		ms,
	)
}

// Stringer interface methods

func (p wProcess) String() string {
	return fmt.Sprintf(
		"goroutines: %d, cpus: %d, cgos: %d",
		p.NumRoutines,
		p.NumCPUs,
		p.NumCGO,
	)
}

// Process interface methods

func (p wProcess) Routines() int {
	return p.NumRoutines
}

func (p wProcess) CPUs() int {
	return p.NumCPUs
}

func (p wProcess) CGO() int {
	return p.NumCGO
}

func (p wProcess) Memory() *runtime.MemStats {
	return p.MemStats
}

func (p wProcess) Break() {
	if IgnoreBreakpoints() {
		return
	}

	runtime.Breakpoint()
}
