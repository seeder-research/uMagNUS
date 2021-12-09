// Add voltage-controlled magnetic anisotropy field to B.
// https://www.nature.com/articles/s42005-019-0189-6.pdf
__kernel void
addvoltagecontrolledanisotropy2(__global real_t* __restrict         Bx, __global real_t* __restrict            By, __global real_t* __restrict Bz,
                                __global real_t* __restrict         mx, __global real_t* __restrict            my, __global real_t* __restrict mz,
                                __global real_t* __restrict        Ms_,                      real_t        Ms_mul,
                                __global real_t* __restrict vcmaCoeff_,                      real_t vcmaCoeff_mul,
                                __global real_t* __restrict   voltage_,                      real_t   voltage_mul,
                                __global real_t* __restrict        ux_,                      real_t        ux_mul,
                                __global real_t* __restrict        uy_,                      real_t        uy_mul,
                                __global real_t* __restrict        uz_,                      real_t        uz_mul,
                                                        int          N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3         u = normalized(vmul(ux_, uy_, uz_, ux_mul, uy_mul, uz_mul, i));
        real_t      invMs = inv_Msat(Ms_, Ms_mul, i);
        real_t  vcmaCoeff = amul(vcmaCoeff_, vcmaCoeff_mul, i) * invMs;
        real_t    voltage = amul(voltage_, voltage_mul, i) * invMs;
        real_t3         m = {mx[i], my[i], mz[i]};
        real_t         mu = dot(m, u);
        real_t3        Ba = 2.0f*vcmaCoeff*voltage*    (mu)*u;

        Bx[i] += Ba.x;
        By[i] += Ba.y;
        Bz[i] += Ba.z;
    }
}
