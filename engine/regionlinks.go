package engine

// data structure used for linking two regions
// used for twosurfaceexchange, etc.

type RegionLink struct {
	regionA byte
	regionB byte
	delX    int
	delY    int
	delZ    int
}

func (r *RegionLink) SetRegionA(rA int) {
	r.regionA = byte(rA)
}

func (r *RegionLink) GetRegionA() int {
	return int(r.regionA)
}

func (r *RegionLink) SetRegionB(rB int) {
	r.regionB = byte(rB)
}

func (r *RegionLink) GetRegionB() int {
	return int(r.regionB)
}

func (r *RegionLink) GetDisplacement() (int, int, int) {
	return r.delX, r.delY, r.delZ
}

func (r *RegionLink) SetDisplacement(x, y, z int) {
	r.delX, r.delY, r.delZ = x, y, z
}
