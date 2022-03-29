package opencl64

import (
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
)

// Slice + scalar multiplier.
type MSlice struct {
	arr *data.Slice
	mul []float64
}

func ToMSlice(s *data.Slice) MSlice {
	return MSlice{
		arr: s,
		mul: ones(s.NComp()),
	}
}

func MakeMSlice(arr *data.Slice, mul []float64) MSlice {
	return MSlice{arr, mul}
}

func (m MSlice) Size() [3]int {
	return m.arr.Size()
}

func (m MSlice) Len() int {
	return m.arr.Len()
}

func (m MSlice) DevPtr(c int) unsafe.Pointer {
	return m.arr.DevPtr(c)
}

func (m MSlice) Mul(c int) float64 {
	return float64(m.mul[c])
}

func (m MSlice) GetSlicePtr(c int) *data.Slice {
	return m.arr
}

func (m MSlice) SetMul(c int, mul float64) {
	m.mul[c] = float64(mul)
}

func (m MSlice) Recycle() {
	if m.arr != nil {
		Recycle(m.arr)
		m.arr = nil
	}
}

func (m MSlice) SetEvent(index int, event *cl.Event) {
	m.arr.SetEvent(index, event)
}

func (m MSlice) GetEvent(index int) *cl.Event {
	return m.arr.GetEvent(index)
}

var _ones = [4]float64{1, 1, 1, 1}

func ones(n int) []float64 {
	return _ones[:n]

}
