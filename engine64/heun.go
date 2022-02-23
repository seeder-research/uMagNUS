package engine64

import (
	opencl "github.com/seeder-research/uMagNUS/opencl64"
	util "github.com/seeder-research/uMagNUS/util"
)

// Adaptive Heun solver.
type Heun struct{}

// Adaptive Heun method, can be used as solver.Step
func (he *Heun) Step() {
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

	Err := opencl.Buffer(3, y.Size())
	defer opencl.Recycle(Err)
	opencl.Madd2(Err, dy0, dy, 1., -1.)

	if simpleController(Err, float64(dt), he.AdvOrder(), he.AdvOrder()+1) {
		// step OK
		opencl.Madd3(y, y, dy, dy0, 1, 0.5*dt, -0.5*dt)
		M.normalize()
		setMaxTorque(dy)
	} else {
                // undo bad step
		Time -= Dt_si
		opencl.Madd2(y, y, dy0, 1, -dt)
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
