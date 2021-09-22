// shift dst by shz cells (positive or negative) along Z-axis.
// new edge value is clampL at left edge or clampR at right edge.
__kernel void
shiftz(__global float* __restrict  dst, __global float* __restrict  src,
       int Nx,  int Ny,  int Nz, int shz, float clampL, float clampR) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if(ix < Nx && iy < Ny && iz < Nz) {
        int iz2 = iz-shz;
        float newval;
        if (iz2 < 0) {
            newval = clampL;
        } else if (iz2 >= Nz) {
            newval = clampR;
        } else {
            newval = src[idx(ix, iy, iz2)];
        }
        dst[idx(ix, iy, iz)] = newval;
    }
}

