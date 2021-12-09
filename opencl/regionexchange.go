package opencl

// Region paired exchange interaction

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl/cl"
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

	var BxPtr unsafe.Pointer
	var ByPtr unsafe.Pointer
	var BzPtr unsafe.Pointer
	var mxPtr unsafe.Pointer
	var myPtr unsafe.Pointer
	var mzPtr unsafe.Pointer
	var MsPtr unsafe.Pointer
	var eventsList []*cl.Event

	if B.DevPtr(X) != nil {
		BxPtr = B.DevPtr(X)
		eventsList = append(eventsList, B.GetEvent(X))
	}
	if B.DevPtr(Y) != nil {
		ByPtr = B.DevPtr(Y)
		eventsList = append(eventsList, B.GetEvent(Y))
	}
	if B.DevPtr(Z) != nil {
		BzPtr = B.DevPtr(Z)
		eventsList = append(eventsList, B.GetEvent(Z))
	}
	if m.DevPtr(X) != nil {
		mxPtr = m.DevPtr(X)
		eventsList = append(eventsList, m.GetEvent(X))
	}
	if m.DevPtr(Y) != nil {
		myPtr = m.DevPtr(Y)
		eventsList = append(eventsList, m.GetEvent(Y))
	}
	if m.DevPtr(Z) != nil {
		mzPtr = m.DevPtr(Z)
		eventsList = append(eventsList, m.GetEvent(Z))
	}
	if Msat.DevPtr(0) != nil {
		MsPtr = Msat.DevPtr(0)
	}

	sig_eff := sig * float32(cellwgt)
	sig2_eff := sig2 * float32(cellwgt)

	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_tworegionexchange_field_async(BxPtr, ByPtr, BzPtr,
		mxPtr, myPtr, mzPtr,
		MsPtr, Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, sig_eff, sig2_eff, N[X], N[Y], N[Z], cfg,
		eventsList)

	if B.DevPtr(X) != nil {
		B.SetEvent(X, event)
	}
	if B.DevPtr(Y) != nil {
		B.SetEvent(Y, event)
	}
	if B.DevPtr(Z) != nil {
		B.SetEvent(Z, event)
	}
	if m.DevPtr(X) != nil {
		m.SetEvent(X, event)
	}
	if m.DevPtr(Y) != nil {
		m.SetEvent(Y, event)
	}
	if m.DevPtr(Z) != nil {
		m.SetEvent(Z, event)
	}

	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
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

	var EdensPtr unsafe.Pointer
	var mxPtr unsafe.Pointer
	var myPtr unsafe.Pointer
	var mzPtr unsafe.Pointer
	var MsPtr unsafe.Pointer
	var eventsList []*cl.Event

	if Edens.DevPtr(0) != nil {
		EdensPtr = Edens.DevPtr(0)
		eventsList = append(eventsList, Edens.GetEvent(0))
	}
	if m.DevPtr(X) != nil {
		mxPtr = m.DevPtr(X)
		eventsList = append(eventsList, m.GetEvent(X))
	}
	if m.DevPtr(Y) != nil {
		myPtr = m.DevPtr(Y)
		eventsList = append(eventsList, m.GetEvent(Y))
	}
	if m.DevPtr(Z) != nil {
		mzPtr = m.DevPtr(Z)
		eventsList = append(eventsList, m.GetEvent(Z))
	}
	if Msat.DevPtr(0) != nil {
		MsPtr = Msat.DevPtr(0)
	}

	sig_eff := sig * float32(cellwgt)
	sig2_eff := sig2 * float32(cellwgt)

	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_tworegionexchange_edens_async(EdensPtr,
		mxPtr, myPtr, mzPtr,
		MsPtr, Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, sig_eff, sig2_eff, N[X], N[Y], N[Z], cfg,
		eventsList)

	if Edens.DevPtr(0) != nil {
		Edens.SetEvent(0, event)
	}
	if m.DevPtr(X) != nil {
		m.SetEvent(X, event)
	}
	if m.DevPtr(Y) != nil {
		m.SetEvent(Y, event)
	}
	if m.DevPtr(Z) != nil {
		m.SetEvent(Z, event)
	}

	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in addtworegionexchange_edens: %+v", err)
	}
}
