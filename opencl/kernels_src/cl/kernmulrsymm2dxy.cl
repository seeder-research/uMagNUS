// 2D XY (in-plane) micromagnetic kernel multiplication:
// |Mx| = |Kxx Kxy| * |Mx|
// |My|   |Kyx Kyy|   |My|
// Using the same symmetries as kernmulrsymm3d.cl
__kernel void
kernmulRSymm2Dxy(__global float* __restrict  fftMx,  __global float* __restrict  fftMy,
                 __global float* __restrict  fftKxx, __global float* __restrict  fftKyy, __global float* __restrict  fftKxy,
                 int Nx, int Ny) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);

	if(ix>= Nx || iy>=Ny) {
		return;
	}

	int I = iy*Nx + ix;
	int e = 2 * I;

	float reMx = fftMx[e  ];
	float imMx = fftMx[e+1];
	float reMy = fftMy[e  ];
	float imMy = fftMy[e+1];

	// symmetry factor
	float fxy = 1.0f;
	if (iy > Ny/2) {
		iy = Ny-iy;
		fxy = -fxy;
	}
	I = iy*Nx + ix;

	float Kxx = fftKxx[I];
	float Kyy = fftKyy[I];
	float Kxy = fxy * fftKxy[I];

	fftMx[e  ] = reMx * Kxx + reMy * Kxy;
	fftMx[e+1] = imMx * Kxx + imMy * Kxy;
	fftMy[e  ] = reMx * Kxy + reMy * Kyy;
	fftMy[e+1] = imMx * Kxy + imMy * Kyy;
}

