package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Add magneto-elasticit coupling field to the effective field.
// see magnetoelasticfield.cl
func AddMagnetoelasticField(Beff, m *data.Slice, exx, eyy, ezz, exy, exz, eyz, B1, B2, Msat MSlice) {
	util.Argument(Beff.Size() == m.Size())
	util.Argument(Beff.Size() == exx.Size())
	util.Argument(Beff.Size() == eyy.Size())
	util.Argument(Beff.Size() == ezz.Size())
	util.Argument(Beff.Size() == exy.Size())
	util.Argument(Beff.Size() == exz.Size())
	util.Argument(Beff.Size() == eyz.Size())

	N := Beff.Len()
	cfg := make1DConf(N)

	eventList := [](*cl.Event){}
	tmpEvt := Beff.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = Beff.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = Beff.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	if exx.GetSlicePtr() != nil {
		tmpEvt = exx.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if eyy.GetSlicePtr() != nil {
		tmpEvt = eyy.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if ezz.GetSlicePtr() != nil {
		tmpEvt = ezz.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if exy.GetSlicePtr() != nil {
		tmpEvt = exy.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if exz.GetSlicePtr() != nil {
		tmpEvt = exz.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if eyz.GetSlicePtr() != nil {
		tmpEvt = eyz.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if B1.GetSlicePtr() != nil {
		tmpEvt = B1.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if B2.GetSlicePtr() != nil {
		tmpEvt = B2.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if Msat.GetSlicePtr() != nil {
		tmpEvt = Msat.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_addmagnetoelasticfield_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		exx.DevPtr(0), exx.Mul(0), eyy.DevPtr(0), eyy.Mul(0), ezz.DevPtr(0), ezz.Mul(0),
		exy.DevPtr(0), exy.Mul(0), exz.DevPtr(0), exz.Mul(0), eyz.DevPtr(0), eyz.Mul(0),
		B1.DevPtr(0), B1.Mul(0), B2.DevPtr(0), B2.Mul(0),
		Msat.DevPtr(0), Msat.Mul(0),
		N, cfg, eventList)

	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	if exx.GetSlicePtr() != nil {
		exx.SetEvent(0, event)
	}
	if eyy.GetSlicePtr() != nil {
		eyy.SetEvent(0, event)
	}
	if ezz.GetSlicePtr() != nil {
		ezz.SetEvent(0, event)
	}
	if exy.GetSlicePtr() != nil {
		exy.SetEvent(0, event)
	}
	if exz.GetSlicePtr() != nil {
		exz.SetEvent(0, event)
	}
	if eyz.GetSlicePtr() != nil {
		eyz.SetEvent(0, event)
	}
	if B1.GetSlicePtr() != nil {
		B1.SetEvent(0, event)
	}
	if B2.GetSlicePtr() != nil {
		B2.SetEvent(0, event)
	}
	if Msat.GetSlicePtr() != nil {
		Msat.SetEvent(0, event)
	}

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in addmagnetoelasticfield failed: %+v \n", err)
		}
	}
}

// Calculate magneto-elasticit force density
// see magnetoelasticforce.cl
func GetMagnetoelasticForceDensity(out, m *data.Slice, B1, B2 MSlice, mesh *data.Mesh) {
	util.Argument(out.Size() == m.Size())

	cellsize := mesh.CellSize()
	N := mesh.Size()
	cfg := make3DConf(N)

	rcsx := float32(1.0 / cellsize[X])
	rcsy := float32(1.0 / cellsize[Y])
	rcsz := float32(1.0 / cellsize[Z])

	eventList := [](*cl.Event){}
	tmpEvt := out.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = out.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = out.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	if B1.GetSlicePtr() != nil {
		tmpEvt = B1.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if B2.GetSlicePtr() != nil {
		tmpEvt = B2.GetEvent(0)
		if tmpEvt != nil {
			eventList = append(eventList, tmpEvt)
		}
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_getmagnetoelasticforce_async(out.DevPtr(X), out.DevPtr(Y), out.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B1.DevPtr(0), B1.Mul(0), B2.DevPtr(0), B2.Mul(0),
		rcsx, rcsy, rcsz,
		N[X], N[Y], N[Z],
		mesh.PBC_code(), cfg, eventList)

	out.SetEvent(X, event)
	out.SetEvent(Y, event)
	out.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	if B1.GetSlicePtr() != nil {
		B1.SetEvent(0, event)
	}
	if B2.GetSlicePtr() != nil {
		B2.SetEvent(0, event)
	}

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in addmagnetoelasticforce failed: %+v \n", err)
		}
	}
}
