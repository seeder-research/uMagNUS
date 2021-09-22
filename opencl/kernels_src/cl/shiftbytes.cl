// shift dst by shx cells (positive or negative) along X-axis.
// new edge value is clampL at left edge or clampR at right edge.
__kernel void
shiftbytes(__global uint8_t* __restrict  dst, __global uint8_t* __restrict  src,
           int Nx,  int Ny,  int Nz, int shx, uint8_t clampV) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

	if(ix < Nx && iy < Ny && iz < Nz) {
		int ix2 = ix-shx;
		uint8_t newval;
		if (ix2 < 0 || ix2 >= Nx) {
			newval = clampV;
		} else {
			newval = src[idx(ix2, iy, iz)];
		}
		dst[idx(ix, iy, iz)] = newval;
	}
}

