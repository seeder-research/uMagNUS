__kernel void
crossproduct(__global real_t* __restrict dstx, __global real_t* __restrict dsty, __global real_t* __restrict dstz,
             __global real_t* __restrict   ax, __global real_t* __restrict   ay, __global real_t* __restrict   az,
             __global real_t* __restrict   bx, __global real_t* __restrict   by, __global real_t* __restrict   bz,
                                     int    N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);
    for (int i = gid; i < N; i += gsize) {
        real_t3   A = {ax[i], ay[i], az[i]};
        real_t3   B = {bx[i], by[i], bz[i]};
        real_t3 AxB = cross(A, B);

        dstx[i] = AxB.x;
        dsty[i] = AxB.y;
        dstz[i] = AxB.z;
    }
}
