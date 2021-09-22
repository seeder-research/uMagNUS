__kernel void
setTheta(__global float* __restrict theta,
        __global float* __restrict mz,
        int Nx, int Ny, int Nz) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

        if (ix >= Nx || iy >= Ny || iz >= Nz)
        {
                return;
        }

        int I = idx(ix, iy, iz);                      // central cell index
        theta[I] = acos(mz[I]);
}
