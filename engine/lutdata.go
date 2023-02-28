package engine

import (
	"sync"
	"unsafe"

	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

// look-up table for region based parameters
type lut struct {
	gpu_buf opencl.LUTPtrs     // gpu copy of cpu buffer, only transferred when needed
	gpu_ok  bool               // gpu cache up-to date with cpu source?
	cpu_buf [][NREGION]float32 // table data on cpu
	source  updater            // updates cpu data
}

type updater interface {
	update() // updates cpu lookup table
}

func (p *lut) init(nComp int, source updater) {
	p.gpu_buf = make(opencl.LUTPtrs, nComp)
	p.cpu_buf = make([][NREGION]float32, nComp)
	p.source = source
}

// get an up-to-date version of the lookup-table on CPU
func (p *lut) cpuLUT() [][NREGION]float32 {
	p.source.update()
	return p.cpu_buf
}

// get an up-to-date version of the lookup-table on GPU
func (p *lut) gpuLUT() opencl.LUTPtrs {
	p.source.update()
	if !p.gpu_ok {
		// upload to GPU
		p.assureAlloc()
		opencl.ClCmdQueue.Finish() // sync previous kernels, may still be using gpu lut
		var wg sync.WaitGroup
		for c := range p.gpu_buf {
			wg.Add(1)
			opencl.MemCpyHtoD(p.gpu_buf[c], unsafe.Pointer(&p.cpu_buf[c][0]), opencl.SIZEOF_FLOAT32*NREGION, &wg)
		}
		p.gpu_ok = true
		wg.Wait()
	}
	return p.gpu_buf
}

// utility for LUT of single-component data
func (p *lut) gpuLUT1() opencl.LUTPtr {
	util.Assert(len(p.gpu_buf) == 1)
	return opencl.LUTPtr(p.gpuLUT()[0])
}

// all data is 0?
func (p *lut) isZero() bool {
	v := p.cpuLUT()
	for c := range v {
		for i := 0; i < NREGION; i++ {
			if v[c][i] != 0 {
				return false
			}
		}
	}
	return true
}

func (p *lut) nonZero() bool { return !p.isZero() }

func (p *lut) assureAlloc() {
	if p.gpu_buf[0] == nil {
		for i := range p.gpu_buf {
			p.gpu_buf[i] = unsafe.Pointer(opencl.MemAlloc(NREGION * opencl.SIZEOF_FLOAT32))
		}
	}
}

func (b *lut) NComp() int { return len(b.cpu_buf) }

// uncompress the table to a full array with parameter values per cell.
func (p *lut) Slice() (*data.Slice, bool) {
	gpu := p.gpuLUT()
	b := opencl.Buffer(p.NComp(), Mesh().Size())
	for c := 0; c < p.NComp(); c++ {
		opencl.RegionDecode(b.Comp(c), opencl.LUTPtr(gpu[c]), regions.Gpu())
	}
	return b, true
}
