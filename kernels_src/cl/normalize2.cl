// normalize vector {vx, vy, vz} to unit length, unless length or vol are zero.
__kernel void
normalize2(__global real_t* __restrict vx, __global real_t* __restrict vy, __global real_t* __restrict vz, __global real_t* __restrict vol, int N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t  v = (vol == NULL) ? (real_t)1.0 : vol[i];
        real_t3 V = {v*vx[i], v*vy[i], v*vz[i]};

        V = normalize(V);
        if (v == (real_t)0.0) {
            vx[i] = 0.0;
            vy[i] = 0.0;
            vz[i] = 0.0;
        } else {
            vx[i] = V.x;
            vy[i] = V.y;
            vz[i] = V.z;
        }
    }
}
