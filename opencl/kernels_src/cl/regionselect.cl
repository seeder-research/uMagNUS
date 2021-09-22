__kernel void
regionselect(__global float* __restrict  dst, __global float* __restrict src, __global uint8_t* regions, uint8_t region, int N) {

	int i = ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (i < N) {
		dst[i] = (regions[i] == region? src[i]: 0.0f);
	}
}

