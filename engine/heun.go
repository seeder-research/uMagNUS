package engine

import (
	"fmt"
	"math"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

// Adaptive Heun solver.
type Heun struct{}

// Adaptive Heun method, can be used as solver.Step
func (he *Heun) Step() {
	// sync in the beginning
	if err1 := opencl.WaitAllQueuesToFinish(); err1 != nil {
		fmt.Printf("error waiting for all queues at start of esdirk32a.step(): %+v \n", err1)
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
	m0 := opencl.Buffer(VECTOR, y.Size())
	defer opencl.Recycle(dy0)
	defer opencl.Recycle(m0)
	data.Copy(m0, y)
	// sync queues to seqQueue
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})

	if FixDt != 0 {
		Dt_si = FixDt
	}

	dt := float32(Dt_si * GammaLL)
	util.Assert(dt > 0)

	// stage 1
	torqueFn(dy0)
	opencl.Madd2(y, y, dy0, 1, dt, queues, nil) // y = y + dt * dy

	// stage 2
	dy := opencl.Buffer(3, y.Size())
	defer opencl.Recycle(dy)
	Time += Dt_si
	torqueFn(dy)

	err := opencl.MaxVecDiff(dy0, dy, seqQueue, nil) * float64(dt)

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		opencl.Madd3(y, y, dy, dy0, 1, 0.5*dt, -0.5*dt, queues, nil)
		// sync queues to seqQueue
		opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
		M.normalize()
		NSteps++
		adaptDt(math.Pow(MaxErr/err, 1./2.))
		setLastErr(err)
		setMaxTorque(dy)
	} else {
		// undo bad step
		util.Assert(FixDt == 0)
		Time -= Dt_si
		data.Copy(y, m0)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./3.))
	}

	// sync before returning
	if err1 := seqQueue.Finish(); err1 != nil {
		fmt.Printf("error waiting for queue to finish after heun.step(): %+v \n", err1)
	}
}

func (_ *Heun) Free() {}

func (_ *Heun) EmType() bool {
	return false
}

func (_ *Heun) AdvOrder() int {
	return 2
}

func (_ *Heun) EmOrder() int {
	return -1
}
