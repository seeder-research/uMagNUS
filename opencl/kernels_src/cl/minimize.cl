// Steepest descent energy minimizer
__kernel void
minimize(__global float* __restrict mx,  __global float* __restrict  my,  __global float* __restrict mz,
         __global float* __restrict m0x, __global float* __restrict  m0y, __global float* __restrict m0z,
         __global float* __restrict tx,  __global float* __restrict  ty,  __global float* __restrict tz,
         float dt, int N) {

	int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
	if (i < N) {

		float3 m0 = {m0x[i], m0y[i], m0z[i]};
		float3 t = {tx[i], ty[i], tz[i]};

		float t2 = dt*dt*dot(t, t);
		float3 result = (4 - t2) * m0 + 4 * dt * t;
		float divisor = 4 + t2;
		
		mx[i] = result.x / divisor;
		my[i] = result.y / divisor;
		mz[i] = result.z / divisor;
	}
}

