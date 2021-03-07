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
	numRoutines int
	numCPUs     int
	numCGO      int
	memStats    *runtime.MemStats
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
		numRoutines: numRoutines,
		numCPUs:     numCPUs,
		numCGO:      numCGO,
		memStats:    memStats,
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
		p.numRoutines,
		p.numCPUs,
		p.numCGO,
	)
}

// Process interface methods

func (p wProcess) Routines() int {
	return p.numRoutines
}

func (p wProcess) CPUs() int {
	return p.numCPUs
}

func (p wProcess) CGO() int {
	return p.numCGO
}

func (p wProcess) Memory() *runtime.MemStats {
	return p.memStats
}

func (p wProcess) Break() {
	if IgnoreBreakpoints() {
		return
	}

	runtime.Breakpoint()
}
