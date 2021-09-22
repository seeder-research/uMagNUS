// dst[i] = a[i] / b[i]

__kernel void
divide(__global float* __restrict  dst, __global float* __restrict  a, __global float* __restrict b, int N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);

	if(i < N) {
		if((a[i] == 0) || (b[i] == 0)) {
			dst[i] = 0.0;
		} else {
			dst[i] = a[i] / b[i];
		}
	}
}

