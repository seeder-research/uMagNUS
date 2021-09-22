// 2D Z (out-of-plane only) micromagnetic kernel multiplication:
// Mz = Kzz * Mz
// Using the same symmetries as kernmulrsymm3d.cl
__kernel void
kernmulRSymm2Dz(__global float* __restrict  fftMz, __global float* __restrict  fftKzz, int Nx, int Ny) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);

	if(ix>= Nx || iy>=Ny) {
		return;
	}

	int I = iy*Nx + ix;
	int e = 2 * I;

	float reMz = fftMz[e  ];
	float imMz = fftMz[e+1];

	if (iy > Ny/2) {
		iy = Ny-iy;
	}
	I = iy*Nx + ix;

	float Kzz = fftKzz[I];

	fftMz[e  ] = reMz * Kzz;
	fftMz[e+1] = imMz * Kzz;
}

