package engine64

import (
	"reflect"

	data "github.com/seeder-research/uMagNUS/data64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
)

// Arbitrary physical quantity.
type Quantity interface {
	NComp() int
	EvalTo(dst *data.Slice)
}

func MeshSize() [3]int {
	return Mesh().Size()
}

func SizeOf(q Quantity) [3]int {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Mesh() *data.Mesh
	}); ok {
		return s.Mesh().Size()
	}
	// otherwise: default mesh
	return MeshSize()
}

func AverageOf(q Quantity) []float64 {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		average() []float64
	}); ok {
		return s.average()
	}
	// otherwise: default mesh
	buf := ValueOf(q)
	defer opencl.Recycle(buf)
	return sAverageMagnet(buf)
}

func NameOf(q Quantity) string {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Name() string
	}); ok {
		return s.Name()
	}
	return "unnamed." + reflect.TypeOf(q).String()
}

func UnitOf(q Quantity) string {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Unit() string
	}); ok {
		return s.Unit()
	}
	return "?"
}

func MeshOf(q Quantity) *data.Mesh {
	// quantity defines its own, custom, implementation:
	if s, ok := q.(interface {
		Mesh() *data.Mesh
	}); ok {
		return s.Mesh()
	}
	return Mesh()
}

func ValueOf(q Quantity) *data.Slice {
	// TODO: check for Buffered() implementation
	buf := opencl.Buffer(q.NComp(), SizeOf(q))
	q.EvalTo(buf)
	return buf
}

// Temporary shim to fit Slice into EvalTo
func EvalTo(q interface {
	Slice() (*data.Slice, bool)
}, dst *data.Slice) {
	v, r := q.Slice()
	if r {
		defer opencl.Recycle(v)
	}
	data.Copy(dst, v)
}
