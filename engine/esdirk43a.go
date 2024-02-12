package engine

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

// Explicit singly diagonal implicit Rnge-Kutta (ESDIRK) solver.
// 3rd order, 4 stages per step, adaptive step.
// John Bagterp Jørgensen, "A family of ESDIRK methods,"
// arXiv:1803.01613
// Advance with y{n+1}
//
//	k1 = f(tn, yn)
//	k2 = f(tn + 0.87173304301691799883 h, yn + 0.43586652150845899942 h k1 + 0.43586652150845899942 h k2)
//	k3 = f(tn + 0.46823874485184439565 h, yn + 0.14073777472470619619 h k1 - 0.1083655513813208000 h k2 + 0.43586652150845899942 h k3)
//	k4 = f(tn + h, yn + 0.10239940061991099768 h k1 - 0.3768784522555561061 h k2 + 0.83861253012718610911 h k3 + 0.43586652150845899942 h k4)
//	y{n+1}  = yn + 0.10239940061991 h k1 - 0.37687845225556 h k2 + 0.83861253012719 h k3 + 0.43586652150846 h k4  // 3rd order
//	z{n+1} = yn + 0.15702489786032493710 h k1 + 0.11733044137043884870 h k2 + 0.61667803039212146434 h k3 + 0.10896663037711474985 h k4) // 4th order
type ESDIRK43A struct {
	k1 *data.Slice // torque at end of step is kept for beginning of next step
}

func (esdirk *ESDIRK43A) Step() {
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
	m_ := opencl.Buffer(3, size)
	defer opencl.Recycle(m0)
	defer opencl.Recycle(m_)
	data.Copy(m0, m)
	// sync queues to seqQueue
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})

	k2, k3, k4 := opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size)
	defer opencl.Recycle(k2)
	defer opencl.Recycle(k3)
	defer opencl.Recycle(k4)

	h := float32(Dt_si * GammaLL) // internal time step = Dt * gammaLL

	// there is no explicit stage 1: k1 from previous step

	// stage 2
	Time = t0 + (0.87173304301691799883)*Dt_si
	opencl.Madd2(m, m0, esdirk.k1, 1, (0.43586652150845899942)*h, queues, nil) // m = m*1 + k1*0.43586652150845899942
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k2, esdirk.k1)
	_, _, _ = fixedPtIterations((0.43586652150845899942)*h, m_, k2)

	// stage 3
	Time = t0 + (0.46823874485184439565)*Dt_si
	opencl.Madd3(m, m0, esdirk.k1, k2, 1, (0.14073777472470619619)*h, (0.1083655513813208000)*h, queues, nil) // m = m0*1 + k1*0.14073777472470619619 + k2*0.1083655513813208000
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k3, k2)
	_, _, _ = fixedPtIterations((0.43586652150845899942)*h, m_, k3)

	// stage 4
	Time = t0 + Dt_si
	opencl.Madd5(m, m0, esdirk.k1, k2, k3, k4, 1, (0.10239940061991099768)*h, (-0.3768784522555561061)*h, (0.83861253012718610911)*h, (0.43586652150845899942)*h, queues, nil) // m = m0*1 + k1*0.10239940061991099768 - k2*0.3768784522555561061 + k3*0.83861253012718610911 + k4*0.43586652150845899942
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k4, k3)
	_, _, _ = fixedPtIterations((0.43586652150845899942)*h, m_, k4)

	// 3rd order solution
	opencl.Madd4(m_, esdirk.k1, k2, k3, k4, (0.10239940061991099768), (-0.3768784522555561061), (0.83861253012718610911), (0.43586652150845899942), queues, nil)
	opencl.Madd2(m, m0, m_, 1, h, queues, nil)
	// sync queues to seqQueue
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	M.normalize()

	// error estimate
	Err := k2 // re-use k2 as error
	// difference of 3rd and 4th order torque without explicitly storing them first
	opencl.Madd4(Err, esdirk.k1, k2, k3, k4, (-0.05462549724041393942), (-0.49420889362599495480), (0.22193449973506464477), (0.32689989113134424957), queues, nil)

	gustafssonController(Err, m_, esdirk.k1, m0, t0, float64(h), esdirk.AdvOrder(), true)
}

func (esdirk *ESDIRK43A) Free() {
	esdirk.k1.Free()
	esdirk.k1 = nil
}

func (_ *ESDIRK43A) EmType() bool {
	return true
}

func (_ *ESDIRK43A) AdvOrder() int {
	return 3
}

func (_ *ESDIRK43A) EmOrder() int {
	return 4
}
