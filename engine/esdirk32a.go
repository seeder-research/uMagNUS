package engine

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

// Explicit singly diagonal implicit Rnge-Kutta (ESDIRK) solver.
// ESDIRK32a
// 3rd order, 4 stages per step, adaptive step.
// Anne Kværnø, Singly diagonal implicit Runge-Kutta methods
// with explicit first stage," BIT Numerical Mathematics vol. 44,
// 489-502, 2004.
//
// Advance using z{n+1}
//
//	k1 = f(tn, yn)
//	k2 = f(tn + 0.871733043 h, yn + 0.4358665215 h k1 + 0.4358665215 h k2)
//	k3 = f(tn + h, yn + 0.490563388419108 h k1 + 0.073570090080892 h k2 + 0.4358665215 h k3)
//	y{n+1}  = yn + 0.490563388419108 h k1 + 0.073570090080892 h k2 + 0.4358665215 h k3 // 2nd order
//	k4 = f(tn + h, yn + 0.308809969973036 h k1 + 1.490563388254108 h k2 - 1.235239879727145 h k3 + 0.4358665215 h k4)
//	z{n+1} = yn + 0.308809969973036 h k1 + 1.490563388254108 h k2 - 1.235239879727145 h k3 + 0.4358665215 h k4 // 3rd order
type ESDIRK32A struct {
	k1 *data.Slice // torque at end of step is kept for beginning of next step
}

func (esdirk *ESDIRK32A) Step() {
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

	m := M.Buffer()
	size := m.Size()

	if FixDt != 0 {
		Dt_si = FixDt
	}

	// upon resize: remove wrongly sized k1
	if esdirk.k1.Size() != m.Size() {
		esdirk.Free()
	}

	// first step ever: one-time k1 init and eval
	if esdirk.k1 == nil {
		esdirk.k1 = opencl.NewSlice(3, size)
		torqueFn(esdirk.k1)
	}

	// FSAL cannot be used with temperature
	if !Temp.isZero() {
		torqueFn(esdirk.k1)
	}

	t0 := Time
	// backup magnetization
	m0 := opencl.Buffer(3, size)
	defer opencl.Recycle(m0)
	data.Copy(m0, m) // D2D copy sync on seqQueue
	// sync queues to seqQueue
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})

	k2, k3, k4 := opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size)
	defer opencl.Recycle(k2)
	defer opencl.Recycle(k3)
	defer opencl.Recycle(k4)

	h := float32(Dt_si * GammaLL) // internal time step = Dt * gammaLL

	// there is no explicit stage 1: k1 from previous step

	// stage 2
	Time = t0 + (0.871733043)*Dt_si
	opencl.Madd2(m, m0, esdirk.k1, 1, (0.871733043)*h, queues, nil) // m{try} = m*1 + 0.871733043*k1
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	m_ := opencl.Buffer(3, size)
	defer opencl.Recycle(m_)
	data.Copy(m_, m)
	data.Copy(k2, esdirk.k1)
	_, _, _ = fixedPtIterations((0.4358665215)*h, m_, k2) // guaranteed sync

	// stage 3
	Time = t0 + Dt_si
	opencl.Madd3(m, m0, esdirk.k1, k2, 1, (0.490563388419108)*h, (0.073570090080892)*h, queues, nil) // m = m0*1 + k1*0.490563388419108 + k2*0.073570090080892
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k3, k2)
	_, _, _ = fixedPtIterations((0.4358665215)*h, m_, k3) // guaranteed sync

	// stage 4
	Time = t0 + Dt_si
	opencl.Madd4(m, m0, esdirk.k1, k2, k3, 1, (0.308809969973036)*h, (1.490563388254108)*h, (-1.235239879727145)*h, queues, nil) // m = m0*1 + k1*0.308809969973036 + k2*1.490563388254108 - k3*1.235239879727145
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k4, k3)
	_, _, _ = fixedPtIterations((0.4358665215)*h, m_, k4) // guaranteed sync

	// 3rd order solution
	opencl.Madd4(m_, esdirk.k1, k2, k3, k4, (0.308809969973036), (1.490563388254108), (-1.235239879727145), (0.4358665215), queues, nil) // m = m0*1 + k1*0.308809969973036 + k2*1.490563388254108 - k3*1.235239879727145 + k4*0.4358665215
	opencl.Madd2(m, m0, m_, 1, h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()

	// error estimate
	Time = t0 + Dt_si
	Err := k2 // re-use k2 as error
	// difference of 3rd and 2nd order torque without explicitly storing them first
	opencl.Madd4(Err, esdirk.k1, k2, k3, k4, (-0.181753418446072), (1.41699329817322), (-1.67110640122714), (0.4358665215), queues, nil)

	gustafssonController(Err, m_, esdirk.k1, m0, t0, float64(h), esdirk.AdvOrder(), true)
}

func (esdirk *ESDIRK32A) Free() {
	esdirk.k1.Free()
	esdirk.k1 = nil
}

func (_ *ESDIRK32A) EmType() bool {
	return true
}

func (_ *ESDIRK32A) AdvOrder() int {
	return 3
}

func (_ *ESDIRK32A) EmOrder() int {
	return 2
}
