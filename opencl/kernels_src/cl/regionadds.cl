// add region-based scalar to dst:
// dst[i] += LUT[region[i]]
__kernel void
regionadds(__global   float* __restrict     dst,
           __global   float* __restrict     LUT,
           __global uint8_t* __restrict regions, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        uint8_t r = regions[i];
        dst[i] += LUT[r];
    }
}
