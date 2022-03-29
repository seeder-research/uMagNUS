package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
)

// Add Slonczewski ST torque to torque (Tesla).
// see slonczewski.cl
func AddSlonczewskiTorque2(torque, m *data.Slice, Msat, J, fixedP, alpha, pol, λ, ε_prime MSlice, thickness MSlice, flp float64, mesh *data.Mesh) {
	N := torque.Len()
	cfg := make1DConf(N)
	meshThickness := mesh.WorldSize()[Z]

	event := k_addslonczewskitorque2_async(
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		J.DevPtr(Z), J.Mul(Z),
		fixedP.DevPtr(X), fixedP.Mul(X),
		fixedP.DevPtr(Y), fixedP.Mul(Y),
		fixedP.DevPtr(Z), fixedP.Mul(Z),
		alpha.DevPtr(0), alpha.Mul(0),
		pol.DevPtr(0), pol.Mul(0),
		λ.DevPtr(0), λ.Mul(0),
		ε_prime.DevPtr(0), ε_prime.Mul(0),
		thickness.DevPtr(0), thickness.Mul(0),
		float64(meshThickness),
		float64(flp),
		N, cfg,
		[](*cl.Event){torque.GetEvent(X), torque.GetEvent(Y), torque.GetEvent(Z),
			m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z),
			fixedP.GetEvent(X), fixedP.GetEvent(Y), fixedP.GetEvent(Z),
			J.GetEvent(Z)})
	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	J.SetEvent(Z, event)
	fixedP.SetEvent(X, event)
	fixedP.SetEvent(Y, event)
	fixedP.SetEvent(Z, event)

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in addslonczewskitorque2: %+v \n", err)
	}
}
