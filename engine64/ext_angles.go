package engine64

import (
	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl"
)

var (
	ext_phi   = NewScalarField("ext_phi", "rad", "Azimuthal angle", SetPhi)
	ext_theta = NewScalarField("ext_theta", "rad", "Polar angle", SetTheta)
)

func SetPhi(dst *data.Slice) {
	opencl.SetPhi(dst, M.Buffer())
}

func SetTheta(dst *data.Slice) {
	opencl.SetTheta(dst, M.Buffer())
}
