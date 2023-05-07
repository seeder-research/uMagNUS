package opencl

import (
	"context"
	"sync"
)

var (
	routinesCtx context.Context // Context for signaling all goroutines to terminate
	routinesFcn chan func() // Channel for passing functions to goroutines for execution
	routinesWG sync.WaitGroup // WaitGroup for tracking if all routines are free
	endRoutines context.CancelFunc
	routinesInit = false
)

// Initialization of parameters for routines
func initRoutinesWithContext(sz int) {
	routinesCtx, endRoutines = context.WithCancel(context.Background())
	routinesFcn = make(chan func(), sz)
	routinesInit = true
}

// Routines will run this function where they continually wait to receive functions
// in the channel unless the context is signalling the routine to return
func supportRoutine(ctx context.Context, fcn <-chan func(), wg *sync.WaitGroup) {
	for {
		select {
		case f := <- fcn:
			wg.Add(1)
			f()
			wg.Done()
		case <-ctx.Done():
			break
		}
	}
	return
}

// Start the routines
func startRoutines() {
	if routinesInit {
		for i := 0; i < cap(routinesFcn); i++ {
			go supportRoutine(routinesCtx, routinesFcn, &routinesWG)
		}
	}
}

func submitFcnToRoutine(fcn func()) {
	if routinesInit {
		routinesFcn <- fcn
	}
}

// Terminate all routines
func terminateRoutines() {
	if routinesInit {
		routinesWG.Wait() // Wait for all routines to be in wait state
		endRoutines() // Send signal to end routines
		close(routinesFcn)
	}
}
