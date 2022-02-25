package engine64

import (
	data "github.com/seeder-research/uMagNUS/data64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
	util "github.com/seeder-research/uMagNUS/util"
	"math"
)

// Time step controller to be used with these solvers:
//   RK4
func simpleController(Err *data.Slice, h float64, accOrder, rejOrder int) bool {

	// determine error
	err := opencl.MaxVecNorm(Err) * h

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

	// determine error
	err := opencl.MaxVecNorm(Err) * h

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// Passed absolute error. Check relative error...
		errnorm := opencl.Buffer(1, size)
		ddtnorm := opencl.Buffer(1, size)
		defer opencl.Recycle(errnorm)
		defer opencl.Recycle(ddtnorm)

		opencl.VecNorm(errnorm, Err)
		opencl.VecNorm(ddtnorm, delM)

		maxdm := opencl.MaxVecNorm(delM)
		fail := 0
		rlerr := float64(0.0)

		if maxdm < MinSlope { // Only step using relerr if dmdt is big enough. Overcomes equilibrium problem
			fail = 0
		} else {
			opencl.Div(errnorm, errnorm, ddtnorm) //re-use errnorm
			rlerr = float64(opencl.MaxAbs(errnorm))
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
			}
		} else {
			// undo bad step
			//util.Println("Bad step at t=", t0, ", err=", err)
			util.Assert(FixDt == 0)
			Time = t0
			data.Copy(m, m0)
			NUndone++
			adaptDt(math.Pow(RelErr/rlerr, 1./float64(rejOrder)))
		}
	} else {
		// undo bad step
		//util.Println("Bad step at t=", t0, ", err=", err)
		util.Assert(FixDt == 0)
		Time = t0
		data.Copy(m, m0)
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

	// determine error
	err := opencl.MaxVecNorm(Err) * h

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// Passed absolute error. Check relative error...
		errnorm := opencl.Buffer(1, size)
		ddtnorm := opencl.Buffer(1, size)
		defer opencl.Recycle(errnorm)
		defer opencl.Recycle(ddtnorm)

		opencl.VecNorm(errnorm, Err)
		opencl.VecNorm(ddtnorm, delM)

		maxdm := opencl.MaxVecNorm(delM)
		fail := 0
		rlerr := float64(0.0)

		if maxdm < MinSlope { // Only step using relerr if dmdt is big enough. Overcomes equilibrium problem
			fail = 0
		} else {
			opencl.Div(errnorm, errnorm, ddtnorm) //re-use errnorm
			rlerr = float64(opencl.MaxAbs(errnorm))
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
			}
		} else {
			// undo bad step
			//util.Println("Bad step at t=", t0, ", err=", err)
			util.Assert(FixDt == 0)
			Time = t0
			data.Copy(m, m0)
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
