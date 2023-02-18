package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// shift dst by shx cells (positive or negative) along X-axis.
// new edge value is clampL at left edge or clampR at right edge.
func ShiftX(dst, src *data.Slice, shiftX int, clampL, clampR float32) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == src.Len())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		shiftx__(dst, src, shiftX, clampL, clampR, wg)
	} else {
		go shiftx__(dst, src, shiftX, clampL, clampR, wg)
	}
	wg.Wait()
}

func shiftx__(dst, src *data.Slice, shiftX int, clampL, clampR float32, wg_ sync.WaitGroup) {
	dst.Lock(0)
	defer dst.Unlock(0)
	src.RLock(0)
	defer src.RUnlock(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("shiftx failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := dst.Size()
	cfg := make3DConf(N)

	event := k_shiftx_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftX, clampL, clampR, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in shiftx failed: %+v \n", err)
	}
}

func ShiftY(dst, src *data.Slice, shiftY int, clampL, clampR float32) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == src.Len())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		shifty__(dst, src, shiftY, clampL, clampR, wg)
	} else {
		go shifty__(dst, src, shiftY, clampL, clampR, wg)
	}
	wg.Wait()
}

func shifty__(dst, src *data.Slice, shiftY int, clampL, clampR float32, wg_ sync.WaitGroup) {
	dst.Lock(0)
	defer dst.Unlock(0)
	src.RLock(0)
	defer src.RUnlock(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("shifty failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := dst.Size()
	cfg := make3DConf(N)

	event := k_shifty_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftY, clampL, clampR, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in shifty failed: %+v \n", err)
	}
}

func ShiftZ(dst, src *data.Slice, shiftZ int, clampL, clampR float32) {
	util.Argument(dst.NComp() == 1 && src.NComp() == 1)
	util.Assert(dst.Len() == src.Len())

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		shiftz__(dst, src, shiftZ, clampL, clampR, wg)
	} else {
		go shiftz__(dst, src, shiftZ, clampL, clampR, wg)
	}
	wg.Wait()
}

func shiftz__(dst, src *data.Slice, shiftZ int, clampL, clampR float32, wg_ sync.WaitGroup) {
	dst.Lock(0)
	defer dst.Unlock(0)
	src.RLock(0)
	defer src.RUnlock(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("shiftz failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := dst.Size()
	cfg := make3DConf(N)

	event := k_shiftz_async(dst.DevPtr(0), src.DevPtr(0), N[X], N[Y], N[Z], shiftZ, clampL, clampR, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in shiftz failed: %+v \n", err)
	}
}

// Like Shift, but for bytes
func ShiftBytes(dst, src *Bytes, m *data.Mesh, shiftX int, clamp byte) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		shiftbytes__(dst, src, m, shiftX, clamp, wg)
	} else {
		go shiftbytes__(dst, src, m, shiftX, clamp, wg)
	}
	wg.Wait()
}

func shiftbytes__(dst, src *Bytes, m *data.Mesh, shiftX int, clamp byte, wg_ sync.WaitGroup) {
	dst.Lock()
	defer dst.Unlock()
	src.RLock()
	defer src.RUnlock()

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("shiftbytes failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := m.Size()
	cfg := make3DConf(N)

	event := k_shiftbytes_async(dst.Ptr, src.Ptr, N[X], N[Y], N[Z], shiftX, clamp, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in shiftbytes failed: %+v \n", err)
	}
}

func ShiftBytesY(dst, src *Bytes, m *data.Mesh, shiftY int, clamp byte) {
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		shiftbytesy__(dst, src, m, shiftY, clamp, wg)
	} else {
		go shiftbytesy__(dst, src, m, shiftY, clamp, wg)
	}
	wg.Wait()
}

func shiftbytesy__(dst, src *Bytes, m *data.Mesh, shiftY int, clamp byte, wg_ sync.WaitGroup) {
	dst.Lock()
	defer dst.Unlock()
	src.RLock()
	defer src.RUnlock()

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("shiftbytes failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	N := m.Size()
	cfg := make3DConf(N)

	event := k_shiftbytesy_async(dst.Ptr, src.Ptr, N[X], N[Y], N[Z], shiftY, clamp, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in shiftbytesy failed: %+v \n", err)
	}
}
