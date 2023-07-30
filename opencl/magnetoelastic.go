package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Add magneto-elasticit coupling field to the effective field.
// see magnetoelasticfield.cl
func AddMagnetoelasticField(Beff, m *data.Slice, exx, eyy, ezz, exy, exz, eyz, B1, B2, Msat MSlice, q *cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(Beff.Size() == m.Size())
	util.Argument(Beff.Size() == exx.Size())
	util.Argument(Beff.Size() == eyy.Size())
	util.Argument(Beff.Size() == ezz.Size())
	util.Argument(Beff.Size() == exy.Size())
	util.Argument(Beff.Size() == exz.Size())
	util.Argument(Beff.Size() == eyz.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	event := k_addmagnetoelasticfield_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		exx.DevPtr(0), exx.Mul(0), eyy.DevPtr(0), eyy.Mul(0), ezz.DevPtr(0), ezz.Mul(0),
		exy.DevPtr(0), exy.Mul(0), exz.DevPtr(0), exz.Mul(0), eyz.DevPtr(0), eyz.Mul(0),
		B1.DevPtr(0), B1.Mul(0), B2.DevPtr(0), B2.Mul(0),
		Msat.DevPtr(0), Msat.Mul(0),
		N, cfg, ewl,
		q)

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)

	glist := []GSlice{m}
	if exx.GetSlicePtr() != nil {
		glist = append(glist, exx)
	}
	if eyy.GetSlicePtr() != nil {
		glist = append(glist, eyy)
	}
	if ezz.GetSlicePtr() != nil {
		glist = append(glist, ezz)
	}
	if exy.GetSlicePtr() != nil {
		glist = append(glist, exy)
	}
	if exz.GetSlicePtr() != nil {
		glist = append(glist, exz)
	}
	if eyz.GetSlicePtr() != nil {
		glist = append(glist, eyz)
	}
	if B1.GetSlicePtr() != nil {
		glist = append(glist, B1)
	}
	if B2.GetSlicePtr() != nil {
		glist = append(glist, B2)
	}
	if Msat.GetSlicePtr() != nil {
		glist = append(glist, Msat)
	}
	InsertEventIntoGSlices(event, glist)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in addmagnetoelasticfield failed: %+v \n", err)
		}
	}

	return
}

// Calculate magneto-elasticit force density
// see magnetoelasticforce.cl
func GetMagnetoelasticForceDensity(out, m *data.Slice, B1, B2 MSlice, mesh *data.Mesh, q *cl.CommandQueue, ewl []*cl.Event) {
	util.Argument(out.Size() == m.Size())

	cellsize := mesh.CellSize()
	N := mesh.Size()
	cfg := make3DConf(N)

	rcsx := float32(1.0 / cellsize[X])
	rcsy := float32(1.0 / cellsize[Y])
	rcsz := float32(1.0 / cellsize[Z])

	event := k_getmagnetoelasticforce_async(out.DevPtr(X), out.DevPtr(Y), out.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B1.DevPtr(0), B1.Mul(0), B2.DevPtr(0), B2.Mul(0),
		rcsx, rcsy, rcsz,
		N[X], N[Y], N[Z],
		mesh.PBC_code(), cfg, ewl,
		q)

	out.SetEvent(X, event)
	out.SetEvent(Y, event)
	out.SetEvent(Z, event)

	glist := []GSlice{m}
	if B1.GetSlicePtr() != nil {
		glist = append(glist, B1)
	}
	if B2.GetSlicePtr() != nil {
		glist = append(glist, B2)
	}
	InsertEventIntoGSlices(event, glist)

	if Synchronous || Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in addmagnetoelasticforce failed: %+v \n", err)
		}
	}

	return
}
