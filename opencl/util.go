package opencl

import (
	"fmt"
)

// OpenCL Launch parameters.
// there might be better choices for recent hardware,
// but it barely makes a difference in the end.
const (
	TileX, TileY = 16, 16
	MaxGridSize  = 65535
)

// opencl launch configuration
type config struct {
	Grid, Block []int
}

// Make a 1D kernel launch configuration suited for N threads.
func make1DConf(N int) *config {
	bl := make([]int, 3)
	bl[0], bl[1], bl[2] = ClPrefWGSz, 1, 1

	n2 := divUp(N, ClPrefWGSz) // N2 blocks left
	nx := divUp(n2, MaxGridSize)
	ny := divUp(n2, nx)
	gr := make([]int, 3)
	gr[0], gr[1], gr[2] = (nx * bl[0]), (ny * bl[1]), bl[2]

	return &config{gr, bl}
}

// Make a 3D kernel launch configuration suited for N threads.
func make3DConf(N [3]int) *config {
	bl := make([]int, 3)
	bl[0], bl[1], bl[2] = TileX, TileY, 1

	nx := divUp(N[X], TileX)
	ny := divUp(N[Y], TileY)
	gr := make([]int, 3)
	gr[0], gr[1], gr[2] = (nx * bl[0]), (ny * bl[1]), (N[Z] * bl[2])

	return &config{gr, bl}
}

// integer minimum
func iMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Integer division rounded up.
func divUp(x, y int) int {
	return ((x - 1) / y) + 1
}

const (
	X = 0
	Y = 1
	Z = 2
)

func checkSize(a interface {
	Size() [3]int
}, b ...interface {
	Size() [3]int
}) {
	sa := a.Size()
	for _, b := range b {
		if b.Size() != sa {
			panic(fmt.Sprintf("size mismatch: %v != %v", sa, b.Size()))
		}
	}
}
