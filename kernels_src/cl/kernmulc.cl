__kernel void
kernmulC(__global real_t* __restrict fftM, __global real_t* __restrict fftK, int Nx, int Ny) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);

    if ((ix>= Nx) || (iy>=Ny)) {
        return;
    }

    int I = iy*Nx + ix;
    int e = 2 * I;

    real_t reM = fftM[e  ];
    real_t imM = fftM[e+1];
    real_t reK = fftK[e  ];
    real_t imK = fftK[e+1];

    fftM[e  ] = reM * reK - imM * imK;
    fftM[e+1] = reM * imK + imM * reK;
}
