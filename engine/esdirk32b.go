package engine

import (
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

// Explicit singly diagonal implicit Rnge-Kutta (ESDIRK) solver.
// ESDIRK32b
// 2nd order, 4 stages per step, adaptive step.
// Anne Kværnø, Singly diagonal implicit Runge-Kutta methods
// with explicit first stage," BIT Numerical Mathematics vol. 44,
// 489-502, 2004.
//
// Advance using y{n+1}
//	k1 = f(tn, yn)
//	k2 = f(tn + 0.5857864376 h, yn + 0.2928932188 h k1 + 0.2928932188 h k2)
//	k3 = f(tn + h, yn + 0.353553390567523 h k1 + 0.353553390632477 h k2 + 0.2928932188 h k3)
//	y{n+1}  = yn + 0.353553390567523 h k1 + 0.353553390632477 h k2 + 0.2928932188 h k3 // 2nd order
//	k4 = f(tn + h, yn + 0.215482203122508 h k1 + 0.686886723913539 h k2 - 0.195262145836047 h k3 + 0.2928932188 h k4)
//	z{n+1} = yn + 0.215482203122508 h k1 + 0.686886723913539 h k2 - 0.195262145836047 h k3 + 0.2928932188 h k4 // 3rd order
type ESDIRK32B struct {
	k1 *data.Slice // torque at end of step is kept for beginning of next step
}

func (esdirk *ESDIRK32B) Step() {
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
	m0 := opencl.Buffer(3, size)
	defer opencl.Recycle(m0)
	data.Copy(m0, m)

	k2, k3, k4 := opencl.Buffer(3, size), opencl.Buffer(3, size), opencl.Buffer(3, size)
	defer opencl.Recycle(k2)
	defer opencl.Recycle(k3)
	defer opencl.Recycle(k4)

	h := float32(Dt_si * GammaLL) // internal time step = Dt * gammaLL

	// there is no explicit stage 1: k1 from previous step

	// stage 2
	Time = t0 + (0.5857864376)*Dt_si
	opencl.Madd2(m, m0, esdirk.k1, 1, (0.2928932188)*h) // m{try} = m*1 + k1*0.2928932188
	M.normalize()
	m_ := opencl.Buffer(3, size)
	defer opencl.Recycle(m_)
	data.Copy(m_, m)
	data.Copy(k2, esdirk.k1)
	_, _, _ = fixedPtIterations((0.2928932188)*h, m_, k2)

	// stage 3
	Time = t0 + Dt_si
	opencl.Madd3(m, m0, esdirk.k1, k2, 1, (0.353553390567523)*h, (0.353553390632477)*h) // m = m0*1 + k1*0.353553390567523 + k2*0.353553390632477
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k3, k2)
	_, _, _ = fixedPtIterations((0.2928932188)*h, m_, k3)

	// stage 4 (for estimating error)
	Time = t0 + Dt_si
	opencl.Madd4(m, m0, esdirk.k1, k2, k3, 1, (0.215482203122508)*h, (0.686886723913539)*h, (-0.195262145836047)*h) // m = m0*1 + k1*0.215482203122508 + k2*0.686886723913539- k3*0.195262145836047
	M.normalize()
	data.Copy(m_, m)
	data.Copy(k4, k3)
	_, _, _ = fixedPtIterations((0.2928932188)*h, m_, k4)

	// 2nd order solution
	opencl.Madd3(m_, esdirk.k1, k2, k3, (0.353553390567523), (0.353553390632477), (0.2928932188)) // m = m0*1 + k1*0.353553390567523 + k2*0.353553390632477 + k3*0.2928932188
	opencl.Madd2(m, m0, m_, 1, h)
	M.normalize()

	// error estimate
	Time = t0 + Dt_si
	Err := k2 // re-use k2 as error
	// difference of 3rd and 2nd order torque without explicitly storing them first
	opencl.Madd4(Err, esdirk.k1, k2, k3, k4, (-0.138071187445015), (0.333333333281062), (-0.488155364636047), (-0.2928932188))

	gustafssonController(Err, m_, esdirk.k1, m0, t0, float64(h), esdirk.AdvOrder(), true)
}

func (esdirk *ESDIRK32B) Free() {
	esdirk.k1.Free()
	esdirk.k1 = nil
}

func (_ *ESDIRK32B) EmType() bool {
	return true
}

func (_ *ESDIRK32B) AdvOrder() int {
	return 2
}

func (_ *ESDIRK32B) EmOrder() int {
	return 3
}
