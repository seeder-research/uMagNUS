__kernel void
reducemaxvecdiff2(         __global real_t*    __restrict      x1,
                           __global real_t*    __restrict      y1,
                           __global real_t*    __restrict      z1,
                           __global real_t*    __restrict      x2,
                           __global real_t*    __restrict      y2,
                           __global real_t*    __restrict      z2,
                  volatile __global realint_t* __restrict     dst,
                                    real_t                initVal,
                                       int                      n,
                  volatile __local  real_t*               scratch) {

    // Calculate indices
    unsigned int    local_idx = get_local_id(0);   // Work-item index within workgroup
    unsigned int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    unsigned int            i = get_group_id(0)*grp_sz + local_idx;
    unsigned int       stride = get_global_size(0);
    real_t               mine = initVal;

    while (i < n) {
        real_t3 v = distance((real_t3){x1[i], y1[i], z1[i]}, (real_t3){x2[i], y2[i], z2[i]});
        mine = fmax(mine, dot(v, v));
        i += stride;
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
    }

    // Store reduction result for each iteration and move to next
    if (local_idx == 0) {
        mine = fmax(scratch[0], scratch[1]);
#if defined(__REAL_IS_DOUBLE__)
        atom_max(dst, as_long(mine));
#else
        atom_max(dst, as_int(mine));
#endif // __REAL_IS_DOUBLE__
    }

}
