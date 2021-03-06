package engine

// Total energy calculation

import (
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

// TODO: Integrate(Edens)
// TODO: consistent naming SetEdensTotal, ...

var (
	energyTerms []func() float64        // all contributions to total energy
	edensTerms  []func(dst *data.Slice) // all contributions to total energy density (add to dst)
	Edens_total = NewScalarField("Edens_total", "J/m3", "Total energy density", SetTotalEdens)
	E_total     = NewScalarValue("E_total", "J", "total energy", GetTotalEnergy)
)

// add energy term to global energy
func registerEnergy(term func() float64, dens func(*data.Slice)) {
	energyTerms = append(energyTerms, term)
	edensTerms = append(edensTerms, dens)
}

// Returns the total energy in J.
func GetTotalEnergy() float64 {
	E := 0.
	for _, f := range energyTerms {
		E += f()
	}
	checkNaN1(E)
	return E
}

// Set dst to total energy density in J/m3
func SetTotalEdens(dst *data.Slice) {
	opencl.Zero(dst)
	for _, addTerm := range edensTerms {
		addTerm(dst)
	}
}

// volume of one cell in m3
func cellVolume() float64 {
	c := Mesh().CellSize()
	return c[0] * c[1] * c[2]
}

// returns a function that adds to dst the energy density:
// 	prefactor * dot (M_full, field)
func makeEdensAdder(field Quantity, prefactor float64) func(*data.Slice) {
	return func(dst *data.Slice) {
		B := ValueOf(field)
		defer opencl.Recycle(B)
		m := ValueOf(M_full)
		defer opencl.Recycle(m)
		factor := float32(prefactor)
		opencl.AddDotProduct(dst, factor, B, m)
	}
}

// vector dot product
func dot(a, b Quantity) float64 {
	A := ValueOf(a)
	defer opencl.Recycle(A)
	B := ValueOf(b)
	defer opencl.Recycle(B)
	return float64(opencl.Dot(A, B))
}
