// Original implementation by Mykola Dvornik for mumax2
// Modified for mumax3 by Arne Vansteenkiste, 2013
__kernel void
addtworegionoommfslonczewskitorque(  __global float* __restrict            tx, __global float* __restrict               ty, __global float* __restrict      tz,
                                     __global float* __restrict            mx, __global float* __restrict               my, __global float* __restrict      mz,
                                     __global float* __restrict           Ms_,                      float           Ms_mul,
                                   __global uint8_t* __restrict       regions,
                                                        uint8_t       regionA,                    uint8_t          regionB,
                                                            int       strideX,                        int          strideY,                        int strideZ,
                                                            int            Nx,                        int               Ny,                        int      Nz,
                                                          float            j_,
                                                          float        alpha_,
                                                          float         pfix_,                      float           pfree_,
                                                          float    lambdafix_,                      float      lambdafree_,
                                                          float epsilonPrime_,
                                                          float          flt_,
                                                            int             N) {

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

    float3 m0 = make_float3(mx[I], my[I], mz[I]);
    float Ms0 = amul(Ms_, Ms_mul, I);
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

    float3 m1 = make_float3(mx[i_], my[i_], mz[i_]); // "neighbor" mag
    float Ms1 = amul(Ms_, Ms_mul, i_);
    if (is0(m1) || (Ms1 == 0.0f)) {
        return;
    }

    if ((j_ == 0.0f) || (Ms0 == 0.0f) || (Ms1 == 0.0f)) {
        return;
    }

    // Calculate for cell belonging to regionA
    float           beta0 = (HBAR / QE) * (j_ / (2.0f *flt_) );
    float            beta = beta0 / Ms0;
    float      lambdafix2 = lambdafix_ * lambdafix_;
    float     lambdafree2 = lambdafree_ * lambdafree_;
    float  lambdafreePlus = sqrt(lambdafree2 + 1.0f);
    float   lambdafixPlus = sqrt(lambdafix2 + 1.0f);
    float lambdafreeMinus = sqrt(lambdafree2 - 1.0f);
    float  lambdafixMinus = sqrt(lambdafix2 - 1.0f);
    float      plus_ratio = lambdafreePlus / lambdafixPlus;
    float     minus_ratio = 1.0f;

    if (lambdafreeMinus > 0.0f) {
        minus_ratio = lambdafixMinus / lambdafreeMinus;
    }

    // Compute q_plus and q_minus
    float  plus_factor = pfix_ * lambdafix2 * plus_ratio;
    float minus_factor = pfree_ * lambdafree2 * minus_ratio;
    float       q_plus = plus_factor + minus_factor;
    float      q_minus = plus_factor - minus_factor;
    float       lplus2 = lambdafreePlus * lambdafixPlus;
    float      lminus2 = lambdafreeMinus * lambdafixMinus;
    float        pdotm = dot(m1, m0);
    float       A_plus = lplus2 + (lminus2 * pdotm);
    float      A_minus = lplus2 - (lminus2 * pdotm);
    float      epsilon = (q_plus / A_plus) - (q_minus / A_minus);

    float A = beta * epsilon;
    float B = beta * epsilonPrime_;

    float     gilb = 1.0f / (1.0f + alpha_ * alpha_);
    float mxpxmFac = gilb * (A + alpha_ * B);
    float   pxmFac = gilb * (B - alpha_ * A);

    float3   pxm = cross(m1, m0);
    float3 mxpxm = cross(m0, pxm);

    tx[I] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
    ty[I] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
    tz[I] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;

    // Now calculate for cell in regionB
    beta        = -1.0f * beta0 / Ms1;
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
