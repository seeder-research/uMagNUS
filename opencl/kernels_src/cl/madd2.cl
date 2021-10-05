// dst[i] = fac1*src1[i] + fac2*src2[i];
__kernel void
madd2(__global float* __restrict  dst,
      __global float* __restrict  src1, float fac1,
      __global float* __restrict  src2, float fac2, int N) {

    int gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = fac1*src1[i] + fac2*src2[i];
    }
}
