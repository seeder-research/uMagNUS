// dst += prefactor * dot(a,b)
__kernel void
dotproduct(__global real_t* __restrict dst,                      real_t prefactor,
           __global real_t* __restrict  ax, __global real_t* __restrict        ay, __global real_t* __restrict az,
           __global real_t* __restrict  bx, __global real_t* __restrict        by, __global real_t* __restrict bz,
                                   int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        real_t3 A = {ax[i], ay[i], az[i]};
        real_t3 B = {bx[i], by[i], bz[i]};

        dst[i] += prefactor * dot(A, B);
    }
}
