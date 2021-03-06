package engine64

import (
	data "github.com/seeder-research/uMagNUS/data64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
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
}

func shiftMag(m *data.Slice, dx int) {
	m2 := opencl.Buffer(1, m.Size())
	defer opencl.Recycle(m2)
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		opencl.ShiftX(m2, comp, dx, float64(ShiftMagL[c]), float64(ShiftMagR[c]))
		data.Copy(comp, m2) // str0 ?
	}
}

// shift the simulation window over dy cells in Y direction
func YShift(dy int) {
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
}

func shiftMagY(m *data.Slice, dy int) {
	m2 := opencl.Buffer(1, m.Size())
	defer opencl.Recycle(m2)
	for c := 0; c < m.NComp(); c++ {
		comp := m.Comp(c)
		opencl.ShiftY(m2, comp, dy, float64(ShiftMagU[c]), float64(ShiftMagD[c]))
		data.Copy(comp, m2) // str0 ?
	}
}
