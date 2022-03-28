package opencl

import (
	"fmt"
	"unsafe"

	"github.com/seeder-research/uMagNUS/cl"
	"github.com/seeder-research/uMagNUS/data"
)

// SetMaxAngle sets dst to the maximum angle of each cells magnetization with all of its neighbors,
// provided the exchange stiffness with that neighbor is nonzero.
func SetMaxAngle(dst, m *data.Slice, Aex_red SymmLUT, regions *Bytes, mesh *data.Mesh) {
	N := mesh.Size()
	pbc := mesh.PBC_code()
	cfg := make3DConf(N)
	event := k_setmaxangle_async(dst.DevPtr(0),
		m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		unsafe.Pointer(Aex_red), regions.Ptr,
		N[X], N[Y], N[Z], pbc, cfg,
		[](*cl.Event){dst.GetEvent(0), m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z)})
	dst.SetEvent(0, event)
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in setmaxangle: %+v \n", err)
	}
}
