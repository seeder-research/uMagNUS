__kernel void
setPhi(__global float* __restrict phi,
        __global float* __restrict mx, __global float* __restrict my,
        int Nx, int Ny, int Nz) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

        if (ix >= Nx || iy >= Ny || iz >= Nz)
        {
                return;
        }

        int I = idx(ix, iy, iz);                      // central cell index
        phi[I] = atan2(my[I], mx[I]);
}
