// Uses 16 registers (if single-precision, then 64B per
// workitem. Twice if double-precision, i.e., 128B per
// workitem) to produce reduction sum in one thread. Allow
// running workgroup to have 32 or 64 workitems
// Single while loop emulating 8 nested for loops
// Single workgroup of 64 workitems can handle up to
// 1024x64x64 = 1024x1024x4 = 128x128x128x2 = 4194304 items
//
// Inputs:
//   src:      pointer to buffer with data to reduce
//   dst:      pointer to buffer where result is stored
//   initVal:  result of reduction will be offset by initVal
//   fac:      signal to tell scheme to write result to dst
//   group_n:  number of items to reduce in buffer
//   n:        total length of src buffer
//   scratch:  local memory for reduction
#ifdef(REDUCE_SUM_THREADS_PER_GROUP)
#define REDUCE_SUM_THREADS_PER_GROUP_OLD REDUCE_SUM_THREADS_PER_GROUP
#endif
#define REDUCE_SUM_THREADS_PER_GROUP 64
#define REDUCE_SUM_NUM_PER_STAGE 8
__kernel void
reducedot(         __global real_t*    __restrict    src1,
                   __global real_t*    __restrict    src2,
          volatile __global real_t*    __restrict     dst,
                            real_t                initVal,
                               int                    fac,
                               int                group_n,
                               int                      n,
          volatile __local  real_t*               scratch) {

    if (get_local_id(0) < REDUCE_SUM_THREADS_PER_GROUP) {
        // Calculate index for workitem
        ulong offset = get_group_id(0) * group_n;
        ulong idx_local = offset + get_local_id(0);
        ulong const stride =  (get_local_size(0) < REDUCE_SUM_THREADS_PER_GROUP) ? get_local_size(0) : REDUCE_SUM_THREADS_PER_GROUP;

        // Initialize registers for reduction
        real_t v0[2];
        v0[0] = (offset + idx_local == 0) ? initVal : real_t(0.0); v0[1] = real_t(0.0);
        real_t v1[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v2[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v3[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v4[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v5[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v6[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v7[2] = {(real_t)(0.0), (real_t)(0.0)};

        // Initialize scratch spaces
        scratch[get_local_id(0)] = (real_t)(0.0);
        scratch[get_local_id(0) + REDUCE_SUM_THREADS_PER_GROUP] = (real_t)(0.0);

        // Begin loop
        uint counterval = 0;
        while ((offset + idx_local < n) && (idx_local < group_n)) {
            // Get inputs from buffer and accumulate into v0 registers
            v0[0] += ((offset + idx_local < n) && (idx_local < group_n)) ? src1[offset + idx_local]*src2[offset + idx_local] : real_t(0.0); idx_local += stride;
            v0[1] += ((offset + idx_local < n) && (idx_local < group_n)) ? src1[offset + idx_local]*src2[offset + idx_local] : real_t(0.0); idx_local += stride;
            v0[0] += ((offset + idx_local < n) && (idx_local < group_n)) ? src1[offset + idx_local]*src2[offset + idx_local] : real_t(0.0); idx_local += stride;
            v0[1] += ((offset + idx_local < n) && (idx_local < group_n)) ? src1[offset + idx_local]*src2[offset + idx_local] : real_t(0.0); idx_local += stride;

            if (!(offset + idx_local < n) || !(idx_local < group_n)) {
                break;
            }

            // Merge results in v0 into v1 registers and clear v0 registers,
            // increment counter after
            v1[(counterval & 0x00000001)] = v0[0] + v0[1];
            v0[0] = real_t(0.0); v0[1] = real_t(0.0); counter++;
            // Determine if v1 counters should be accumulated into v2 counters.
            // If so, accumulate and clear v1 counters.
            if ((counterval ^ 0x00000003) == 0) {
                v2[((counterval >> 3) & 0x00000001)] = v1[0] + v1[1];
                v1[0] = real_t(0.0); v1[1] = real_t(0.0);
            }

            // Determine if v2 counters should be accumulated into v3 counters.
            // If so, accumulate and clear v2 counters.
            if ((counterval ^ 0x0000000f) == 0) {
                v3[((counterval >> 5) & 0x00000001)] = v2[0] + v2[1];
                v2[0] = real_t(0.0); v2[1] = real_t(0.0);
            }

            // Determine if v3 counters should be accumulated into v4 counters.
            // If so, accumulate and clear v3 counters.
            if ((counterval ^ 0x0000003f) == 0) {
                v4[((counterval >> 7) & 0x00000001)] = v3[0] + v3[1];
                v3[0] = real_t(0.0); v3[1] = real_t(0.0);
            }

            // Determine if v4 counters should be accumulated into v5 counters.
            // If so, accumulate and clear v4 counters.
            if ((counterval ^ 0x000000ff) == 0) {
                v5[((counterval >> 9) & 0x00000001)] = v4[0] + v4[1];
                v4[0] = real_t(0.0); v4[1] = real_t(0.0);
            }

            // Determine if v5 counters should be accumulated into v6 counters.
            // If so, accumulate and clear v5 counters.
            if ((counterval ^ 0x000003ff) == 0) {
                v6[((counterval >> 11) & 0x00000001)] = v5[0] + v5[1];
                v5[0] = real_t(0.0); v5[1] = real_t(0.0);
            }

            // Determine if v6 counters should be accumulated into v7 counters.
            // If so, accumulate and clear v6 counters.
            if ((counterval ^ 0x00000fff) == 0) {
                v7[((counterval >> 13) & 0x00000001)] = v6[0] + v6[1];
                v6[0] = real_t(0.0); v6[1] = real_t(0.0);
            }
        }

        // Main while loop over
        // Accumulate registers in scratch
        v1[1] += v0[0] + v0[1];
        v2[1] += v1[0] + v1[1];
        v3[1] += v2[0] + v2[1];
        v4[1] += v3[0] + v3[1];
        v5[1] += v4[0] + v4[1];
        v6[1] += v5[0] + v5[1];
        v7[1] += v6[0] + v6[1];
        scratch[get_local_id(0)] += (v7[0] + v7[1]);
        
    }

    // Force all workitems to execute here for synchronization
    scratch[get_local_id(0)] += real_t(0.0);
    barrier(CLK_LOCAL_MEM_FENCE);

    // Reduce in scratch
    if ((get_local_size(0) > 32) && (get_local_id(0) < 32)) {
        scratch[get_local_id(0)] += scratch[get_local_id(0) + 32];
    }
    if (get_local_id(0) < 16) {
        scratch[get_local_id(0)] += scratch[get_local_id(0) + 16];
        scratch[get_local_id(0)] += scratch[get_local_id(0) + 8];
        scratch[get_local_id(0)] += scratch[get_local_id(0) + 4];
        scratch[get_local_id(0)] += scratch[get_local_id(0) + 2];
        scratch[get_local_id(0)] += scratch[get_local_id(0) + 1];
    }

    // Add atomically to global buffer
    if ((get_local_id(0) == 0) && (scratch[0] != 0)){
        switch (fac)
        {
            // if every workgroup has its dedicated atomic buffer
            case 1:
                atomicAdd_r(dst[get_group_id(0)], scratch[0]);
                break;

            // every two workgroups shares one atomic buffer
            case 2:
                atomicAdd_r(dst[get_group_id(0) << 1], scratch[0]);
                break;

            // only one global atomic buffer available
            default:
                atomicAdd_r(dst, scratch[0]);
        }
    }
}

//////////////////////////////////////////////////////////////////
//
// For reducesum and reducedot, the input data buffer is
// partitioned as evenly as possible for each workgroup to
// process. Depending on the number of elements to reduce, the
// number of workgroups used will be different. Each workgroup
// will have either 64 workitems (32 allowed but avoided).
// Each workitem can process up to 4^8 = 65,536 items with optimum
// truncation errors.
//
//////////////////////////////////////////////////////////////////
