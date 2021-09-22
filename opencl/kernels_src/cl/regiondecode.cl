// decode the regions+LUT pair into an uncompressed array
__kernel void
regiondecode(__global float* __restrict  dst, __global float* __restrict LUT, __global uint8_t* regions, int N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (i < N) {

		dst[i] = LUT[regions[i]];

	}
}

