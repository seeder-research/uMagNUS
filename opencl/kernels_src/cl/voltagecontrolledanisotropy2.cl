// Add voltage-controlled magnetic anisotropy field to B.
// https://www.nature.com/articles/s42005-019-0189-6.pdf
__kernel void
addvoltagecontrolledanisotropy2(__global float* __restrict         Bx, __global float* __restrict            By, __global float* __restrict Bz,
                                __global float* __restrict         mx, __global float* __restrict            my, __global float* __restrict mz,
                                __global float* __restrict        Ms_,                      float        Ms_mul,
                                __global float* __restrict vcmaCoeff_,                      float vcmaCoeff_mul,
                                __global float* __restrict   voltage_,                      float   voltage_mul,
                                __global float* __restrict        ux_,                      float        ux_mul,
                                __global float* __restrict        uy_,                      float        uy_mul,
                                __global float* __restrict        uz_,                      float        uz_mul,
                                                       int          N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        float3         u = normalized(vmul(ux_, uy_, uz_, ux_mul, uy_mul, uz_mul, i));
        float      invMs = inv_Msat(Ms_, Ms_mul, i);
        float  vcmaCoeff = amul(vcmaCoeff_, vcmaCoeff_mul, i) * invMs;
        float    voltage = amul(voltage_, voltage_mul, i) * invMs;
        float3         m = {mx[i], my[i], mz[i]};
        float         mu = dot(m, u);
        float3        Ba = 2.0f*vcmaCoeff*voltage*    (mu)*u;

        Bx[i] += Ba.x;
        By[i] += Ba.y;
        Bz[i] += Ba.z;
    }
}
