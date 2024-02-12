package engine

import (
	"fmt"
	"math"

	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

var (
	Ext_TopologicalChargeLattice        = NewScalarValue("ext_topologicalchargelattice", "", "2D topological charge according to Berg and Lüscher", GetTopologicalChargeLattice)
	Ext_TopologicalChargeDensityLattice = NewScalarField("ext_topologicalchargedensitylattice", "1/m2",
		"2D topological charge density according to Berg and Lüscher", SetTopologicalChargeDensityLattice)
)

func SetTopologicalChargeDensityLattice(dst *data.Slice) {
	Refer("Berg1981")
	// sync in the beginning
	seqQueue := opencl.ClCmdQueue[0]
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queues to finish in settopologicalchargelattice: %+v \n", err)
	}
	opencl.SetTopologicalChargeLattice(dst, M.Buffer(), M.Mesh(), seqQueue, nil)
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queues to finish after settopologicalchargelattice: %+v \n", err)
	}
}

func GetTopologicalChargeLattice() float64 {
	// sync in the beginning
	seqQueue := opencl.ClCmdQueue[0]
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queues to finish in gettopologicalchargelattice: %+v \n", err)
	}
	s := ValueOf(Ext_TopologicalChargeDensityLattice)
	defer opencl.Recycle(s)
	c := Mesh().CellSize()
	N := Mesh().Size()

	return (0.25 * c[X] * c[Y] / math.Pi / float64(N[Z])) * float64(opencl.Sum(s, seqQueue, nil))
}
