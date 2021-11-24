__kernel void
crossproduct(__global float* __restrict dstx, __global float* __restrict dsty, __global float* __restrict dstz,
             __global float* __restrict   ax, __global float* __restrict   ay, __global float* __restrict   az,
             __global float* __restrict   bx, __global float* __restrict   by, __global float* __restrict   bz,
                                    int    N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);
    for (int i = gid; i < N; i += gsize) {
        float3   A = {ax[i], ay[i], az[i]};
        float3   B = {bx[i], by[i], bz[i]};
        float3 AxB = cross(A, B);

        dstx[i] = AxB.x;
        dsty[i] = AxB.y;
        dstz[i] = AxB.z;
    }
}
