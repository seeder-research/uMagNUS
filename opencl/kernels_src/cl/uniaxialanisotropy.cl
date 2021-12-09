// Add uniaxial magnetocrystalline anisotropy field to B.
// http://www.southampton.ac.uk/~fangohr/software/oxs_uniaxial4.html
__kernel void
adduniaxialanisotropy(__global real_t* __restrict  Bx, __global real_t* __restrict     By, __global real_t* __restrict  Bz,
                      __global real_t* __restrict  mx, __global real_t* __restrict     my, __global real_t* __restrict  mz,
                      __global real_t* __restrict Ms_,                      real_t Ms_mul,
                      __global real_t* __restrict K1_,                      real_t K1_mul,
                      __global real_t* __restrict ux_,                      real_t ux_mul,
                      __global real_t* __restrict uy_,                      real_t uy_mul,
                      __global real_t* __restrict uz_,                      real_t uz_mul,
                                              int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3     u = normalized(vmul(ux_, uy_, uz_, ux_mul, uy_mul, uz_mul, i));
        real_t  invMs = inv_Msat(Ms_, Ms_mul, i);
        real_t     K1 = amul(K1_, K1_mul, i);

        K1  *= invMs;

        real_t3  m = {mx[i], my[i], mz[i]};
        real_t  mu = dot(m, u);
        real_t3 Ba = 2.0f*K1*(mu)*u;

        Bx[i] += Ba.x;
        By[i] += Ba.y;
        Bz[i] += Ba.z;
    }
}
