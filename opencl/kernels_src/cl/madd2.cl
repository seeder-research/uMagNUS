// dst[i] = fac1*src1[i] + fac2*src2[i];
__kernel void
madd2(__global real_t* __restrict   dst,
      __global real_t* __restrict  src1, real_t fac1,
      __global real_t* __restrict  src2, real_t fac2, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = fac1*src1[i] + fac2*src2[i];
    }
}
