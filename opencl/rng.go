package opencl

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl/RNGmtgp"
	"github.com/seeder-research/uMagNUS/opencl/cl"
	"github.com/seeder-research/uMagNUS/timer"
	"math/rand"
	"time"
)

type Prng_ interface {
	Init(uint32, []*cl.Event)
	GenerateUniform(unsafe.Pointer, int, []*cl.Event) *cl.Event
	GenerateNormal(unsafe.Pointer, int, []*cl.Event) *cl.Event
	GetGroupSize() int
	GetGroupCount() int
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

const MTGP32_MEXP = RNGmtgp.MTGPDC_MEXP
const MTGP32_N = RNGmtgp.MTGPDC_N
const MTGP32_FLOOR_2P = RNGmtgp.MTGPDC_FLOOR_2P
const MTGP32_CEIL_2P = RNGmtgp.MTGPDC_CEIL_2P
const MTGP32_TN = RNGmtgp.MTGPDC_TN
const MTGP32_LS = RNGmtgp.MTGPDC_LS
const MTGP32_TS = RNGmtgp.MTGPDC_TS
const MTGP32_PARAM_NUM = RNGmtgp.MTGPDC_PARAMS_NUM

type mtgp32_params RNGmtgp.MTGP32dc_params_array_ptr

func NewGenerator(name string) *Generator {
	switch name {
	case "mtgp":
		var prng_ptr Prng_
		prng_ptr = NewMTGPRNGParams()
		return &Generator{Name: "mtgp", PRNG: prng_ptr}
	case "mrg32k3a":
		fmt.Println("mrg32k3a not yet implemented")
		return nil
	default:
		fmt.Println("RNG not implemented: ", name)
		return nil
	}
}

func (g *Generator) CreatePNG() {
	switch g.Name {
	case "mtgp":
		var prng_ptr Prng_
		prng_ptr = NewMTGPRNGParams()
		g.PRNG = prng_ptr
	case "mrg32k3a":
		fmt.Println("mrg32k3a not yet implemented")
	default:
		fmt.Println("RNG not implemented: ", g.Name)
	}
}

func (g *Generator) Init(seed *uint32, events []*cl.Event) {
	g.buf_size = 8 * g.PRNG.GetGroupSize() * MTGP32_TN
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
			event, err = ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(g.buf.DevPtr(0)), (*cl.MemObject)(data), SIZEOF_FLOAT32*g.sup_offset, SIZEOF_FLOAT32*demand_offset, SIZEOF_FLOAT32*demand, nil)
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
			event, err = ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(g.buf.DevPtr(0)), (*cl.MemObject)(data), SIZEOF_FLOAT32*g.sup_offset, SIZEOF_FLOAT32*demand_offset, SIZEOF_FLOAT32*g.supply, nil)
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
			event, err = ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(g.buf.DevPtr(0)), (*cl.MemObject)(data), SIZEOF_FLOAT32*g.sup_offset, SIZEOF_FLOAT32*demand_offset, SIZEOF_FLOAT32*demand, nil)
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
			event, err = ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(g.buf.DevPtr(0)), (*cl.MemObject)(data), SIZEOF_FLOAT32*g.sup_offset, SIZEOF_FLOAT32*demand_offset, SIZEOF_FLOAT32*g.supply, nil)
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

func NewMTGPRNGParams() *mtgp32_params {
	var err error
	var events_list []*cl.Event
	var event *cl.Event
	tmp := RNGmtgp.NewMTGPParams()
	maxNumGroups, max_size := ClCUnits, MTGP32_PARAM_NUM
	if maxNumGroups > max_size {
		maxNumGroups = max_size
	}
	tmp.SetGroupCount(maxNumGroups)
	if ClWGSize < MTGP32_TN {
		log.Fatalln("Unable to use PRNG on device! Insufficient resources for parallel work-items")
	}
	local_item := MTGP32_N
	if local_item > ClWGSize {
		local_item = MTGP32_TN
	}
	tmp.SetGroupSize(local_item)
	tmp.GetMTGPArrays()
	tmp.CreateParamBuffers(ClCtx)
	events_list, err = tmp.LoadAllParamBuffersToDevice(ClCmdQueue, nil)
	if err != nil {
		log.Fatalln("Unable to load RNG parameters to device")
	}
	event, err = tmp.LoadStatusBuffersToDevice(ClCmdQueue, nil)
	if err != nil {
		log.Fatalln("Unable to load RNG status to device")
	}
	err = cl.WaitForEvents(append(events_list, event))
	return (*mtgp32_params)(tmp)
}

func initRNG() uint32 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Uint32()
}

func (p *mtgp32_params) Init(seed uint32, events []*cl.Event) {

	event := k_mtgp32_init_seed_kernel_async(unsafe.Pointer(p.Rec_buf), unsafe.Pointer(p.Temper_buf), unsafe.Pointer(p.Flt_temper_buf), unsafe.Pointer(p.Pos_buf),
		unsafe.Pointer(p.Sh1_buf), unsafe.Pointer(p.Sh2_buf), unsafe.Pointer(p.Status_buf), seed,
		&config{[]int{p.GetGroupCount() * p.GetGroupSize()}, []int{p.GetGroupSize()}}, events)

	p.Ini = true
	err := cl.WaitForEvents([]*cl.Event{event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in InitRNG: %+v \n", err)
	}

}

func (p *mtgp32_params) GenerateUniform(d_data unsafe.Pointer, data_size int, events []*cl.Event) *cl.Event {

	if p.Ini == false {
		log.Fatalln("Generator has not been initialized!")
	}

	item_num := MTGP32_TN * p.GetGroupCount()
	min_size := MTGP32_LS * p.GetGroupCount()
	tmpSize := data_size
	if data_size%min_size != 0 {
		tmpSize = (data_size/min_size + 1) * min_size
	}

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Start("mtgp32_uniform")
	}

	event := k_mtgp32_uniform_async(unsafe.Pointer(p.Rec_buf), unsafe.Pointer(p.Temper_buf), unsafe.Pointer(p.Flt_temper_buf), unsafe.Pointer(p.Pos_buf),
		unsafe.Pointer(p.Sh1_buf), unsafe.Pointer(p.Sh2_buf), unsafe.Pointer(p.Status_buf), d_data, tmpSize/p.GetGroupCount(),
		&config{[]int{item_num}, []int{MTGP32_TN}}, events)

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Stop("mtgp32_uniform")
	}

	return event
}

func (p *mtgp32_params) GenerateNormal(d_data unsafe.Pointer, data_size int, events []*cl.Event) *cl.Event {

	if p.Ini == false {
		log.Fatalln("Generator has not been initialized!")
	}

	item_num := MTGP32_TN * p.GetGroupCount()
	min_size := MTGP32_LS * p.GetGroupCount()
	tmpSize := data_size
	if data_size%min_size != 0 {
		tmpSize = (data_size/min_size + 1) * min_size
	}

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Start("mtgp32_normal")
	}

	event := k_mtgp32_normal_async(unsafe.Pointer(p.Rec_buf), unsafe.Pointer(p.Temper_buf), unsafe.Pointer(p.Flt_temper_buf), unsafe.Pointer(p.Pos_buf),
		unsafe.Pointer(p.Sh1_buf), unsafe.Pointer(p.Sh2_buf), unsafe.Pointer(p.Status_buf), d_data, tmpSize/p.GetGroupCount(),
		&config{[]int{item_num}, []int{MTGP32_TN}}, events)

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Stop("mtgp32_normal")
	}

	return event
}

func (p *mtgp32_params) SetRecursionArray(arr []uint32) {
	p.Rec = arr
}

func (p *mtgp32_params) GetRecursionArray() []uint32 {
	return p.Rec
}

func (p *mtgp32_params) SetPositionArray(arr []int) {
	p.Pos = arr
}

func (p *mtgp32_params) GetPositionArray() []int {
	return p.Pos
}

func (p *mtgp32_params) SetSH1Array(arr []int) {
	p.Sh1 = arr
}

func (p *mtgp32_params) GetSH1Array() []int {
	return p.Sh1
}

func (p *mtgp32_params) SetSH2Array(arr []int) {
	p.Sh2 = arr
}

func (p *mtgp32_params) GetSH2Array() []int {
	return p.Sh2
}

func (p *mtgp32_params) SetStatusArray(arr []uint32) {
	p.Status = arr
}

func (p *mtgp32_params) GetStatusArray() []uint32 {
	return p.Status
}

func (p *mtgp32_params) SetGroupSize(in int) {
	p.GroupSize = in
}

func (p *mtgp32_params) GetGroupSize() int {
	return p.GroupSize
}

func (p *mtgp32_params) SetGroupCount(in int) {
	p.GroupCount = in
}

func (p *mtgp32_params) GetGroupCount() int {
	return p.GroupCount
}
