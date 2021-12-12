// set dst to zero in cells where mask != 0
__kernel void
zeromask(__global real_t* __restrict dst, __global real_t* maskLUT, __global uint8_t* regions, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        if (maskLUT[regions[i]] != 0){
            dst[i] = (real_t)0.0;
        }
    }
}
