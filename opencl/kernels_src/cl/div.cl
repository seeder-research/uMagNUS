// dst[i] = a[i] / b[i]
__kernel void
pointwise_div(__global float* __restrict  dst, __global float* __restrict  a, __global float* __restrict b, int N) {

    int gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        if (b[i] != 0.0f) {
            dst[i] = a[i] / b[i];
        } else {
            dst[i] = 0.0f;
        }
    }
}
