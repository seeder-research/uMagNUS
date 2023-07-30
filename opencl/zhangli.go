package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Add Zhang-Li ST torque (Tesla) to torque.
// see zhangli2.cl
func AddZhangLiTorque(torque, m *data.Slice, Msat, J, alpha, xi, pol MSlice, mesh *data.Mesh, q *cl.CommandQueue, ewl []*cl.Event) {
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
		ewl, q)

	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)

	glist := []GSlice{m}
	if J.GetSlicePtr() != nil {
		glist = append(glist, J)
	}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	if alpha.GetSlicePtr() != nil {
		glist = append(glist, alpha)
	}
	if xi.GetSlicePtr() != nil {
		glist = append(glist, xi)
	}
	if pol.GetSlicePtr() != nil {
		glist = append(glist, pol)
	}
	InsertEventIntoGSlices(event, glist)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents failed in addzhanglitorque: %+v \n", err)
		}
		WaitAndUpdateDataSliceEvents(event, glist, false)
	}

	return
}
