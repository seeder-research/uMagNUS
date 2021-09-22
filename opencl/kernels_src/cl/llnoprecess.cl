// Landau-Lifshitz torque without precession
__kernel void
llnoprecess(__global float* __restrict  tx, __global float* __restrict  ty, __global float* __restrict  tz,
            __global float* __restrict  mx, __global float* __restrict  my, __global float* __restrict  mz,
            __global float* __restrict  hx, __global float* __restrict  hy, __global float* __restrict  hz, int N) {

	int i = ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (i < N) {

		float3 m = {mx[i], my[i], mz[i]};
		float3 H = {hx[i], hy[i], hz[i]};

		float3 mxH = cross(m, H);
		float3 torque = -cross(m, mxH);

		tx[i] = torque.x;
		ty[i] = torque.y;
		tz[i] = torque.z;
	}
}

