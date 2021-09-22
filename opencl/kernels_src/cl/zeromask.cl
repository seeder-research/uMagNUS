// set dst to zero in cells where mask != 0
__kernel void
zeromask(__global float* __restrict  dst, __global float* maskLUT, __global uint8_t* regions, int N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (i < N) {
		if (maskLUT[regions[i]] != 0){
			dst[i] = 0;
		}
	}
}

