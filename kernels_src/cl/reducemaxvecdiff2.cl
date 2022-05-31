__kernel void
reducemaxvecdiff2(__global real_t* __restrict      x1,
                  __global real_t* __restrict      y1,
                  __global real_t* __restrict      z1,
                  __global real_t* __restrict      x2,
                  __global real_t* __restrict      y2,
                  __global real_t* __restrict      z2,
                  __global real_t* __restrict     dst,
                           real_t             initVal,
                              int                   n,
                  __local  real_t*            scratch) {

    // Calculate indices
    int    local_idx = get_local_id(0);   // Work-item index within workgroup
    int       grp_id = get_group_id(0);   // ID of workgroup
    int      num_grp = get_num_groups(0); // Number of workgroups launched
    int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int            i = grp_id * grp_sz + local_idx;
    int       stride = num_grp * grp_sz;
    real_t      mine = initVal;

    while (i < n) {
        real_t3 dv = (real_t3){x1[i] - x2[i], y1[i] - y2[i], z1[i] - z2[i]};
        real_t   v = dot(v, v);
        mine = fmax(mine, v);
        i += stride;
    }

    // Load value into local buffer to reduce
    scratch[local_idx] = mine;

    // Sync all workitems before reducing
    barrier(CLK_LOCAL_MEM_FENCE);

    // For loop reduction
    for (int s = (grp_sz >> 1); s > 32; s >>= 1) {
        if (local_idx < s) {
            scratch[local_idx] = fmax(scratch[local_idx], scratch[local_idx + s]);
        }

        // Add barrier to sync all threads before next iteration
        barrier(CLK_LOCAL_MEM_FENCE);

    }

    // Unroll loop
    if (local_idx < 32) {
        __local volatile real_t* smem = scratch;
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 32]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx + 16]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  8]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  4]);
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  2]);
    }

    // Write back to global buffer
    if (local_idx == 0) {
        dst[grp_id] = fmax(scratch[0], scratch[1]);
    }

}
