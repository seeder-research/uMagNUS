package engine

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

// declare a global structure that stores pairwise couplings
// other extensions can invoke this data structure to calculate
// exchange or other phenomena between the coupled cells

var (
	regionexchangelinks map[byte]RegionExchange // global links map
	B_tworegionexch     = NewVectorField("B_tworegionexch", "T", "Two Region Exchange field", AddRegionExchangeField)
	E_tworegionexch     = NewScalarValue("E_tworegionexch", "J", "Total two region exchange energy", GetRegionExchangeEnergy)
	Edens_tworegionexch = NewScalarField("Edens_tworegionexch", "J/m3", "Total two region exchange energy density", AddRegionExchangeEnergyDensity)
)

const NREGIONLINKS = 256 // maximum number of region links, limited by size of byte.

func init() {
	registerEnergy(GetRegionExchangeEnergy, AddRegionExchangeEnergyDensity)
	regionexchangelinks = make(map[byte]RegionExchange)
	DeclFunc("ext_NewRegionExchange", DefRegionExchange, "Define a cell link with given index (0-255) and shape")
	DeclFunc("ext_DeleteRegionExchange", DeleteRegionExchange, "Remove a defined cell link with given link id")
	DeclROnly("regionexchangelinks", &regionexchangelinks, "Outputs the link index for each cell")
}

type RegionExchange struct {
	link RegionLink
	sig  float32
	sig2 float32
}

// Define a region with id (0-255) to be inside the Shape.
func DefRegionExchange(id, regionA, regionB, delX, delY, delZ int, sig, sig2 float64) {
	defRegionExchangeId(id)
	if regionA < 0 || regionA > NREGION {
		util.Fatalf("regionA id should be 0 -%v, have: %v", NREGION, regionA)
	}
	if regionB < 0 || regionB > NREGION {
		util.Fatalf("regionA id should be 0 -%v, have: %v", NREGION, regionB)
	}
	createRegionExchangeLink(byte(id), regionA, regionB, delX, delY, delZ, sig, sig2)
}

func defRegionExchangeId(id int) {
	if id < 0 || id > NREGIONLINKS {
		util.Fatalf("regionexchange id should be 0 -%v, have: %v", NREGIONLINKS, id)
	}
	checkMesh()
}

func DeleteRegionExchange(id int) {
	delete(regionexchangelinks, byte(id))
}

func createRegionExchangeLink(id byte, regionA, regionB, delX, delY, delZ int, sig, sig2 float64) {
	regionexchangelinks[id] = RegionExchange{}
	rr := regionexchangelinks[id]
	rPtr := &rr
	rPtr.SetRegionA(regionA)
	rPtr.SetRegionB(regionB)
	rPtr.SetSig(float32(sig))
	rPtr.SetSig2(float32(sig2))
	rPtr.SetDisplacement(delX, delY, delZ)
}

// Adds the current region exchange field to dst
func AddRegionExchangeField(dst *data.Slice) {
	ms := Msat.MSlice()
	defer ms.Recycle()

	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queues to finish in addregionexchangefield: %+v \n", err)
	}

	buf := opencl.Buffer(3, Mesh().Size())
	defer opencl.Recycle(buf)

	opencl.Zero(buf)
	seqQueue := opencl.ClCmdQueue[0]
	for _, linkpair := range regionexchangelinks {
		linkPtr := &linkpair
		sX, sY, sZ := linkPtr.GetDisplacement()
		opencl.AddRegionExchangeField(buf, M.Buffer(), ms, regions.Gpu(), uint8(linkPtr.GetRegionA()), uint8(linkPtr.GetRegionB()), sX, sY, sZ, linkPtr.GetSig(), linkPtr.GetSig2(), M.Mesh(), seqQueue, nil)
	}

	// checkout queues for 3-component vector addition
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})

	// add vectors
	opencl.Add(dst, dst, buf, queues, nil)

	// sync before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for all queues to finish in addregionlinkspintorque: %+v \n", err)
	}
}

// Adds the region exchange energy densities
func AddRegionExchangeEnergyDensity(dst *data.Slice) {
	ms := Msat.MSlice()
	defer ms.Recycle()

	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queues to finish in addregionexchangeenergydensity: %+v \n", err)
	}

	buf := opencl.Buffer(1, Mesh().Size())
	defer opencl.Recycle(buf)
	opencl.Zero(buf)
	seqQueue := opencl.ClCmdQueue[0]
	for _, linkpair := range regionexchangelinks {
		linkPtr := &linkpair
		sX, sY, sZ := linkPtr.GetDisplacement()
		opencl.AddRegionExchangeEdens(buf, M.Buffer(), ms, regions.Gpu(), uint8(linkPtr.GetRegionA()), uint8(linkPtr.GetRegionB()), sX, sY, sZ, linkPtr.GetSig(), linkPtr.GetSig2(), M.Mesh(), seqQueue, nil)
	}
	opencl.Add(dst, dst, buf, []*cl.CommandQueue{seqQueue}, nil)

	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for sequential queue to finish after addregionexchangeenergydensity: %+v \n", err)
	}
}

func GetRegionExchangeEnergy() float64 {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queues to finish in addregionexchangeenergydensity: %+v \n", err)
	}

	buf := opencl.Buffer(1, Mesh().Size())
	defer opencl.Recycle(buf)
	opencl.Zero(buf)
	AddRegionExchangeEnergyDensity(buf)
	seqQueue := opencl.ClCmdQueue[0]
	return cellVolume() * float64(opencl.Sum(buf, seqQueue, nil))
}

func (r *RegionExchange) SetSig(s float32) {
	r.sig = s
}

func (r *RegionExchange) GetSig() float32 {
	return r.sig
}

func (r *RegionExchange) SetSig2(s float32) {
	r.sig2 = s
}

func (r *RegionExchange) GetSig2() float32 {
	return r.sig2
}

func (r *RegionExchange) SetRegionA(rA int) {
	ptr := &r.link
	ptr.SetRegionA(rA)
}

func (r *RegionExchange) GetRegionA() int {
	ptr := &r.link
	return ptr.GetRegionA()
}

func (r *RegionExchange) SetRegionB(rB int) {
	ptr := &r.link
	ptr.SetRegionB(rB)
}

func (r *RegionExchange) GetRegionB() int {
	ptr := &r.link
	return ptr.GetRegionB()
}

func (r *RegionExchange) SetDisplacement(x, y, z int) {
	ptr := &r.link
	ptr.SetDisplacement(x, y, z)
}

func (r *RegionExchange) GetDisplacement() (int, int, int) {
	ptr := &r.link
	return ptr.GetDisplacement()
}
