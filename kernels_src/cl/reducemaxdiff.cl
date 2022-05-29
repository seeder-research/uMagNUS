__kernel void
reducemaxdiff(__global real_t* __restrict    src1,
              __global real_t* __restrict    src2,
              __global real_t* __restrict     dst,
                       real_t             initVal,
                          int                   n,
              __local  real_t*            scratch) {

    // Calculate indices
    int    local_idx = get_local_id(0);   // Work-item index within workgroup
    int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    real_t       res = 0.0;

    for (int idx_base = 0; idx_base < n; idx_base += grp_sz) {
        int global_idx = idx_base + local_idx;
        scratch[local_idx] = 0.0;
        if (global_idx < n) {
            scratch[local_idx] = fabs(src1[global_idx] - src2[global_idx]);
        }

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

        // Store reduction result for each iteration and move to next
        if (local_idx == 0) {
            res = fmax(scratch[0], scratch[1]);
        }

    }

    // Store reduction result for each iteration and move to next
    if (local_idx == 0) {
        dst[0] = fmax(res, initVal);
    }

}
