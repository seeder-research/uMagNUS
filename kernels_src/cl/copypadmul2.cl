// Copy src (size S, smaller) into dst (size D, larger),
// and multiply by Bsat * vol
__kernel void
copypadmul2(__global real_t* __restrict dst,    int     Dx, int Dy, int Dz,
            __global real_t* __restrict src,    int     Sx, int Sy, int Sz,
            __global real_t* __restrict Ms_, real_t Ms_mul,
            __global real_t* __restrict vol) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix<Sx) && (iy<Sy) && (iz<Sz)) {
        int        sI = index(ix, iy, iz, Sx, Sy, Sz);  // source index
        real_t tmpFac = amul(Ms_, Ms_mul, sI);
        real_t   Bsat = MU0 * tmpFac;
        real_t      v = amul(vol, 1.0f, sI);

        dst[index(ix, iy, iz, Dx, Dy, Dz)] = Bsat * v * src[sI];
    }
}
