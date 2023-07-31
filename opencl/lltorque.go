package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Landau-Lifshitz torque divided by gamma0:
//   - 1/(1+α²) [ m x B +  α m x (m x B) ]
//     torque in Tesla
//     m normalized
//     B in Tesla
//
// see lltorque.cl
func LLTorque(torque, m, B *data.Slice, alpha MSlice, q *cl.CommandQueue, ewl []*cl.Event) {
	N := torque.Len()
	cfg := make1DConf(N)

	// Launch kernel
	event := k_lltorque2_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		alpha.DevPtr(0), alpha.Mul(0), N, cfg,
		ewl, q)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in lltorque: %+v \n", err)
		}
	}

	return
}

// Landau-Lifshitz torque with precession disabled.
// Used by engine.Relax().
func LLNoPrecess(torque, m, B *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) {
	N := torque.Len()
	cfg := make1DConf(N)

	// Launch kernel
	event := k_llnoprecess_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z), N, cfg,
		ewl, q)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in llnoprecess: %+v \n", err)
		}
	}

	return
}
