// add region-based vector to dst:
// dst[i] += LUT[region[i]]
__kernel void
regionaddv(__global float* __restrict dstx, __global float* __restrict dsty, __global float* __restrict dstz,
           __global float* __restrict LUTx, __global float* __restrict LUTy, __global float* __restrict LUTz,
           __global uint8_t* regions, int N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (i < N) {

		uint8_t r = regions[i];
		dstx[i] += LUTx[r];
		dsty[i] += LUTy[r];
		dstz[i] += LUTz[r];
	}
}

