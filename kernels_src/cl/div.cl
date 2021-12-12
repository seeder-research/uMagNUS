// dst[i] = a[i] / b[i]
__kernel void
pointwise_div(__global real_t* __restrict dst, __global real_t* __restrict a, __global real_t* __restrict b, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = (b[i] == (real_t)0.0) ? (real_t)0.0 : (a[i] / b[i]);
    }
}
