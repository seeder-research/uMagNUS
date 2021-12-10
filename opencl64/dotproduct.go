package opencl64

import (
	"fmt"

	data "github.com/seeder-research/uMagNUS/data64"
	"github.com/seeder-research/uMagNUS/cl"
	"github.com/seeder-research/uMagNUS/util"
)

// dst += prefactor * dot(a, b), as used for energy density
func AddDotProduct(dst *data.Slice, prefactor float64, a, b *data.Slice) {
	util.Argument(dst.NComp() == 1 && a.NComp() == 3 && b.NComp() == 3)
	util.Argument(dst.Len() == a.Len() && dst.Len() == b.Len())

	N := dst.Len()
	cfg := make1DConf(N)
	event := k_dotproduct_async(dst.DevPtr(0), prefactor,
		a.DevPtr(X), a.DevPtr(Y), a.DevPtr(Z),
		b.DevPtr(X), b.DevPtr(Y), b.DevPtr(Z),
		N, cfg, [](*cl.Event){dst.GetEvent(0), a.GetEvent(X), a.GetEvent(Y), a.GetEvent(Z),
			b.GetEvent(X), b.GetEvent(Y), b.GetEvent(Z)})

	dst.SetEvent(0, event)
	a.SetEvent(X, event)
	a.SetEvent(Y, event)
	a.SetEvent(Z, event)
	b.SetEvent(X, event)
	b.SetEvent(Y, event)
	b.SetEvent(Z, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in adddotproduct: %+v \n", err)
	}
}
