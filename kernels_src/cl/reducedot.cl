__kernel void
reducedot(         __global real_t* __restrict       x1, __global real_t* __restrict       x2,
          volatile __global real_t* __restrict      dst,
                                        real_t  initVal,
                                           int        n,
                               __local real_t* scratch1,             __local real_t* scratch2) {

    // Initialize indices
    int  local_idx = get_local_id(0);
    int    grp_idx = get_group_id(0);
    int grp_offset = get_local_size(0);
    int global_idx = grp_idx * grp_offset + local_idx;
    grp_offset *= get_num_groups(0);

    // Initialize memory
    real_t currVal = 0.0f;
    real_t currErr = 0.0f;
    real_t   tmpR0 = 0.0f;
    real_t   tmpR1 = 0.0f;
    real_t   tmpR2 = 0.0f;
    real_t   tmpR3 = 0.0f;
    real_t   tmpR4 = 0.0f;
    
    // Set the accumulator value to initVal for the first work-item only
    if (global_idx == 0) {
        currVal = initVal;
    }

    // Loop over input elements in chunks and accumulate each chunk into local memory
    while (global_idx < n) {
        tmpR0 = x1[global_idx];
        tmpR1 = x2[global_idx];
        tmpR2 = fma(tmpR0, tmpR1, (real_t)0.0);
        tmpR3 = currVal + tmpR2;
        tmpR0 = tmpR3 - currVal;
        tmpR1 = tmpR3 - tmpR2;
        tmpR4 = tmpR0 - tmpR2;
        tmpR0 = tmpR1 - currVal;
        currVal = tmpR3;
        currErr += tmpR4 + tmpR0;
        global_idx += grp_offset;
    }

    // At this point, accumulated values on chunks are in local memory. Perform parallel reduction
    scratch1[local_idx] = currVal;
    scratch2[local_idx] = currErr;

    // Add barrier to sync all threads
    barrier(CLK_LOCAL_MEM_FENCE);
    for (int offset = get_local_size(0) / 2; offset > 0; offset = offset / 2) {
        if (local_idx < offset) {
            tmpR0 = scratch1[local_idx];
            tmpR1 = scratch1[local_idx + offset];
            currErr = scratch2[local_idx] + scratch2[local_idx + offset];
            currVal = tmpR0 + tmpR1;
            tmpR3 = currVal - tmpR0;
            tmpR4 = currVal - tmpR1;
            tmpR2 = tmpR3 - tmpR1;
            tmpR3 = tmpR4 - tmpR1;
            currErr += tmpR2 + tmpR3;
            scratch1[local_idx] = currVal;
            scratch2[local_idx] = currErr;
        }
        // barrier for syncing work group
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    if (local_idx == 0) {
        dst[grp_idx] = scratch1[0];
    }
}
