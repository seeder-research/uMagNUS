package opencl

import (
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
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

func (m MSlice) NComp() int {
	return m.arr.NComp()
}

func (m MSlice) DevPtr(c int) unsafe.Pointer {
	return m.arr.DevPtr(c)
}

func (m MSlice) Mul(c int) float32 {
	return float32(m.mul[c])
}

func (m MSlice) GetSlicePtr() *data.Slice {
	return m.arr
}

func (m MSlice) SetMul(c int, mul float32) {
	m.mul[c] = float64(mul)
}

func (m MSlice) Recycle() {
	if m.arr != nil {
		Recycle(m.arr)
		m.arr = nil
	}
}

func (m MSlice) RLock() {
	if m.arr == nil {
		return
	}
	for c := 0; c < m.NComp(); c++ {
		m.arr.RLock(c)
	}
}

func (m MSlice) RUnlock() {
	if m.arr == nil {
		return
	}
	for c := 0; c < m.NComp(); c++ {
		m.arr.RUnlock(c)
	}
}

func (m MSlice) Lock() {
	if m.arr == nil {
		return
	}
	for c := 0; c < m.NComp(); c++ {
		m.arr.Lock(c)
	}
}

func (m MSlice) Unlock() {
	if m.arr == nil {
		return
	}
	for c := 0; c < m.NComp(); c++ {
		m.arr.Unlock(c)
	}
}

func (m MSlice) SetEvent(index int, event *cl.Event) {
	m.arr.SetEvent(index, event)
}

func (m MSlice) GetEvent(index int) *cl.Event {
	return m.arr.GetEvent(index)
}

// Sets the rdEvent of the slice
func (m MSlice) SetReadEvents(index int, eventList []*cl.Event) {
	m.arr.SetReadEvents(index, eventList)
}

// Insert a cl.Event to rdEvent of the slice
func (m MSlice) InsertReadEvent(index int, event *cl.Event) {
	m.arr.InsertReadEvent(index, event)
}

// Remove a cl.Event from rdEvent of the slice
func (m MSlice) RemoveReadEvent(index int, event *cl.Event) {
	m.arr.RemoveReadEvent(index, event)
}

// Returns rdEvent of the slice as a slice
func (m MSlice) GetReadEvents(index int) []*cl.Event {
	return m.arr.GetReadEvents(index)
}

// Returns all events of the slice (for syncing kernels writing to the slice)
func (m MSlice) GetAllEvents(index int) []*cl.Event {
	return m.arr.GetAllEvents(index)
}

var _ones = [4]float64{1, 1, 1, 1}

func ones(n int) []float64 {
	return _ones[:n]

}
