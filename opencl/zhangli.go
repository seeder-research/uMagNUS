package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Add Zhang-Li ST torque (Tesla) to torque.
// see zhangli2.cl
func AddZhangLiTorque(torque, m *data.Slice, Msat, J, alpha, xi, pol MSlice, mesh *data.Mesh) {
	c := mesh.CellSize()
	N := mesh.Size()
	cfg := make3DConf(N)

	event := k_addzhanglitorque2_async(
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		J.DevPtr(X), J.Mul(X),
		J.DevPtr(Y), J.Mul(Y),
		J.DevPtr(Z), J.Mul(Z),
		alpha.DevPtr(0), alpha.Mul(0),
		xi.DevPtr(0), xi.Mul(0),
		pol.DevPtr(0), pol.Mul(0),
		float32(c[X]), float32(c[Y]), float32(c[Z]),
		N[X], N[Y], N[Z], mesh.PBC_code(), cfg,
		[](*cl.Event){torque.GetEvent(X), torque.GetEvent(Y), torque.GetEvent(Z),
			m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z),
			J.GetEvent(X), J.GetEvent(Y), J.GetEvent(Z)})
	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	J.SetEvent(X, event)
	J.SetEvent(Y, event)
	J.SetEvent(Z, event)
	Msat.SetEvent(0, event)
	alpha.SetEvent(0, event)
	xi.SetEvent(0, event)
	pol.SetEvent(0, event)

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in addzhanglitorque: %+v \n", err)
	}
}
