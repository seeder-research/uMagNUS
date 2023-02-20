package opencl

import (
	"fmt"
	"sync"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// multiply: dst[i] = a[i] * b[i]
// a and b must have the same number of components
func Mul(dst, a, b *data.Slice) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(a.Len() == N && a.NComp() == nComp && b.Len() == N && b.NComp() == nComp)
	var wg sync.WaitGroup
	wg.Add(nComp)
	for c := 0; c < nComp; c++ {
		if Synchronous {
			mul__(dst, a, b, c, &wg)
		} else {
			go mul__(dst, a, b, c, &wg)
		}
	}
	wg.Wait()
}

func mul__(dst, a, b *data.Slice, idx int, wg_ *sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	if dst.DevPtr(idx) != a.DevPtr(idx) {
		a.RLock(idx)
		defer a.RUnlock(idx)
	}
	if dst.DevPtr(idx) != b.DevPtr(idx) {
		b.RLock(idx)
		defer b.RUnlock(idx)
	}

	N := dst.Len()
	cfg := make1DConf(N)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	defer cmdqueue.Release()
	if err != nil {
		fmt.Printf("mul failed to create command queue: %+v \n", err)
		return
	}

	event := k_mul_async(dst.DevPtr(idx), a.DevPtr(idx), b.DevPtr(idx), N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in mul: %+v \n", err)
	}
}

// divide: dst[i] = a[i] / b[i]
// divide-by-zero yields zero.
func Div(dst, a, b *data.Slice) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(a.Len() == N && a.NComp() == nComp && b.Len() == N && b.NComp() == nComp)

	var wg sync.WaitGroup
	wg.Add(nComp)
	for c := 0; c < nComp; c++ {
		if Synchronous {
			div__(dst, a, b, c, &wg)
		} else {
			go div__(dst, a, b, c, &wg)
		}
	}
	wg.Wait()
}

func div__(dst, a, b *data.Slice, idx int, wg_ *sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	if dst.DevPtr(idx) != a.DevPtr(idx) {
		a.RLock(idx)
		defer a.RUnlock(idx)
	}
	if dst.DevPtr(idx) != b.DevPtr(idx) {
		b.RLock(idx)
		defer b.RUnlock(idx)
	}

	N := dst.Len()
	cfg := make1DConf(N)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("div failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	event := k_pointwise_div_async(dst.DevPtr(idx), a.DevPtr(idx), b.DevPtr(idx), N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in div: %+v \n", err)
	}
}

// Add: dst = src1 + src2.
func Add(dst, src1, src2 *data.Slice) {
	Madd2(dst, src1, src2, 1, 1)
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2
func Madd2(dst, src1, src2 *data.Slice, factor1, factor2 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp)

	var wg sync.WaitGroup
	for c := 0; c < nComp; c++ {
		wg.Add(1)
		if Synchronous {
			madd2__(dst, src1, src2, factor1, factor2, c, &wg)
		} else {
			go madd2__(dst, src1, src2, factor1, factor2, c, &wg)
		}
	}
	wg.Wait()
}

func madd2__(dst, src1, src2 *data.Slice, factor1, factor2 float32, idx int, wg_ *sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	if dst.DevPtr(idx) != src1.DevPtr(idx) {
		src1.RLock(idx)
		defer src1.RUnlock(idx)
	}
	if dst.DevPtr(idx) != src2.DevPtr(idx) {
		src2.RLock(idx)
		defer src2.RUnlock(idx)
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("madd2 failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_madd2_async(dst.DevPtr(idx),
		src1.DevPtr(idx), factor1,
		src2.DevPtr(idx), factor2, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents in madd2 failed: %+v", err)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3 * factor3
func Madd3(dst, src1, src2, src3 *data.Slice, factor1, factor2, factor3 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp)

	var wg sync.WaitGroup
	wg.Add(nComp)
	for c := 0; c < nComp; c++ {
		if Synchronous {
			madd3__(dst, src1, src2, src3, factor1, factor2, factor3, c, &wg)
		} else {
			go madd3__(dst, src1, src2, src3, factor1, factor2, factor3, c, &wg)
		}
	}
	wg.Wait()
}

func madd3__(dst, src1, src2, src3 *data.Slice, factor1, factor2, factor3 float32, idx int, wg_ *sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	if dst.DevPtr(idx) != src1 {
		src1.RLock(idx)
		defer src1.RUnlock(idx)
	}
	if dst != src2 {
		src2.RLock(idx)
		defer src2.RUnlock(idx)
	}
	if dst != src3 {
		src3.RLock(idx)
		defer src3.RUnlock(idx)
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("madd3 failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_madd3_async(dst.DevPtr(idx),
		src1.DevPtr(idx), factor1,
		src2.DevPtr(idx), factor2,
		src3.DevPtr(idx), factor3, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents in madd3 failed: %+v", err)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4
func Madd4(dst, src1, src2, src3, src4 *data.Slice, factor1, factor2, factor3, factor4 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp)

	var wg sync.WaitGroup
	for c := 0; c < nComp; c++ {
		wg.Add(1)
		if Synchronous {
			madd4__(dst, src1, src2, src3, src4, factor1, factor2, factor3, factor4, c, &wg)
		} else {
			go madd4__(dst, src1, src2, src3, src4, factor1, factor2, factor3, factor4, c, &wg)
		}
	}
	wg.Wait()
}

func madd4__(dst, src1, src2, src3, src4 *data.Slice, factor1, factor2, factor3, factor4 float32, idx int, wg_ *sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	if dst != src1 {
		src1.RLock(idx)
		defer src1.RUnlock(idx)
	}
	if dst != src2 {
		src2.RLock(idx)
		defer src2.RUnlock(idx)
	}
	if dst != src3 {
		src3.RLock(idx)
		defer src3.RUnlock(idx)
	}
	if dst != src4 {
		src4.RLock(idx)
		defer src4.RUnlock(idx)
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("madd4 failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_madd4_async(dst.DevPtr(idx),
		src1.DevPtr(idx), factor1,
		src2.DevPtr(idx), factor2,
		src3.DevPtr(idx), factor3,
		src4.DevPtr(idx), factor4, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents in madd4 failed: %+v", err)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4 + src5[i] * factor5
func Madd5(dst, src1, src2, src3, src4, src5 *data.Slice, factor1, factor2, factor3, factor4, factor5 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N && src5.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp && src5.NComp() == nComp)

	var wg sync.WaitGroup
	for c := 0; c < nComp; c++ {
		wg.Add(1)
		if Synchronous {
			madd5__(dst, src1, src2, src3, src4, src5, factor1, factor2, factor3, factor4, factor5, c, &wg)
		} else {
			go madd5__(dst, src1, src2, src3, src4, src5, factor1, factor2, factor3, factor4, factor5, c, &wg)
		}
	}
	wg.Wait()
}

func madd5__(dst, src1, src2, src3, src4, src5 *data.Slice, factor1, factor2, factor3, factor4, factor5 float32, idx int, wg_ *sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	if dst != src1 {
		src1.RLock(idx)
		defer src1.RUnlock(idx)
	}
	if dst != src2 {
		src2.RLock(idx)
		defer src2.RUnlock(idx)
	}
	if dst != src3 {
		src3.RLock(idx)
		defer src3.RUnlock(idx)
	}
	if dst != src4 {
		src4.RLock(idx)
		defer src4.RUnlock(idx)
	}
	if dst != src5 {
		src5.RLock(idx)
		defer src5.RUnlock(idx)
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("madd5 failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_madd5_async(dst.DevPtr(idx),
		src1.DevPtr(idx), factor1,
		src2.DevPtr(idx), factor2,
		src3.DevPtr(idx), factor3,
		src4.DevPtr(idx), factor4,
		src5.DevPtr(idx), factor5, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents in madd5 failed: %+v", err)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4 + src5[i] * factor5 + src6[i] * factor6
func Madd6(dst, src1, src2, src3, src4, src5, src6 *data.Slice, factor1, factor2, factor3, factor4, factor5, factor6 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N && src5.Len() == N && src6.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp && src5.NComp() == nComp && src6.NComp() == nComp)

	var wg sync.WaitGroup
	for c := 0; c < nComp; c++ {
		wg.Add(1)
		if Synchronous {
			madd6__(dst, src1, src2, src3, src4, src5, src6, factor1, factor2, factor3, factor4, factor5, factor6, c, &wg)
		} else {
			go madd6__(dst, src1, src2, src3, src4, src5, src6, factor1, factor2, factor3, factor4, factor5, factor6, c, &wg)
		}
	}
	wg.Wait()
}

func madd6__(dst, src1, src2, src3, src4, src5, src6 *data.Slice, factor1, factor2, factor3, factor4, factor5, factor6 float32, idx int, wg_ *sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	if dst != src1 {
		src1.RLock(idx)
		defer src1.RUnlock(idx)
	}
	if dst != src2 {
		src2.RLock(idx)
		defer src2.RUnlock(idx)
	}
	if dst != src3 {
		src3.RLock(idx)
		defer src3.RUnlock(idx)
	}
	if dst != src4 {
		src4.RLock(idx)
		defer src4.RUnlock(idx)
	}
	if dst != src5 {
		src5.RLock(idx)
		defer src5.RUnlock(idx)
	}
	if dst != src6 {
		src6.RLock(idx)
		defer src6.RUnlock(idx)
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("madd6 failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_madd6_async(dst.DevPtr(idx),
		src1.DevPtr(idx), factor1,
		src2.DevPtr(idx), factor2,
		src3.DevPtr(idx), factor3,
		src4.DevPtr(idx), factor4,
		src5.DevPtr(idx), factor5,
		src6.DevPtr(idx), factor6, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents in madd6 failed: %+v", err)
	}
}

// multiply-add: dst[i] = src1[i] * factor1 + src2[i] * factor2 + src3[i] * factor3 + src4[i] * factor4 + src5[i] * factor5 + src6[i] * factor6 + src7[i] * factor7
func Madd7(dst, src1, src2, src3, src4, src5, src6, src7 *data.Slice, factor1, factor2, factor3, factor4, factor5, factor6, factor7 float32) {
	N := dst.Len()
	nComp := dst.NComp()
	util.Assert(src1.Len() == N && src2.Len() == N && src3.Len() == N && src4.Len() == N && src5.Len() == N && src6.Len() == N && src7.Len() == N)
	util.Assert(src1.NComp() == nComp && src2.NComp() == nComp && src3.NComp() == nComp && src4.NComp() == nComp && src5.NComp() == nComp && src6.NComp() == nComp && src7.NComp() == nComp)

	var wg sync.WaitGroup
	for c := 0; c < nComp; c++ {
		wg.Add(1)
		if Synchronous {
			madd7__(dst, src1, src2, src3, src4, src5, src6, src7, factor1, factor2, factor3, factor4, factor5, factor6, factor7, c, &wg)
		} else {
			go madd7__(dst, src1, src2, src3, src4, src5, src6, src7, factor1, factor2, factor3, factor4, factor5, factor6, factor7, c, &wg)
		}
	}
	wg.Wait()
}

func madd7__(dst, src1, src2, src3, src4, src5, src6, src7 *data.Slice, factor1, factor2, factor3, factor4, factor5, factor6, factor7 float32, idx int, wg_ *sync.WaitGroup) {
	dst.Lock(idx)
	defer dst.Unlock(idx)
	if dst != src1 {
		src1.RLock(idx)
		defer src1.RUnlock(idx)
	}
	if dst != src2 {
		src2.RLock(idx)
		defer src2.RUnlock(idx)
	}
	if dst != src3 {
		src3.RLock(idx)
		defer src3.RUnlock(idx)
	}
	if dst != src4 {
		src4.RLock(idx)
		defer src4.RUnlock(idx)
	}
	if dst != src5 {
		src5.RLock(idx)
		defer src5.RUnlock(idx)
	}
	if dst != src6 {
		src6.RLock(idx)
		defer src6.RUnlock(idx)
	}
	if dst != src7 {
		src7.RLock(idx)
		defer src7.RUnlock(idx)
	}

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("madd7 failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	N := dst.Len()
	cfg := make1DConf(N)

	event := k_madd7_async(dst.DevPtr(idx),
		src1.DevPtr(idx), factor1,
		src2.DevPtr(idx), factor2,
		src3.DevPtr(idx), factor3,
		src4.DevPtr(idx), factor4,
		src5.DevPtr(idx), factor5,
		src6.DevPtr(idx), factor6,
		src7.DevPtr(idx), factor7, N, cfg,
		cmdqueue, nil)

	wg_.Done()

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents in madd7 failed: %+v", err)
	}
}
