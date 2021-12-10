package opencl64

import (
	"fmt"

	data "github.com/seeder-research/uMagNUS/data64"
	"github.com/seeder-research/uMagNUS/cl"
	"github.com/seeder-research/uMagNUS/util"
)

// Crop stores in dst a rectangle cropped from src at given offset position.
// dst size may be smaller than src.
func Crop(dst, src *data.Slice, offX, offY, offZ int) {
	D := dst.Size()
	S := src.Size()
	util.Argument(dst.NComp() == src.NComp())
	util.Argument(D[X]+offX <= S[X] && D[Y]+offY <= S[Y] && D[Z]+offZ <= S[Z])

	cfg := make3DConf(D)

	eventList := make([](*cl.Event), dst.NComp())
	for c := 0; c < dst.NComp(); c++ {
		eventList[c] = k_crop_async(dst.DevPtr(c), D[X], D[Y], D[Z],
			src.DevPtr(c), S[X], S[Y], S[Z],
			offX, offY, offZ, cfg, [](*cl.Event){dst.GetEvent(c), src.GetEvent(c)})
		dst.SetEvent(c, eventList[c])
		src.SetEvent(c, eventList[c])
	}
	err := cl.WaitForEvents(eventList)
	if err != nil {
		fmt.Printf("WaitForEvents failed in crop: %+v \n", err)
	}
}
