// Add magneto-elastic coupling field to B.
// H = - δUmel / δM, 
// where Umel is magneto-elastic energy denstiy given by the eq. (12.18) of Gurevich&Melkov "Magnetization Oscillations and Waves", CRC Press, 1996
__kernel void
addmagnetoelasticfield(__global real_t* __restrict   Bx, __global real_t* __restrict      By, __global real_t* __restrict  Bz,
                       __global real_t* __restrict   mx, __global real_t* __restrict      my, __global real_t* __restrict  mz,
                       __global real_t* __restrict exx_,                      real_t exx_mul,
                       __global real_t* __restrict eyy_,                      real_t eyy_mul,
                       __global real_t* __restrict ezz_,                      real_t ezz_mul,
                       __global real_t* __restrict exy_,                      real_t exy_mul,
                       __global real_t* __restrict exz_,                      real_t exz_mul,
                       __global real_t* __restrict eyz_,                      real_t eyz_mul,
                       __global real_t* __restrict  B1_,                      real_t  B1_mul, 
                       __global real_t* __restrict  B2_,                      real_t  B2_mul,
                       __global real_t* __restrict  Ms_,                      real_t  Ms_mul,
                                               int    N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int I = gid; I < N; I += gsize) {

        real_t Exx = amul(exx_, exx_mul, I);
        real_t Eyy = amul(eyy_, eyy_mul, I);
        real_t Ezz = amul(ezz_, ezz_mul, I);

        real_t Exy = amul(exy_, exy_mul, I);
        real_t Eyx = Exy;

        real_t Exz = amul(exz_, exz_mul, I);
        real_t Ezx = Exz;

        real_t Eyz = amul(eyz_, eyz_mul, I);
        real_t Ezy = Eyz;

        real_t invMs = inv_Msat(Ms_, Ms_mul, I);

        real_t B1 = amul(B1_, B1_mul, I) * invMs;
        real_t B2 = amul(B2_, B2_mul, I) * invMs;

        real_t3 m = {mx[I], my[I], mz[I]};

        Bx[I] += -((real_t)2.0*B1*m.x*Exx + B2*(m.y*Exy + m.z*Exz));
        By[I] += -((real_t)2.0*B1*m.y*Eyy + B2*(m.x*Eyx + m.z*Eyz));
        Bz[I] += -((real_t)2.0*B1*m.z*Ezz + B2*(m.x*Ezx + m.y*Ezy));
    }
}
