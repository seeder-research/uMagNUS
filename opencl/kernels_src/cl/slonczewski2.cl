// Original implementation by Mykola Dvornik for mumax2
// Modified for mumax3 by Arne Vansteenkiste, 2013, 2016

__kernel void
addslonczewskitorque2(__global real_t* __restrict                tx, __global real_t* __restrict             ty, __global real_t* __restrict tz,
                      __global real_t* __restrict                mx, __global real_t* __restrict             my, __global real_t* __restrict mz,
                      __global real_t* __restrict               Ms_,                      real_t         Ms_mul,
                      __global real_t* __restrict               jz_,                      real_t         jz_mul,
                      __global real_t* __restrict               px_,                      real_t         px_mul,
                      __global real_t* __restrict               py_,                      real_t         py_mul,
                      __global real_t* __restrict               pz_,                      real_t         pz_mul,
                      __global real_t* __restrict            alpha_,                      real_t      alpha_mul,
                      __global real_t* __restrict              pol_,                      real_t        pol_mul,
                      __global real_t* __restrict           lambda_,                      real_t     lambda_mul,
                      __global real_t* __restrict         epsPrime_,                      real_t   epsPrime_mul,
                      __global real_t* __restrict        thickness_,                      real_t  thickness_mul,
                                           real_t     meshThickness,
                                           real_t freeLayerPosition,
                                              int                 N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3            m = make_float3(mx[i], my[i], mz[i]);
        real_t             J = amul(jz_, jz_mul, i);
        real_t3            p = normalized(vmul(px_, py_, pz_, px_mul, py_mul, pz_mul, i));
        real_t            Ms = amul(Ms_, Ms_mul, i);
        real_t         alpha = amul(alpha_, alpha_mul, i);
        real_t           pol = amul(pol_, pol_mul, i);
        real_t        lambda = amul(lambda_, lambda_mul, i);
        real_t  epsilonPrime = amul(epsPrime_, epsPrime_mul, i);
        real_t     thickness = amul(thickness_, thickness_mul, i);

        if (thickness == 0.0) { // if thickness is not set, use the thickness of the mesh instead
            thickness = meshThickness;
        }
        thickness *= freeLayerPosition; // switch sign if fixedlayer is at the bottom

        if (J == 0.0f || Ms == 0.0f) {
            return;
        }

        real_t    beta = (HBAR / QE) * (J / (thickness*Ms) );
        real_t lambda2 = lambda * lambda;
        real_t epsilon = pol * lambda2 / ((lambda2 + 1.0f) + (lambda2 - 1.0f) * dot(p, m));

        real_t A = beta * epsilon;
        real_t B = beta * epsilonPrime;

        real_t     gilb = 1.0f / (1.0f + alpha * alpha);
        real_t mxpxmFac = gilb * (A + alpha * B);
        real_t   pxmFac = gilb * (B - alpha * A);

        real_t3   pxm = cross(p, m);
        real_t3 mxpxm = cross(m, pxm);

        tx[i] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
        ty[i] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
        tz[i] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;
    }
}
