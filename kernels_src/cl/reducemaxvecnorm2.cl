__kernel void
reducemaxvecnorm2(__global real_t* __restrict       x,
                  __global real_t* __restrict       y,
                  __global real_t* __restrict       z,
                  __global real_t* __restrict     dst,
                           real_t             initVal,
                              int                   n,
                  __local  real_t*            scratch) {

    // Calculate indices
    int  local_idx = get_local_id(0);   // Work-item index within workgroup
    int     grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int    grp_cnt = grp_sz << 4;       // Maximum number of workgroups to emulate
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

        // Initialize value to track
        real_t currVal = initVal;

        while (global_idx < n) {
            real_t element = (x[global_idx]*x[global_idx]) + (y[global_idx]*y[global_idx]) + (z[global_idx]*z[global_idx]);
            currVal = fmax(currVal, element);
            global_idx += grp_offset;
        }

        // At this point, max values on chunks are in local memory. Perform parallel reduction
        scratch[local_idx] = currVal;

        // Add barrier to sync all threads
        barrier(CLK_LOCAL_MEM_FENCE);

        for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
            if (local_idx < s) {
                real_t other = scratch[local_idx + s];
                real_t  mine = scratch[local_idx];
                scratch[local_idx] = fmax(mine, other);
            }

            // Synchronize work-group
            barrier(CLK_LOCAL_MEM_FENCE);
        }

        if (local_idx == 0) {
            dst[grp_id] = fmax(scratch[0], scratch[1]);
        }
    }
}
