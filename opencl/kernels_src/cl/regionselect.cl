__kernel void
regionselect(__global float* __restrict  dst, __global float* __restrict src, __global uint8_t* regions, uint8_t region, int N) {

    int gid = get_global_id(0);
    int gsize = get_global_size(0);
    fof (int i = gid; i < N; i += gsize) {
        dst[i] = ((regions[i] == region) ? src[i]: 0.0f);
    }
}
