package data64

// Slice stores N-component GPU or host data.

import (
	"bytes"
	"fmt"
	"github.com/seeder-research/uMagNUS/cl"
	"github.com/seeder-research/uMagNUS/util"
	"log"
	"reflect"
	"unsafe"
)

const SIZEOF_FLOAT32 = 4
const SIZEOF_FLOAT64 = 8

// Slice is like a [][]float64, but may be stored in GPU or host memory.
type Slice struct {
	ptrs    []unsafe.Pointer
	size    [3]int
	memType int8
	event   []*cl.Event
}

// this package must not depend on OpenCL.
// NOTE: cpyDtoH and cpuHtoD are only needed to support 32-bit builds,
// otherwise, it could be removed in favor of memCpy only.
var (
	memFree, memFreeHost           func(unsafe.Pointer)
	memCpy, memCpyDtoH, memCpyHtoD func(dst, src unsafe.Pointer, bytes int) []*cl.Event
)

// Internal: enables slices on GPU. Called upon opencl init.
func EnableGPU(free, freeHost func(unsafe.Pointer),
	cpy, cpyDtoH, cpyHtoD func(dst, src unsafe.Pointer, bytes int) []*cl.Event) {
	memFree = free
	memFreeHost = freeHost
	memCpy = cpy
	memCpyDtoH = cpyDtoH
	memCpyHtoD = cpyHtoD
}

// Make a CPU Slice with nComp components of size length.
func NewSlice(nComp int, size [3]int) *Slice {
	length := prod(size)
	ptrs := make([]unsafe.Pointer, nComp)
	for i := range ptrs {
		ptrs[i] = unsafe.Pointer(&(make([]float64, length)[0]))
	}
	return SliceFromPtrs(size, CPUMemory, ptrs)
}

func SliceFromArray(data [][]float64, size [3]int) *Slice {
	nComp := len(data)
	length := prod(size)
	ptrs := make([]unsafe.Pointer, nComp)
	for i := range ptrs {
		if len(data[i]) != length {
			panic("size mismatch")
		}
		ptrs[i] = unsafe.Pointer(&data[i][0])
	}
	return SliceFromPtrs(size, CPUMemory, ptrs)
}

// Return a slice without underlying storage. Used to represent a mask containing all 1's.
func NilSlice(nComp int, size [3]int) *Slice {
	return SliceFromPtrs(size, GPUMemory, make([]unsafe.Pointer, nComp))
}

// Internal: construct a Slice using bare memory pointers.
func SliceFromPtrs(size [3]int, memType int8, ptrs []unsafe.Pointer) *Slice {
	length := prod(size)
	nComp := len(ptrs)
	util.Argument(nComp > 0 && length > 0)
	s := new(Slice)
	s.ptrs = make([]unsafe.Pointer, nComp)
	s.size = size
	s.event = make([]*cl.Event, nComp)
	for c := range ptrs {
		s.ptrs[c] = ptrs[c]
		s.event[c] = nil
	}
	s.memType = memType
	return s
}

// Frees the underlying storage and zeros the Slice header to avoid accidental use.
// Slices sharing storage will be invalid after Free. Double free is OK.
func (s *Slice) Free() {
	if s == nil {
		return
	}
	// free storage
	switch s.memType {
	case 0:
		return // already freed
	case GPUMemory:
		for _, ptr := range s.ptrs {
			memFree(ptr)
		}
	//case UnifiedMemory:
	//	for _, ptr := range s.ptrs {
	//		memFreeHost(ptr)
	//	}
	case CPUMemory:
		// nothing to do
	default:
		panic("invalid memory type")
	}
	s.Disable()
}

// INTERNAL. Overwrite struct fields with zeros to avoid
// accidental use after Free.
func (s *Slice) Disable() {
	s.ptrs = s.ptrs[:0]
	s.size = [3]int{0, 0, 0}
	s.memType = 0
}

// value for Slice.memType
const (
	CPUMemory = 1 << 0
	GPUMemory = 1 << 1
	//UnifiedMemory = CPUMemory | GPUMemory
)

// MemType returns the memory type of the underlying storage:
// CPUMemory, GPUMemory or UnifiedMemory
func (s *Slice) MemType() int {
	return int(s.memType)
}

// GPUAccess returns whether the Slice is accessible by the GPU.
// true means it is either stored on GPU or in unified host memory.
func (s *Slice) GPUAccess() bool {
	return s.memType&GPUMemory != 0
}

// CPUAccess returns whether the Slice is accessible by the CPU.
// true means it is stored in host memory.
func (s *Slice) CPUAccess() bool {
	return s.memType&CPUMemory != 0
}

// NComp returns the number of components.
func (s *Slice) NComp() int {
	return len(s.ptrs)
}

// Len returns the number of elements per component.
func (s *Slice) Len() int {
	return prod(s.size)
}

func (s *Slice) Size() [3]int {
	if s == nil {
		return [3]int{0, 0, 0}
	}
	return s.size
}

// Comp returns a single component of the Slice.
func (s *Slice) Comp(i int) *Slice {
	sl := new(Slice)
	sl.ptrs = make([]unsafe.Pointer, 1)
	sl.ptrs[0] = s.ptrs[i]
	sl.size = s.size
	sl.memType = s.memType
	sl.event = []*cl.Event{s.event[i]}
	return sl
}

// DevPtr returns a OpenCL memory object handle to a component.
// Slice must have GPUAccess.
// It is safe to call on a nil slice, returns NULL.
func (s *Slice) DevPtr(component int) unsafe.Pointer {
	if s == nil {
		return nil
	}
	if !s.GPUAccess() {
		panic("slice not accessible by GPU")
	}
	return s.ptrs[component]
}

// Host returns the Slice as a [][]float64 indexed by component, cell number.
// It should have CPUAccess() == true.
func (s *Slice) Host() [][]float64 {
	if !s.CPUAccess() {
		log.Panic("slice not accessible by CPU")
	}
	list := make([][]float64, s.NComp())
	for c := range list {
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&list[c]))
		hdr.Data = uintptr(s.ptrs[c])
		hdr.Len = s.Len()
		hdr.Cap = hdr.Len
	}
	return list
}

// Returns a copy of the Slice, allocated on CPU.
func (s *Slice) HostCopy() *Slice {
	cpy := NewSlice(s.NComp(), s.Size())
	Copy(cpy, s)
	return cpy
}

//Associate a list of events to the slice
func (s *Slice) SetEvents(events []*cl.Event) {
	if s.NComp() != len(events) {
		log.Panic("size of event list does not match number of components in slice")
	}
	s.event = make([]*cl.Event, len(events))
	for idx, event := range events {
		s.event[idx] = event
	}
}

// Associate a cl.Event to the slice
func (s *Slice) SetEvent(index int, event *cl.Event) {
	s.event[index] = event
}

// Returns cl.Event associated with the slice
func (s *Slice) GetEvent(index int) *cl.Event {
	return s.event[index]
}

func Copy(dst, src *Slice) {
	if dst.NComp() != src.NComp() || dst.Len() != src.Len() {
		panic(fmt.Sprintf("slice copy: illegal sizes: dst: %vx%v, src: %vx%v", dst.NComp(), dst.Len(), src.NComp(), src.Len()))
	}
	d, s := dst.GPUAccess(), src.GPUAccess()
	bytes := SIZEOF_FLOAT64 * dst.Len()
	switch {
	default:
		panic("bug")
	case d && s:
		for c := 0; c < dst.NComp(); c++ {
			eventsList := memCpy(dst.DevPtr(c), src.DevPtr(c), bytes)
			dst.SetEvent(c, eventsList[0])
			src.SetEvent(c, eventsList[1])
		}
	case s && !d:
		for c := 0; c < dst.NComp(); c++ {
			eventsList := memCpyDtoH(dst.ptrs[c], src.DevPtr(c), bytes)
			src.SetEvent(c, eventsList[0])
		}
	case !s && d:
		for c := 0; c < dst.NComp(); c++ {
			eventsList := memCpyHtoD(dst.DevPtr(c), src.ptrs[c], bytes)
			dst.SetEvent(c, eventsList[0])
		}
	case !d && !s:
		dst, src := dst.Host(), src.Host()
		for c := range dst {
			copy(dst[c], src[c])
		}
	}
}

// Floats returns the data as 3D array,
// indexed by cell position. Data should be
// scalar (1 component) and have CPUAccess() == true.
func (f *Slice) Scalars() [][][]float64 {
	x := f.Tensors()
	if len(x) != 1 {
		panic(fmt.Sprintf("expecting 1 component, got %v", f.NComp()))
	}
	return x[0]
}

// Vectors returns the data as 4D array,
// indexed by component, cell position. Data should have
// 3 components and have CPUAccess() == true.
func (f *Slice) Vectors() [3][][][]float64 {
	x := f.Tensors()
	if len(x) != 3 {
		panic(fmt.Sprintf("expecting 3 components, got %v", f.NComp()))
	}
	return [3][][][]float64{x[0], x[1], x[2]}
}

// Tensors returns the data as 4D array,
// indexed by component, cell position.
// Requires CPUAccess() == true.
func (f *Slice) Tensors() [][][][]float64 {
	tensors := make([][][][]float64, f.NComp())
	host := f.Host()
	for i := range tensors {
		tensors[i] = reshape(host[i], f.Size())
	}
	return tensors
}

// IsNil returns true if either s is nil or s.pointer[0] == nil
func (s *Slice) IsNil() bool {
	if s == nil {
		return true
	}
	return s.ptrs[0] == nil
}

func (s *Slice) String() string {
	if s == nil {
		return "nil"
	}
	var buf bytes.Buffer
	util.Fprint(&buf, s.Tensors())
	return buf.String()
}

func (s *Slice) Set(comp, ix, iy, iz int, value float64) {
	s.checkComp(comp)
	s.Host()[comp][s.Index(ix, iy, iz)] = float64(value)
}

func (s *Slice) SetVector(ix, iy, iz int, v Vector) {
	i := s.Index(ix, iy, iz)
	for c := range v {
		s.Host()[c][i] = float64(v[c])
	}
}

func (s *Slice) SetScalar(ix, iy, iz int, v float64) {
	s.Host()[0][s.Index(ix, iy, iz)] = float64(v)
}

func (s *Slice) Get(comp, ix, iy, iz int) float64 {
	s.checkComp(comp)
	return float64(s.Host()[comp][s.Index(ix, iy, iz)])
}

func (s *Slice) checkComp(comp int) {
	if comp < 0 || comp >= s.NComp() {
		panic(fmt.Sprintf("slice: invalid component index: %v (number of components=%v)\n", comp, s.NComp()))
	}
}

func (s *Slice) Index(ix, iy, iz int) int {
	return Index(s.Size(), ix, iy, iz)
}

func Index(size [3]int, ix, iy, iz int) int {
	if ix < 0 || ix >= size[X] || iy < 0 || iy >= size[Y] || iz < 0 || iz >= size[Z] {
		panic(fmt.Sprintf("Slice index out of bounds: %v,%v,%v (bounds=%v)\n", ix, iy, iz, size))
	}
	return (iz*size[Y]+iy)*size[X] + ix
}
