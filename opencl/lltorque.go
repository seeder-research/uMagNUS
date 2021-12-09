package opencl

import (
	"fmt"

	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/cl"
)

// Landau-Lifshitz torque divided by gamma0:
// 	- 1/(1+α²) [ m x B +  α m x (m x B) ]
// 	torque in Tesla
// 	m normalized
// 	B in Tesla
// see lltorque.cl
func LLTorque(torque, m, B *data.Slice, alpha MSlice) {
	N := torque.Len()
	cfg := make1DConf(N)

	event := k_lltorque2_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		alpha.DevPtr(0), alpha.Mul(0), N, cfg,
		[](*cl.Event){torque.GetEvent(X), torque.GetEvent(Y), torque.GetEvent(Z),
			m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z),
			B.GetEvent(X), B.GetEvent(Y), B.GetEvent(Z),
			alpha.GetEvent(0)})
	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	B.SetEvent(X, event)
	B.SetEvent(Y, event)
	B.SetEvent(Z, event)
	alpha.SetEvent(0, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in lltorque: %+v \n", err)
	}
}

// Landau-Lifshitz torque with precession disabled.
// Used by engine.Relax().
func LLNoPrecess(torque, m, B *data.Slice) {
	N := torque.Len()
	cfg := make1DConf(N)

	event := k_llnoprecess_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z), N, cfg,
		[](*cl.Event){torque.GetEvent(X), torque.GetEvent(Y), torque.GetEvent(Z),
			m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z),
			B.GetEvent(X), B.GetEvent(Y), B.GetEvent(Z)})
	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	B.SetEvent(X, event)
	B.SetEvent(Y, event)
	B.SetEvent(Z, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in llnoprecess: %+v \n", err)
	}
}
