package engine

import (
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	"github.com/seeder-research/uMagNUS/util"
	"math"
)

// Explicit singly diagonal implicit Rnge-Kutta (ESDIRK) solver.
// 3rd order, 5 stages per step, adaptive step.
// Anne Kværnø, Singly diagonal implicit Runge-Kutta methods
// with explicit first stage," BIT Numerical Mathematics vol. 44,
// 489-502, 2004.
// Advance with y{n+1}
// 	k1 = f(tn, yn)
// 	k2 = f(tn + 0.871733043 h, yn + 0.4358665215 h k1 + 0.4358665215 h k2)
// 	k3 = f(tn + 0.468238744853136 h, yn + 0.140737774731968 h k1 - 0.108365551378832 h k2 + 0.43586652150 h k3)
// 	k4 = f(tn + h, yn + 0.102399400616089 h k1 - 0.376878452267324 h k2 + 0.838612530151233 h k3 + 0.4358665215 h k4)
// 	y{n+1}  = yn + 0.102399400616089 h k1 - 0.376878452267324 h k2 + 0.838612530151233 h k3 + 0.4358665215 h k4 // 3rd order
// 	k5 = f(tn + h, yn + 0.157024897860995 h k1 + 0.117330441357768 h k2 + 0.616678030391680 h k3 - 0.326899891110444 h k4 + 0.4358665215 h k5)
// 	z{n+1} = yn + 0.157024897860995 h k1 + 0.117330441357768 h k2 + 0.616678030391680 h k3 - 0.326899891110444 h k4 + 0.4358665215 h k5) // 4th order
type ESDIRK43B struct {
	k1       *data.Slice // torque at end of step is kept for beginning of next step
	EmType   bool
	AdvOrder int
	EmOrder  int
}

func (esdirk *ESDIRK43B) Step() {
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
	m0, m_ := opencl.Buffer(3, size), opencl.Buffer(3, size)
	defer opencl.Recycle(m0)
	defer opencl.Recycle(m_)
	data.Copy(m0, m)

	k2, k3, k4, k5 := opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size)
	defer opencl.Recycle(k2)
	defer opencl.Recycle(k3)
	defer opencl.Recycle(k4)
	defer opencl.Recycle(k5)

	h := float32(Dt_si * GammaLL) // internal time step = Dt * gammaLL

	// there is no explicit stage 1: k1 from previous step

	// stage 2
	Time = t0 + (0.871733043)*Dt_si
	opencl.Madd2(m, m0, esdirk.k1, 1, (0.4358665215)*h) // m = m*1 + k1*0.4358665215
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k2, esdirk.k1)
	_, _, _ = fixedPtIterations((0.4358665215)*h, m_, k2)

	// stage 3
	Time = t0 + (0.468238744853136)*Dt_si
	opencl.Madd3(m, m0, esdirk.k1, k2, 1, (0.140737774731968)*h, (-0.108365551378832)*h) // m = m0*1 + k1*0.140737774731968 - k2*0.108365551378832
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k3, k2)
	_, _, _ = fixedPtIterations((0.4358665215)*h, m_, k3)

	// stage 4
	Time = t0 + Dt_si
	opencl.Madd4(m, m0, esdirk.k1, k2, k3, 1, (0.102399400616089)*h, (-0.376878452267324)*h, (0.838612530151233)*h) // m = m0*1 + k1*0.102399400616089 - k2*0.376878452267324 + k3*0.838612530151233
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k4, k3)
	_, _, _ = fixedPtIterations((0.4358665215)*h, m_, k4)

	// stage 5
	Time = t0 + Dt_si
	opencl.Madd5(m, m0, esdirk.k1, k2, k3, k4, 1, (0.157024897860995)*h, (0.117330441357768)*h, (0.616678030391680)*h, (-0.326899891110444)*h) // m = m0*1 + k1*0.157024897860995 + k2*0.117330441357768 + k3*0.616678030391680 - k4*0.326899891110444
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k5, k4)
	_, _, _ = fixedPtIterations((0.4358665215)*h, m_, k5)

	// 3rd order solution
	opencl.Madd4(m_, esdirk.k1, k2, k3, k4, (0.102399400616089), (-0.376878452267324), (0.838612530151233), (0.4358665215))
	opencl.Madd2(m, m0, m_, 1, h)
	M.normalize()

	// error estimate
	Time = t0 + Dt_si
	Err := k2 // re-use k2 as error
	// difference of 3rd and 4th order torque without explicitly storing them first
	opencl.Madd5(Err, esdirk.k1, k2, k3, k4, k5, (0.0546254972449057), (0.494208893625092), (-0.221934499759553), (-0.762766412610444), (0.4358665215))

	// determine error
	err := opencl.MaxVecNorm(Err) * float64(h)

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// Passed absolute error. Check relative error...
		errnorm := opencl.Buffer(1, size)
		defer opencl.Recycle(errnorm)
		opencl.VecNorm(errnorm, Err)
		ddtnorm := opencl.Buffer(1, size)
		defer opencl.Recycle(ddtnorm)
		opencl.VecNorm(ddtnorm, m_)
		maxdm := opencl.MaxVecNorm(m_)
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
			setMaxTorque(m_)
			NSteps++
			Time = t0 + Dt_si
			if fail == 0 {
				adaptDt(math.Pow(MaxErr/err, 1./3.))
			} else {
				adaptDt(math.Pow(RelErr/rlerr, 1./3.))
			}
			data.Copy(esdirk.k1, k4) // FSAL
		} else {
			// undo bad step
			//util.Println("Bad step at t=", t0, ", err=", err)
			util.Assert(FixDt == 0)
			Time = t0
			data.Copy(m, m0)
			NUndone++
			adaptDt(math.Pow(RelErr/rlerr, 1./4.))
		}
	} else {
		// undo bad step
		//util.Println("Bad step at t=", t0, ", err=", err)
		util.Assert(FixDt == 0)
		Time = t0
		data.Copy(m, m0)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./4.))
	}
}

func (esdirk *ESDIRK43B) Free() {
	esdirk.k1.Free()
	esdirk.k1 = nil
}
