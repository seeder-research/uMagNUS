// shift dst by shy cells (positive or negative) along Y-axis.
__kernel void
shiftbytesy(__global uint8_t* __restrict  dst, __global uint8_t* __restrict  src,
            int Nx,  int Ny,  int Nz, int shy, uint8_t clampV) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if(ix < Nx && iy < Ny && iz < Nz) {
        int iy2 = iy-shy;
        uint8_t newval;
        if (iy2 < 0 || iy2 >= Ny) {
            newval = clampV;
        } else {
            newval = src[idx(ix, iy2, iz)];
        }
        dst[idx(ix, iy, iz)] = newval;
    }
}

