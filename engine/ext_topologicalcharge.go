package engine

import (
	"fmt"
	"math"

	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

var (
	Ext_TopologicalCharge        = NewScalarValue("ext_topologicalcharge", "", "2D topological charge", GetTopologicalCharge)
	Ext_TopologicalChargeDensity = NewScalarField("ext_topologicalchargedensity", "1/m2",
		"2D topological charge density m·(∂m/∂x ❌ ∂m/∂y)", SetTopologicalChargeDensity)
)

func SetTopologicalChargeDensity(dst *data.Slice) {
	// sync in the beginning
	seqQueue := opencl.ClCmdQueue[0]
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error in waiting for queues to finish in settopologicalchargedensity: %+v \n", err)
	}
	opencl.SetTopologicalCharge(dst, M.Buffer(), M.Mesh(), seqQueue, nil)
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error in waiting for queues to finish after settopologicalchargedensity: %+v \n", err)
	}
}

func GetTopologicalCharge() float64 {
	// sync in the beginning
	seqQueue := opencl.ClCmdQueue[0]
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error in waiting for queues to finish in gettopologicalcharge: %+v \n", err)
	}
	s := ValueOf(Ext_TopologicalChargeDensity)
	defer opencl.Recycle(s)
	c := Mesh().CellSize()
	N := Mesh().Size()
	return (0.25 * c[X] * c[Y] / math.Pi / float64(N[Z])) * float64(opencl.Sum(s, seqQueue, nil))
}
