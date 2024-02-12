package engine

import (
	"fmt"

	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

var (
	MaxAngle  = NewScalarValue("MaxAngle", "rad", "maximum angle between neighboring spins", GetMaxAngle)
	SpinAngle = NewScalarField("spinAngle", "rad", "Angle between neighboring spins", SetSpinAngle)
)

func SetSpinAngle(dst *data.Slice) {
	// sync in the beginning
	seqQueue := opencl.ClCmdQueue[0]
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish in setspinangle: %+v \n", err)
	}
	opencl.SetMaxAngle(dst, M.Buffer(), lex2.Gpu(), regions.Gpu(), M.Mesh(), seqQueue, nil)
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish after setspinangle: %+v \n", err)
	}
}

func GetMaxAngle() float64 {
	// sync in the beginning
	seqQueue := opencl.ClCmdQueue[0]
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish in setspinangle: %+v \n", err)
	}
	s := ValueOf(SpinAngle)
	defer opencl.Recycle(s)
	return float64(opencl.MaxAbs(s, seqQueue, nil)) // just a max would be fine, but not currently implemented
}
