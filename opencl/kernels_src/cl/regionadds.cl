// add region-based scalar to dst:
// dst[i] += LUT[region[i]]
__kernel void
regionadds(__global float* __restrict dst,
           __global float* __restrict LUT,
           __global uint8_t* regions, int N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (i < N) {

		uint8_t r = regions[i];
		dst[i] += LUT[r];
	}
}

