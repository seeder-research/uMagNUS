// Add two region exchange field to Beff.
// The cells of the regions are separated
// by the displacement vector defined by
// real_t3{strideX*cellsize[X], strideY*cellsize[Y], strideZ*cellsize[Z]}
//        m: normalized magnetization
//        B: effective field in Tesla
//  sig_eff: bilinear exchange coefficient (with cell discretization) in J / m^3
// sig2_eff: biquadratic exchange coefficient (with cell discretization) in J / m^3
__kernel void
tworegionexchange_field( __global real_t* __restrict      Bx, __global real_t* __restrict       By, __global real_t* __restrict      Bz,
                         __global real_t* __restrict      mx, __global real_t* __restrict       my, __global real_t* __restrict      mz,
                         __global real_t* __restrict     Ms_,                      real_t   Ms_mul,
                        __global uint8_t* __restrict regions,
                                             uint8_t regionA,                     uint8_t  regionB,
                                                 int strideX,                         int  strideY,                         int strideZ,
                                              real_t sig_eff,                      real_t sig2_eff,
                                                 int      Nx,                         int       Ny,                         int      Nz) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    // central cell
    int I = idx(ix, iy, iz);
    if (regions[I] != regionA) {
        return;
    }

    real_t3  m0 = make_float3(mx[I], my[I], mz[I]);
    real_t  Ms0 = amul(Ms_, Ms_mul, I);

    if (is0(m0) || (Ms0 == 0.0f)) {
        return;
    }

    int cX = ix + strideX;
    int cY = iy + strideY;
    int cZ = iz + strideZ;

    if ((cX >= Nx) || (cY >= Ny) || (cZ >= Nz)) {
        return;
    }

    int i_ = idx(cX, cY, cZ); // "neighbor" index
    if (regions[i_] != regionB) {
        return;
    }

    real_t3  m1 = make_float3(mx[i_], my[i_], mz[i_]); // "neighbor" mag
    real_t  Ms1 = amul(Ms_, Ms_mul, i_);
    if (is0(m1) || (Ms1 == 0.0f)) {
        return;
    }

    real_t3     B = m0 - m1;
    real_t   dot1 = dot(m0, m1);
    real_t    fac = 2.0f * (sig_eff + 2.0f * sig2_eff * dot1);
    real_t  invMs = inv_Msat(Ms_, Ms_mul, I);

    if (Bx != NULL) {
        Bx[I]  -= B.x * (fac*invMs);
        Bx[i_] += B.x * (fac*invMs);
    }
    if (By != NULL) {
        By[I]  -= B.y * (fac*invMs);
        By[i_] += B.y * (fac*invMs);
    }
    if (Bz != NULL) {
        Bz[I]  -= B.z * (fac*invMs);
        Bz[i_] += B.z * (fac*invMs);
    }
}
