// Original implementation by Mykola Dvornik for mumax2
// Modified for mumax3 by Arne Vansteenkiste, 2013

__kernel void
addoommfslonczewskitorque(__global float* __restrict tx, __global float* __restrict ty, __global float* __restrict tz,
						  __global float* __restrict mx, __global float* __restrict my, __global float* __restrict mz,
						  __global float* __restrict Ms_,      		float  Ms_mul,
						  __global float* __restrict jz_,      		float  jz_mul,
						  __global float* __restrict px_,      		float  px_mul,
						  __global float* __restrict py_,      		float  py_mul,
						  __global float* __restrict pz_,      		float  pz_mul,
						  __global float* __restrict alpha_,   		float  alpha_mul,
						  __global float* __restrict pfix_,    		float  pfix_mul,
						  __global float* __restrict pfree_,   		float  pfree_mul,
						  __global float* __restrict lambdafix_,    float  lambdafix_mul,
						  __global float* __restrict lambdafree_,   float  lambdafree_mul,
						  __global float* __restrict epsilonPrime_, float  epsilonPrime_mul,
						  __global float* __restrict flt_,          float  flt_mul,
						  int N) {

	int I =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (I < N) {

		float3 m = make_float3(mx[I], my[I], mz[I]);
        float  J = amul(jz_, jz_mul, I);
        float3 p = normalized(vmul(px_, py_, pz_, px_mul, py_mul, pz_mul, I));
        float  Ms           = amul(Ms_, Ms_mul, I);
        float  alpha        = amul(alpha_, alpha_mul, I);
        float  flt          = amul(flt_, flt_mul, I);
        float  pfix         = amul(pfix_, pfix_mul, I);
        float  pfree        = amul(pfree_, pfree_mul, I);
        float  lambdafix    = amul(lambdafix_, lambdafix_mul, I);
        float  lambdafree   = amul(lambdafree_, lambdafix_mul, I);
        float  epsilonPrime = amul(epsilonPrime_, epsilonPrime_mul, I);

		if (J == 0.0f || Ms == 0.0f) {
			return;
		}

		float beta    = (HBAR / QE) * (J / (2.0f *flt*Ms) );
		float lambdafix2 = lambdafix * lambdafix;
		float lambdafree2 = lambdafree * lambdafree;
		float lambdafreePlus = sqrt(lambdafree2 + 1.0f);
		float lambdafixPlus = sqrt(lambdafix2 + 1.0f);
		float lambdafreeMinus = sqrt(lambdafree2 - 1.0f);
		float lambdafixMinus = sqrt(lambdafix2 - 1.0f);
		float plus_ratio = lambdafreePlus / lambdafixPlus;
		float minus_ratio = 1.0f;
		if (lambdafreeMinus > 0) {
		   	minus_ratio = lambdafixMinus / lambdafreeMinus;
		}
		// Compute q_plus and q_minus
		float plus_factor = pfix * lambdafix2 * plus_ratio;
		float minus_factor = pfree * lambdafree2 * minus_ratio;
		float q_plus = plus_factor + minus_factor;
		float q_minus = plus_factor - minus_factor;
		float lplus2 = lambdafreePlus * lambdafixPlus;
		float lminus2 = lambdafreeMinus * lambdafixMinus;
		float pdotm = dot(p, m);
		float A_plus = lplus2 + (lminus2 * pdotm);
		float A_minus = lplus2 - (lminus2 * pdotm);
		float epsilon = (q_plus / A_plus) - (q_minus / A_minus);

		float A = beta * epsilon;
		float B = beta * epsilonPrime;

		float gilb     = 1.0f / (1.0f + alpha * alpha);
		float mxpxmFac = gilb * (A + alpha * B);
		float pxmFac   = gilb * (B - alpha * A);

		float3 pxm      = cross(p, m);
		float3 mxpxm    = cross(m, pxm);

		tx[I] += mxpxmFac * mxpxm.x + pxmFac * pxm.x;
		ty[I] += mxpxmFac * mxpxm.y + pxmFac * pxm.y;
		tz[I] += mxpxmFac * mxpxm.z + pxmFac * pxm.z;
	}
}
