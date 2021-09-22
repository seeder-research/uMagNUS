// dst[i] = src1[i] * fac1 + src2[i] * fac2 + src3[i] * fac3 + src4[i] * fac4 + src5[i] * fac5 + src6[i] * fac6 + src7[i] * fac7
__kernel void
madd7(__global float* __restrict__ dst,
      __global float* __restrict src1, float fac1,
      __global float* __restrict src2, float fac2,
      __global float* __restrict src3, float fac3,
      __global float* __restrict src4, float fac4,
      __global float* __restrict src5, float fac5,
      __global float* __restrict src6, float fac6,
      __global float* __restrict src7, float fac7, int N) {

	int i = ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);

	if(i < N) {
		dst[i] = (fac1 * src1[i]) + (fac2 * src2[i]) + (fac3 * src3[i]) + (fac4 * src4[i]) + (fac5 * src5[i]) + (fac6 * src6[i]) + (fac7 * src7[i]);
	}
}

