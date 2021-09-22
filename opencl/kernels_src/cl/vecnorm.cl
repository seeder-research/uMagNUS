__kernel void
vecnorm(__global float* __restrict dst,
        __global float* __restrict ax, __global float* __restrict ay, __global float* __restrict az,
        int N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (i < N) {
		float3 A = {ax[i], ay[i], az[i]};
		dst[i] = sqrt(dot(A, A));
	}
}

