package opencl

import (
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
func LLTorque(torque, m, B *data.Slice, alpha MSlice, queue *cl.CommandQueue, events []*cl.Event) {
	N := torque.Len()
	cfg := make1DConf(N)

	k_lltorque2_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		alpha.DevPtr(0), alpha.Mul(0), N, cfg,
		queue, events)
}

// Landau-Lifshitz torque with precession disabled.
// Used by engine.Relax().
func LLNoPrecess(torque, m, B *data.Slice, queue *cl.CommandQueue, events []*cl.Event) {
	N := torque.Len()
	cfg := make1DConf(N)

	k_llnoprecess_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z), N, cfg,
		queue, events)
}
