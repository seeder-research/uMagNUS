// dst[i] = fac1 * src1[i] + fac2 * src2[i] + fac3 * src3[i]
__kernel void
madd3(__global float* __restrict__ dst,
      __global float* __restrict src1, float fac1,
      __global float* __restrict src2, float fac2,
      __global float* __restrict src3, float fac3, int N) {

	int i = ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);

	if(i < N) {
		dst[i] = (fac1 * src1[i]) + (fac2 * src2[i] + fac3 * src3[i]);
		// parens for better accuracy heun solver.
	}
}

