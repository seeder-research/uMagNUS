// 2D XY (in-plane) micromagnetic kernel multiplication:
// |Mx| = |Kxx Kxy| * |Mx|
// |My|   |Kyx Kyy|   |My|
// Using the same symmetries as kernmulrsymm3d.cl
__kernel void
kernmulRSymm2Dxy(__global real_t* __restrict  fftMx, __global real_t* __restrict  fftMy,
                 __global real_t* __restrict fftKxx, __global real_t* __restrict fftKyy, __global real_t* __restrict fftKxy,
                                         int     Nx,                         int     Ny) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);

    if ((ix>= Nx) || (iy>=Ny)) {
        return;
    }

    int I = iy*Nx + ix;
    int e = 2 * I;

    real_t reMx = fftMx[e  ];
    real_t imMx = fftMx[e+1];
    real_t reMy = fftMy[e  ];
    real_t imMy = fftMy[e+1];

    // symmetry factor
    real_t fxy = (real_t)1.0;
    if (iy > Ny/2) {
         iy = Ny-iy;
        fxy = -fxy;
    }
    I = iy*Nx + ix;

    real_t Kxx = fftKxx[I];
    real_t Kyy = fftKyy[I];
    real_t Kxy = fxy * fftKxy[I];

    fftMx[e  ] = reMx * Kxx + reMy * Kxy;
    fftMx[e+1] = imMx * Kxx + imMy * Kxy;
    fftMy[e  ] = reMx * Kxy + reMy * Kyy;
    fftMy[e+1] = imMx * Kxy + imMy * Kyy;
}
