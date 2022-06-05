__kernel void
reducesum(         __global real_t*    __restrict     src,
          volatile __global real_t*    __restrict     dst,
                            real_t                initVal,
                               int                      n,
          volatile __local  real_t*               scratch){

    // Calculate indices
    unsigned int    local_idx = get_local_id(0);   // Work-item index within workgroup
    unsigned int       grp_id = get_group_id(0);   // ID of workgroup
    unsigned int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    unsigned int        grp_i = grp_id(0)*grp_sz;
    unsigned int       stride = get_global_size(0);

    // Initialize ring accumulator for intermediate results
    real_t accum = (real_t)(0.0);
    if (get_global_id(0) == 0) {
        accum = initVal;
    }
    unsigned int itr = 0;

    while (grp_i < (unsigned int)(n)) {
        unsigned int i = grp_i + local_idx;

        // Read from global memory and accumulate local memory
        scratch[local_idx] = (real_t)(0.0);
        if (i < n) {
                scratch[local_idx] = src[i];
        }

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

        if (local_idx == itr) {
            accum += scratch[0];
        }

        // Move pointer to ring accumulator
        itr++;
        if (itr >= 32) {
            itr = 0;
        }

        // Update pointer to next global value
        grp_i += stride;
    }

    // All elements in global buffer have been picked up
    // Reduce intermediate results and add atomically to global buffer
    if (local_idx < 32) {
        scratch[local_idx] = accum;

        // Unroll loop for remaining 32 workitems
        if (local_idx < 16) {
            volatile __local real_t* smem = scratch;
            smem[local_idx] += smem[local_idx + 16];
            smem[local_idx] += smem[local_idx +  8];
            smem[local_idx] += smem[local_idx +  4];
            smem[local_idx] += smem[local_idx +  2];
            smem[local_idx] += smem[local_idx +  1];
        }

        // Add atomically to global buffer
        if (local_idx == 0) {
    //        atomicAdd_r(dst, scratch[0]);
            dst[grp_id] = scratch[0];
        }
    }

}
