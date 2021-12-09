// Landau-Lifshitz torque without precession
__kernel void
llnoprecess(__global real_t* __restrict tx, __global real_t* __restrict ty, __global real_t* __restrict tz,
            __global real_t* __restrict mx, __global real_t* __restrict my, __global real_t* __restrict mz,
            __global real_t* __restrict hx, __global real_t* __restrict hy, __global real_t* __restrict hz,
                                    int  N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t3 m = {mx[i], my[i], mz[i]};
        real_t3 H = {hx[i], hy[i], hz[i]};

        real_t3    mxH = cross(m, H);
        real_t3 torque = -cross(m, mxH);

        tx[i] = torque.x;
        ty[i] = torque.y;
        tz[i] = torque.z;
    }
}
