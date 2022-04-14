/*
engine does the simulation bookkeeping, I/O and GUI.

space-dependence:
value: space-independent
param: region-dependent parameter (always input)
field: fully space-dependent field

TODO: godoc everything

*/
package engine64

import (
	"os"
	"runtime"
	"sync"
	"time"

	opencl "github.com/seeder-research/uMagNUS/opencl64"
	timer "github.com/seeder-research/uMagNUS/timer"
)

const VERSION = "uMagNUS 2.2.1"

var UNAME = VERSION + " " + runtime.GOOS + "_" + runtime.GOARCH + " " + runtime.Version() + " (" + runtime.Compiler + ")"

var StartTime = time.Now()

var (
	busyLock sync.Mutex
	busy     bool // are we so busy we can't respond from run loop? (e.g. calc kernel)
)

// We set SetBusy(true) when the simulation is too busy too accept GUI input on Inject channel.
// E.g. during kernel init.
func SetBusy(b bool) {
	busyLock.Lock()
	defer busyLock.Unlock()
	busy = b
}

func GetBusy() bool {
	busyLock.Lock()
	defer busyLock.Unlock()
	return busy
}

// Cleanly exits the simulation, assuring all output is flushed.
func Close() {
	drainOutput()
	LogUsedRefs()
	Table.flush()
	if logfile != nil {
		logfile.Close()
	}
	if bibfile != nil {
		bibfile.Close()
	}
	if *Flag_sync {
		timer.Print(os.Stdout)
	}
	opencl.ReleaseAndClean()
}
