/*
Alternative reducesum kernels

__kernel void
reducesum_onestage(__global real_t* __restrict      src,
                   __global real_t* __restrict      dst,
                            real_t              initVal,
                               int                    n,
                               int               nbatch,
                    __local real_t*            scratch1) {

    // Calculate indices
    // Each work-group to process 2*nbatch number of items
    int  local_idx = get_local_id(0);                 // Work-item index within workgroup
    int     grp_sz = get_local_size(0);               // Total number of work-items in each workgroup
    int     grp_id = get_group_id(0);                 // Index of workgroup
    int global_idx = grp_id * 2 * nbatch + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * 2 * nbatch;  // Offset for memory access

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
    for (int ii = grp_sz >> 1; ii > 0; ii >>= 1) {
        if (local_idx < ii) {
            scratch1[local_idx] += scratch1[local_idx + ii];
        }
        // Barrier synchronization
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Output to separate global buffer
    if (local_idx == 0) {
        dst[grp_idx] = scratch1[0] + scratch1[1];
    }
}

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

__kernel void
reducesum_twostages(__global real_t* __restrict      src, __global real_t* __restrict      dst, real_t initVal, int n, int nbatch
                      __local real_t* scratch1,             __local real_t* scratch2) {
    // Calculate indices
    int  local_idx = get_local_id(0); // Work-item index within workgroup
    int     grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int     grp_id = get_group_id(0); // Index of workgroup
    int global_idx = grp_id * nbatch + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * nbatch; // Offset for memory access

    // Initialize an intermediate group sum
    real_t grpSum = 0.0f;

    // loop for every work-group
    for (int cnt = 0; cnt < grp_sz; cnt++) {
        // Grab data from global memory
        if ((global_idx < n) && (local_idx < nbatch)) {
            if (global_idx + grp_offset < n) {
                scratch1[local_idx] = src[global_idx] + src[global_idx + grp_offset];
            } else {
                scratch1[local_idx] = src[global_idx];
            }
        } else {
            scratch1[local_idx] = 0.0;
        }

        // Barrier synchronization
        barrier(CLK_LOCAL_MEM_FENCE);

        // Reduce items in local memory
        for (int ii = grp_sz >> 1; ii > 1; ii >>= 1) {
            if (local_idx < ii) {
                scratch1[local_idx] += scratch1[local_idx + ii];
            }

            // Barrier synchronization
            barrier(CLK_LOCAL_MEM_FENCE);
        }

        // accumulate group sum in work-item specific register
        // i.e., distribute storage of first stage results
        if (local_idx == cnt) {
            grpSum = scratch1[0] + scratch1[1];
        }

        global_idx += grp_offset + grp_offset;

        // Barrier synchronization
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // reduce results of distributed storage of group sums...
    scratch1[local_idx] = grpSum;

    // Barrier synchronization
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce items in local memory
    for (int ii = grp_sz >> 1; ii > 1; ii >>= 1) {
        if (local_idx < ii) {
            scratch1[local_idx] += scratch1[local_idx + ii];
        }
        // Barrier synchronization
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    // Output sum after 2 stages to separate global buffer
    if (local_idx == 0) {
        dst[grp_idx] = scratch1[0] + scratch1[1];
    }
}
*/
__kernel void
reducesum(__global real_t* __restrict      src, __global real_t* __restrict      dst, real_t initVal, int n,
                      __local real_t* scratch1,             __local real_t* scratch2) {
    // Calculate indices
    int  local_idx = get_local_id(0); // Work-item index within workgroup
    int     grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int     grp_id = get_group_id(0); // Index of workgroup
    int global_idx = grp_id * grp_sz + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * grp_sz; // Offset for memory access

    // Initialize registers for work-item
    real_t grpSum = 0.0f; // Accumulator for workgroup
    real_t grpErr = 0.0f; // Error for workgroup accumulator
    real_t   aVal = 0.0f; // Register to track operand A
    real_t   bVal = 0.0f; // Register to track operand B
    real_t   lsum = 0.0f; // Register to temporarily store A + B
    real_t   lerr = 0.0f; // Register to temporarily store error from A + B
    real_t  lerr2 = 0.0f;
    real_t2 tmpR0 = 0.0f; // Temporary register
    real_t2 tmpR1 = 0.0f; // Temporary register
    real_t2 tmpR2 = 0.0f; // Temporary register
    real_t2 tmpR3 = 0.0f; // Temporary register
    real_t2 tmpR4 = 0.0f; // Temporary register
    
    // Set the accumulator value to initVal for the first work-item only
    if (global_idx == 0) {
        grpSum = initVal;
    }

/*
    // During each loop iteration, we:
    // 1) load source data into scratch1. If global index exceeds the position, then we load 0
    // 2) perform a reduction sum over the values in scratch1
    // 3) Accumulate into accumulator in work-item with local_idx = 0
*/

    for (int ii = 0; ii < n; ii += grp_offset) {
        // Get source data and load into local memory
        scratch1[local_idx] = (global_idx < n) ? src[global_idx] : 0.0f ;
        scratch2[local_idx] = 0.0f;

        // Add barrier to sync all threads
        barrier(CLK_LOCAL_MEM_FENCE);

        // Reduce sum of data in local memory using divide and conquer strategy
        for (int offset = get_local_size(0) / 2; offset > 0; offset >>= 1) {
            if (local_idx < offset) {
                aVal = scratch1[local_idx]; // Load accumulator
                bVal = scratch1[local_idx + offset]; // Load accumulator
                lsum = aVal + bVal; // Temporary sum

                // Write sum back into workgroup scratch memory
                scratch1[local_idx] = lsum;

                // Calculate error in summing error
                tmpR0.x = lsum; tmpR0.y = lsum;
                tmpR1.x = aVal; tmpR1.y = bVal;
                tmpR2.y = aVal; tmpR2.x = bVal;

                tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
                tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
                lerr = tmpR4.x + tmpR4.y; // Combine the errors

                // Retrieve the errors from scratch memory
                aVal = scratch2[local_idx];
                bVal = scratch2[local_idx + offset];
                lsum = aVal + bVal;

                // Calculate error in summing error
                tmpR0.x = lsum; tmpR0.y = lsum;
                tmpR1.x = aVal; tmpR1.y = bVal;
                tmpR2.y = aVal; tmpR2.x = bVal;
                tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
                tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
                lerr2 = tmpR4.x + tmpR4.y; // Combine the errors

                aVal = lerr; bVal = lsum;
                lsum = aVal + bVal;

                // Calculate error in summing error
                tmpR0.x = lsum; tmpR0.y = lsum;
                tmpR1.x = aVal; tmpR1.y = bVal;
                tmpR2.y = aVal; tmpR2.x = bVal;
                tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
                tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
                lerr = tmpR4.x + tmpR4.y; // Combine the errors
                lerr -= lerr2;

                scratch2[local_idx] = lerr;
            }
            // barrier for syncing workgroup
            barrier(CLK_LOCAL_MEM_FENCE);
        }

        // barrier for syncing workgroup
        barrier(CLK_LOCAL_MEM_FENCE);
        if (local_idx == 0) {
            aVal = scratch1[0]; bVal = grpSum;
            lsum = aVal + bVal;

            // Calculate error in summing error
            tmpR0.x = lsum; tmpR0.y = lsum;
            tmpR1.x = aVal; tmpR1.y = bVal;
            tmpR2.y = aVal; tmpR2.x = bVal;

            tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
            tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
            lerr = tmpR4.x + tmpR4.y; // Combine the errors

            grpSum = lsum;

            aVal = scratch2[0]; bVal = grpErr;
            lsum = aVal + bVal;
            // Calculate error in summing error
            tmpR0.x = lsum; tmpR0.y = lsum;
            tmpR1.x = aVal; tmpR1.y = bVal;
            tmpR2.y = aVal; tmpR2.x = bVal;

            tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
            tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
            lerr2 = tmpR4.x + tmpR4.y; // Combine the errors

            aVal = lerr; bVal = lerr2;
            lsum = aVal + bVal;
            // Calculate error in summing error
            tmpR0.x = lsum; tmpR0.y = lsum;
            tmpR1.x = aVal; tmpR1.y = bVal;
            tmpR2.y = aVal; tmpR2.x = bVal;

            tmpR3 = tmpR0 - tmpR1; // Calculate the operands of the sum from temporary sum
            tmpR4 = tmpR3 - tmpR2; // Calculate error between calculated operands and actual operands
            lerr = tmpR4.x + tmpR4.y; // Combine the errors
            grpErr = lsum - lerr;
        }
        global_idx += grp_offset;
    }

    if (local_idx == 0) {
        dst[grp_id] = grpSum - grpErr;
    }
}
