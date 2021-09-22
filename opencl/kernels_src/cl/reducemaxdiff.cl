__kernel void
reducemaxdiff(__global float* __restrict src1, __global float* __restrict  src2, __global float* __restrict dst,
	      float initVal, int n, __local float* scratch) {
	// Initialize memory
	int global_idx =  get_global_id(0);
	int local_idx = get_local_id(0);
	float currVal = initVal;

	// Loop over input elements in chunks and accumulate each chunk into local memory
	while (global_idx < n) {
		float element = fabs(src1[global_idx] - src2[global_idx]);
		currVal = fmax(currVal, element);
		global_idx += get_global_size(0);
	}

	// At this point, accumulated values on chunks are in local memory. Perform parallel reduction
	scratch[local_idx] = currVal;
	// Add barrier to sync all threads
	barrier(CLK_LOCAL_MEM_FENCE);
	for (int offset = get_local_size(0) / 2; offset > 0; offset = offset / 2) {
		if (local_idx < offset) {
			float other = scratch[local_idx + offset];
			float mine = scratch[local_idx];
			scratch[local_idx] = fmax(mine, other);
		}
		// barrier for syncing work group
		barrier(CLK_LOCAL_MEM_FENCE);
	}

	if (local_idx == 0) {
		dst[get_group_id(0)] = scratch[0];
	}
}

