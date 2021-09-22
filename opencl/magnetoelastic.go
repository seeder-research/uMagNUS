package opencl

import (
	"fmt"
	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl/cl"
	"github.com/seeder-research/uMagNUS/util"
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

	event := k_addmagnetoelasticfield_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		exx.DevPtr(0), exx.Mul(0), eyy.DevPtr(0), eyy.Mul(0), ezz.DevPtr(0), ezz.Mul(0),
		exy.DevPtr(0), exy.Mul(0), exz.DevPtr(0), exz.Mul(0), eyz.DevPtr(0), eyz.Mul(0),
		B1.DevPtr(0), B1.Mul(0), B2.DevPtr(0), B2.Mul(0),
		Msat.DevPtr(0), Msat.Mul(0),
		N, cfg,
		[](*cl.Event){Beff.GetEvent(X), Beff.GetEvent(Y), Beff.GetEvent(Z),
			m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z),
			exx.GetEvent(0), eyy.GetEvent(0), ezz.GetEvent(0),
			exy.GetEvent(0), exz.GetEvent(0), eyz.GetEvent(0),
			B1.GetEvent(0), B2.GetEvent(0), Msat.GetEvent(0)})
	Beff.SetEvent(X, event)
	Beff.SetEvent(Y, event)
	Beff.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	exx.SetEvent(0, event)
	eyy.SetEvent(0, event)
	ezz.SetEvent(0, event)
	exy.SetEvent(0, event)
	exz.SetEvent(0, event)
	eyz.SetEvent(0, event)
	B1.SetEvent(0, event)
	B2.SetEvent(0, event)
	Msat.SetEvent(0, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents in addmagnetoelasticfield failed: %+v \n", err)
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

	event := k_getmagnetoelasticforce_async(out.DevPtr(X), out.DevPtr(Y), out.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		B1.DevPtr(0), B1.Mul(0), B2.DevPtr(0), B2.Mul(0),
		rcsx, rcsy, rcsz,
		N[X], N[Y], N[Z],
		mesh.PBC_code(), cfg,
		[](*cl.Event){out.GetEvent(X), out.GetEvent(Y), out.GetEvent(Z),
			m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z),
			B1.GetEvent(0), B2.GetEvent(0)})

	out.SetEvent(X, event)
	out.SetEvent(Y, event)
	out.SetEvent(Z, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	B1.SetEvent(0, event)
	B2.SetEvent(0, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents in addmagnetoelasticforce failed: %+v \n", err)
	}
}
