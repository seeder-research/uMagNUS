__kernel void
reducedot(         __global real_t* __restrict     src1,
                   __global real_t* __restrict     src2,
          volatile __global real_t* __restrict      dst,
                            real_t              initVal,
                               int                    n,
          volatile __local  real_t*            scratch){

    // Calculate indices
    int    local_idx = get_local_id(0);   // Work-item index within workgroup
    int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int            i = get_group_id(0)*grp_sz + local_idx;
    int       stride = get_global_size(0);
    // Initialize ring accumulator for intermediate results
    real_t accum[__REDUCE_REG_COUNT__];
    for (unsigned int s = 0; s < __REDUCE_REG_COUNT__; s++) {
        accum[s] = 0.0;
    }
    accum[0] = initVal;
    unsigned int itr = 0;

    // Read from global memory and accumulate in workitem ring accumulator
    while (i < n) {
        accum[itr] += src1[i] * src2[i]; // Load value from global buffer into ring accumulator

        // Update pointer to ring accumulator
        itr++;
        if (itr >= __REDUCE_REG_COUNT__) {
            itr = 0;
        }

        // Update pointer to next global value
        i += stride;
    }

    // All elements in global buffer have been picked up
    // Reduce intermediate results and add atomically to global buffer

    // Reduce value in ring buffer
    for (unsigned int s1 = (__REDUCE_REG_COUNT__ >> 1); s1 > 1; s1 >>= 1) {
        for (unsigned int s2 = 0; s2 < s1; s2++) {
            accum[s2] += accum[s2+s1];
        }
    }

    // Reduce in local buffer
    scratch[local_idx] = accum[0] + accum[1];

    // Synchronize workgroup before reduction in local buffer
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce in local buffer
    for (unsigned int s = (grp_sz >> 1); s > 32; s >>= 1 ) {
        if (local_idx < s) {
            scratch[local_idx] += scratch[local_idx + s];
        }

        // Synchronize workgroup before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Unroll loop for remaining 32 workitems
    if (local_idx < 32) {
        volatile __local real_t* smem = scratch;
        smem[local_idx] += smem[local_idx + 32];
        smem[local_idx] += smem[local_idx + 16];
        smem[local_idx] += smem[local_idx +  8];
        smem[local_idx] += smem[local_idx +  4];
        smem[local_idx] += smem[local_idx +  2];
        smem[local_idx] += smem[local_idx +  1];
    }

    // Add atomically to global buffer
    if (local_idx == 0) {
        dst[get_group_id(0)] = scratch[0];
    }

}
