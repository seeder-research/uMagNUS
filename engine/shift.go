package engine

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

var (
	TotalShift, TotalYShift                    float64                        // accumulated window shift (X and Y) in meter
	ShiftMagL, ShiftMagR, ShiftMagU, ShiftMagD data.Vector                    // when shifting m, put these value at the left/right edge.
	ShiftM, ShiftGeom, ShiftRegions            bool        = true, true, true // should shift act on magnetization, geometry, regions?
)

func init() {
	DeclFunc("Shift", Shift, "Shifts the simulation by +1/-1 cells along X")
	DeclVar("ShiftMagL", &ShiftMagL, "Upon shift, insert this magnetization from the left")
	DeclVar("ShiftMagR", &ShiftMagR, "Upon shift, insert this magnetization from the right")
	DeclVar("ShiftMagU", &ShiftMagU, "Upon shift, insert this magnetization from the top")
	DeclVar("ShiftMagD", &ShiftMagD, "Upon shift, insert this magnetization from the bottom")
	DeclVar("ShiftM", &ShiftM, "Whether Shift() acts on magnetization")
	DeclVar("ShiftGeom", &ShiftGeom, "Whether Shift() acts on geometry")
	DeclVar("ShiftRegions", &ShiftRegions, "Whether Shift() acts on regions")
	DeclVar("TotalShift", &TotalShift, "Amount by which the simulation has been shifted (m).")
}

// position of the window lab frame
func GetShiftPos() float64  { return -TotalShift }
func GetShiftYPos() float64 { return -TotalYShift }

// shift the simulation window over dx cells in X direction
func Shift(dx int) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues in yshift: %+v \n", err)
	}
	TotalShift += float64(dx) * Mesh().CellSize()[X] // needed to re-init geom, regions
	if ShiftM {
		shiftMag(M.Buffer(), dx) // TODO: M.shift?
	}
	if ShiftRegions {
		regions.shift(dx)
	}
	if ShiftGeom {
		geometry.shift(dx)
	}
	M.normalize()
	// sync before returning
	seqQueue := opencl.ClCmdQueue[0]
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish after shift: %+v \n", err)
	}
}

func shiftMag(m *data.Slice, dx int) {
	m2 := make([]*data.Slice, m.NComp())
	for c := 0; c < m.NComp(); c++ {
		m2[c] = opencl.Buffer(1, m.Size())
		defer opencl.Recycle(m2[c])
	}
	seqQueue := opencl.ClCmdQueue[0]
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		opencl.ShiftX(m2[c], comp, dx, float32(ShiftMagL[c]), float32(ShiftMagR[c]), queues[c], nil)
		opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{queues[c]})
		data.Copy(comp, m2[c]) // str0 ?
	}
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish after shiftmag: %+v \n", err)
	}
}

// shift the simulation window over dy cells in Y direction
func YShift(dy int) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues in yshift: %+v \n", err)
	}
	TotalYShift += float64(dy) * Mesh().CellSize()[Y] // needed to re-init geom, regions
	if ShiftM {
		shiftMagY(M.Buffer(), dy)
	}
	if ShiftRegions {
		regions.shiftY(dy)
	}
	if ShiftGeom {
		geometry.shiftY(dy)
	}
	M.normalize()
	// sync before returning
	seqQueue := opencl.ClCmdQueue[0]
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish after yshift: %+v \n", err)
	}
}

func shiftMagY(m *data.Slice, dy int) {
	m2 := make([]*data.Slice, m.NComp())
	for c := 0; c < m.NComp(); c++ {
		m2[c] = opencl.Buffer(1, m.Size())
		defer opencl.Recycle(m2[c])
	}
	seqQueue := opencl.ClCmdQueue[0]
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		opencl.ShiftY(m2[c], comp, dy, float32(ShiftMagU[c]), float32(ShiftMagD[c]), queues[c], nil)
		opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{queues[c]})
		data.Copy(comp, m2[c]) // str0 ?
	}
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish after shiftmagy: %+v \n", err)
	}
}
