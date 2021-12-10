// TODO: this could act on x,y,z, so that we need to call it only once.
__kernel void
settemperature2(__global real_t* __restrict      B, __global real_t* __restrict     noise, real_t kB2_VgammaDt,
                __global real_t* __restrict    Ms_,                      real_t    Ms_mul,
                __global real_t* __restrict  temp_,                      real_t  temp_mul,
                __global real_t* __restrict alpha_,                      real_t alpha_mul,
                                        int      N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {
        real_t invMs = inv_Msat(Ms_, Ms_mul, i);
        real_t  temp = amul(temp_, temp_mul, i);
        real_t alpha = amul(alpha_, alpha_mul, i);

        B[i] = noise[i] * sqrt((kB2_VgammaDt * alpha * temp * invMs ));
    }
}
