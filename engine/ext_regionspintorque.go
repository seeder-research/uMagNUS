package engine

import (
//	"github.com/seeder-research/uMagNUS/data"
//	"github.com/seeder-research/uMagNUS/opencl"
	"github.com/seeder-research/uMagNUS/util"
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
	J       float32
	alpha   float32
	pfix    float32
	pfree   float32
	λfix    float32
	λfree   float32
	ε_prime float32
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
	rPtr.SetJ(float32(J))
	rPtr.SetAlpha(float32(alpha))
	rPtr.SetPfix(float32(pfix))
	rPtr.SetPfree(float32(pfree))
	rPtr.SetLambdafix(float32(λfix))
	rPtr.SetLambdafree(float32(λfree))
	rPtr.SetEPrime(float32(ε_prime))
	rPtr.SetDisplacement(delX, delY, delZ)
}

func (r *RegionSpinTorque) SetJ(s float32) {
	r.J = s
}

func (r *RegionSpinTorque) GetJ() float32 {
	return r.J
}

func (r *RegionSpinTorque) SetAlpha(s float32) {
	r.alpha = s
}

func (r *RegionSpinTorque) SetPfix(s float32) {
	r.pfix = s
}

func (r *RegionSpinTorque) GetPfix() float32 {
	return r.pfix
}

func (r *RegionSpinTorque) SetPfree(s float32) {
	r.pfree = s
}

func (r *RegionSpinTorque) GetPfree() float32 {
	return r.pfree
}

func (r *RegionSpinTorque) SetLambdafix(s float32) {
	r.λfix = s
}

func (r *RegionSpinTorque) GetLambdafix() float32 {
	return r.λfix
}

func (r *RegionSpinTorque) SetLambdafree(s float32) {
	r.λfree = s
}

func (r *RegionSpinTorque) GetLambdafree() float32 {
	return r.λfree
}

func (r *RegionSpinTorque) SetEPrime(s float32) {
	r.ε_prime = s
}

func (r *RegionSpinTorque) GetEPrime() float32 {
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

func (r *RegionSpinTorque) SetDisplacement(x int, y int, z int) {
	ptr := &r.link
	ptr.SetDisplacement(x, y, z)
}

func (r *RegionSpinTorque) GetDisplacement() (int, int, int) {
	ptr := &r.link
	return ptr.GetDisplacement()
}
