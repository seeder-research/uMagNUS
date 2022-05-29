__kernel void
reducesum(__global real_t* __restrict     src,
          __global real_t* __restrict     dst,
                   real_t             initVal,
                      int                   n,
          __local  real_t*            scratch){

    // Calculate indices
    int    local_idx = get_local_id(0);   // Work-item index within workgroup
    int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    real_t       res = init_val;

    // Accumulators for intermediate results
    real_t data1 = 0.0;
    real_t data2 = 0.0;
    real_t data3 = 0.0;
    real_t data4 = 0.0;

    // Indices for accumulators for intermediate results
    unsigned int id1 = 0;
    unsigned int id2 = 0;
    unsigned int id3 = 0;
    unsigned int id4 = 0;

    for (int base_gid = 0; base_gid < n; base_gid += grp_sz) {
        // Load data from global buffer to local buffer
        scratch[local_idx] = 0.0;
        int global_idx = base_gid + local_idx;
        if (global_idx < n) {
            scratch[local_idx] = src[global_idx];
        }

        // Synchronize workgroup before reduction in local buffer
        barrier(CLK_LOCAL_MEM_FENCE);

        // Reduce in local buffer
        for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
            if (local_idx < s) {
                scratch[local_idx] += scratch[local_idx + s];
            }

            // Synchronize workgroup before next iteration
            barrier(CLK_LOCAL_MEM_FENCE);
        }

        // Accumulate intermediate result in register of
        // corresponding workitem
        if (local_idx == id1) {
            data1 = scratch[0] + scratch[1];
        }

        // Synchronize workgroup before continuing
        barrier(CLK_LOCAL_MEM_FENCE);

        // Increment index of first stage, and reduce if needed,
        // accumulating result in second stage
        id1++;
        if (id1 >= grp_sz) {
            // Reset index of first stage
            id1 = 0;

            // Reduce intermediate results of first stage and
            // accumulate in second stage
            scratch[local_idx] = data1;

            // Synchronize workgroup before reduction in local buffer
            barrier(CLK_LOCAL_MEM_FENCE);

            // Reduce in local buffer
            for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
                if (local_idx < s) {
                    scratch[local_idx] += scratch[local_idx + s];
                }

                // Synchronize workgroup before next iteration
                barrier(CLK_LOCAL_MEM_FENCE);
            }

            // Accumulate intermediate result in register of
            // corresponding workitem (stage 2)
            if (local_idx == id2) {
                data2 = scratch[0] + scratch[1];
            }

            // Synchronize workgroup before continuing
            barrier(CLK_LOCAL_MEM_FENCE);

            // Increment index of second stage, and reduce if needed,
            // accumulating result in third stage
            id2++;
            if (id2 >= grp_sz) {
                // Reset index of second stage
                id2 = 0;

                // Reduce intermediate results of second stage and
                // accumulate in third stage
                scratch[local_idx] = data2;

                // Synchronize workgroup before reduction in local buffer
                barrier(CLK_LOCAL_MEM_FENCE);

                // Reduce in local buffer
                for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
                    if (local_idx < s) {
                        scratch[local_idx] += scratch[local_idx + s];
                    }

                    // Synchronize workgroup before next iteration
                    barrier(CLK_LOCAL_MEM_FENCE);
                }

                // Accumulate intermediate result in register of
                // corresponding workitem (stage 3)
                if (local_idx == id3) {
                    data3 = scratch[0] + scratch[1];
                }

                // Synchronize workgroup before continuing
                barrier(CLK_LOCAL_MEM_FENCE);

                // Increment index of third stage, and reduce if needed,
                // accumulating result in fourth stage
                id3++;
                if (id3 >= grp_sz) {
                    // Reset index of third stage
                    id3 = 0;

                    // Reduce intermediate results of third stage and
                    // accumulate in fourth stage
                    scratch[local_idx] = data3;

                    // Synchronize workgroup before reduction in local buffer
                    barrier(CLK_LOCAL_MEM_FENCE);

                    // Reduce in local buffer
                    for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
                        if (local_idx < s) {
                            scratch[local_idx] += scratch[local_idx + s];
                        }

                        // Synchronize workgroup before next iteration
                        barrier(CLK_LOCAL_MEM_FENCE);
                    }

                    // Accumulate intermediate result in register of
                    // corresponding workitem (stage 4)
                    if (local_idx == id4) {
                        data4 = scratch[0] + scratch[1];
                    }

                    // Synchronize workgroup before continuing
                    barrier(CLK_LOCAL_MEM_FENCE);

                    // Increment index of fourth stage, and reduce if needed,
                    // accumulating result in final register
                    id4++;
                    if (id4 >= grp_sz) {
                        // Reset index of fourth stage
                        id4 = 0;

                        // Reduce intermediate results of fourth stage and
                        // accumulate in final register
                        scratch[local_idx] = data4;

                        // Synchronize workgroup before reduction in local buffer
                        barrier(CLK_LOCAL_MEM_FENCE);

                        // Reduce in local buffer
                        for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
                            if (local_idx < s) {
                                scratch[local_idx] += scratch[local_idx + s];
                            }

                            // Synchronize workgroup before next iteration
                            barrier(CLK_LOCAL_MEM_FENCE);
                        }

                        // Accumulate intermediate result in register of
                        // corresponding workitem (final stage)
                        if (local_idx == id4) {
                            res += scratch[0] + scratch[1];
                        }

                        // Synchronize workgroup before continuing
                        barrier(CLK_LOCAL_MEM_FENCE);

                        // Reset registers of fourth stage
                        data4 = 0.0;
                    }

                    // Reset registers of third stage
                    data3 = 0.0;
                }

                // Reset registers of second stage
                data2 = 0.0;
            }

            // Reset registers of first stage
            data1 = 0.0;
        }
    }

    // All elements in global buffer have been picked up
    // Reduce intermediate results in each stage and accumulate
    // in final stage

    // Synchronize workgroup before reduction in local buffer
    barrier(CLK_LOCAL_MEM_FENCE);

    // Stage 1
    scratch[local_idx] = data1;

    // Synchronize workgroup before reduction in local buffer
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce in local buffer
    for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
        if (local_idx < s) {
            scratch[local_idx] += scratch[local_idx + s];
        }

        // Synchronize workgroup before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Accumulate intermediate result in register of
    // corresponding workitem (stage 2)
    if (local_idx == id2) {
        data2 = scratch[0] + scratch[1];
    }

    // Synchronize workgroup before continuing
    barrier(CLK_LOCAL_MEM_FENCE);

    // Stage 2
    scratch[local_idx] = data2;

    // Synchronize workgroup before reduction in local buffer
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce in local buffer
    for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
        if (local_idx < s) {
            scratch[local_idx] += scratch[local_idx + s];
        }

        // Synchronize workgroup before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Accumulate intermediate result in register of
    // corresponding workitem (stage 3)
    if (local_idx == id3) {
        data3 = scratch[0] + scratch[1];
    }

    // Synchronize workgroup before continuing
    barrier(CLK_LOCAL_MEM_FENCE);

    // Stage 3
    scratch[local_idx] = data3;

    // Synchronize workgroup before reduction in local buffer
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce in local buffer
    for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
        if (local_idx < s) {
            scratch[local_idx] += scratch[local_idx + s];
        }

        // Synchronize workgroup before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Accumulate intermediate result in register of
    // corresponding workitem (stage 4)
    if (local_idx == id4) {
        data4 = scratch[0] + scratch[1];
    }

    // Synchronize workgroup before continuing
    barrier(CLK_LOCAL_MEM_FENCE);

    // Stage 4
    scratch[local_idx] = data4;

    // Synchronize workgroup before reduction in local buffer
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce in local buffer
    for (unsigned int s = (grp_sz >> 1); s > 1; s >>= 1 ) {
        if (local_idx < s) {
            scratch[local_idx] += scratch[local_idx + s];
        }

        // Synchronize workgroup before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Accumulate intermediate result in register of
    // corresponding workitem (final stage)
    if (local_idx == 0) {
        res += scratch[0] + scratch[1];
        dst[0] = res;
    }

}
