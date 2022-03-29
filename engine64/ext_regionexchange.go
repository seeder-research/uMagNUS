package engine64

import (
	data "github.com/seeder-research/uMagNUS/data64"
	opencl "github.com/seeder-research/uMagNUS/opencl64"
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
	sig  float64
	sig2 float64
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
	rPtr.SetSig(float64(sig))
	rPtr.SetSig2(float64(sig2))
	rPtr.SetDisplacement(delX, delY, delZ)
}

// Adds the current region exchange field to dst
func AddRegionExchangeField(dst *data.Slice) {
	ms := Msat.MSlice()
	defer ms.Recycle()
	buf := opencl.Buffer(3, Mesh().Size())
	defer opencl.Recycle(buf)
	opencl.Zero(buf)
	for _, linkpair := range regionexchangelinks {
		linkPtr := &linkpair
		sX, sY, sZ := linkPtr.GetDisplacement()
		opencl.AddRegionExchangeField(buf, M.Buffer(), ms, regions.Gpu(), uint8(linkPtr.GetRegionA()), uint8(linkPtr.GetRegionB()), sX, sY, sZ, linkPtr.GetSig(), linkPtr.GetSig2(), M.Mesh())
	}
	opencl.Add(dst, dst, buf)
}

// Adds the region exchange energy densities
func AddRegionExchangeEnergyDensity(dst *data.Slice) {
	ms := Msat.MSlice()
	defer ms.Recycle()
	buf := opencl.Buffer(1, Mesh().Size())
	defer opencl.Recycle(buf)
	opencl.Zero(buf)
	for _, linkpair := range regionexchangelinks {
		linkPtr := &linkpair
		sX, sY, sZ := linkPtr.GetDisplacement()
		opencl.AddRegionExchangeEdens(buf, M.Buffer(), ms, regions.Gpu(), uint8(linkPtr.GetRegionA()), uint8(linkPtr.GetRegionB()), sX, sY, sZ, linkPtr.GetSig(), linkPtr.GetSig2(), M.Mesh())
	}
	opencl.Add(dst, dst, buf)
}

func GetRegionExchangeEnergy() float64 {
	buf := opencl.Buffer(1, Mesh().Size())
	defer opencl.Recycle(buf)
	opencl.Zero(buf)
	AddRegionExchangeEnergyDensity(buf)
	return cellVolume() * float64(opencl.Sum(buf))
}

func (r *RegionExchange) SetSig(s float64) {
	r.sig = s
}

func (r *RegionExchange) GetSig() float64 {
	return r.sig
}

func (r *RegionExchange) SetSig2(s float64) {
	r.sig2 = s
}

func (r *RegionExchange) GetSig2() float64 {
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
