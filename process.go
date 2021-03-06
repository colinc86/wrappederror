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

// Process implementation

// Values to use when we can't get components of the process.
const (
	processRoutinesNumberUnknown int = -1
	processCPUsNumberUnknown     int = -1
	processCGONumberUnknown      int = 0
)

// A type containing process information.
type process struct {
	numRoutines int
	numCPUs     int
	numCGO      int
	memStats    *runtime.MemStats
}

// Initializers

// newProcess creates a new process with the specified components.
func newProcess(
	numRoutines int,
	numCPUs int,
	numCGO int,
	memStats *runtime.MemStats,
) *process {
	return &process{
		numRoutines: numRoutines,
		numCPUs:     numCPUs,
		numCGO:      numCGO,
		memStats:    memStats,
	}
}

// Methods

// currentProcess gets the current process.
func currentProcess() *process {
	ms := new(runtime.MemStats)
	runtime.ReadMemStats(ms)

	return newProcess(
		runtime.NumGoroutine(),
		runtime.NumCPU(),
		int(runtime.NumCgoCall()),
		ms,
	)
}

// String interface methods

func (p process) String() string {
	return fmt.Sprintf(
		"goroutines: %d, cpus: %d, cgos: %d",
		p.numRoutines,
		p.numCPUs,
		p.numCGO,
	)
}

// Process interface methods

func (p process) Routines() int {
	return p.numRoutines
}

func (p process) CPUs() int {
	return p.numCPUs
}

func (p process) CGO() int {
	return p.numCGO
}

func (p process) Memory() *runtime.MemStats {
	return p.memStats
}

func (p process) Break() {
	runtime.Breakpoint()
}
