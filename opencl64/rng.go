package opencl64

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	"github.com/seeder-research/uMagNUS/opencl64/oclRAND"
	"math/rand"
	"time"
)

type Prng_ interface {
	Init(uint64, []*cl.Event)
	GenerateUniform(unsafe.Pointer, int, []*cl.Event) *cl.Event
	GenerateNormal(unsafe.Pointer, int, []*cl.Event) *cl.Event
	GetGroupSize() int
	GetGroupCount() int
	RecommendSize() int
}

type Generator struct {
	Name       string
	PRNG       Prng_
	r_seed     *uint32
	buf_size   int
	buf        *data.Slice
	supply     int
	sup_offset int
}

func NewGenerator(name string) *Generator {
	switch name {
	case "threefry":
		oclRAND.Init(ClCmdQueue, Synchronous, KernList)
		var prng_ptr Prng_
		prng_ptr = NewTHREEFRYRNGParams()
		return &Generator{Name: "threefry", PRNG: prng_ptr}
	case "xorwow":
		oclRAND.Init(ClCmdQueue, Synchronous, KernList)
		var prng_ptr Prng_
		prng_ptr = NewXORWOWRNGParams()
		return &Generator{Name: "xorwow", PRNG: prng_ptr}
	default:
		fmt.Println("RNG not implemented: ", name)
		return nil
	}
}

func (g *Generator) CreatePNG() {
	switch g.Name {
	case "threefry":
		oclRAND.Init(ClCmdQueue, Synchronous, KernList)
		var prng_ptr Prng_
		prng_ptr = NewTHREEFRYRNGParams()
		g.PRNG = prng_ptr
	case "xorwow":
		oclRAND.Init(ClCmdQueue, Synchronous, KernList)
		var prng_ptr Prng_
		prng_ptr = NewXORWOWRNGParams()
		g.PRNG = prng_ptr
	default:
		fmt.Println("RNG not implemented: ", g.Name)
	}
}

func (g *Generator) Init(seed *uint64, events []*cl.Event) {
	g.buf_size = g.PRNG.RecommendSize()
	if seed == nil {
		g.PRNG.Init(initRNG(), events)
	} else {
		g.PRNG.Init(*seed, events)
	}
	if g.buf == nil {
		g.buf = Buffer(1, [3]int{g.buf_size, 1, 1})
	} else {
		if g.buf.NComp() != 1 {
			log.Fatalln("Bad buffer for RNG \n")
		} else {
			bufferSize := g.buf.Size()
			if bufferSize[0] != g.buf_size {
				g.buf.Free()
				g.buf = Buffer(1, [3]int{g.buf_size, 1, 1})
			}
		}
	}
	g.supply = 0
}

func (g *Generator) Uniform(data unsafe.Pointer, d_size int, events []*cl.Event) {
	var event *cl.Event
	var err error
	demand, demand_offset := d_size, 0
	err = cl.WaitForEvents(events)
	if err != nil {
		fmt.Printf("WaitForEvents prior to generating random numbers failed: %+v \n", err)
	}
	err = ClCmdQueue.Finish()
	if err != nil {
		fmt.Printf("Waiting for Command Queue to empty prior to generating random numbers failed: %+v \n", err)
	}
	for demand > 0 {
		if g.supply <= 0 {
			event = g.PRNG.GenerateUniform(g.buf.DevPtr(0), g.buf_size, events)
			err = cl.WaitForEvents([]*cl.Event{event})
			if err != nil {
				fmt.Printf("WaitForEvents in generating uniform random numbers failed: %+v \n", err)
			}
			bufferSize := g.buf.Size()
			if bufferSize[0] != g.buf_size {
				fmt.Printf("Error in buffer size variables! \n")
			}
			g.supply = bufferSize[0]
			g.sup_offset = 0
		}
		if g.supply >= demand {
			event, err = ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(g.buf.DevPtr(0)), (*cl.MemObject)(data), SIZEOF_FLOAT64*g.sup_offset, SIZEOF_FLOAT64*demand_offset, SIZEOF_FLOAT64*demand, nil)
			if err != nil {
				fmt.Printf("EnqueueCopyBuffer in copying uniform random numbers failed: %+v \n", err)
			}
			err = cl.WaitForEvents([]*cl.Event{event})
			if err != nil {
				fmt.Printf("WaitForEvents in copying uniform random numbers failed: %+v \n", err)
			}
			g.sup_offset += demand
			g.supply -= demand
			demand = 0
		} else {
			event, err = ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(g.buf.DevPtr(0)), (*cl.MemObject)(data), SIZEOF_FLOAT64*g.sup_offset, SIZEOF_FLOAT64*demand_offset, SIZEOF_FLOAT64*g.supply, nil)
			if err != nil {
				fmt.Printf("EnqueueCopyBuffer in copying uniform random numbers failed: %+v \n", err)
			}
			err = cl.WaitForEvents([]*cl.Event{event})
			if err != nil {
				fmt.Printf("WaitForEvents in copying uniform random numbers failed: %+v \n", err)
			}
			demand -= g.supply
			demand_offset += g.supply
			g.supply = 0
		}
	}
}

func (g *Generator) Normal(data unsafe.Pointer, d_size int, events []*cl.Event) {
	var event *cl.Event
	var err error
	demand, demand_offset := d_size, 0
	if events != nil {
		err = cl.WaitForEvents(events)
		if err != nil {
			fmt.Printf("WaitForEvents prior to generating random numbers failed: %+v \n", err)
		}
	}
	for demand > 0 {
		if g.supply <= 0 {
			event = g.PRNG.GenerateNormal(g.buf.DevPtr(0), g.buf_size, events)
			err = cl.WaitForEvents([]*cl.Event{event})
			if err != nil {
				fmt.Printf("WaitForEvents in generating normally distributed random numbers failed: %+v \n", err)
			}
			bufferSize := g.buf.Size()
			if bufferSize[0] != g.buf_size {
				fmt.Printf("Error in buffer size variables! \n")
			}
			g.supply = bufferSize[0]
			g.sup_offset = 0
		}
		if g.supply >= demand {
			event, err = ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(g.buf.DevPtr(0)), (*cl.MemObject)(data), SIZEOF_FLOAT64*g.sup_offset, SIZEOF_FLOAT64*demand_offset, SIZEOF_FLOAT64*demand, nil)
			if err != nil {
				fmt.Printf("EnqueueCopyBuffer in copying normally distributed random numbers failed: %+v \n", err)
			}
			err = cl.WaitForEvents([]*cl.Event{event})
			if err != nil {
				fmt.Printf("WaitForEvents in copying normally distributed random numbers failed: %+v \n", err)
			}
			g.sup_offset += demand
			g.supply -= demand
			demand = 0
		} else {
			event, err = ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(g.buf.DevPtr(0)), (*cl.MemObject)(data), SIZEOF_FLOAT64*g.sup_offset, SIZEOF_FLOAT64*demand_offset, SIZEOF_FLOAT64*g.supply, nil)
			if err != nil {
				fmt.Printf("EnqueueCopyBuffer in copying normally distributed random numbers failed: %+v \n", err)
			}
			err = cl.WaitForEvents([]*cl.Event{event})
			if err != nil {
				fmt.Printf("WaitForEvents in copying uniform normally distributed numbers failed: %+v \n", err)
			}
			demand -= g.supply
			demand_offset += g.supply
			g.supply = 0
		}
	}
}

func NewTHREEFRYRNGParams() *oclRAND.THREEFRY_status_array_ptr {
	tmp := oclRAND.NewTHREEFRYStatus()
	tmp.SetGroupCount(ClCUnits)
	tmp.SetGroupSize(ClPrefWGSz)
	tmp.SetStatusSize(ClCUnits * ClPrefWGSz)
	tmp.CreateStatusBuffer(ClCtx)

	return tmp
}

func NewXORWOWRNGParams() *oclRAND.XORWOW_status_array_ptr {
	tmp := oclRAND.NewXORWOWStatus()
	tmp.SetGroupCount(ClCUnits)
	tmp.SetGroupSize(ClPrefWGSz)
	tmp.SetStatusSize(ClCUnits * ClPrefWGSz)
	tmp.CreateStatusBuffer(ClCtx)

	return tmp
}

func initRNG() uint64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Uint64()
}
