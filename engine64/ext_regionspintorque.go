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
	regionspintorquelinks map[byte]RegionSpinTorque // global links map
)

const NREGIONSPINTORQUELINKS = 256 // maximum number of region links, limited by size of byte.

func init() {
	regionspintorquelinks = make(map[byte]RegionSpinTorque)
	DeclFunc("ext_NewRegionSpinTorque", DefRegionSpinTorque, "Define a cell link (for spin torque) with given index (0-255) and shape")
	DeclFunc("ext_DeleteRegionSpinTorque", DeleteRegionSpinTorque, "Remove a defined cell link (for spin torque) with given link id")
	DeclROnly("regionspintorquelinks", &regionspintorquelinks, "Outputs the link index for each cell")
}

type RegionSpinTorque struct {
	link    RegionLink
	J       float64
	alpha   float64
	pfix    float64
	pfree   float64
	λfix    float64
	λfree   float64
	ε_prime float64
}

// Define a region with id (0-255) to be inside the Shape.
func DefRegionSpinTorque(id, regionA, regionB, delX, delY, delZ int, J, alpha, pfix, pfree, λfix, λfree, ε_prime float64) {
	defRegionSpinTorqueId(id)
	if regionA < 0 || regionA > NREGION {
		util.Fatalf("regionA id should be 0 -%v, have: %v", NREGION, regionA)
	}
	if regionB < 0 || regionB > NREGION {
		util.Fatalf("regionA id should be 0 -%v, have: %v", NREGION, regionB)
	}
	createRegionSpinTorqueLink(byte(id), regionA, regionB, delX, delY, delZ, J, alpha, pfix, pfree, λfix, λfree, ε_prime)
}

func defRegionSpinTorqueId(id int) {
	if id < 0 || id > NREGIONSPINTORQUELINKS {
		util.Fatalf("regionexchange id should be 0 -%v, have: %v", NREGIONLINKS, id)
	}
	checkMesh()
}

func DeleteRegionSpinTorque(id int) {
	delete(regionspintorquelinks, byte(id))
}

func createRegionSpinTorqueLink(id byte, regionA, regionB, delX, delY, delZ int, J, alpha, pfix, pfree, λfix, λfree, ε_prime float64) {
	regionspintorquelinks[id] = RegionSpinTorque{}
	rr := regionspintorquelinks[id]
	rPtr := &rr
	rPtr.SetRegionA(regionA)
	rPtr.SetRegionB(regionB)
	rPtr.SetJ(float64(J))
	rPtr.SetAlpha(float64(alpha))
	rPtr.SetPfix(float64(pfix))
	rPtr.SetPfree(float64(pfree))
	rPtr.SetLambdafix(float64(λfix))
	rPtr.SetLambdafree(float64(λfree))
	rPtr.SetEPrime(float64(ε_prime))
	rPtr.SetDisplacement(delX, delY, delZ)
}

func AddRegionLinkSpinTorque(dst *data.Slice) {
	ms := Msat.MSlice()
	defer ms.Recycle()
	buf := opencl.Buffer(3, Mesh().Size())
	defer opencl.Recycle(buf)
	opencl.Zero(buf)
	for _, linkpair := range regionspintorquelinks {
		linkPtr := &linkpair
		sX, sY, sZ := linkPtr.GetDisplacement()
		opencl.AddRegionSpinTorque(buf, M.Buffer(), ms, regions.Gpu(), uint8(linkPtr.GetRegionA()), uint8(linkPtr.GetRegionB()), sX, sY, sZ, linkPtr.GetJ(), linkPtr.GetAlpha(), linkPtr.GetPfix(), linkPtr.GetPfree(), linkPtr.GetLambdafix(), linkPtr.GetLambdafree(), linkPtr.GetEPrime(), M.Mesh())
	}
	opencl.Add(dst, dst, buf)
}

func (r *RegionSpinTorque) SetJ(s float64) {
	r.J = s
}

func (r *RegionSpinTorque) GetJ() float64 {
	return r.J
}

func (r *RegionSpinTorque) SetAlpha(s float64) {
	r.alpha = s
}

func (r *RegionSpinTorque) GetAlpha() float64 {
	return r.alpha
}

func (r *RegionSpinTorque) SetPfix(s float64) {
	r.pfix = s
}

func (r *RegionSpinTorque) GetPfix() float64 {
	return r.pfix
}

func (r *RegionSpinTorque) SetPfree(s float64) {
	r.pfree = s
}

func (r *RegionSpinTorque) GetPfree() float64 {
	return r.pfree
}

func (r *RegionSpinTorque) SetLambdafix(s float64) {
	r.λfix = s
}

func (r *RegionSpinTorque) GetLambdafix() float64 {
	return r.λfix
}

func (r *RegionSpinTorque) SetLambdafree(s float64) {
	r.λfree = s
}

func (r *RegionSpinTorque) GetLambdafree() float64 {
	return r.λfree
}

func (r *RegionSpinTorque) SetEPrime(s float64) {
	r.ε_prime = s
}

func (r *RegionSpinTorque) GetEPrime() float64 {
	return r.ε_prime
}

func (r *RegionSpinTorque) SetRegionA(rA int) {
	ptr := &r.link
	ptr.SetRegionA(rA)
}

func (r *RegionSpinTorque) GetRegionA() int {
	ptr := &r.link
	return ptr.GetRegionA()
}

func (r *RegionSpinTorque) SetRegionB(rB int) {
	ptr := &r.link
	ptr.SetRegionB(rB)
}

func (r *RegionSpinTorque) GetRegionB() int {
	ptr := &r.link
	return ptr.GetRegionB()
}

func (r *RegionSpinTorque) SetDisplacement(x, y, z int) {
	ptr := &r.link
	ptr.SetDisplacement(x, y, z)
}

func (r *RegionSpinTorque) GetDisplacement() (int, int, int) {
	ptr := &r.link
	return ptr.GetDisplacement()
}
