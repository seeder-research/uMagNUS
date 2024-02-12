package engine

import (
	"fmt"
	"reflect"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

var M magnetization // reduced magnetization (unit length)

func init() { DeclLValue("m", &M, `Reduced magnetization (unit length)`) }

// Special buffered quantity to store magnetization
// makes sure it's normalized etc.
type magnetization struct {
	buffer_ *data.Slice
}

func (m *magnetization) Mesh() *data.Mesh    { return Mesh() }
func (m *magnetization) NComp() int          { return 3 }
func (m *magnetization) Name() string        { return "m" }
func (m *magnetization) Unit() string        { return "" }
func (m *magnetization) Buffer() *data.Slice { return m.buffer_ } // todo: rename Gpu()?

func (m *magnetization) Comp(c int) ScalarField  { return Comp(m, c) }
func (m *magnetization) SetValue(v interface{})  { m.SetInShape(nil, v.(Config)) }
func (m *magnetization) InputType() reflect.Type { return reflect.TypeOf(Config(nil)) }
func (m *magnetization) Type() reflect.Type      { return reflect.TypeOf(new(magnetization)) }
func (m *magnetization) Eval() interface{}       { return m }
func (m *magnetization) average() []float64      { return sAverageMagnet(M.Buffer()) }
func (m *magnetization) Average() data.Vector    { return unslice(m.average()) }

func (m *magnetization) normalize() {
	// sync in the beginning
	seqQueue := opencl.ClCmdQueue[0]
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting all queues to finish before magnetization.normalize: %+v \n", err)
	}
	opencl.Normalize(m.Buffer(), geometry.Gpu(), seqQueue, nil)
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting queue to finish after magnetization.normalize: %+v \n", err)
	}
}

// allocate storage (not done by init, as mesh size may not yet be known then)
func (m *magnetization) alloc() {
	m.buffer_ = opencl.NewSlice(3, m.Mesh().Size())
	m.Set(RandomMag()) // sane starting config
}

func (b *magnetization) SetArray(src *data.Slice) {
	if src.Size() != b.Mesh().Size() {
		src = data.Resample(src, b.Mesh().Size())
	}
	data.Copy(b.Buffer(), src)
	seqQueue := opencl.ClCmdQueue[0]
	if src.CPUAccess() {
		opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, []*cl.CommandQueue{opencl.H2DQueue})
	}
	b.normalize()
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting queue to finish after magnetization.setarray: %+v \n", err)
	}
}

func (m *magnetization) Set(c Config) {
	checkMesh()
	m.SetInShape(nil, c)
}

func (m *magnetization) LoadFile(fname string) {
	m.SetArray(LoadFile(fname))
}

func (m *magnetization) Slice() (s *data.Slice, recycle bool) {
	return m.Buffer(), false
}

func (m *magnetization) EvalTo(dst *data.Slice) {
	data.Copy(dst, m.buffer_)
	// sync before returning
	if err0, err1, err2 := opencl.WaitAllQueuesToFinish(), opencl.H2DQueue.Finish(), opencl.D2HQueue.Finish(); (err0 != nil) || (err1 != nil) || (err2 != nil) {
		fmt.Printf("error waiting for all queues to finish in magnetization.evalto() 0: %+v \n", err0)
		fmt.Printf("error waiting for all queues to finish in magnetization.evalto() 1: %+v \n", err1)
		fmt.Printf("error waiting for all queues to finish in magnetization.evalto() 2: %+v \n", err2)
	}
}

func (m *magnetization) Region(r int) *vOneReg { return vOneRegion(m, r) }

func (m *magnetization) String() string {
	tmp := m.Buffer().HostCopy()
	// sync before continuing
	if m.Buffer().GPUAccess() {
		if err := opencl.D2HQueue.Finish(); err != nil {
			fmt.Printf("error waiting for queue to finish in magnetization.string(): %+v \n", err)
		}
	}
	return util.Sprint(tmp)
}

// Set the value of one cell.
func (m *magnetization) SetCell(ix, iy, iz int, v data.Vector) {
	for c := 0; c < 3; c++ {
		opencl.SetCell(m.Buffer(), c, ix, iy, iz, float32(v[c]))
	}
}

// Get the value of one cell.
func (m *magnetization) GetCell(ix, iy, iz int) data.Vector {
	mx := float64(opencl.GetCell(m.Buffer(), X, ix, iy, iz))
	my := float64(opencl.GetCell(m.Buffer(), Y, ix, iy, iz))
	mz := float64(opencl.GetCell(m.Buffer(), Z, ix, iy, iz))
	return Vector(mx, my, mz)
}

func (m *magnetization) Quantity() []float64 { return slice(m.Average()) }

// Sets the magnetization inside the shape
func (m *magnetization) SetInShape(region Shape, conf Config) {
	checkMesh()

	if region == nil {
		region = universe
	}
	host := m.Buffer().HostCopy()
	if m.Buffer().GPUAccess() {
		if err := opencl.D2HQueue.Finish(); err != nil {
			fmt.Printf("error waiting for queue to finish in beginning of magnetization.setregion(): %+v \n", err)
		}
	}
	h := host.Vectors()
	n := m.Mesh().Size()

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := 0; ix < n[X]; ix++ {
				r := Index2Coord(ix, iy, iz)
				x, y, z := r[X], r[Y], r[Z]
				if region(x, y, z) { // inside
					m := conf(x, y, z)
					h[X][iz][iy][ix] = float32(m[X])
					h[Y][iz][iy][ix] = float32(m[Y])
					h[Z][iz][iy][ix] = float32(m[Z])
				}
			}
		}
	}
	m.SetArray(host)
	// sync before returning
	if m.Buffer().GPUAccess() {
		if err := opencl.H2DQueue.Finish(); err != nil {
			fmt.Printf("error waiting for queue to finish after magnetization.setinshape(): %+v \n", err)
		}
	}
}

// set m to config in region
func (m *magnetization) SetRegion(region int, conf Config) {
	host := m.Buffer().HostCopy()
	if m.Buffer().GPUAccess() {
		if err := opencl.D2HQueue.Finish(); err != nil {
			fmt.Printf("error waiting for queue to finish in beginning of magnetization.setregion(): %+v \n", err)
		}
	}
	h := host.Vectors()
	n := m.Mesh().Size()
	r := byte(region)

	regionsArr := regions.HostArray()

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := 0; ix < n[X]; ix++ {
				pos := Index2Coord(ix, iy, iz)
				x, y, z := pos[X], pos[Y], pos[Z]
				if regionsArr[iz][iy][ix] == r {
					m := conf(x, y, z)
					h[X][iz][iy][ix] = float32(m[X])
					h[Y][iz][iy][ix] = float32(m[Y])
					h[Z][iz][iy][ix] = float32(m[Z])
				}
			}
		}
	}
	m.SetArray(host)
	// sync before returning
	if m.Buffer().GPUAccess() {
		if err := opencl.H2DQueue.Finish(); err != nil {
			fmt.Printf("error waiting for queue to finish after magnetization.setinshape(): %+v \n", err)
		}
	}
}

func (m *magnetization) resize() {
	backup := m.Buffer().HostCopy()
	// sync before continuing
	if m.Buffer().GPUAccess() {
		if err := opencl.D2HQueue.Finish(); err != nil {
			fmt.Printf("error waiting for queue to finish in beginning of magnetization.resize(): %+v \n", err)
		}
	}
	s2 := Mesh().Size()
	resized := data.Resample(backup, s2)
	m.buffer_.Free()
	m.buffer_ = opencl.NewSlice(VECTOR, s2)
	data.Copy(m.buffer_, resized)
	if resized.CPUAccess() {
		if err := opencl.H2DQueue.Finish(); err != nil {
			fmt.Printf("error waiting for h2d queue to finish in magnetization.resize(): %+v \n", err)
		}
	} else {
		seqQueue := opencl.ClCmdQueue[0]
		if err := seqQueue.Finish(); err != nil {
			fmt.Printf("error waiting for seqQueue to finish in magnetization.resize(): %+v \n", err)
		}
	}
}
