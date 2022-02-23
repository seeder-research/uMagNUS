package engine64

import (
	//"fmt"
	data "github.com/seeder-research/uMagNUS/data64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
	util "github.com/seeder-research/uMagNUS/util"
)

type BackwardEuler struct {
	dy1 *data.Slice
}

// Euler method, can be used as solver.Step.
func (s *BackwardEuler) Step() {
	util.AssertMsg(MaxErr > 0, "Backward euler solver requires MaxErr > 0")

	t0 := Time

	y := M.Buffer()

	y0 := opencl.Buffer(VECTOR, y.Size())
	defer opencl.Recycle(y0)
	data.Copy(y0, y)

	//	dy0 := opencl.Buffer(VECTOR, y.Size())
	//	defer opencl.Recycle(dy0)
	if s.dy1 == nil {
		s.dy1 = opencl.Buffer(VECTOR, M.Buffer().Size())
	}
	//	dy1 := s.dy1

	Dt_si = FixDt
	dt := float64(Dt_si * GammaLL)
	util.AssertMsg(dt > 0, "Backward Euler solver requires fixed time step > 0")

	// First guess
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

	Time = t0 + Dt_si
	abserr, _, _ := fixedPtIterations(dt, y0, s.dy1)
	// adjust next time step
	//if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
	// step OK
	NSteps++
	setLastErr(abserr)
	setMaxTorque(s.dy1)
	//} else {
	// undo bad step
	//	util.Assert(FixDt == 0)
	//	Time = t0
	//	data.Copy(y, y0)
	//	NUndone++
	//}
}

func (s *BackwardEuler) Free() {
	s.dy1.Free()
	s.dy1 = nil
}

func (_ *BackwardEuler) EmType() bool {
	return false
}

func (_ *BackwardEuler) AdvOrder() int {
	return 1
}

func (_ *BackwardEuler) EmOrder() int {
	return -1
}
