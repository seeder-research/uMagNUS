package opencl64

// Region paired spin torque calculations

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/seeder-research/uMagNUS/cl"
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

	var TxPtr unsafe.Pointer
	var TyPtr unsafe.Pointer
	var TzPtr unsafe.Pointer
	var mxPtr unsafe.Pointer
	var myPtr unsafe.Pointer
	var mzPtr unsafe.Pointer
	var MsPtr unsafe.Pointer
	var eventsList []*cl.Event

	if torque.DevPtr(X) != nil {
		TxPtr = torque.DevPtr(X)
		eventsList = append(eventsList, torque.GetEvent(X))
	}
	if torque.DevPtr(Y) != nil {
		TyPtr = torque.DevPtr(Y)
		eventsList = append(eventsList, torque.GetEvent(Y))
	}
	if torque.DevPtr(Z) != nil {
		TzPtr = torque.DevPtr(Z)
		eventsList = append(eventsList, torque.GetEvent(Z))
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

	if len(eventsList) == 0 {
		eventsList = nil
	}

	event := k_addtworegionoommfslonczewskitorque_async(TxPtr, TyPtr, TzPtr,
		mxPtr, myPtr, mzPtr,
		MsPtr, Msat.Mul(0),
		regions.Ptr, regionA, regionB,
		sX, sY, sZ, N[X], N[Y], N[Z],
		J, alpha, pfix, pfree, λfix, λfree, ε_prime, float64(cellwgt),
		cfg,
		eventsList)

	if torque.DevPtr(X) != nil {
		torque.SetEvent(X, event)
	}
	if torque.DevPtr(Y) != nil {
		torque.SetEvent(Y, event)
	}
	if torque.DevPtr(Z) != nil {
		torque.SetEvent(Z, event)
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
		fmt.Printf("WaitForEvents failed in addtworegionoommfslonczewskitorque: %+v", err)
	}
}
