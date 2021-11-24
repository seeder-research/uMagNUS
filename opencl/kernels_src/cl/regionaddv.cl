// add region-based vector to dst:
// dst[i] += LUT[region[i]]
__kernel void
regionaddv(__global   float* __restrict    dstx, __global float* __restrict dsty, __global float* __restrict dstz,
           __global   float* __restrict    LUTx, __global float* __restrict LUTy, __global float* __restrict LUTz,
           __global uint8_t* __restrict regions,
                                    int       N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        uint8_t r = regions[i];
        dstx[i] += LUTx[r];
        dsty[i] += LUTy[r];
        dstz[i] += LUTz[r];
    }
}
