__kernel void
reducedot(__global real_t* __restrict     src1,
          __global real_t* __restrict     src2,
          __global real_t* __restrict      dst,
                   real_t              initVal,
                      int                    n,
          __local  real_t*            scratch1){

    // Calculate indices
    int  local_idx = get_local_id(0); // Work-item index within workgroup
    int     grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int grp_offset = grp_sz * grp_sz; // Offset for memory access

    // loop through groups
    for (int grp_id = get_group_id(0); grp_id < grp_sz; grp_id += get_num_groups(0)) {
        // Early termination if work group is noop
        int global_idx = grp_id * grp_sz;
        if (global_idx >= n) {
            break;
	}
        global_idx += local_idx; // Calculate global index of work-item

        // Use 8 local resisters to track work-item sum to reduce truncation errors
        real_t mine[8] = {0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0};
        uint itr = 0;
        while (global_idx < n) {
            itr = itr & 0x00000007;
            mine[itr] = fma(src1[global_idx], src2[global_idx], mine[itr]);
            global_idx += grp_offset;
            itr++;
        }

        // Merge work-item sums
        mine[0] += mine[4];
        mine[1] += mine[5];
        mine[2] += mine[6];
        mine[3] += mine[7];
        mine[0] += mine[2];
        mine[1] += mine[3];

        // Load work-item sums into local shared memory
        scratch1[local_idx] = mine[0] + mine[1];

        // Synchronize work-group
        barrier(CLK_LOCAL_MEM_FENCE);

        for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
            if (local_idx < s) {
                scratch1[local_idx] += scratch1[local_idx + s];
            }

            // Synchronize work-group
            barrier(CLK_LOCAL_MEM_FENCE);
        }
        if (local_idx == 0) {
            dst[grp_id] = (scratch1[0] + scratch1[1]) + initVal;
        }
    }
}
