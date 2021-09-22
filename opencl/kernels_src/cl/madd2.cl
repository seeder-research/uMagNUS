// dst[i] = fac1*src1[i] + fac2*src2[i];
__kernel void
madd2(__global float* __restrict  dst,
      __global float* __restrict  src1, float fac1,
      __global float* __restrict  src2, float fac2, int N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);

	if(i < N) {
		dst[i] = fac1*src1[i] + fac2*src2[i];
	}
}

