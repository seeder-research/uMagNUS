package opencl64

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

var config1DSize int

// Make a 1D kernel launch configuration suited for N threads.
func make1DConf(N int) *config {

	gr := make([]int, 3)
	gr[0], gr[1], gr[2] = config1DSize, 1, 1

	return &config{Grid: gr, Block: nil}
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

func UpdateLaunchConfigs(c []int) {
	numItems := c[0] * c[1] * c[2] // total number of size of main data arrays

	// Work-items per Work-group
	groupSize := ClPrefWGSz

	// Find first multiple of groupSize larger than numItems
	if numItems >= ClTotalPE-groupSize {
		config1DSize = ClTotalPE
	} else {
		for i0 := groupSize; i0 < numItems; i0 += groupSize {
			config1DSize = i0
		}
	}

	// Find reduce config for intermediate reduce step
	if numItems <= reduceSingleSize {
		reduceintcfg = nil
	} else {
		if numItems >= ClTotalPE {
			reduceintcfg = &config{Grid: []int{ClTotalPE, 1, 1}, Block: []int{groupSize, 1, 1}}
		} else {
			for ii0 := groupSize; ii0 < numItems; ii0 += groupSize {
				reduceintcfg = &config{Grid: []int{ii0, 1, 1}, Block: []int{groupSize, 1, 1}}
			}
		}
	}
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
