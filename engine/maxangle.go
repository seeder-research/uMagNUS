package engine

import (
	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl"
)

var (
	MaxAngle  = NewScalarValue("MaxAngle", "rad", "maximum angle between neighboring spins", GetMaxAngle)
	SpinAngle = NewScalarField("spinAngle", "rad", "Angle between neighboring spins", SetSpinAngle)
)

func SetSpinAngle(dst *data.Slice) {
	opencl.SetMaxAngle(dst, M.Buffer(), lex2.Gpu(), regions.Gpu(), M.Mesh())
}

func GetMaxAngle() float64 {
	s := ValueOf(SpinAngle)
	defer opencl.Recycle(s)
	return float64(opencl.MaxAbs(s)) // just a max would be fine, but not currently implemented
}
