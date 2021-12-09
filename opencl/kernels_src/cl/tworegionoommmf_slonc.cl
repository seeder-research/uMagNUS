// Original implementation by Mykola Dvornik for mumax2
// Modified for mumax3 by Arne Vansteenkiste, 2013
__kernel void
addtworegionoommfslonczewskitorque( __global real_t* __restrict            tx, __global real_t* __restrict               ty, __global real_t* __restrict      tz,
                                    __global real_t* __restrict            mx, __global real_t* __restrict               my, __global real_t* __restrict      mz,
                                    __global real_t* __restrict           Ms_,                      real_t           Ms_mul,
                                   __global uint8_t* __restrict       regions,
                                                        uint8_t       regionA,                     uint8_t          regionB,
                                                            int       strideX,                         int          strideY,                         int strideZ,
                                                            int            Nx,                         int               Ny,                         int      Nz,
                                                         real_t            j_,
                                                         real_t        alpha_,
                                                         real_t         pfix_,                      real_t           pfree_,
                                                         real_t    lambdafix_,                      real_t      lambdafree_,
                                                         real_t epsilonPrime_,
                                                         real_t          flt_) {

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

    if ((j_ == 0.0f) || (Ms0 == 0.0f) || (Ms1 == 0.0f)) {
        return;
    }

    // Calculate for cell belonging to regionA
    real_t           beta0 = (HBAR / QE) * (j_ / (2.0f *flt_) );
    real_t            beta = beta0 / Ms0;
    real_t      lambdafix2 = lambdafix_ * lambdafix_;
    real_t     lambdafree2 = lambdafree_ * lambdafree_;
    real_t  lambdafreePlus = sqrt(lambdafree2 + 1.0f);
    real_t   lambdafixPlus = sqrt(lambdafix2 + 1.0f);
    real_t lambdafreeMinus = sqrt(lambdafree2 - 1.0f);
    real_t  lambdafixMinus = sqrt(lambdafix2 - 1.0f);
    real_t      plus_ratio = lambdafreePlus / lambdafixPlus;
    real_t     minus_ratio = 1.0f;

    if (lambdafreeMinus > 0.0f) {
        minus_ratio = lambdafixMinus / lambdafreeMinus;
    }

    // Compute q_plus and q_minus
    real_t  plus_factor = pfix_ * lambdafix2 * plus_ratio;
    real_t minus_factor = pfree_ * lambdafree2 * minus_ratio;
    real_t       q_plus = plus_factor + minus_factor;
    real_t      q_minus = plus_factor - minus_factor;
    real_t       lplus2 = lambdafreePlus * lambdafixPlus;
    real_t      lminus2 = lambdafreeMinus * lambdafixMinus;
    real_t        pdotm = dot(m1, m0);
    real_t       A_plus = lplus2 + (lminus2 * pdotm);
    real_t      A_minus = lplus2 - (lminus2 * pdotm);
    real_t      epsilon = (q_plus / A_plus) - (q_minus / A_minus);

    real_t A = beta * epsilon;
    real_t B = beta * epsilonPrime_;

    real_t     gilb = 1.0f / (1.0f + alpha_ * alpha_);
    real_t mxpxmFac = gilb * (A + alpha_ * B);
    real_t   pxmFac = gilb * (B - alpha_ * A);

    real_t3   pxm = cross(m1, m0);
    real_t3 mxpxm = cross(m0, pxm);

    tx[I] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
    ty[I] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
    tz[I] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;

    // Now calculate for cell in regionB
    beta        = beta0 / Ms1;
    plus_ratio  = lambdafixPlus / lambdafreePlus;
    minus_ratio = 1.0f;

    if (lambdafixMinus > 0.0f) {
        minus_ratio = lambdafreeMinus / lambdafixMinus;
    }

    // Compute q_plus and q_minus
    plus_factor  = pfree_ * lambdafree2 * plus_ratio;
    minus_factor = pfix_ * lambdafix2 * minus_ratio;
    q_plus       = plus_factor + minus_factor;
    q_minus      = plus_factor - minus_factor;
    epsilon      = (q_plus / A_plus) - (q_minus / A_minus);

    A = beta * epsilon;
    B = beta * epsilonPrime_;

    mxpxmFac = gilb * (A + alpha_ * B);
    pxmFac   = gilb * (B - alpha_ * A);

    mxpxm = cross(m1, pxm);

    tx[i_] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
    ty[i_] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
    tz[i_] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;
}
