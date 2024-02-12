package engine

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

type RK45DP struct {
	k1 *data.Slice // torque at end of step is kept for beginning of next step
}

func (rk *RK45DP) Step() {
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

	// FSAL cannot be used with finite temperature
	if !Temp.isZero() {
		torqueFn(rk.k1)
	}

	t0 := Time
	// backup magnetization
	m0 := opencl.Buffer(3, size)
	defer opencl.Recycle(m0)
	data.Copy(m0, m)
	// sync queues to seqQueue
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})

	k2, k3, k4, k5, k6 := opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size)
	defer opencl.Recycle(k2)
	defer opencl.Recycle(k3)
	defer opencl.Recycle(k4)
	defer opencl.Recycle(k5)
	defer opencl.Recycle(k6)
	// k2 will be re-used as k7

	h := float32(Dt_si * GammaLL) // internal time step = Dt * gammaLL

	// there is no explicit stage 1: k1 from previous step

	// stage 2
	Time = t0 + (1./5.)*Dt_si
	opencl.Madd2(m, m, rk.k1, 1, (1./5.)*h, queues, nil) // m = m*1 + k1*h/5
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k2)

	// stage 3
	Time = t0 + (3./10.)*Dt_si
	opencl.Madd3(m, m0, rk.k1, k2, 1, (3./40.)*h, (9./40.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k3)

	// stage 4
	Time = t0 + (4./5.)*Dt_si
	opencl.Madd4(m, m0, rk.k1, k2, k3, 1, (44./45.)*h, (-56./15.)*h, (32./9.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k4)

	// stage 5
	Time = t0 + (8./9.)*Dt_si
	opencl.Madd5(m, m0, rk.k1, k2, k3, k4, 1, (19372./6561.)*h, (-25360./2187.)*h, (64448./6561.)*h, (-212./729.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k5)

	// stage 6
	Time = t0 + (1.)*Dt_si
	opencl.Madd6(m, m0, rk.k1, k2, k3, k4, k5, 1, (9017./3168.)*h, (-355./33.)*h, (46732./5247.)*h, (49./176.)*h, (-5103./18656.)*h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	torqueFn(k6)

	// stage 7: 5th order solution
	Time = t0 + (1.)*Dt_si
	// no k2
	opencl.Madd6(m, m0, rk.k1, k3, k4, k5, k6, 1, (35./384.)*h, (500./1113.)*h, (125./192.)*h, (-2187./6784.)*h, (11./84.)*h, queues, nil) // 5th
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	k7 := k2     // re-use k2
	torqueFn(k7) // next torque if OK

	// error estimate
	Err := opencl.Buffer(3, size) //k3 // re-use k3 as error estimate
	defer opencl.Recycle(Err)
	opencl.Madd6(Err, rk.k1, k3, k4, k5, k6, k7, (35./384.)-(5179./57600.), (500./1113.)-(7571./16695.), (125./192.)-(393./640.), (-2187./6784.)-(-92097./339200.), (11./84.)-(187./2100.), (0.)-(1./40.), queues, nil)

	integralController(Err, k7, rk.k1, m0, t0, float64(h), rk.AdvOrder(), rk.AdvOrder()+1, true)
}

func (rk *RK45DP) Free() {
	rk.k1.Free()
	rk.k1 = nil
}

func (_ *RK45DP) EmType() bool {
	return true
}

func (_ *RK45DP) AdvOrder() int {
	return 5
}

func (_ *RK45DP) EmOrder() int {
	return 4
}
