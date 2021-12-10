// Steepest descent energy minimizer
__kernel void
minimize(__global real_t* __restrict  mx, __global real_t* __restrict  my, __global real_t* __restrict  mz,
         __global real_t* __restrict m0x, __global real_t* __restrict m0y, __global real_t* __restrict m0z,
         __global real_t* __restrict  tx, __global real_t* __restrict  ty, __global real_t* __restrict  tz,
                              real_t  dt,                         int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3 m0 = {m0x[i], m0y[i], m0z[i]};
        real_t3  t = {tx[i], ty[i], tz[i]};

        real_t       t2 = dt*dt*dot(t, t);
        real_t3  result = (4.0f - t2) * m0 + 4.0f * dt * t;
        real_t  divisor = 4.0f + t2;

        mx[i] = result.x / divisor;
        my[i] = result.y / divisor;
        mz[i] = result.z / divisor;
    }
}
