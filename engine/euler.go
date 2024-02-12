package engine

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

type Euler struct{}

// Euler method, can be used as solver.Step.
func (_ *Euler) Step() {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues at start of esdirk32a.step(): %+v \n", err)
	}

	// Get queues
	seqQueue := opencl.ClCmdQueue[0]
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}

	y := M.Buffer()
	dy0 := opencl.Buffer(VECTOR, y.Size())
	defer opencl.Recycle(dy0)

	torqueFn(dy0)
	setMaxTorque(dy0)

	// Adaptive time stepping: treat MaxErr as the maximum magnetization delta
	// (proportional to the error, but an overestimation for sure)
	var dt float32
	if FixDt != 0 {
		Dt_si = FixDt
		dt = float32(Dt_si * GammaLL)
	} else {
		dt = float32(MaxErr / LastTorque)
		Dt_si = float64(dt) / GammaLL
	}
	util.AssertMsg(dt > 0, "Euler solver requires fixed time step > 0")
	setLastErr(float64(dt) * LastTorque)

	opencl.Madd2(y, y, dy0, 1, dt, queues, nil) // y = y + dt * dy
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	// sync before returning
	if err1 := seqQueue.Finish(); err1 != nil {
		fmt.Printf("error waiting for all queues at start of esdirk32a.step(): %+v \n", err1)
	}

	Time += Dt_si
	NSteps++
}

func (_ *Euler) Free() {}

func (_ *Euler) EmType() bool {
	return false
}

func (_ *Euler) AdvOrder() int {
	return 1
}

func (_ *Euler) EmOrder() int {
	return -1
}
