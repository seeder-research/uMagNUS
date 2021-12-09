package engine64

import (
	"github.com/seeder-research/uMagNUS/opencl"
	"github.com/seeder-research/uMagNUS/util"
	"math"
)

// Adaptive Heun solver.
type Heun struct{}

// Adaptive Heun method, can be used as solver.Step
func (_ *Heun) Step() {
	y := M.Buffer()
	dy0 := opencl.Buffer(VECTOR, y.Size())
	defer opencl.Recycle(dy0)

	if FixDt != 0 {
		Dt_si = FixDt
	}

	dt := float64(Dt_si * GammaLL)
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
		opencl.Madd2(y, y, dy0, 1, -dt)
		NUndone++
		adaptDt(math.Pow(MaxErr/err, 1./3.))
	}
}

func (_ *Heun) Free() {}
