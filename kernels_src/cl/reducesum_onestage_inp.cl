// In-place reduction sum
// Ideally, the number of work-items per work-group is the smallest
// power of 2 that is larger than nbatch
__kernel void
reducesum_onestage_inp(__global real_t* __restrict      src,
                       __global real_t* __restrict      dst,
                                real_t              initVal,
                                   int                    n,
                                   int               nbatch,
                                  bool            overwrite,
                        __local real_t*            scratch1) {

    // Calculate indices
    // Each work-group to process 2*nbatch number of items
    // Work-group will iterate in for loop to emulate other groups
    int  local_idx = get_local_id(0);   // Work-item index within workgroup
    int     grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int grp_stride = get_num_groups(0); // Update group id at every iteration
    int   base_idx = get_group_id(0);   // Use the actual work-group id
    int    itr_cnt = 0;                 // Track the number of iterations of outer loop

    // grp_id: Base index of workgroup that gets updated every iteration
    for (int grp_id = get_group_id(0); grp_id*(nbatch << 1) < n; grp_id += grp_stride) {
        int global_idx = grp_id * (nbatch << 1) + local_idx; // Calculate global index of work-item
        // Grab data from global memory
        scratch1[local_idx] = 0.0; // Unsure scratch is zeroed at the beginning
        if ((global_idx < n) && (local_idx < nbatch)) { // Execute only if work-item is valid
            if (global_idx + nbatch < n) { // If work-item has two valid inputs
                scratch1[local_idx] = src[global_idx] + src[global_idx + nbatch];
                if (overwrite) {
                    src[global_idx] = 0.0;
                    src[global_idx + nbatch] = 0.0;
                }
            } else { // If work-item has only one valid input
                scratch1[local_idx] = src[global_idx];
                if (overwrite) {
                    src[global_idx] = 0.0;
                }
            }
        }

        // Barrier synchronization
        barrier(CLK_LOCAL_MEM_FENCE);

        // Reduce items in local memory
        for (int i0 = grp_sz >> 1; i0 > 0; i0 >>= 1) {
            if (local_idx < i0) {
                scratch1[local_idx] += scratch1[local_idx + i0];
            }

            // Barrier synchronization
            barrier(CLK_LOCAL_MEM_FENCE);
        }

        // Output to global buffer in location previous processed by
        // the current work-group
        if (local_idx == 0) {
            dst[base_idx + itr_cnt] = scratch1[0] + scratch1[1];
            itr_cnt += 1;
            if (itr_cnt == nbatch) {
                base_idx += (nbatch << 1) * grp_stride;
                itr_cnt = 0;
            }
        }
    }
}
