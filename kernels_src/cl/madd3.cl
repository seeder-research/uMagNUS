// dst[i] = fac1 * src1[i] + fac2 * src2[i] + fac3 * src3[i]
__kernel void
madd3(__global real_t* __restrict  dst,
      __global real_t* __restrict src1, real_t fac1,
      __global real_t* __restrict src2, real_t fac2,
      __global real_t* __restrict src3, real_t fac3, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (fac1 * src1[i]) + (fac2 * src2[i] + fac3 * src3[i]);
        // parens for better accuracy heun solver.
    }
}
