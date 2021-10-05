__kernel void
vecnorm(__global float* __restrict dst,
        __global float* __restrict ax, __global float* __restrict ay, __global float* __restrict az,
        int N) {

    int gid = get_global_id(0);
    int gsize = get_global_size(0);
    for (int i = gid; i < N; i += gsize) {
        float3 A = {ax[i], ay[i], az[i]};
        dst[i] = sqrt(dot(A, A));
    }
}
