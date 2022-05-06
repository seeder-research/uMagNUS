__kernel void
reducesum(__global real_t* __restrict      src,
          __global real_t* __restrict      dst,
                   real_t              initVal,
                      int                    n,
          __local  real_t*            scratch1){

    // Calculate indices
    int  local_idx = get_local_id(0);   // Work-item index within workgroup
    int     grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int    grp_cnt = grp_sz << 4;       // Maximum number of workgroups emulated
    int grp_offset = grp_sz;            // Offset for memory access (if sole workgroup)
    int      nGrps = get_num_groups(0); // Total number of workgroups launched

    // If this is not the final stage reduction, need to
    // change the stride for memory accesses
    if (nGrps > 1) {
        grp_offset *= grp_cnt;
    }

    // Loop through groups
    for (int grp_id = get_group_id(0); grp_id < grp_cnt; grp_id += nGrps) {
        // Early termination if work-group is noop
        int global_idx = grp_id * grp_sz; // Calculate global_idx for work-item 0 of group
        if (global_idx >= n) { // Entire work-group is noop
            break;
        }

        global_idx += local_idx; // Calculate global index of work-item

        // Use 8 local resisters to track work-item sum to reduce truncation errors
        real_t4 data1 = {0.0, 0.0, 0.0, 0.0};
        real_t4 data2 = {0.0, 0.0, 0.0, 0.0};
        while (global_idx < n) {
            data1.x += src[global_idx];
            global_idx += grp_offset;
            if (global_idx >= n) {
                break;
            }
            data1.y += src[global_idx];
            global_idx += grp_offset;
            if (global_idx >= n) {
                break;
            }
            data1.z += src[global_idx];
            global_idx += grp_offset;
            if (global_idx >= n) {
                break;
            }
            data1.w += src[global_idx];
            global_idx += grp_offset;
            if (global_idx >= n) {
                break;
            }
            data2.x += src[global_idx];
            global_idx += grp_offset;
            if (global_idx >= n) {
                break;
            }
            data2.y += src[global_idx];
            global_idx += grp_offset;
            if (global_idx >= n) {
                break;
            }
            data2.z += src[global_idx];
            global_idx += grp_offset;
            if (global_idx >= n) {
                break;
            }
            data2.w += src[global_idx];
            global_idx += grp_offset;
        }

        // Merge work-item partial sums
        data1 += data2;
        data1.x += data1.z;
        data1.y += data1.w;

        // Load work-item sums into local shared memory
        scratch1[local_idx] = data1.x + data1.y;

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
