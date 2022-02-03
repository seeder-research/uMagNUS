package engine64

import (
	//"fmt"
	data "github.com/seeder-research/uMagNUS/data64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
	"github.com/seeder-research/uMagNUS/util"
)

type BackwardEuler struct {
	dy1 *data.Slice
}

// Euler method, can be used as solver.Step.
func (s *BackwardEuler) Step() {
	util.AssertMsg(MaxErr > 0, "Backward euler solver requires MaxErr > 0")

	//	t0 := Time

	//	y := M.Buffer()

	//	y0 := opencl.Buffer(VECTOR, y.Size())
	//	defer opencl.Recycle(y0)
	//	data.Copy(y0, y)

	//	dy0 := opencl.Buffer(VECTOR, y.Size())
	//	defer opencl.Recycle(dy0)
	if s.dy1 == nil {
		s.dy1 = opencl.Buffer(VECTOR, M.Buffer().Size())
	}
	//	dy1 := s.dy1

	Dt_si = FixDt
	dt := float64(Dt_si * GammaLL)
	util.AssertMsg(dt > 0, "Backward Euler solver requires fixed time step > 0")

	// Fist guess
	//	Time = t0 + 0.5*Dt_si // 0.5 dt makes it implicit midpoint method

	// with temperature, previous torque cannot be used as predictor
	//	if Temp.isZero() {
	//		opencl.Madd2(y, y0, dy1, 1, dt) // predictor euler step with previous torque
	//		M.normalize()
	//	}

	//	torqueFn(dy0)
	//	opencl.Madd2(y, y0, dy0, 1, dt) // y = y0 + dt * dy
	//	M.normalize()

	// One iteration
	//	torqueFn(dy1)
	//	opencl.Madd2(y, y0, dy1, 1, dt) // y = y0 + dt * dy1
	//	M.normalize()

	//	Time = t0 + Dt_si

	//	err := opencl.MaxVecDiff(dy0, dy1) * float64(dt)

	err := bwEulerFixedPtIterations(Dt_si, dt, s.dy1)
	// adjust next time step
	//if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
	// step OK
	NSteps++
	setLastErr(err)
	setMaxTorque(s.dy1)
	//} else {
	// undo bad step
	//	util.Assert(FixDt == 0)
	//	Time = t0
	//	data.Copy(y, y0)
	//	NUndone++
	//}
}

func bwEulerFixedPtIterations(Dt_si, dt float64, dy1 *data.Slice) float64 {
	// For backward Euler, need to solve for
	// y_{n+1} = y_{n} + h * f(t_{n+1}, y_{n+1})
	// This function implements y_{n+1} = g(y_{n+1}) to find
	// the solution to the backward Euler step using
	// fixed point method, using the forward Euler result
	// as the initial guess

	t0 := Time

	y := M.Buffer()

	y0 := opencl.Buffer(VECTOR, y.Size())
	defer opencl.Recycle(y0)
	data.Copy(y0, y)

	yprev := opencl.Buffer(VECTOR, y.Size())
	defer opencl.Recycle(yprev)

	errVector := opencl.Buffer(VECTOR, y.Size())
	defer opencl.Recycle(errVector)

	dy0 := opencl.Buffer(VECTOR, y.Size())
	defer opencl.Recycle(dy0)

	torqueFn(dy0)

	// Initial guess
	opencl.Madd2(y, y0, dy0, 1.0, dt) // y = y0 + dt * dy
	M.normalize()
	data.Copy(yprev, y)

	Time = t0 + Dt_si

	errIter := float64(100.0)
	Niters := 0

	// fixed point iterations until converence criterion reached
	for (errIter > ErrConv) && (Niters < NConv) {
		torqueFn(dy1)
		opencl.Madd2(y, y0, dy1, 1.0, dt) // Backward Euler step
		M.normalize()

		// Calculate error as the difference in calculated predictions
		// in consecutive fixed point iterations
		errIter = opencl.MaxVecDiff(yprev, y)
		Niters++

		// Record fixed point result for next iteration
		data.Copy(yprev, y)
	}
	if Niters == NConv {
		util.Log("backward Euler fixed point iterations exceeded limit!")
	}
	return errIter
}

func (s *BackwardEuler) Free() {
	s.dy1.Free()
	s.dy1 = nil
}
