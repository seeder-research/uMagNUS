__kernel void
reducesum_twophase(__global real_t* __restrict      src,
                   __global real_t* __restrict      dst,
                            real_t              initVal,
                               int                    n,
                               int               nbatch,
                    __local real_t*            scratch1) {

    // Calculate indices
    // Each work-group to process 2*nbatch*grp_sz number of items
    int  local_idx = get_local_id(0);                          // Work-item index within workgroup
    int     grp_sz = get_local_size(0);                        // Total number of work-items in each workgroup
    int     grp_id = get_group_id(0);                          // Index of workgroup
    int global_idx = grp_id * 2 * nbatch * grp_sz + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * 2 * nbatch * grp_sz;  // Offset for memory access

    // Initialize an intermediate group sum
    real_t grpSum[4];
    grpSum[0] = 0.0; grpSum[1] = 0.0; grpSum[2] = 0.0; grpSum[3] = 0.0;

    // use flag to terminate nested loop early
    bool termFlag = false;

    // use loop to emulate separate workgroups
    for (int i0 = 0; i0 < 4; i0++) {
        for (int i1 = 0; i1 < grp_sz; i1++) {
            if (global_idx - local_idx >= n) { // early termination
                termFlag = true;
                break;
            }
            // Grab data from global memory
            scratch1[local_idx] = 0.0;
            if ((global_idx < n) && (local_idx < nbatch)) {
                if (global_idx + nbatch < n) {
                    scratch1[local_idx] = src[global_idx] + src[global_idx + nbatch];
                } else {
                    scratch1[local_idx] = src[global_idx];
                }
            }

            // Barrier synchronization
            barrier(CLK_LOCAL_MEM_FENCE);

            // Reduce items in local memory
            for (int ii = (grp_sz >> 1); ii > 0; ii >>= 1) {
               if (local_idx < ii) {
                    scratch1[local_idx] += scratch1[local_idx + ii];
                }
                // Barrier synchronization
                barrier(CLK_LOCAL_MEM_FENCE);
            }

            // Store to register in a work-item
            if (local_idx == i1) {
                grpSum[i0] = scratch1[0] + scratch1[1];
            }
        }
        if (termFlag) {
            break;
        }
    }

    // store group partial sums into local memory before reduction
    grpSum[0] += grpSum[2];
    grpSum[1] += grpSum[3];
    scratch1[local_idx] = grpSum[0] + grpSum[1];

    // Barrier synchronization
    barrier(CLK_LOCAL_MEM_FENCE);

    // begin to merge partial sums
    for (int ii = (grp_sz >> 1); ii > 0; ii >>= 1) {
        if (local_idx < ii) {
            scratch1[local_idx] += scratch1[local_idx + ii];
        }
        // Barrier synchronization
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Store to global memory
    if (local_idx == 0) {
        dst[grp_id] = scratch1[0] + scratch1[1];
    }
}
