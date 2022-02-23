package engine64

// Minimize follows the steepest descent method as per Exl et al., JAP 115, 17D118 (2014).

import (
	data "github.com/seeder-research/uMagNUS/data64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
)

var (
	DmSamples int     = 10   // number of dm to keep for convergence check
	StopMaxDm float64 = 1e-6 // stop minimizer if sampled dm is smaller than this
)

func init() {
	DeclFunc("Minimize", Minimize, "Use steepest conjugate gradient method to minimize the total energy")
	DeclVar("MinimizerStop", &StopMaxDm, "Stopping max dM for Minimize")
	DeclVar("MinimizerSamples", &DmSamples, "Number of max dM to collect for Minimize convergence check.")
}

// fixed length FIFO. Items can be added but not removed
type fifoRing struct {
	count int
	tail  int // index to put next item. Will loop to 0 after exceeding length
	data  []float64
}

func FifoRing(length int) fifoRing {
	return fifoRing{data: make([]float64, length)}
}

func (r *fifoRing) Add(item float64) {
	r.data[r.tail] = item
	r.count++
	r.tail = (r.tail + 1) % len(r.data)
	if r.count > len(r.data) {
		r.count = len(r.data)
	}
}

func (r *fifoRing) Max() float64 {
	max := r.data[0]
	for i := 1; i < r.count; i++ {
		if r.data[i] > max {
			max = r.data[i]
		}
	}
	return max
}

type Minimizer struct {
	k      *data.Slice // torque saved to calculate time step
	lastDm fifoRing
	h      float64
}

func (mini *Minimizer) Step() {
	m := M.Buffer()
	size := m.Size()

	if mini.k == nil {
		mini.k = opencl.Buffer(3, size)
		torqueFn(mini.k)
	}

	k := mini.k
	h := mini.h

	// save original magnetization
	m0 := opencl.Buffer(3, size)
	defer opencl.Recycle(m0)
	data.Copy(m0, m)

	// make descent
	opencl.Minimize(m, m0, k, h)

	// calculate new torque for next step
	k0 := opencl.Buffer(3, size)
	defer opencl.Recycle(k0)
	data.Copy(k0, k)
	torqueFn(k)
	setMaxTorque(k) // report to user

	// just to make the following readable
	dm := m0
	dk := k0

	// calculate step difference of m and k
	opencl.Madd2(dm, m, m0, 1., -1.)
	opencl.Madd2(dk, k, k0, -1., 1.) // reversed due to LLNoPrecess sign

	// get maxdiff and add to list
	max_dm := opencl.MaxVecNorm(dm)
	mini.lastDm.Add(max_dm)
	setLastErr(mini.lastDm.Max()) // report maxDm to user as LastErr

	// adjust next time step
	var nom, div float64
	if NSteps%2 == 0 {
		nom = opencl.Dot(dm, dm)
		div = opencl.Dot(dm, dk)
	} else {
		nom = opencl.Dot(dm, dk)
		div = opencl.Dot(dk, dk)
	}
	if div != 0. {
		mini.h = nom / div
	} else { // in case of division by zero
		mini.h = 1e-4
	}

	M.normalize()

	// as a convention, time does not advance during relax
	NSteps++
}

func (mini *Minimizer) Free() {
	if mini.k != nil {
		opencl.Recycle(mini.k)
		mini.k = nil
	}
}

func Minimize() {
	Refer("exl2014")
	SanityCheck()
	// Save the settings we are changing...
	prevType := solvertype
	prevFixDt := FixDt
	prevPrecess := Precess
	t0 := Time

	relaxing = true // disable temperature noise

	// ...to restore them later
	defer func() {
		SetSolver(prevType)
		FixDt = prevFixDt
		Precess = prevPrecess
		Time = t0

		relaxing = false
	}()

	Precess = false // disable precession for torque calculation
	// remove previous stepper
	if stepper != nil {
		stepper.Free()
	}

	// set stepper to the minimizer
	mini := Minimizer{
		h:      1e-4,
		k:      nil,
		lastDm: FifoRing(DmSamples)}
	stepper = &mini

	cond := func() bool {
		return (mini.lastDm.count < DmSamples || mini.lastDm.Max() > StopMaxDm)
	}

	RunWhile(cond)
	pause = true
}

func (_ *Minimizer) EmType() bool {
	return false
}

func (_ *Minimizer) AdvOrder() int {
	return -1
}

func (_ *Minimizer) EmOrder() int {
	return -1
}
