// Original implementation by Mykola Dvornik for mumax2
// Modified for mumax3 by Arne Vansteenkiste, 2013

__kernel void
addoommfslonczewskitorque(__global real_t* __restrict            tx, __global real_t* __restrict              ty, __global real_t* __restrict tz,
                          __global real_t* __restrict            mx, __global real_t* __restrict              my, __global real_t* __restrict mz,
                          __global real_t* __restrict           Ms_,                      real_t          Ms_mul,
                          __global real_t* __restrict           jz_,                      real_t          jz_mul,
                          __global real_t* __restrict           px_,                      real_t          px_mul,
                          __global real_t* __restrict           py_,                      real_t          py_mul,
                          __global real_t* __restrict           pz_,                      real_t           pz_mul,
                          __global real_t* __restrict        alpha_,                      real_t        alpha_mul,
                          __global real_t* __restrict         pfix_,                      real_t         pfix_mul,
                          __global real_t* __restrict        pfree_,                      real_t        pfree_mul,
                          __global real_t* __restrict    lambdafix_,                      real_t    lambdafix_mul,
                          __global real_t* __restrict   lambdafree_,                      real_t   lambdafree_mul,
                          __global real_t* __restrict epsilonPrime_,                      real_t epsilonPrime_mul,
                          __global real_t* __restrict          flt_,                      real_t          flt_mul,
                                                  int             N) {

    int     I = ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int I = gid; I < N; I += gsize) {

        real_t3 m = make_float3(mx[I], my[I], mz[I]);
        real_t  J = amul(jz_, jz_mul, I);
        real_t3 p = normalized(vmul(px_, py_, pz_, px_mul, py_mul, pz_mul, I));

        real_t  Ms           = amul(Ms_, Ms_mul, I);
        real_t  alpha        = amul(alpha_, alpha_mul, I);
        real_t  flt          = amul(flt_, flt_mul, I);
        real_t  pfix         = amul(pfix_, pfix_mul, I);
        real_t  pfree        = amul(pfree_, pfree_mul, I);
        real_t  lambdafix    = amul(lambdafix_, lambdafix_mul, I);
        real_t  lambdafree   = amul(lambdafree_, lambdafix_mul, I);
        real_t  epsilonPrime = amul(epsilonPrime_, epsilonPrime_mul, I);

        if ((J == 0.0f) || (Ms == 0.0f)) {
            return;
        }

        real_t            beta = (HBAR / QE) * (J / (2.0f *flt*Ms) );
        real_t      lambdafix2 = lambdafix * lambdafix;
        real_t     lambdafree2 = lambdafree * lambdafree;
        real_t  lambdafreePlus = sqrt(lambdafree2 + 1.0f);
        real_t   lambdafixPlus = sqrt( lambdafix2 + 1.0f);
        real_t lambdafreeMinus = sqrt(lambdafree2 - 1.0f);
        real_t  lambdafixMinus = sqrt( lambdafix2 - 1.0f);
        real_t      plus_ratio = lambdafreePlus / lambdafixPlus;
        real_t     minus_ratio = 1.0f;

        if (lambdafreeMinus > 0) {
            minus_ratio = lambdafixMinus / lambdafreeMinus;
        }

        // Compute q_plus and q_minus
        real_t  plus_factor = pfix * lambdafix2 * plus_ratio;
        real_t minus_factor = pfree * lambdafree2 * minus_ratio;
        real_t       q_plus = plus_factor + minus_factor;
        real_t      q_minus = plus_factor - minus_factor;
        real_t       lplus2 = lambdafreePlus * lambdafixPlus;
        real_t      lminus2 = lambdafreeMinus * lambdafixMinus;
        real_t        pdotm = dot(p, m);
        real_t       A_plus = lplus2 + (lminus2 * pdotm);
        real_t      A_minus = lplus2 - (lminus2 * pdotm);
        real_t      epsilon = (q_plus / A_plus) - (q_minus / A_minus);

        real_t A = beta * epsilon;
        real_t B = beta * epsilonPrime;

        real_t gilb     = 1.0f / (1.0f + alpha * alpha);
        real_t mxpxmFac = gilb * (A + alpha * B);
        real_t pxmFac   = gilb * (B - alpha * A);

        real_t3 pxm      = cross(p, m);
        real_t3 mxpxm    = cross(m, pxm);

        tx[I] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
        ty[I] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
        tz[I] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;
    }
}
