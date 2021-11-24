// decode the regions+LUT pair into an uncompressed array
__kernel void
regiondecode(__global float* __restrict dst, __global float* __restrict LUT, __global uint8_t* regions, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        dst[i] = LUT[regions[i]];
    }
}
