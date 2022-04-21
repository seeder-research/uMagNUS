package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	util "github.com/seeder-research/uMagNUS/util"
)

// multiply: dst[i] = a[i] * b[i]
// a and b must have the same number of components
func Mul(dst, a, b *data.Slice) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(a.Len() == N && a.NComp() == nComp && b.Len() == N && b.NComp() == nComp)
	cfg := make1DConf(N)
	eventList := make([]*cl.Event, nComp)
	for c := 0; c < nComp; c++ {
		intEventList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = a.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = b.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		if len(intEventList) == 0 {
			intEventList = nil
		}
		eventList[c] = k_mul_async(dst.DevPtr(c), a.DevPtr(c), b.DevPtr(c), N, cfg,
			intEventList)

		dst.SetEvent(c, eventList[c])
		a.SetEvent(c, eventList[c])
		b.SetEvent(c, eventList[c])
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in mul: %+v \n", err)
		}
	}
}

// divide: dst[i] = a[i] / b[i]
// divide-by-zero yields zero.
func Div(dst, a, b *data.Slice) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(a.Len() == N && a.NComp() == nComp && b.Len() == N && b.NComp() == nComp)
	cfg := make1DConf(N)
	eventList := make([]*cl.Event, nComp)
	for c := 0; c < nComp; c++ {
		intEventList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = a.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = b.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		if len(intEventList) == 0 {
			intEventList = nil
		}
		eventList[c] = k_pointwise_div_async(dst.DevPtr(c), a.DevPtr(c), b.DevPtr(c), N, cfg,
			intEventList)

		dst.SetEvent(c, eventList[c])
		a.SetEvent(c, eventList[c])
		b.SetEvent(c, eventList[c])
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in div: %+v \n", err)
		}
	}
}

// Add: dst = src1 + src2.
func Add(dst, src1, src2 *data.Slice) {
	Madd2(dst, src1, src2, 1, 1)
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2
func Madd2(dst, src1, src2 *data.Slice, factor1, factor2 float64) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp)
	cfg := make1DConf(N)
	eventList := make([]*cl.Event, nComp)
	for c := 0; c < nComp; c++ {
		intEventList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src1.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src2.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		if len(intEventList) == 0 {
			intEventList = nil
		}
		eventList[c] = k_madd2_async(dst.DevPtr(c), src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2, N, cfg,
			intEventList)

		dst.SetEvent(c, eventList[c])
		src1.SetEvent(c, eventList[c])
		src2.SetEvent(c, eventList[c])
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in madd2: %+v \n", err)
		}
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3 * factor3
func Madd3(dst, src1, src2, src3 *data.Slice, factor1, factor2, factor3 float64) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp)
	cfg := make1DConf(N)
	eventList := make([]*cl.Event, nComp)
	for c := 0; c < nComp; c++ {
		intEventList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src1.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src2.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src3.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		if len(intEventList) == 0 {
			intEventList = nil
		}
		eventList[c] = k_madd3_async(dst.DevPtr(c), src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2, src3.DevPtr(c), factor3, N, cfg,
			intEventList)

		dst.SetEvent(c, eventList[c])
		src1.SetEvent(c, eventList[c])
		src2.SetEvent(c, eventList[c])
		src3.SetEvent(c, eventList[c])
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in madd3: %+v \n", err)
		}
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4
func Madd4(dst, src1, src2, src3, src4 *data.Slice, factor1, factor2, factor3, factor4 float64) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp)
	cfg := make1DConf(N)
	eventList := make([]*cl.Event, nComp)
	for c := 0; c < nComp; c++ {
		intEventList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src1.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src2.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src3.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src4.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		if len(intEventList) == 0 {
			intEventList = nil
		}
		eventList[c] = k_madd4_async(dst.DevPtr(c),
			src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2,
			src3.DevPtr(c), factor3,
			src4.DevPtr(c), factor4, N, cfg,
			intEventList)

		dst.SetEvent(c, eventList[c])
		src1.SetEvent(c, eventList[c])
		src2.SetEvent(c, eventList[c])
		src3.SetEvent(c, eventList[c])
		src4.SetEvent(c, eventList[c])
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in madd4: %+v \n", err)
		}
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4 + src5[i] * factor5
func Madd5(dst, src1, src2, src3, src4, src5 *data.Slice, factor1, factor2, factor3, factor4, factor5 float64) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N && src5.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp && src5.NComp() == nComp)
	cfg := make1DConf(N)
	eventList := make([]*cl.Event, nComp)
	for c := 0; c < nComp; c++ {
		intEventList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src1.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src2.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src3.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src4.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src5.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		if len(intEventList) == 0 {
			intEventList = nil
		}
		eventList[c] = k_madd5_async(dst.DevPtr(c),
			src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2,
			src3.DevPtr(c), factor3,
			src4.DevPtr(c), factor4,
			src5.DevPtr(c), factor5, N, cfg,
			intEventList)

		dst.SetEvent(c, eventList[c])
		src1.SetEvent(c, eventList[c])
		src2.SetEvent(c, eventList[c])
		src3.SetEvent(c, eventList[c])
		src4.SetEvent(c, eventList[c])
		src5.SetEvent(c, eventList[c])
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in madd5: %+v \n", err)
		}
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4 + src5[i] * factor5 + src6[i] * factor6
func Madd6(dst, src1, src2, src3, src4, src5, src6 *data.Slice, factor1, factor2, factor3, factor4, factor5, factor6 float64) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N && src5.Len() == N && src6.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp && src5.NComp() == nComp && src6.NComp() == nComp)
	cfg := make1DConf(N)
	eventList := make([]*cl.Event, nComp)
	for c := 0; c < nComp; c++ {
		intEventList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src1.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src2.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src3.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src4.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src5.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src6.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		if len(intEventList) == 0 {
			intEventList = nil
		}
		eventList[c] = k_madd6_async(dst.DevPtr(c),
			src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2,
			src3.DevPtr(c), factor3,
			src4.DevPtr(c), factor4,
			src5.DevPtr(c), factor5,
			src6.DevPtr(c), factor6, N, cfg,
			intEventList)

		dst.SetEvent(c, eventList[c])
		src1.SetEvent(c, eventList[c])
		src2.SetEvent(c, eventList[c])
		src3.SetEvent(c, eventList[c])
		src4.SetEvent(c, eventList[c])
		src5.SetEvent(c, eventList[c])
		src6.SetEvent(c, eventList[c])
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in madd6: %+v \n", err)
		}
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4 + src5[i] * factor5 + src6[i] * factor6 + src7[i] * factor7
func Madd7(dst, src1, src2, src3, src4, src5, src6, src7 *data.Slice, factor1, factor2, factor3, factor4, factor5, factor6, factor7 float64) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N && src5.Len() == N && src6.Len() == N && src7.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp && src5.NComp() == nComp && src6.NComp() == nComp && src7.NComp() == nComp)
	cfg := make1DConf(N)
	eventList := make([]*cl.Event, nComp)
	for c := 0; c < nComp; c++ {
		intEventList := []*cl.Event{}
		tmpEvt := dst.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src1.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src2.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src3.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src4.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src5.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src6.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		tmpEvt = src7.GetEvent(c)
		if tmpEvt != nil {
			intEventList = append(intEventList, tmpEvt)
		}
		if len(intEventList) == 0 {
			intEventList = nil
		}
		eventList[c] = k_madd7_async(dst.DevPtr(c),
			src1.DevPtr(c), factor1,
			src2.DevPtr(c), factor2,
			src3.DevPtr(c), factor3,
			src4.DevPtr(c), factor4,
			src5.DevPtr(c), factor5,
			src6.DevPtr(c), factor6,
			src7.DevPtr(c), factor7, N, cfg,
			intEventList)

		dst.SetEvent(c, eventList[c])
		src1.SetEvent(c, eventList[c])
		src2.SetEvent(c, eventList[c])
		src3.SetEvent(c, eventList[c])
		src4.SetEvent(c, eventList[c])
		src5.SetEvent(c, eventList[c])
		src6.SetEvent(c, eventList[c])
		src7.SetEvent(c, eventList[c])
	}
	if Debug {
		if err := cl.WaitForEvents(eventList); err != nil {
			fmt.Printf("WaitForEvents failed in madd7: %+v \n", err)
		}
	}
}
