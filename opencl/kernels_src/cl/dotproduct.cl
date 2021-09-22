// dst += prefactor * dot(a,b)
__kernel void
dotproduct(__global float* __restrict  dst, float  prefactor,
           __global float* __restrict  ax, __global float* __restrict  ay, __global float* __restrict  az,
           __global float* __restrict  bx, __global float* __restrict  by, __global float* __restrict  bz,
           int  N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (i < N) {
		float3 A = {ax[i], ay[i], az[i]};
		float3 B = {bx[i], by[i], bz[i]};
		dst[i] += prefactor * dot(A, B);
	}
}

