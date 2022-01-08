// dst[i] = a[i] / b[i]
__kernel void
divide(__global float* __restrict dst, __global float* __restrict a, __global float* __restrict b, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);
    for (int i = gid; i < N; i += gsize) {
        dst[i] = ((a[i] == 0) || (b[i] == 0)) ? 0.0f : (a[i] / b[i]);
    }
}