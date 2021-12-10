package opencl64

import (
	"math"
	"os"
	"testing"
	"unsafe"

	data "github.com/seeder-research/uMagNUS/data64"
	"github.com/seeder-research/uMagNUS/util"
)

// test input data
var in1, in2, in3 *data.Slice

func initTest() {
	inh1 := make([]float64, 1000)
	for i := range inh1 {
		inh1[i] = float64(i)
	}
	in1 = toGPU(inh1)

	inh2 := make([]float64, 100000)
	for i := range inh2 {
		inh2[i] = -float64(i) / 100
	}
	in2 = toGPU(inh2)
}

func toGPU(list []float64) *data.Slice {
	mesh := [3]int{1, 1, len(list)}
	h := sliceFromList([][]float64{list}, mesh)
	d := NewSlice(1, mesh)
	data.Copy(d, h)
	return d
}

func TestMain(m *testing.M) {
	Init(0)
	initTest()
	code := m.Run()
	ReleaseAndClean()
	os.Exit(code)
}

func TestReduceSum(t *testing.T) {
	result := Sum(in1)
	if result != 499500 {
		t.Error("got:", result)
	}
}

func TestReduceDot(t *testing.T) {
	// test for 1 comp
	a := toGPU([]float64{1, 2, 3, 4, 5})
	b := toGPU([]float64{5, 4, 3, -1, 2})
	result := Dot(a, b)
	if result != 5+8+9-4+10 {
		t.Error("got:", result)
	}

	// test for 3 comp
	const N = 32
	mesh := [3]int{1, 1, N}
	c := NewSlice(3, mesh)
	d := NewSlice(3, mesh)
	Memset(c, 1, 2, 3)
	Memset(d, 4, 5, 6)
	result = Dot(c, d)
	if result != N*(4+10+18) {
		t.Error("got:", result)
	}
}

func TestReduceMaxAbs(t *testing.T) {
	result := MaxAbs(in1)
	if result != 999 {
		t.Error("got:", result)
	}
	result = MaxAbs(in2)
	if result != 999.99 {
		t.Error("got:", result)
	}
}

func TestReduceMaxDiff(t *testing.T) {
	// Test on a 1-D array first
	ah1 := make([]float64, 1000)
	bh1 := make([]float64, 1000)
	for i := range ah1 {
		ah1[i] = float64(i)
		bh1[i] = float64(i + i)
	}
	a1 := toGPU(ah1)
	b1 := toGPU(bh1)
	result := MaxDiff(a1, b1)
	if len(result) != 1 {
		t.Error("unexpected result length:", len(result))
	}
	if result[0] != (bh1[999] - ah1[999]) {
		t.Error("got:", result)
	}
	// Check on same 1-D array but change specific
	// entry in b1 first
	SetElem(b1, 0, 1, 10000)
	result = MaxDiff(a1, b1)
	if result[0] != 9999 {
		t.Error("got:", result)
	}

	// Check on 3-D array but each element in list of 32
	// has the same value
	const N = 32
	mesh := [3]int{1, 1, N}
	a1 = NewSlice(3, mesh)
	b1 = NewSlice(3, mesh)
	Memset(a1, 1, 2, 3)
	Memset(b1, 3, 6, 9)
	result = MaxDiff(a1, b1)
	if (result[0] != 2) || (result[1] != 4) || (result[2] != 6) {
		t.Error("got:")
		t.Error("result[0]: ", result[0])
		t.Error("result[1]: ", result[1])
		t.Error("result[2]: ", result[2])
	}
	SetElem(b1, 0, 19, 325)
	SetElem(b1, 1, 19, 48)
	SetElem(b1, 2, 19, 831)
	result = MaxDiff(a1, b1)
	if (result[0] != 324) || (result[1] != 46) || (result[2] != 828) {
		t.Error("got:")
		t.Error("result[0]: ", result[0])
		t.Error("result[1]: ", result[1])
		t.Error("result[2]: ", result[2])
	}
}

func TestReduceMaxVecDiff(t *testing.T) {
	const N = 32
	mesh := [3]int{1, 1, N}
	a := NewSlice(3, mesh)
	b := NewSlice(3, mesh)
	// Set all elements in a and b
	Memset(a, 1, 2, 3)
	Memset(b, 4, 5, 6)
	result := MaxVecDiff(a, b)
	if result != math.Sqrt(9+9+9) {
		t.Error("got:", result)
	}
	// Change only one element in b
	SetElem(b, 0, 5, 6)
	SetElem(b, 1, 5, 7)
	SetElem(b, 2, 5, 9)
	result = MaxVecDiff(a, b)
	if result != math.Sqrt(25+25+36) {
		t.Error("got:", result)
	}
}

func TestReduceMaxNorm(t *testing.T) {
	const N = 32
	mesh := [3]int{1, 1, N}
	a := NewSlice(3, mesh)
	Memset(a, 1, 2, 3)
	result := MaxVecNorm(a)
	if result != math.Sqrt(1+4+9) {
		t.Error("got:", result)
	}
	// Change only one element in a
	SetElem(a, 0, 5, 6)
	SetElem(a, 1, 5, 7)
	SetElem(a, 2, 5, 8)
	result = MaxVecNorm(a)
	if result != math.Sqrt(36+49+64) {
		t.Error("got:", result)
	}
}

func sliceFromList(arr [][]float64, size [3]int) *data.Slice {
	ptrs := make([]unsafe.Pointer, len(arr))
	for i := range ptrs {
		util.Argument(len(arr[i]) == prod(size))
		ptrs[i] = unsafe.Pointer(&arr[i][0])
	}
	return data.SliceFromPtrs(size, data.CPUMemory, ptrs)
}
