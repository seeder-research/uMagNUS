// dst[i] = src1[i] * fac1 + src2[i] * fac2 + src3[i] * fac3 + src4[i] * fac4 + src5[i] * fac5 + src6[i] * fac6
__kernel void
madd6(__global real_t* __restrict  dst,
      __global real_t* __restrict src1, real_t fac1,
      __global real_t* __restrict src2, real_t fac2,
      __global real_t* __restrict src3, real_t fac3,
      __global real_t* __restrict src4, real_t fac4,
      __global real_t* __restrict src5, real_t fac5,
      __global real_t* __restrict src6, real_t fac6, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (fac1 * src1[i]) + (fac2 * src2[i]) + (fac3 * src3[i]) + (fac4 * src4[i]) + (fac5 * src5[i]) + (fac6 * src6[i]);
    }
}
