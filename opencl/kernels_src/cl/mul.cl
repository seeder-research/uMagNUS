// dst[i] = a[i] * b[i]
__kernel void
mul(__global float* __restrict  dst, __global float* __restrict  a, __global float* __restrict b, int N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);

	if(i < N) {
		dst[i] = a[i] * b[i];
	}
}

