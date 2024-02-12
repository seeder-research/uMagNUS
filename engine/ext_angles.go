package engine

import (
	"fmt"

	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

var (
	ext_phi   = NewScalarField("ext_phi", "rad", "Azimuthal angle", SetPhi)
	ext_theta = NewScalarField("ext_theta", "rad", "Polar angle", SetTheta)
)

func SetPhi(dst *data.Slice) {
	// sync in the beginning
	seqQueue := opencl.ClCmdQueue[0]
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queues to finish in setphi: %+v \n", err)
	}
	opencl.SetPhi(dst, M.Buffer(), seqQueue, nil)
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queues to finish after setphi: %+v \n", err)
	}
}

func SetTheta(dst *data.Slice) {
	// sync in the beginning
	seqQueue := opencl.ClCmdQueue[0]
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queues to finish in settheta: %+v \n", err)
	}
	opencl.SetTheta(dst, M.Buffer(), seqQueue, nil)
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queues to finish after settheta: %+v \n", err)
	}
}
