package engine

import (
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
	"math"
)

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
			if (FSAL && (k1 != nil) && (delM != nil)) {
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
