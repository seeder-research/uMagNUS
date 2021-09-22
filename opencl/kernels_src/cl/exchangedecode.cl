// Finds the average exchange strength around each cell, for debugging.
__kernel void
exchangedecode(__global float* __restrict dst, __global float* __restrict aLUT2d, __global uint8_t* __restrict regions,
               float wx, float wy, float wz, int Nx, int Ny, int Nz, uint8_t PBC) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

	if (ix >= Nx || iy >= Ny || iz >= Nz) {
		return;
	}

	// central cell
	int I = idx(ix, iy, iz);
	uint8_t r0 = regions[I];

	int i_;    // neighbor index
	float avg = 0.0f;

	// left neighbor
	i_  = idx(lclampx(ix-1), iy, iz);           // clamps or wraps index according to PBC
	avg += aLUT2d[symidx(r0, regions[i_])];

	// right neighbor
	i_  = idx(hclampx(ix+1), iy, iz);
	avg += aLUT2d[symidx(r0, regions[i_])];

	// back neighbor
	i_  = idx(ix, lclampy(iy-1), iz);
	avg += aLUT2d[symidx(r0, regions[i_])];

	// front neighbor
	i_  = idx(ix, hclampy(iy+1), iz);
	avg += aLUT2d[symidx(r0, regions[i_])];

	// only take vertical derivative for 3D sim
	if (Nz != 1) {
		// bottom neighbor
		i_  = idx(ix, iy, lclampz(iz-1));
		avg += aLUT2d[symidx(r0, regions[i_])];

		// top neighbor
		i_  = idx(ix, iy, hclampz(iz+1));
		avg += aLUT2d[symidx(r0, regions[i_])];
	}

	dst[I] = avg;
}

