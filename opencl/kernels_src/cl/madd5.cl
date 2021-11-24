// dst[i] = src1[i] * fac1 + src2[i] * fac2 + src3[i] * fac3 + src4[i] * fac4 + src5[i] * fac5
__kernel void
madd5(__global float* __restrict  dst,
      __global float* __restrict src1, float fac1,
      __global float* __restrict src2, float fac2,
      __global float* __restrict src3, float fac3,
      __global float* __restrict src4, float fac4,
      __global float* __restrict src5, float fac5, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (fac1 * src1[i]) + (fac2 * src2[i]) + (fac3 * src3[i]) + (fac4 * src4[i]) + (fac5 * src5[i]);
    }
}
