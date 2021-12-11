package engine64

// Mangeto-elastic coupling.

import (
	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl"
	"github.com/seeder-research/uMagNUS/util"
)

var (
	B1        = NewScalarParam("B1", "J/m3", "First magneto-elastic coupling constant")
	B2        = NewScalarParam("B2", "J/m3", "Second magneto-elastic coupling constant")
	exx       = NewScalarExcitation("exx", "", "exx component of the strain tensor")
	eyy       = NewScalarExcitation("eyy", "", "eyy component of the strain tensor")
	ezz       = NewScalarExcitation("ezz", "", "ezz component of the strain tensor")
	exy       = NewScalarExcitation("exy", "", "exy component of the strain tensor")
	exz       = NewScalarExcitation("exz", "", "exz component of the strain tensor")
	eyz       = NewScalarExcitation("eyz", "", "eyz component of the strain tensor")
	B_mel     = NewVectorField("B_mel", "T", "Magneto-elastic filed", AddMagnetoelasticField)
	F_mel     = NewVectorField("F_mel", "N/m3", "Magneto-elastic force density", GetMagnetoelasticForceDensity)
	Edens_mel = NewScalarField("Edens_mel", "J/m3", "Magneto-elastic energy density", AddMagnetoelasticEnergyDensity)
	E_mel     = NewScalarValue("E_mel", "J", "Magneto-elastic energy", GetMagnetoelasticEnergy)
)

var (
	zeroMel = NewScalarParam("_zeroMel", "", "utility zero parameter")
)

func init() {
	registerEnergy(GetMagnetoelasticEnergy, AddMagnetoelasticEnergyDensity)
}

func AddMagnetoelasticField(dst *data.Slice) {
	haveMel := B1.nonZero() || B2.nonZero()
	if !haveMel {
		return
	}

	Exx := exx.MSlice()
	defer Exx.Recycle()

	Eyy := eyy.MSlice()
	defer Eyy.Recycle()

	Ezz := ezz.MSlice()
	defer Ezz.Recycle()

	Exy := exy.MSlice()
	defer Exy.Recycle()

	Exz := exz.MSlice()
	defer Exz.Recycle()

	Eyz := eyz.MSlice()
	defer Eyz.Recycle()

	b1 := B1.MSlice()
	defer b1.Recycle()

	b2 := B2.MSlice()
	defer b2.Recycle()

	ms := Msat.MSlice()
	defer ms.Recycle()

	opencl.AddMagnetoelasticField(dst, M.Buffer(),
		Exx, Eyy, Ezz,
		Exy, Exz, Eyz,
		b1, b2, ms)
}

func GetMagnetoelasticForceDensity(dst *data.Slice) {
	haveMel := B1.nonZero() || B2.nonZero()
	if !haveMel {
		return
	}

	util.AssertMsg(B1.IsUniform() && B2.IsUniform(), "Magnetoelastic: B1, B2 must be uniform")

	b1 := B1.MSlice()
	defer b1.Recycle()

	b2 := B2.MSlice()
	defer b2.Recycle()

	opencl.GetMagnetoelasticForceDensity(dst, M.Buffer(),
		b1, b2, M.Mesh())
}

func AddMagnetoelasticEnergyDensity(dst *data.Slice) {
	haveMel := B1.nonZero() || B2.nonZero()
	if !haveMel {
		return
	}

	buf := opencl.Buffer(B_mel.NComp(), B_mel.Mesh().Size())
	defer opencl.Recycle(buf)

	// unnormalized magnetization:
	Mf := ValueOf(M_full)
	defer opencl.Recycle(Mf)

	Exx := exx.MSlice()
	defer Exx.Recycle()

	Eyy := eyy.MSlice()
	defer Eyy.Recycle()

	Ezz := ezz.MSlice()
	defer Ezz.Recycle()

	Exy := exy.MSlice()
	defer Exy.Recycle()

	Exz := exz.MSlice()
	defer Exz.Recycle()

	Eyz := eyz.MSlice()
	defer Eyz.Recycle()

	b1 := B1.MSlice()
	defer b1.Recycle()

	b2 := B2.MSlice()
	defer b2.Recycle()

	ms := Msat.MSlice()
	defer ms.Recycle()

	zeromel := zeroMel.MSlice()
	defer zeromel.Recycle()

	// 1st
	opencl.Zero(buf)
	opencl.AddMagnetoelasticField(buf, M.Buffer(),
		Exx, Eyy, Ezz,
		Exy, Exz, Eyz,
		b1, zeromel, ms)
	opencl.AddDotProduct(dst, -1./2., buf, Mf)

	// 1nd
	opencl.Zero(buf)
	opencl.AddMagnetoelasticField(buf, M.Buffer(),
		Exx, Eyy, Ezz,
		Exy, Exz, Eyz,
		zeromel, b2, ms)
	opencl.AddDotProduct(dst, -1./1., buf, Mf)
}

// Returns magneto-ell energy in joules.
func GetMagnetoelasticEnergy() float64 {
	haveMel := B1.nonZero() || B2.nonZero()
	if !haveMel {
		return float64(0.0)
	}

	buf := opencl.Buffer(1, Mesh().Size())
	defer opencl.Recycle(buf)

	opencl.Zero(buf)
	AddMagnetoelasticEnergyDensity(buf)
	return cellVolume() * float64(opencl.Sum(buf))
}