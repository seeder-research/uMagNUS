// Original implementation by Mykola Dvornik for mumax2
// Modified for mumax3 by Arne Vansteenkiste, 2013, 2016

__kernel void
addslonczewskitorque2(__global float* __restrict tx, __global float* __restrict ty, __global float* __restrict tz,
                      __global float* __restrict mx, __global float* __restrict my, __global float* __restrict mz,
                      __global float* __restrict Ms_,         float  Ms_mul,
                      __global float* __restrict jz_,         float  jz_mul,
                      __global float* __restrict px_,         float  px_mul,
                      __global float* __restrict py_,         float  py_mul,
                      __global float* __restrict pz_,         float  pz_mul,
                      __global float* __restrict alpha_,      float  alpha_mul,
                      __global float* __restrict pol_,        float  pol_mul,
                      __global float* __restrict lambda_,     float  lambda_mul,
                      __global float* __restrict epsPrime_,   float  epsPrime_mul,
                      __global float* __restrict thickness_,  float  thickness_mul,
                      float meshThickness,
                      float freeLayerPosition,
                      int N) {

    int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
    if (i < N) {

        float3 m = make_float3(mx[i], my[i], mz[i]);
        float  J = amul(jz_, jz_mul, i);
        float3 p = normalized(vmul(px_, py_, pz_, px_mul, py_mul, pz_mul, i));
        float  Ms           = amul(Ms_, Ms_mul, i);
        float  alpha        = amul(alpha_, alpha_mul, i);
        float  pol          = amul(pol_, pol_mul, i);
        float  lambda       = amul(lambda_, lambda_mul, i);
        float  epsilonPrime = amul(epsPrime_, epsPrime_mul, i);

        float thickness = amul(thickness_, thickness_mul, i);
        if (thickness == 0.0) { // if thickness is not set, use the thickness of the mesh instead
            thickness = meshThickness;
        }
        thickness *= freeLayerPosition; // switch sign if fixedlayer is at the bottom

        if (J == 0.0f || Ms == 0.0f) {
            return;
        }

        float beta    = (HBAR / QE) * (J / (thickness*Ms) );
        float lambda2 = lambda * lambda;
        float epsilon = pol * lambda2 / ((lambda2 + 1.0f) + (lambda2 - 1.0f) * dot(p, m));

        float A = beta * epsilon;
        float B = beta * epsilonPrime;

        float gilb     = 1.0f / (1.0f + alpha * alpha);
        float mxpxmFac = gilb * (A + alpha * B);
        float pxmFac   = gilb * (B - alpha * A);

        float3 pxm      = cross(p, m);
        float3 mxpxm    = cross(m, pxm);

        tx[i] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
        ty[i] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
        tz[i] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;
    }
}

