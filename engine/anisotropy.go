package engine

// Magnetocrystalline anisotropy.

import (
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

func addUniaxialAnisotropyFrom(dst *data.Slice, M magnetization, Msat, Ku1, Ku2 *RegionwiseScalar, AnisU *RegionwiseVector) {
	if Ku1.nonZero() || Ku2.nonZero() {
		ms := Msat.MSlice()
		defer ms.Recycle()
		ku1 := Ku1.MSlice()
		defer ku1.Recycle()
		ku2 := Ku2.MSlice()
		defer ku2.Recycle()
		u := AnisU.MSlice()
		defer u.Recycle()

		opencl.AddUniaxialAnisotropy2(dst, M.Buffer(), ms, ku1, ku2, u)
	}
}

func addCubicAnisotropyFrom(dst *data.Slice, M magnetization, Msat, Kc1, Kc2, Kc3 *RegionwiseScalar, AnisC1, AnisC2 *RegionwiseVector) {
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
		opencl.AddCubicAnisotropy2(dst, M.Buffer(), ms, kc1, kc2, kc3, c1, c2)
	}
}

func addTensorAnisotropyFrom(dst *data.Slice, M magnetization, Msat, KuXX, KuYY, KuZZ *RegionwiseScalar, AnisUxx, AnisUyy, AnisUzz *RegionwiseVector) {
	if KuXX.nonZero() {
		ms := Msat.MSlice()
		defer ms.Recycle()
		kuXX := KuXX.MSlice()
		defer kuXX.Recycle()
		uXX := AnisUxx.MSlice()
		defer uXX.Recycle()

		opencl.AddUniaxialAnisotropy(dst, M.Buffer(), ms, kuXX, uXX)
	}
	if KuYY.nonZero() {
		ms := Msat.MSlice()
		defer ms.Recycle()
		kuYY := KuYY.MSlice()
		defer kuYY.Recycle()
		uYY := AnisUyy.MSlice()
		defer uYY.Recycle()

		opencl.AddUniaxialAnisotropy(dst, M.Buffer(), ms, kuYY, uYY)
	}
	if KuZZ.nonZero() {
		ms := Msat.MSlice()
		defer ms.Recycle()
		kuZZ := KuZZ.MSlice()
		defer kuZZ.Recycle()
		uZZ := AnisUzz.MSlice()
		defer uZZ.Recycle()

		opencl.AddUniaxialAnisotropy(dst, M.Buffer(), ms, kuZZ, uZZ)
	}
}

func addVCAnisotropyFrom(dst *data.Slice, M magnetization, Msat, VcmaCoeff1, VcmaCoeff2 *RegionwiseScalar, AnisVCMAU1, AnisVCMAU2 *RegionwiseVector) {
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

			opencl.AddVoltageControlledAnisotropy(dst, M.Buffer(), ms, vcmaCoeff, Vapp, c1)
		}
	}
	if VcmaCoeff2.nonZero() {
		if !Vint2.isZero() {
			ms := Msat.MSlice()
			defer ms.Recycle()

			vcmaCoeff := VcmaCoeff2.MSlice()
			defer vcmaCoeff.Recycle()

			c1 := AnisVCMAU2.MSlice()
			defer c1.Recycle()

			vapp, rec := Vint2.Slice()
			if rec {
				defer opencl.Recycle(vapp)
			}
			Vapp := opencl.ToMSlice(vapp)
			defer Vapp.Recycle()

			opencl.AddVoltageControlledAnisotropy(dst, M.Buffer(), ms, vcmaCoeff, Vapp, c1)
		}
	}
}

// Add the anisotropy field to dst
func AddAnisotropyField(dst *data.Slice) {
	addUniaxialAnisotropyFrom(dst, M, Msat, Ku1, Ku2, AnisU)
	addCubicAnisotropyFrom(dst, M, Msat, Kc1, Kc2, Kc3, AnisC1, AnisC2)
	addTensorAnisotropyFrom(dst, M, Msat, KuXX, KuYY, KuZZ, AnisUxx, AnisUyy, AnisUzz)
	addVCAnisotropyFrom(dst, M, Msat, VcmaCoeff1, VcmaCoeff2, AnisVCMAU1, AnisVCMAU2)
}

// Add the anisotropy energy density to dst
func AddAnisotropyEnergyDensity(dst *data.Slice) {
	haveUnixial := Ku1.nonZero() || Ku2.nonZero()
	haveCubic := Kc1.nonZero() || Kc2.nonZero() || Kc3.nonZero()

	if !haveUnixial && !haveCubic {
		return
	}

	buf := opencl.Buffer(B_anis.NComp(), Mesh().Size())
	defer opencl.Recycle(buf)

	// unnormalized magnetization:
	Mf := ValueOf(M_full)
	defer opencl.Recycle(Mf)

	if haveUnixial {
		// 1st
		opencl.Zero(buf)
		addUniaxialAnisotropyFrom(buf, M, Msat, Ku1, sZero, AnisU)
		opencl.AddDotProduct(dst, -1./2., buf, Mf)

		// 2nd
		opencl.Zero(buf)
		addUniaxialAnisotropyFrom(buf, M, Msat, sZero, Ku2, AnisU)
		opencl.AddDotProduct(dst, -1./4., buf, Mf)
	}

	if haveCubic {
		// 1st
		opencl.Zero(buf)
		addCubicAnisotropyFrom(buf, M, Msat, Kc1, sZero, sZero, AnisC1, AnisC2)
		opencl.AddDotProduct(dst, -1./4., buf, Mf)

		// 2nd
		opencl.Zero(buf)
		addCubicAnisotropyFrom(buf, M, Msat, sZero, Kc2, sZero, AnisC1, AnisC2)
		opencl.AddDotProduct(dst, -1./6., buf, Mf)

		// 3nd
		opencl.Zero(buf)
		addCubicAnisotropyFrom(buf, M, Msat, sZero, sZero, Kc3, AnisC1, AnisC2)
		opencl.AddDotProduct(dst, -1./8., buf, Mf)
	}
}

// Returns anisotropy energy in joules.
func GetAnisotropyEnergy() float64 {
	buf := opencl.Buffer(1, Mesh().Size())
	defer opencl.Recycle(buf)

	opencl.Zero(buf)
	AddAnisotropyEnergyDensity(buf)
	return cellVolume() * float64(opencl.Sum(buf))
}
