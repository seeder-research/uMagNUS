package opencl

// Region paired exchange interaction

import (
	"fmt"
	"math"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
)

// Add exchange field to Beff.
// 	m: normalized magnetization
// 	B: effective field in Tesla
func AddRegionExchangeField(B, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, sig, sig2 float32, mesh *data.Mesh) {
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

	eventsList := []*cl.Event{B.GetEvent(X), B.GetEvent(Y), B.GetEvent(Z),
		m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z)}

	sig_eff := sig * float32(cellwgt)
	sig2_eff := sig2 * float32(cellwgt)

	event := k_tworegionexchange_field_async(B.DevPtr(X), B.DevPtr(Y), B.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, sig_eff, sig2_eff, N[X], N[Y], N[Z], cfg,
		eventsList)

	B.SetEvent(X, event)
	B.SetEvent(Y, event)
	B.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addtworegionexchange_field: %+v", err)
	}
}

func AddRegionExchangeEdens(Edens, m *data.Slice, Msat MSlice, regions *Bytes, regionA, regionB uint8, sX, sY, sZ int, sig, sig2 float32, mesh *data.Mesh) {
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

	eventsList := []*cl.Event{Edens.GetEvent(0), m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z)}

	sig_eff := sig * float32(cellwgt)
	sig2_eff := sig2 * float32(cellwgt)

	event := k_tworegionexchange_edens_async(Edens.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		Msat.DevPtr(0), Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, sig_eff, sig2_eff, N[X], N[Y], N[Z], cfg,
		eventsList)

	Edens.SetEvent(0, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)

	if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in addtworegionexchange_edens: %+v", err)
	}
}
