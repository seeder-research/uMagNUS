package engine

import (
	"math"

	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

// Adaptive Heun solver.
type Heun struct{}

// Adaptive Heun method, can be used as solver.Step
func (he *Heun) Step() {
	y := M.Buffer()
	dy0 := opencl.Buffer(VECTOR, y.Size())
	m0 := opencl.Buffer(VECTOR, y.Size())
	defer opencl.Recycle(dy0)
	defer opencl.Recycle(m0)
	data.Copy(m0, y)

	if FixDt != 0 {
		Dt_si = FixDt
	}

	dt := float32(Dt_si * GammaLL)
	util.Assert(dt > 0)

	// stage 1
	torqueFn(dy0)
	opencl.Madd2(y, y, dy0, 1, dt) // y = y + dt * dy

	// stage 2
	dy := opencl.Buffer(3, y.Size())
	defer opencl.Recycle(dy)
	Time += Dt_si
	torqueFn(dy)

	err := opencl.MaxVecDiff(dy0, dy) * float64(dt)

	// adjust next time step
	if err < MaxErr || Dt_si <= MinDt || FixDt != 0 { // mindt check to avoid infinite loop
		// step OK
		opencl.Madd3(y, y, dy, dy0, 1, 0.5*dt, -0.5*dt)
		M.normalize()
		NSteps++
		adaptDt(math.Pow(MaxErr/err, 1./2.))
		setLastErr(err)
		setMaxTorque(dy)
	} else {
		// undo bad step
		util.Assert(FixDt == 0)
		Time -= Dt_si
		data.Copy(y, m0)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./3.))
	}
}

func (_ *Heun) Free() {}

func (_ *Heun) EmType() bool {
	return false
}

func (_ *Heun) AdvOrder() int {
	return 2
}

func (_ *Heun) EmOrder() int {
	return -1
}
