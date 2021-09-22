// Crop stores in dst a rectangle cropped from src at given offset position.
// dst size may be smaller than src.
__kernel void
crop(__global float* __restrict  dst, int Dx, int Dy, int Dz,
     __global float* __restrict  src, int Sx, int Sy, int Sz,
     int Offx, int Offy, int Offz) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

	if (ix<Dx && iy<Dy && iz<Dz) {
		dst[index(ix, iy, iz, Dx, Dy, Dz)] = src[index(ix+Offx, iy+Offy, iz+Offz, Sx, Sy, Sz)];
	}
}

