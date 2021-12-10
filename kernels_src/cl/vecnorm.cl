__kernel void
vecnorm(__global real_t* __restrict dst,
        __global real_t* __restrict  ax, __global real_t* __restrict ay, __global real_t* __restrict az,
                                int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        real_t3 A = {ax[i], ay[i], az[i]};
        dst[i] = sqrt(dot(A, A));
    }
}
