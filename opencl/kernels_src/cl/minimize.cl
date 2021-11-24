// Steepest descent energy minimizer
__kernel void
minimize(__global float* __restrict  mx, __global float* __restrict  my, __global float* __restrict  mz,
         __global float* __restrict m0x, __global float* __restrict m0y, __global float* __restrict m0z,
         __global float* __restrict  tx, __global float* __restrict  ty, __global float* __restrict  tz,
                              float  dt,                        int   N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        float3 m0 = {m0x[i], m0y[i], m0z[i]};
        float3  t = {tx[i], ty[i], tz[i]};

        float       t2 = dt*dt*dot(t, t);
        float3  result = (4.0f - t2) * m0 + 4.0f * dt * t;
        float  divisor = 4.0f + t2;

        mx[i] = result.x / divisor;
        my[i] = result.y / divisor;
        mz[i] = result.z / divisor;
    }
}
