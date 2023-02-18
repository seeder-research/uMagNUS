package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// Crop stores in dst a rectangle cropped from src at given offset position.
// dst size may be smaller than src.
func Crop(dst, src *data.Slice, offX, offY, offZ int) {
	D := dst.Size()
	S := src.Size()
	util.Argument(dst.NComp() == src.NComp())
	util.Argument(D[X]+offX <= S[X] && D[Y]+offY <= S[Y] && D[Z]+offZ <= S[Z])

	var wg sync.WaitGroup
	for c := 0; c < dst.NComp(); c++ {
		wg.Add(1)
		if Synchronous {
			crop__(dst, src, offX, offY, offZ, c, wg)
		} else {
			go crop__(dst, src, offX, offY, offZ, c, wg)
		}
	}
	wg.Wait()
}

func crop__(dst, src *data.Slice, offX, offY, offZ, idx int, wg_ sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	src.RLock(idx)
	defer src.RUnlock(idx)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("crop failed to create command queue: %+v \n", err)
		return nil
	}
	defer cmdqueue.Release()

	D := dst.Size()
	S := src.Size()
	cfg := make3DConf(D)

	ev := k_crop_async(dst.DevPtr(idx), D[X], D[Y], D[Z],
		src.DevPtr(idx), S[X], S[Y], S[Z],
		offX, offY, offZ, cfg, cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{ev}); err != nil {
		fmt.Printf("WaitForEvents failed in crop: %+v \n", err)
	}
}
