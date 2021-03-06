// 3D micromagnetic kernel multiplication:
//
// |Mx|   |Kxx Kxy Kxz|   |Mx|
// |My| = |Kxy Kyy Kyz| * |My|
// |Mz|   |Kxz Kyz Kzz|   |Mz|
//
// ~kernel has mirror symmetry along Y and Z-axis,
// apart form first row,
// and is only stored (roughly) half:
//
// K11, K22, K02:
// xxxxx
// aaaaa
// bbbbb
// ....
// bbbbb
// aaaaa
//
// K12:
// xxxxx
// aaaaa
// bbbbb
// ...
// -bbbb
// -aaaa

__kernel void
kernmulRSymm3D(__global real_t* __restrict  fftMx, __global real_t* __restrict  fftMy, __global real_t* __restrict  fftMz,
               __global real_t* __restrict fftKxx, __global real_t* __restrict fftKyy, __global real_t* __restrict fftKzz,
               __global real_t* __restrict fftKyz, __global real_t* __restrict fftKxz, __global real_t* __restrict fftKxy,
                                       int     Nx,                         int     Ny,                         int     Nz) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix>= Nx) || (iy>= Ny) || (iz>=Nz)) {
        return;
    }

    // fetch (complex) FFT'ed magnetization
    int I = (iz*Ny + iy)*Nx + ix;
    int e = 2 * I;

    real_t reMx = fftMx[e  ];
    real_t imMx = fftMx[e+1];
    real_t reMy = fftMy[e  ];
    real_t imMy = fftMy[e+1];
    real_t reMz = fftMz[e  ];
    real_t imMz = fftMz[e+1];

    // fetch kernel

    // minus signs are added to some elements if
    // reconstructed from symmetry.
    real_t signYZ = (real_t)1.0;
    real_t signXZ = (real_t)1.0;
    real_t signXY = (real_t)1.0;

    // use symmetry to fetch from redundant parts:
    // mirror index into first quadrant and set signs.
    if (iy > Ny/2) {
            iy = Ny-iy;
        signYZ = -signYZ;
        signXY = -signXY;
    }
    if (iz > Nz/2) {
            iz = Nz-iz;
        signYZ = -signYZ;
        signXZ = -signXZ;
    }

    // fetch kernel element from non-redundant part
    // and apply minus signs for mirrored parts.
    I = (iz*(Ny/2+1) + iy)*Nx + ix; // Ny/2+1: only half is stored

    real_t Kxx = fftKxx[I];
    real_t Kyy = fftKyy[I];
    real_t Kzz = fftKzz[I];
    real_t Kyz = fftKyz[I] * signYZ;
    real_t Kxz = fftKxz[I] * signXZ;
    real_t Kxy = fftKxy[I] * signXY;

    // m * K matrix multiplication, overwrite m with result.
    fftMx[e  ] = reMx * Kxx + reMy * Kxy + reMz * Kxz;
    fftMx[e+1] = imMx * Kxx + imMy * Kxy + imMz * Kxz;
    fftMy[e  ] = reMx * Kxy + reMy * Kyy + reMz * Kyz;
    fftMy[e+1] = imMx * Kxy + imMy * Kyy + imMz * Kyz;
    fftMz[e  ] = reMx * Kxz + reMy * Kyz + reMz * Kzz;
    fftMz[e+1] = imMx * Kxz + imMy * Kyz + imMz * Kzz;
}
