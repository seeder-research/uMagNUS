__kernel void
reducemaxdiff(         __global real_t* __restrict    src1,
                       __global real_t* __restrict    src2,
              volatile __global real_t* __restrict     dst,
                                real_t             initVal,
                                   int                   n,
              volatile __local  real_t*            scratch) {

    // Calculate indices
    unsigned int    local_idx = get_local_id(0);   // Work-item index within workgroup
    unsigned int       grp_id = get_group_id(0);   // ID of workgroup
    unsigned int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    unsigned int        grp_i = grp_id*grp_sz;
    unsigned int       stride = get_global_size(0);
    real_t               mine = initVal;

    while (grp_i < (unsigned int)(n)) {
        unsigned int i = grp_i + local_idx;
        if (i < (unsigned int)(n)) {
            mine = fmax(mine, fabs(src1[i] - src2[i]));
        }

        // Load workitem value into local buffer and synchronize
        scratch[local_idx] = mine;
        barrier(CLK_LOCAL_MEM_FENCE);

        // Reduce using lor loop
        for (unsigned int s = (grp_sz >> 1); s > 32; s >>= 1) {
            if (local_idx < s) {
                scratch[local_idx] = fmax(scratch[local_idx], scratch[local_idx + s]);
            }

            // Synchronize workitems before next iteration
            barrier(CLK_LOCAL_MEM_FENCE);
        }

        // Unroll for loop that executes within one unit that works on 32 workitems
        if (local_idx < 32) {
            volatile __local real_t* smem = scratch;
            smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 32]);
            smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 16]);
            smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  8]);
            smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  4]);
            smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  2]);
            smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  1]);
            mine = scratch[local_idx];
        }

        grp_i += stride;
    }

    // Store reduction result for each iteration and move to next
    if (local_idx == 0) {
        mine = fmax(scratch[0], scratch[1]);
        atomicMax_r(dst, mine);
//        dst[grp_id] = fmax(scratch[0], scratch[1]);
//        dst[grp_id] = mine;
    }

}
