package engine64

import (
	data "github.com/seeder-research/uMagNUS/data64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
)

// Classical 4th order RK solver.
type RK4 struct{}

func (rk *RK4) Step() {
	m := M.Buffer()
	size := m.Size()

	if FixDt != 0 {
		Dt_si = FixDt
	}

	t0 := Time
	// backup magnetization
	m0 := opencl.Buffer(3, size)
	defer opencl.Recycle(m0)
	data.Copy(m0, m)

	k1, k2, k3, k4 := opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size)

	defer opencl.Recycle(k1)
	defer opencl.Recycle(k2)
	defer opencl.Recycle(k3)
	defer opencl.Recycle(k4)

	h := float64(Dt_si * GammaLL) // internal time step = Dt * gammaLL

	// stage 1
	torqueFn(k1)

	// stage 2
	Time = t0 + (1./2.)*Dt_si
	opencl.Madd2(m, m, k1, 1, (1./2.)*h) // m = m*1 + k1*h/2
	M.normalize()
	torqueFn(k2)

	// stage 3
	opencl.Madd2(m, m0, k2, 1, (1./2.)*h) // m = m0*1 + k2*1/2
	M.normalize()
	torqueFn(k3)

	// stage 4
	Time = t0 + Dt_si
	opencl.Madd2(m, m0, k3, 1, 1.*h) // m = m0*1 + k3*1
	M.normalize()
	torqueFn(k4)

	Err := opencl.Buffer(3, size)
	defer opencl.Recycle(Err)
	opencl.Madd2(Err, k1, k4, 1., -1.)

	if simpleController(Err, float64(h), rk.AdvOrder(), rk.AdvOrder()+1) { // mindt check to avoid infinite loop
		// step OK
		// 4th order solution
		opencl.Madd5(m, m0, k1, k2, k3, k4, 1, (1./6.)*h, (1./3.)*h, (1./3.)*h, (1./6.)*h)
		M.normalize()
		setMaxTorque(k4)
	} else {
		// undo bad step
		Time = t0
		data.Copy(m, m0)
	}
}

func (_ *RK4) Free() {}

func (_ *RK4) EmType() bool {
	return false
}

func (_ *RK4) AdvOrder() int {
	return 4
}

func (_ *RK4) EmOrder() int {
	return -1
}
