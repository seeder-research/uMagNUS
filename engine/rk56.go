package engine

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

type RK56 struct {
	k1 *data.Slice // torque at end of step is kept for beginning of next step
}

func (rk *RK56) Step() {
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
	if rk.k1.Size() != m.Size() {
		rk.Free()
	}

	// first step ever: one-time k1 init and eval
	if rk.k1 == nil {
		rk.k1 = opencl.NewSlice(3, size)
		torqueFn(rk.k1)
	}

	t0 := Time
	// backup magnetization
	m0 := opencl.Buffer(3, size)
	defer opencl.Recycle(m0)
	data.Copy(m0, m)
	// sync queues to seqQueue
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})

	k2, k3, k4, k5, k6, k7, k8 := opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size)
	defer opencl.Recycle(k2)
	defer opencl.Recycle(k3)
	defer opencl.Recycle(k4)
	defer opencl.Recycle(k5)
	defer opencl.Recycle(k6)
	defer opencl.Recycle(k7)
	defer opencl.Recycle(k8)
	//k2 will be recyled as k9

	h := float32(Dt_si * GammaLL) // internal time step = Dt * gammaLL

	// stage 1
	torqueFn(rk.k1)

	// stage 2
	Time = t0 + (1./6.)*Dt_si
	opencl.Madd2(m, m, rk.k1, 1, (1./6.)*h, queues, nil) // m = m*1 + k1*h/6
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k2)

	// stage 3
	Time = t0 + (4./15.)*Dt_si
	opencl.Madd3(m, m0, rk.k1, k2, 1, (4./75.)*h, (16./75.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k3)

	// stage 4
	Time = t0 + (2./3.)*Dt_si
	opencl.Madd4(m, m0, rk.k1, k2, k3, 1, (5./6.)*h, (-8./3.)*h, (5./2.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k4)

	// stage 5
	Time = t0 + (4./5.)*Dt_si
	opencl.Madd5(m, m0, rk.k1, k2, k3, k4, 1, (-8./5.)*h, (144./25.)*h, (-4.)*h, (16./25.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k5)

	// stage 6
	Time = t0 + (1.)*Dt_si
	opencl.Madd6(m, m0, rk.k1, k2, k3, k4, k5, 1, (361./320.)*h, (-18./5.)*h, (407./128.)*h, (-11./80.)*h, (55./128.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k6)

	// stage 7
	Time = t0
	opencl.Madd5(m, m0, rk.k1, k3, k4, k5, 1, (-11./640.)*h, (11./256.)*h, (-11/160.)*h, (11./256.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k7)

	// stage 8
	Time = t0 + (1.)*Dt_si
	opencl.Madd7(m, m0, rk.k1, k2, k3, k4, k5, k7, 1, (93./640.)*h, (-18./5.)*h, (803./256.)*h, (-11./160.)*h, (99./256.)*h, (1.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k8)

	// stage 9: 6th order solution
	Time = t0 + (1.)*Dt_si
	//madd6(m, m0, k1, k3, k4, k5, k6, 1, (31./384.)*h, (1125./2816.)*h, (9./32.)*h, (125./768.)*h, (5./66.)*h)
	opencl.Madd7(m, m0, rk.k1, k3, k4, k5, k7, k8, 1, (7./1408.)*h, (1125./2816.)*h, (9./32.)*h, (125./768.)*h, (5./66.)*h, (5./66.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k2) // re-use k2

	// error estimate
	Err := opencl.Buffer(3, size)
	defer opencl.Recycle(Err)
	opencl.Madd4(Err, rk.k1, k6, k7, k8, (-5. / 66.), (-5. / 66.), (5. / 66.), (5. / 66.), queues, nil)

	integralController(Err, k2, rk.k1, m0, t0, float64(h), rk.AdvOrder(), rk.AdvOrder()+1, true)
}

func (rk *RK56) Free() {
	rk.k1.Free()
	rk.k1 = nil
}

func (_ *RK56) EmType() bool {
	return true
}

func (_ *RK56) AdvOrder() int {
	return 6
}

func (_ *RK56) EmOrder() int {
	return 5
}
