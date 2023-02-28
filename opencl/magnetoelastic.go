package opencl

import (
	"fmt"
	"sync"

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

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		addmagnetoelasticfield__(Beff, m, exx, eyy, ezz, exy, exz, eyz, B1, B2, Msat, &wg)
	} else {
		go func() {
			addmagnetoelasticfield__(Beff, m, exx, eyy, ezz, exy, exz, eyz, B1, B2, Msat, &wg)
		}()
	}
	wg.Wait()
}

func addmagnetoelasticfield__(Beff, m *data.Slice, exx, eyy, ezz, exy, exz, eyz, B1, B2, Msat MSlice, wg_ *sync.WaitGroup) {
	Beff.Lock(X)
	Beff.Lock(Y)
	Beff.Lock(Z)
	defer Beff.Unlock(X)
	defer Beff.Unlock(Y)
	defer Beff.Unlock(Z)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer m.RUnlock(X)
	defer m.RUnlock(Y)
	defer m.RUnlock(Z)
	if exx.GetSlicePtr() != nil {
		exx.RLock()
		defer exx.RUnlock()
	}
	if eyy.GetSlicePtr() != nil {
		eyy.RLock()
		defer eyy.RUnlock()
	}
	if ezz.GetSlicePtr() != nil {
		ezz.RLock()
		defer ezz.RUnlock()
	}
	if exy.GetSlicePtr() != nil {
		exy.RLock()
		defer exy.RUnlock()
	}
	if exz.GetSlicePtr() != nil {
		exz.RLock()
		defer exz.RUnlock()
	}
	if eyz.GetSlicePtr() != nil {
		eyz.RLock()
		defer eyz.RUnlock()
	}
	if B1.GetSlicePtr() != nil {
		B1.RLock()
		defer B1.RUnlock()
	}
	if B2.GetSlicePtr() != nil {
		B2.RLock()
		defer B2.RUnlock()
	}
	if Msat.GetSlicePtr() != nil {
		Msat.RLock()
		defer Msat.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("addmagnetoelasticfield failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	N := Beff.Len()
	cfg := make1DConf(N)

	event := k_addmagnetoelasticfield_async(Beff.DevPtr(X), Beff.DevPtr(Y), Beff.DevPtr(Z),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		exx.DevPtr(0), exx.Mul(0), eyy.DevPtr(0), eyy.Mul(0), ezz.DevPtr(0), ezz.Mul(0),
		exy.DevPtr(0), exy.Mul(0), exz.DevPtr(0), exz.Mul(0), eyz.DevPtr(0), eyz.Mul(0),
		B1.DevPtr(0), B1.Mul(0), B2.DevPtr(0), B2.Mul(0),
		Msat.DevPtr(0), Msat.Mul(0),
		N, cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in addmagnetoelasticfield: %+v \n", err)
	}
}

// Calculate magneto-elasticit force density
// see magnetoelasticforce.cl
func GetMagnetoelasticForceDensity(out, m *data.Slice, B1, B2 MSlice, mesh *data.Mesh) {
	util.Argument(out.Size() == m.Size())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		getmagnetoelasticforcedensity__(out, m, B1, B2, mesh, &wg)
	} else {
		go func() {
			getmagnetoelasticforcedensity__(out, m, B1, B2, mesh, &wg)
		}()
	}
	wg.Wait()
}

func getmagnetoelasticforcedensity__(out, m *data.Slice, B1, B2 MSlice, mesh *data.Mesh, wg_ *sync.WaitGroup) {
	out.Lock(X)
	out.Lock(Y)
	out.Lock(Z)
	defer out.Unlock(X)
	defer out.Unlock(Y)
	defer out.Unlock(Z)
	m.RLock(X)
	m.RLock(Y)
	m.RLock(Z)
	defer out.RUnlock(X)
	defer out.RUnlock(Y)
	defer out.RUnlock(Z)
	if B1.GetSlicePtr() != nil {
		B1.RLock()
		defer B1.RUnlock()
	}
	if B2.GetSlicePtr() != nil {
		B2.RLock()
		defer B2.RUnlock()
	}

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	fmt.Printf("getmagnetoelasticforcedensity failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

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
		mesh.PBC_code(), cfg, cmdqueue, nil)

	wg_.Done()

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in getmagnetoelasticforcedensity: %+v \n", err)
	}
}
