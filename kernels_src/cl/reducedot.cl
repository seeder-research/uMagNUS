__kernel void
reducedot(__global real_t*          __restrict     src1,
          __global real_t*          __restrict     src2,
          __global real_t* volatile __restrict      dst,
                   real_t              initVal,
                      int                    n,
          __local  real_t*             scratch){

    // Calculate indices
    int    local_idx = get_local_id(0);   // Work-item index within workgroup
    int       grp_id = get_group_id(0);   // ID of workgroup
    int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int        grp_i = grp_id * grp_sz;
    int       stride = get_global_size(0);

    // Accumulators for intermediate results
    real_t data1 = initVal;
    int      itr = 0;

    // Read into local buffer, reduce and accumulate in workitem registers
    while (grp_i < n) {
        int i = grp_i + local_idx;
        scratch[local_idx] = 0.0;
        if (i < n) {
            scratch[local_idx] = src1[i] * src2[i];
        }

        // Sync all workitems before reducing
        barrier(CLK_LOCAL_MEM_FENCE);

        // For loop reduction
        for (int s = (grp_sz >> 1); s > 32; s >>= 1) {
            if (local_idx < s) {
                scratch[local_idx] += scratch[local_idx + s];
            }

            // Add barrier to sync all threads before next iteration
            barrier(CLK_LOCAL_MEM_FENCE);

        }

        // Unroll loop
        if (local_idx < 32) {
            __local real_t* volatile smem = scratch;
            smem[local_idx] += smem[local_idx + 32];
            smem[local_idx] += smem[local_idx + 16];
            smem[local_idx] += smem[local_idx +  8];
            smem[local_idx] += smem[local_idx +  4];
            smem[local_idx] += smem[local_idx +  2];
            smem[local_idx] += smem[local_idx +  1];
        }

        // Sync all workitems before reducing
        barrier(CLK_LOCAL_MEM_FENCE);

        // Write back to global buffer
        if (local_idx == itr) {
            data1 += scratch[0];
        }

        // Sync all workitems before reducing
        barrier(CLK_LOCAL_MEM_FENCE);

        itr++;
        if (itr >= grp_sz) {
            itr = 0;
        }

        grp_i += stride;

    }

    // Reduce to global buffer
    // Load accumulator values to local buffer and reduce
    scratch[local_idx] = data1;

    // Sync all workitems before reducing
    barrier(CLK_LOCAL_MEM_FENCE);

    // For loop reduction
    for (int s = (grp_sz >> 1); s > 32; s >>= 1) {
        if (local_idx < s) {
            scratch[local_idx] += scratch[local_idx + s];
        }

        // Add barrier to sync all threads before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);

    }

    // Unroll loop
    if (local_idx < 32) {
        __local real_t* volatile smem = scratch;
        smem[local_idx] += smem[local_idx + 32];
        smem[local_idx] += smem[local_idx + 16];
        smem[local_idx] += smem[local_idx +  8];
        smem[local_idx] += smem[local_idx +  4];
        smem[local_idx] += smem[local_idx +  2];
        smem[local_idx] += smem[local_idx +  1];
    }

    // Write back to global buffer
    if (local_idx == 0) {
          real_t tmp = scratch[0];
#if defined(__REAL_IS_DOUBLE__)
          while ((tmp = atomic_xchg(dst, as_long(as_double(atomic_xchg(dst, 0.0)) + tmp))) != 0.0);
#else
          while ((tmp = atomic_xchg(dst, as_int(as_float(atomic_xchg(dst, 0.0)) + tmp))) != 0.0);
#endif // __REAL_IS_DOUBLE__
    }
}
