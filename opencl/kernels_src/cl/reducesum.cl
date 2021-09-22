__kernel void
reducesum(__global float* __restrict src, __global float* __restrict dst, float initVal, int n, __local float* scratch1, __local float* scratch2) {
	// Calculate indices
	int local_idx = get_local_id(0); // Work-item index within workgroup
	int grp_sz = get_local_size(0); // Total number of work-items in each workgroup
	int grp_id = get_group_id(0); // Index of workgroup
	int global_idx = grp_id * grp_sz + local_idx; // Calculate global index of work-item
	int grp_offset = get_num_groups(0) * grp_sz; // Offset for memory access

	// Initialize registers for work-item
	float grpSum = 0; // Accumulator for workgroup
	float grpErr = 0; // Error for workgroup accumulator
	float aVal = 0; // Register to track operand A
	float bVal = 0; // Register to track operand B
	float lsum = 0; // Register to temporarily store A + B
	float lerr = 0; // Register to temporarily store error from A + B
	float lerr2 = 0;
	float2 tmpR0 = 0; // Temporary register
	float2 tmpR1 = 0; // Temporary register
	float2 tmpR2 = 0; // Temporary register
	float2 tmpR3 = 0; // Temporary register
	float2 tmpR4 = 0; // Temporary register
	
	// Set the accumulator value to initVal for the first work-item only
	if (global_idx == 0) {
		grpSum = initVal;
	}

/*
	// During each loop iteration, we:
	// 1) load source data into scratch1. If global index exceeds the position, then we load 0
	// 2) perform a reduction sum over the values in scratch1
	// 3) Accumulate into accumulator in work-item with local_idx = 0
*/

	for (int ii = 0; ii < n; ii += grp_offset) {
		// Get source data and load into local memory
		scratch1[local_idx] = (global_idx < n) ? src[global_idx] : 0.0f ;
		scratch2[local_idx] = 0.0f;

		// Add barrier to sync all threads
		barrier(CLK_LOCAL_MEM_FENCE);

		// Reduce sum of data in local memory using divide and conquer strategy
		for (int offset = get_local_size(0) / 2; offset > 0; offset = offset / 2) {
			if (local_idx < offset) {
				aVal = scratch1[local_idx]; // Load accumulator
				bVal = scratch1[local_idx + offset]; // Load accumulator
				lsum = aVal + bVal; // Temporary sum

				// Write sum back into workgroup scratch memory
				scratch1[local_idx] = lsum;

				// Calculate error in summing error
				tmpR0.x = lsum; tmpR0.y = lsum;
				tmpR1.x = aVal; tmpR1.y = bVal;
				tmpR2.y = aVal; tmpR2.x = bVal;

				tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
				tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
				lerr = tmpR4.x + tmpR4.y; // Combine the errors

				// Retrieve the errors from scratch memory
				aVal = scratch2[local_idx];
				bVal = scratch2[local_idx + offset];
				lsum = aVal + bVal;

				// Calculate error in summing error
				tmpR0.x = lsum; tmpR0.y = lsum;
				tmpR1.x = aVal; tmpR1.y = bVal;
				tmpR2.y = aVal; tmpR2.x = bVal;
				tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
				tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
				lerr2 = tmpR4.x + tmpR4.y; // Combine the errors

				aVal = lerr; bVal = lsum;
				lsum = aVal + bVal;

				// Calculate error in summing error
				tmpR0.x = lsum; tmpR0.y = lsum;
				tmpR1.x = aVal; tmpR1.y = bVal;
				tmpR2.y = aVal; tmpR2.x = bVal;
				tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
				tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
				lerr = tmpR4.x + tmpR4.y; // Combine the errors
				lerr -= lerr2;

				scratch2[local_idx] = lerr;
			}
			// barrier for syncing workgroup
			barrier(CLK_LOCAL_MEM_FENCE);
		}

		// barrier for syncing workgroup
		barrier(CLK_LOCAL_MEM_FENCE);
		if (local_idx == 0) {
			aVal = scratch1[0]; bVal = grpSum;
			lsum = aVal + bVal;

			// Calculate error in summing error
			tmpR0.x = lsum; tmpR0.y = lsum;
			tmpR1.x = aVal; tmpR1.y = bVal;
			tmpR2.y = aVal; tmpR2.x = bVal;

			tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
			tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
			lerr = tmpR4.x + tmpR4.y; // Combine the errors

			grpSum = lsum;

			aVal = scratch2[0]; bVal = grpErr;
			lsum = aVal + bVal;
			// Calculate error in summing error
			tmpR0.x = lsum; tmpR0.y = lsum;
			tmpR1.x = aVal; tmpR1.y = bVal;
			tmpR2.y = aVal; tmpR2.x = bVal;

			tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
			tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
			lerr2 = tmpR4.x + tmpR4.y; // Combine the errors

			aVal = lerr; bVal = lerr2;
			lsum = aVal + bVal;
			// Calculate error in summing error
			tmpR0.x = lsum; tmpR0.y = lsum;
			tmpR1.x = aVal; tmpR1.y = bVal;
			tmpR2.y = aVal; tmpR2.x = bVal;

			tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
			tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
			lerr = tmpR4.x + tmpR4.y; // Combine the errors
			grpErr = lsum - lerr;
		}
		global_idx += grp_offset;
	}

	if (local_idx == 0) {
		dst[grp_id] = grpSum - grpErr;
	}
}

