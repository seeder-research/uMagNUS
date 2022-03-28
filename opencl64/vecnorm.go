package opencl64

import (
	"fmt"

	"github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	"github.com/seeder-research/uMagNUS/util"
)

// dst = sqrt(dot(a, a)),
func VecNorm(dst *data.Slice, a *data.Slice) {
	util.Argument(dst.NComp() == 1 && a.NComp() == 3)
	util.Argument(dst.Len() == a.Len())

	N := dst.Len()
	cfg := make1DConf(N)
	event := k_vecnorm_async(dst.DevPtr(0),
		a.DevPtr(X), a.DevPtr(Y), a.DevPtr(Z),
		N, cfg, [](*cl.Event){dst.GetEvent(0), a.GetEvent(X), a.GetEvent(Y), a.GetEvent(Z)})
	dst.SetEvent(0, event)
	a.SetEvent(X, event)
	a.SetEvent(Y, event)
	a.SetEvent(Z, event)

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in vecnorm: %+v \n", err)
	}
}
