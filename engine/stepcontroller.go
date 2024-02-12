package engine

import (
	"fmt"
	"math"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

// Time step controller to be used with these solvers:
//
//	RK4
func simpleController(Err *data.Slice, h float64, accOrder, rejOrder int) bool {

	// sync in the beginning
	if err1 := opencl.WaitAllQueuesToFinish(); err1 != nil {
		fmt.Printf("error waiting for all queues to finish in simplecontroller: %+v \n", err1)
	}
	seqQueue := opencl.ClCmdQueue[0]

	// determine error
	err := opencl.MaxVecNorm(Err, seqQueue, nil) * h

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		NSteps++
		adaptDt(math.Pow(MaxErr/err, 1./float64(accOrder)))
		setLastErr(err)
		return true
	} else {
		// undo bad step
		util.Assert(FixDt == 0)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./float64(rejOrder)))
		return false
	}
}

// Time step controllers to be used with embedded solvers
func integralController(Err, delM, k1, m0 *data.Slice, t0, h float64, accOrder, rejOrder int, FSAL bool) {

	m := M.Buffer()
	size := m.Size()

	// sync in the beginning
	if err1 := opencl.WaitAllQueuesToFinish(); err1 != nil {
		fmt.Printf("error waiting for all queues to finish in integralcontroller: %+v \n", err1)
	}
	seqQueue := opencl.ClCmdQueue[0]

	// determine error
	err := opencl.MaxVecNorm(Err, seqQueue, nil) * h

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// checkout queues to execute div
		q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
		defer opencl.CheckinQueue(q1idx)
		defer opencl.CheckinQueue(q2idx)
		defer opencl.CheckinQueue(q3idx)
		queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}

		// Passed absolute error. Check relative error...
		maxdm := opencl.MaxVecNorm(delM, seqQueue, nil)
		fail := 0
		rlerr := float64(0.0)

		if maxdm < MinSlope { // Only step using relerr if dmdt is big enough. Overcomes equilibrium problem
			fail = 0
		} else {
			errnorm := opencl.Buffer(1, size)
			ddtnorm := opencl.Buffer(1, size)
			defer opencl.Recycle(errnorm)
			defer opencl.Recycle(ddtnorm)

			opencl.VecNorm(errnorm, Err, queues[0], nil)
			opencl.VecNorm(ddtnorm, delM, queues[1], nil)

			// sync all queues
			if err1 := opencl.WaitAllQueuesToFinish(); err1 != nil {
				fmt.Printf("error waiting for all queues to finish in integralcontroller after maxdm: %+v \n", err1)
			}

			opencl.Div(errnorm, errnorm, ddtnorm, queues, nil) //re-use errnorm
			// sync queues and seqQueue
			opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)

			rlerr = float64(opencl.MaxAbs(errnorm, seqQueue, nil))
			fail = 1
		}

		if fail == 0 || RelErr <= 0.0 || rlerr < RelErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
			// step OK
			setLastErr(err)
			setMaxTorque(delM)
			NSteps++
			Time = t0 + Dt_si
			if fail == 0 {
				adaptDt(math.Pow(MaxErr/err, 1./float64(accOrder)))
			} else {
				adaptDt(math.Pow(RelErr/rlerr, 1./float64(accOrder)))
			}
			if FSAL && (k1 != nil) && (delM != nil) {
				data.Copy(k1, delM) // FSAL
				// sync seqQueue
				if err1 := seqQueue.Finish(); err1 != nil {
					fmt.Printf("error waiting for seqQueue to finish in integralcontroller: %+v \n", err1)
				}
			}
		} else {
			// undo bad step
			//util.Println("Bad step at t=", t0, ", err=", err)
			util.Assert(FixDt == 0)
			Time = t0
			data.Copy(m, m0)
			// sync seqQueue
			if err1 := seqQueue.Finish(); err1 != nil {
				fmt.Printf("error waiting for seqQueue to finish in integralcontroller: %+v \n", err1)
			}
			NUndone++
			adaptDt(math.Pow(RelErr/rlerr, 1./float64(rejOrder)))
		}
	} else {
		// undo bad step
		//util.Println("Bad step at t=", t0, ", err=", err)
		util.Assert(FixDt == 0)
		Time = t0
		data.Copy(m, m0)
		// sync seqQueue
		if err1 := seqQueue.Finish(); err1 != nil {
			fmt.Printf("error waiting for seqQueue to finish in integralcontroller: %+v \n", err1)
		}
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./float64(rejOrder)))
	}
}

// Gustafsson accelerated time stepper
// Used for ESDIRK steppers
// K. Gustafsson, "Control-theoretic techniques for stepsize selection in
// implicit Ruge-Kutta methods," ACM Trans. Mathematical Sotware vol. 20,
// No. 4, pp. 496-517, Dec. 1994.
func gustafssonController(Err, delM, k1, m0 *data.Slice, t0, h float64, kOrder int, FSAL bool) {
	//	if {current-step-accept} then
	//		if {first-step | first-after-succ-rej | conv-restrict } then
	//			h_r <= pow((tol / r), (1. / k)) * h
	//		else
	//			h_r <= (h / h_acc) * pow((tol / r), (k_2 / k)) * pow((r_acc / r), (k_1 / k)) * h
	//		r_acc <= r
	//		h_acc < h
	//	else
	//		if {successive-rej} then
	//			k_est <= log(r / r_rej) / log(h / h_rej)
	//			h_r <= pow((tol / r), (1. / k_est)) * h
	//		else
	//			h_r <= pow((tol / r), (1. / k)) * h
	//		r_rej <= r
	//		h_rej <= h
	//	h <= restrict(h_r)
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	m := M.Buffer()
	size := m.Size()

	// sync in the beginning
	if err1 := opencl.WaitAllQueuesToFinish(); err1 != nil {
		fmt.Printf("error waiting for all queues to finish in gustafssoncontroller: %+v \n", err1)
	}
	seqQueue := opencl.ClCmdQueue[0]

	// determine error
	err := opencl.MaxVecNorm(Err, seqQueue, nil) * h

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// checkout queues to execute div
		q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
		defer opencl.CheckinQueue(q1idx)
		defer opencl.CheckinQueue(q2idx)
		defer opencl.CheckinQueue(q3idx)
		queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}

		// Passed absolute error. Check relative error...
		maxdm := opencl.MaxVecNorm(delM, seqQueue, nil)
		fail := 0
		rlerr := float64(0.0)

		if maxdm < MinSlope { // Only step using relerr if dmdt is big enough. Overcomes equilibrium problem
			fail = 0
		} else {
			errnorm := opencl.Buffer(1, size)
			ddtnorm := opencl.Buffer(1, size)
			defer opencl.Recycle(errnorm)
			defer opencl.Recycle(ddtnorm)

			opencl.VecNorm(errnorm, Err, queues[0], nil)
			opencl.VecNorm(ddtnorm, delM, queues[1], nil)

			// sync all queues
			if err1 := opencl.WaitAllQueuesToFinish(); err1 != nil {
				fmt.Printf("error waiting for all queues to finish in gustafssoncontroller after maxdm: %+v \n", err1)
			}

			opencl.Div(errnorm, errnorm, ddtnorm, queues, nil) //re-use errnorm
			// sync queues and seqQueue
			opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)

			rlerr = float64(opencl.MaxAbs(errnorm, seqQueue, nil))
			fail = 1
		}

		if fail == 0 || RelErr <= 0.0 || rlerr < RelErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
			// step OK
			setLastErr(err)
			setMaxTorque(delM)
			NSteps++
			Time = t0 + Dt_si
			if fail == 0 {
				if accErr < 0. {
					dtScale := math.Pow(MaxErr/err, 1./float64(kOrder))
					accDt = dtScale * Dt_si
					adaptDt(dtScale)
					accErr = err
				} else {
					dtScale := math.Pow(MaxErr/err, gustafssonk1/float64(kOrder)) * math.Pow(accErr/err, gustafssonk2/float64(kOrder)) * (Dt_si / accDt)
					accDt = dtScale * Dt_si
					adaptDt(dtScale)
					accErr = err
				}
				rejErr = -1.0
				rejRelErr = 1.0
				accRelErr = -1.0
			} else {
				if accRelErr < 0. {
					dtScale := math.Pow(RelErr/rlerr, 1./float64(kOrder))
					accDt = dtScale * Dt_si
					adaptDt(dtScale)
					accRelErr = rlerr
				} else {
					dtScale := math.Pow(RelErr/rlerr, gustafssonk1/float64(kOrder)) * math.Pow(accErr/rlerr, gustafssonk2/float64(kOrder)) * (Dt_si / accDt)
					accDt = dtScale * Dt_si
					adaptDt(dtScale)
					accRelErr = rlerr
				}
				rejErr = -1.0
				rejRelErr = 1.0
				accErr = -1.0
			}
			if FSAL && (k1 != nil) && (delM != nil) {
				data.Copy(k1, delM) // FSAL
				// sync seqQueue
				if err1 := seqQueue.Finish(); err1 != nil {
					fmt.Printf("error waiting for seqQueue to finish in gustafssoncontroller: %+v \n", err1)
				}
			}
		} else {
			// undo bad step
			//util.Println("Bad step at t=", t0, ", err=", err)
			util.Assert(FixDt == 0)
			Time = t0
			data.Copy(m, m0)
			// sync seqQueue
			if err1 := seqQueue.Finish(); err1 != nil {
				fmt.Printf("error waiting for seqQueue to finish in gustafssoncontroller: %+v \n", err1)
			}
			NUndone++
			if rejRelErr <= 0. {
				dtScale := math.Pow(RelErr/rlerr, 1./float64(kOrder))
				rejDt = dtScale * Dt_si
				adaptDt(dtScale)
				rejRelErr = rlerr
			} else {
				k_est := math.Log(rlerr/rejRelErr) / math.Log(Dt_si/rejDt)
				if k_est > float64(kOrder) {
					k_est = float64(kOrder)
				}
				if k_est < 0.1 {
					k_est = 0.1
				}
				dtScale := math.Pow(RelErr/rlerr, 1./k_est)
				rejDt = dtScale * Dt_si
				adaptDt(dtScale)
				rejRelErr = rlerr
			}
			accErr = -1.0
			accRelErr = -1.0
			rejErr = -1.0
		}
	} else {
		// undo bad step
		//util.Println("Bad step at t=", t0, ", err=", err)
		util.Assert(FixDt == 0)
		Time = t0
		data.Copy(m, m0)
		// sync seqQueue
		if err1 := seqQueue.Finish(); err1 != nil {
			fmt.Printf("error waiting for seqQueue to finish in gustafssoncontroller: %+v \n", err1)
		}
		NUndone++
		if rejErr <= 0. {
			dtScale := math.Pow(MaxErr/err, 1./float64(kOrder))
			rejDt = dtScale * Dt_si
			adaptDt(dtScale)
			rejErr = err
		} else {
			k_est := math.Log(err/rejErr) / math.Log(Dt_si/rejDt)
			if k_est > float64(kOrder) {
				k_est = float64(kOrder)
			}
			if k_est < 0.1 {
				k_est = 0.1
			}
			dtScale := math.Pow(MaxErr/err, 1./k_est)
			rejDt = dtScale * Dt_si
			adaptDt(dtScale)
			rejErr = err
		}
		accErr = -1.0
		accRelErr = -1.0
		rejRelErr = -1.0
	}
}
