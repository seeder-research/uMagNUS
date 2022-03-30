package opencl64

// Region paired spin torque calculations

import (
	"fmt"
	"math"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
)

func AddRegionSpinTorque(torque, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, J, alpha, pfix, pfree, λfix, λfree, ε_prime float64, mesh *data.Mesh) {
	c := mesh.CellSize()
	dX := float64(sX) * c[X]
	dY := float64(sY) * c[Y]
	dZ := float64(sZ) * c[Z]

	distsq := dX*dX + dY*dY + dZ*dZ
	cellwgt := math.Abs(dX*c[X]) + math.Abs(dY*c[Y]) + math.Abs(dZ*c[Z])
	if cellwgt > 0.0 {
		cellwgt = math.Sqrt(distsq) / cellwgt
	}

	N := mesh.Size()
	cfg := make3DConf(N)

	eventsList := []*cl.Event{torque.GetEvent(X), torque.GetEvent(Y), torque.GetEvent(Z), m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z)}

	event := k_addtworegionoommfslonczewskitorque_async(torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, N[X], N[Y], N[Z],
		J, alpha, pfix, pfree, λfix, λfree, ε_prime, float64(cellwgt),
		cfg,
		eventsList)

	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addtworegionoommfslonczewskitorque: %+v", err)
	}
}
