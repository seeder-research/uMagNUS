package engine

// Magnetocrystalline anisotropy.

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

// Anisotropy variables
var (
	Ku1        = NewScalarParam("Ku1", "J/m3", "1st order uniaxial anisotropy constant")
	Ku2        = NewScalarParam("Ku2", "J/m3", "2nd order uniaxial anisotropy constant")
	Kc1        = NewScalarParam("Kc1", "J/m3", "1st order cubic anisotropy constant")
	Kc2        = NewScalarParam("Kc2", "J/m3", "2nd order cubic anisotropy constant")
	Kc3        = NewScalarParam("Kc3", "J/m3", "3rd order cubic anisotropy constant")
	AnisU      = NewVectorParam("anisU", "", "Uniaxial anisotropy direction")
	AnisC1     = NewVectorParam("anisC1", "", "Cubic anisotropy direction #1")
	AnisC2     = NewVectorParam("anisC2", "", "Cubic anisotorpy directon #2")
	B_anis     = NewVectorField("B_anis", "T", "Anisotropy field", AddAnisotropyField)
	Edens_anis = NewScalarField("Edens_anis", "J/m3", "Anisotropy energy density", AddAnisotropyEnergyDensity)
	E_anis     = NewScalarValue("E_anis", "J", "total anisotropy energy", GetAnisotropyEnergy)
	KuXX       = NewScalarParam("KuXX", "J/m3", "1st order uniaxial anisotropy constant (XX) for monodomain shape anisotropy")
	KuYY       = NewScalarParam("KuYY", "J/m3", "1st order uniaxial anisotropy constant (YY) for monodomain shape anisotropy")
	KuZZ       = NewScalarParam("KuZZ", "J/m3", "1st order uniaxial anisotropy constant (ZZ) for monodomain shape anisotropy")
	AnisUxx    = NewVectorParam("anisUxx", "", "Uniaxial anisotropy direction (XX) for monodomain shape anisotropy")
	AnisUyy    = NewVectorParam("anisUyy", "", "Uniaxial anisotropy direction (YY) for monodomain shape anisotropy")
	AnisUzz    = NewVectorParam("anisUzz", "", "Uniaxial anisotropy direction (ZZ) for monodomain shape anisotropy")
)

var (
	sZero = NewScalarParam("_zero", "", "utility zero parameter")
)

func init() {
	registerEnergy(GetAnisotropyEnergy, AddAnisotropyEnergyDensity)
}

func addUniaxialAnisotropyFrom(dst *data.Slice, M magnetization, Msat, Ku1, Ku2 *RegionwiseScalar, AnisU *RegionwiseVector, q *cl.CommandQueue, events []*cl.Event) {
	// prefer to be synchronous due to the recycle operations??
	if Ku1.nonZero() || Ku2.nonZero() {
		ms := Msat.MSlice()
		defer ms.Recycle()
		ku1 := Ku1.MSlice()
		defer ku1.Recycle()
		ku2 := Ku2.MSlice()
		defer ku2.Recycle()
		u := AnisU.MSlice()
		defer u.Recycle()

		opencl.AddUniaxialAnisotropy2(dst, M.Buffer(), ms, ku1, ku2, u, q, events)
	}
}

func addCubicAnisotropyFrom(dst *data.Slice, M magnetization, Msat, Kc1, Kc2, Kc3 *RegionwiseScalar, AnisC1, AnisC2 *RegionwiseVector, q *cl.CommandQueue, events []*cl.Event) {
	// prefer to be synchronous due to the recycle operations??
	if Kc1.nonZero() || Kc2.nonZero() || Kc3.nonZero() {
		ms := Msat.MSlice()
		defer ms.Recycle()

		kc1 := Kc1.MSlice()
		defer kc1.Recycle()

		kc2 := Kc2.MSlice()
		defer kc2.Recycle()

		kc3 := Kc3.MSlice()
		defer kc3.Recycle()

		c1 := AnisC1.MSlice()
		defer c1.Recycle()

		c2 := AnisC2.MSlice()
		defer c2.Recycle()
		opencl.AddCubicAnisotropy2(dst, M.Buffer(), ms, kc1, kc2, kc3, c1, c2, q, events)
	}
}

func addTensorAnisotropyFrom(dst *data.Slice, M magnetization, Msat, KuXX, KuYY, KuZZ *RegionwiseScalar, AnisUxx, AnisUyy, AnisUzz *RegionwiseVector, q *cl.CommandQueue, events []*cl.Event) {
	// prefer to be synchronous due to the recycle operations??
	if KuXX.nonZero() {
		ms := Msat.MSlice()
		defer ms.Recycle()
		kuXX := KuXX.MSlice()
		defer kuXX.Recycle()
		uXX := AnisUxx.MSlice()
		defer uXX.Recycle()

		opencl.AddUniaxialAnisotropy(dst, M.Buffer(), ms, kuXX, uXX, q, events)
	}
	if KuYY.nonZero() {
		ms := Msat.MSlice()
		defer ms.Recycle()
		kuYY := KuYY.MSlice()
		defer kuYY.Recycle()
		uYY := AnisUyy.MSlice()
		defer uYY.Recycle()

		opencl.AddUniaxialAnisotropy(dst, M.Buffer(), ms, kuYY, uYY, q, events)
	}
	if KuZZ.nonZero() {
		ms := Msat.MSlice()
		defer ms.Recycle()
		kuZZ := KuZZ.MSlice()
		defer kuZZ.Recycle()
		uZZ := AnisUzz.MSlice()
		defer uZZ.Recycle()

		opencl.AddUniaxialAnisotropy(dst, M.Buffer(), ms, kuZZ, uZZ, q, events)
	}
}

func addVCAnisotropyFrom(dst *data.Slice, M magnetization, Msat, VcmaCoeff1, VcmaCoeff2 *RegionwiseScalar, AnisVCMAU1, AnisVCMAU2 *RegionwiseVector, q *cl.CommandQueue, events []*cl.Event) {
	// prefer to be synchronous due to the recycle operations??
	if VcmaCoeff1.nonZero() {
		if !Vint1.isZero() {
			ms := Msat.MSlice()
			defer ms.Recycle()

			vcmaCoeff := VcmaCoeff1.MSlice()
			defer vcmaCoeff.Recycle()

			c1 := AnisVCMAU1.MSlice()
			defer c1.Recycle()

			vapp, rec := Vint1.Slice()
			if rec {
				defer opencl.Recycle(vapp)
			}
			Vapp := opencl.ToMSlice(vapp)
			defer Vapp.Recycle()

			opencl.AddVoltageControlledAnisotropy(dst, M.Buffer(), ms, vcmaCoeff, Vapp, c1, q, events)
		}
	}
	if VcmaCoeff2.nonZero() {
		if !Vint2.isZero() {
			ms := Msat.MSlice()
			defer ms.Recycle() // perform in AddVoltageControlledAnisotropy()??

			vcmaCoeff := VcmaCoeff2.MSlice()
			defer vcmaCoeff.Recycle() // perform in AddVoltageControlledAnisotropy()??

			c1 := AnisVCMAU2.MSlice()
			defer c1.Recycle() // perform in AddVoltageControlledAnisotropy()??

			vapp, rec := Vint2.Slice()
			if rec {
				defer opencl.Recycle(vapp) // perform in AddVoltageControlledAnisotropy()??
			}
			Vapp := opencl.ToMSlice(vapp)
			defer Vapp.Recycle() // perform in AddVoltageControlledAnisotropy()??

			opencl.AddVoltageControlledAnisotropy(dst, M.Buffer(), ms, vcmaCoeff, Vapp, c1, q, events)
		}
	}
}

// Add the anisotropy field to dst
func AddAnisotropyField(dst *data.Slice) {
	// Parallelism is exposed here (i.e., all 4 calls compute results independently, which needs
	// to be merged into dst at the end

	// sync all queues in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish in the beginning of addanisotropyfield: %+v \n", err)
	}

	// checkout queues for executing kernels
	seqQueue := opencl.ClCmdQueue[0]
	uniQueueIdx := opencl.CheckoutQueue()
	uniQueue := opencl.ClCmdQueue[uniQueueIdx]
	defer opencl.CheckinQueue(uniQueueIdx)
	cubicQueueIdx := opencl.CheckoutQueue()
	cubicQueue := opencl.ClCmdQueue[cubicQueueIdx]
	defer opencl.CheckinQueue(cubicQueueIdx)
	vcmaQueueIdx := opencl.CheckoutQueue()
	vcmaQueue := opencl.ClCmdQueue[vcmaQueueIdx]
	defer opencl.CheckinQueue(vcmaQueueIdx)

	// checkout buffers for temporarily storing results
	uniBuffer := opencl.Buffer(dst.NComp(), dst.Size())
	defer opencl.Recycle(uniBuffer)
	cubicBuffer := opencl.Buffer(dst.NComp(), dst.Size())
	defer opencl.Recycle(cubicBuffer)
	vcmaBuffer := opencl.Buffer(dst.NComp(), dst.Size())
	defer opencl.Recycle(vcmaBuffer)

	// addUniaxialAnisotropyFrom(dst, M, Msat, Ku1, Ku2, AnisU)
	ms := Msat.MSlice()
	defer ms.Recycle()
	if Ku1.nonZero() || Ku2.nonZero() {
		ku1 := Ku1.MSlice()
		defer ku1.Recycle()
		ku2 := Ku2.MSlice()
		defer ku2.Recycle()
		u := AnisU.MSlice()
		defer u.Recycle()

		opencl.AddUniaxialAnisotropy2(uniBuffer, M.Buffer(), ms, ku1, ku2, u, uniQueue, nil)
	}

	// addCubicAnisotropyFrom(dst, M, Msat, Kc1, Kc2, Kc3, AnisC1, AnisC2)
	if Kc1.nonZero() || Kc2.nonZero() || Kc3.nonZero() {
		kc1 := Kc1.MSlice()
		defer kc1.Recycle()

		kc2 := Kc2.MSlice()
		defer kc2.Recycle()

		kc3 := Kc3.MSlice()
		defer kc3.Recycle()

		c1 := AnisC1.MSlice()
		defer c1.Recycle()

		c2 := AnisC2.MSlice()
		defer c2.Recycle()
		opencl.AddCubicAnisotropy2(cubicBuffer, M.Buffer(), ms, kc1, kc2, kc3, c1, c2, cubicQueue, nil)
	}

	// addTensorAnisotropyFrom(dst, M, Msat, KuXX, KuYY, KuZZ, AnisUxx, AnisUyy, AnisUzz)
	if KuXX.nonZero() {
		kuXX := KuXX.MSlice()
		defer kuXX.Recycle()
		uXX := AnisUxx.MSlice()
		defer uXX.Recycle()

		opencl.AddUniaxialAnisotropy(vcmaBuffer, M.Buffer(), ms, kuXX, uXX, vcmaQueue, nil)
	}
	if KuYY.nonZero() {
		kuYY := KuYY.MSlice()
		defer kuYY.Recycle()
		uYY := AnisUyy.MSlice()
		defer uYY.Recycle()

		opencl.AddUniaxialAnisotropy(uniBuffer, M.Buffer(), ms, kuYY, uYY, uniQueue, nil)
	}
	if KuZZ.nonZero() {
		kuZZ := KuZZ.MSlice()
		defer kuZZ.Recycle()
		uZZ := AnisUzz.MSlice()
		defer uZZ.Recycle()

		opencl.AddUniaxialAnisotropy(cubicBuffer, M.Buffer(), ms, kuZZ, uZZ, cubicQueue, nil)
	}

	// addVCAnisotropyFrom(dst, M, Msat, VcmaCoeff1, VcmaCoeff2, AnisVCMAU1, AnisVCMAU2)
	if VcmaCoeff1.nonZero() {
		if !Vint1.isZero() {
			vcmaCoeff1 := VcmaCoeff1.MSlice()
			defer vcmaCoeff1.Recycle()

			cvcma1 := AnisVCMAU1.MSlice()
			defer cvcma1.Recycle()

			vapp1, rec := Vint1.Slice()
			if rec {
				defer opencl.Recycle(vapp1)
			}
			Vapp1 := opencl.ToMSlice(vapp1)
			defer Vapp1.Recycle()

			opencl.AddVoltageControlledAnisotropy(vcmaBuffer, M.Buffer(), ms, vcmaCoeff1, Vapp1, cvcma1, vcmaQueue, nil)
		}
	}
	if VcmaCoeff2.nonZero() {
		if !Vint2.isZero() {
			vcmaCoeff2 := VcmaCoeff2.MSlice()
			defer vcmaCoeff2.Recycle()

			cvcma2 := AnisVCMAU2.MSlice()
			defer cvcma2.Recycle()

			vapp2, rec := Vint2.Slice()
			if rec {
				defer opencl.Recycle(vapp2)
			}
			Vapp2 := opencl.ToMSlice(vapp2)
			defer Vapp2.Recycle()

			opencl.AddVoltageControlledAnisotropy(uniBuffer, M.Buffer(), ms, vcmaCoeff2, Vapp2, cvcma2, uniQueue, nil)
		}
	}

	// sync execution queues before merging results
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish at end of addanisotropyfield: %+v \n", err)
	}

	// merge results
	opencl.Madd4(dst, dst, uniBuffer, cubicBuffer, vcmaBuffer, 1.0, 1.0, 1.0, 1.0, []*cl.CommandQueue{uniQueue, cubicQueue, vcmaQueue}, nil)

	// sync queues before returning
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{uniQueue, cubicQueue, vcmaQueue})
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for seqQueue to finish at end of addanisotropyfield: %+v \n", err)
	}
}

// Add the anisotropy energy density to dst
func AddAnisotropyEnergyDensity(dst *data.Slice) {
	haveUnixial := Ku1.nonZero() || Ku2.nonZero()
	haveCubic := Kc1.nonZero() || Kc2.nonZero() || Kc3.nonZero()

	if !haveUnixial && !haveCubic {
		return
	}

	// unnormalized magnetization:
	Mf := ValueOf(M_full)
	defer opencl.Recycle(Mf)

	ms := Msat.MSlice()
	defer ms.Recycle()

	szero := sZero.MSlice()
	defer szero.Recycle()

	uniQueueIdx := opencl.CheckoutQueue()
	uniQueue := opencl.ClCmdQueue[uniQueueIdx]
	defer opencl.CheckinQueue(uniQueueIdx)
	cubicQueueIdx := opencl.CheckoutQueue()
	cubicQueue := opencl.ClCmdQueue[cubicQueueIdx]
	defer opencl.CheckinQueue(cubicQueueIdx)
	seqQueue := opencl.ClCmdQueue[0]

	// sync all queues in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish at the beginning of addanisotropyenergydensity: %+v \n", err)
	}

	if haveUnixial {
		bufA := opencl.Buffer(B_anis.NComp(), Mesh().Size())
		defer opencl.Recycle(bufA)

		u := AnisU.MSlice()
		defer u.Recycle()

		// 1st
		// addUniaxialAnisotropyFrom(buf, M, Msat, Ku1, sZero, AnisU, q1, nil)
		if Ku1.nonZero() {
			ku1 := Ku1.MSlice()
			defer ku1.Recycle()

			opencl.AddUniaxialAnisotropy2(bufA, M.Buffer(), ms, ku1, szero, u, uniQueue, nil)
		}
		opencl.AddDotProduct(dst, -1./2., bufA, Mf, uniQueue, nil)

		// 2nd
		opencl.Zero(bufA)
		opencl.SyncQueues([]*cl.CommandQueue{uniQueue}, []*cl.CommandQueue{seqQueue})
		// addUniaxialAnisotropyFrom(buf, M, Msat, sZero, Ku2, AnisU, q1, nil)
		if Ku2.nonZero() {
			ku2 := Ku2.MSlice()
			defer ku2.Recycle()

			opencl.AddUniaxialAnisotropy2(bufA, M.Buffer(), ms, szero, ku2, u, uniQueue, nil)
		}
		opencl.AddDotProduct(dst, -1./4., bufA, Mf, uniQueue, nil)
		// sync uniQueue to seqQueue
		opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{uniQueue})
	}

	if haveCubic {
		bufB := opencl.Buffer(B_anis.NComp(), Mesh().Size())
		defer opencl.Recycle(bufB)

		c1 := AnisC1.MSlice()
		defer c1.Recycle()

		c2 := AnisC2.MSlice()
		defer c2.Recycle()

		// 1st
		// addCubicAnisotropyFrom(buf, M, Msat, Kc1, sZero, sZero, AnisC1, AnisC2, q1, nil)
		if Kc1.nonZero() {
			kc1 := Kc1.MSlice()
			defer kc1.Recycle()

			opencl.AddCubicAnisotropy2(bufB, M.Buffer(), ms, kc1, szero, szero, c1, c2, cubicQueue, nil)
		}
		opencl.SyncQueues([]*cl.CommandQueue{cubicQueue}, []*cl.CommandQueue{seqQueue})
		opencl.AddDotProduct(dst, -1./4., bufB, Mf, cubicQueue, nil)

		// 2nd
		opencl.Zero(bufB)
		opencl.SyncQueues([]*cl.CommandQueue{cubicQueue}, []*cl.CommandQueue{seqQueue})
		// addCubicAnisotropyFrom(buf, M, Msat, sZero, Kc2, sZero, AnisC1, AnisC2, q1, nil)
		if Kc2.nonZero() {
			kc2 := Kc2.MSlice()
			defer kc2.Recycle()

			opencl.AddCubicAnisotropy2(bufB, M.Buffer(), ms, szero, kc2, szero, c1, c2, cubicQueue, nil)
		}
		opencl.AddDotProduct(dst, -1./6., bufB, Mf, cubicQueue, nil)
		// sync q1 queue

		// 3nd
		opencl.Zero(bufB)
		opencl.SyncQueues([]*cl.CommandQueue{cubicQueue}, []*cl.CommandQueue{seqQueue})
		// addCubicAnisotropyFrom(buf, M, Msat, sZero, sZero, Kc3, AnisC1, AnisC2, q1, nil)
		if Kc3.nonZero() {
			kc3 := Kc3.MSlice()
			defer kc3.Recycle()

			opencl.AddCubicAnisotropy2(bufB, M.Buffer(), ms, szero, szero, kc3, c1, c2, cubicQueue, nil)
		}
		opencl.AddDotProduct(dst, -1./8., bufB, Mf, cubicQueue, nil)
		// sync cubicQueue to seqQueue
		opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{cubicQueue})
	}
	// sync queues before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for seqQueue to finish at end of addanisotropyenergydensity: %+v \n", err)
	}
}

// Returns anisotropy energy in joules.
func GetAnisotropyEnergy() float64 {
	buf := opencl.Buffer(1, Mesh().Size())
	defer opencl.Recycle(buf)
	seqQueue := opencl.ClCmdQueue[0]

	opencl.Zero(buf)
	AddAnisotropyEnergyDensity(buf)
	return cellVolume() * float64(opencl.Sum(buf, seqQueue, nil))
}
