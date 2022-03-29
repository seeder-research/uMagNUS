package engine

import (
	"math"

	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

var (
	Ext_TopologicalChargeLattice        = NewScalarValue("ext_topologicalchargelattice", "", "2D topological charge according to Berg and Lüscher", GetTopologicalChargeLattice)
	Ext_TopologicalChargeDensityLattice = NewScalarField("ext_topologicalchargedensitylattice", "1/m2",
		"2D topological charge density according to Berg and Lüscher", SetTopologicalChargeDensityLattice)
)

func SetTopologicalChargeDensityLattice(dst *data.Slice) {
	Refer("Berg1981")
	opencl.SetTopologicalChargeLattice(dst, M.Buffer(), M.Mesh())
}

func GetTopologicalChargeLattice() float64 {
	s := ValueOf(Ext_TopologicalChargeDensityLattice)
	defer opencl.Recycle(s)
	c := Mesh().CellSize()
	N := Mesh().Size()

	return (0.25 * c[X] * c[Y] / math.Pi / float64(N[Z])) * float64(opencl.Sum(s))
}
