// shift dst by shy cells (positive or negative) along Y-axis.
// new edge value is clampL at left edge or clampR at right edge.
__kernel void
shifty(__global float* __restrict  dst, __global float* __restrict  src,
       int Nx,  int Ny,  int Nz, int shy, float clampL, float clampR) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if(ix < Nx && iy < Ny && iz < Nz) {
        int iy2 = iy-shy;
        float newval;
        if (iy2 < 0) {
            newval = clampL;
        } else if (iy2 >= Ny) {
            newval = clampR;
        } else {
            newval = src[idx(ix, iy2, iz)];
        }
        dst[idx(ix, iy, iz)] = newval;
    }
}

