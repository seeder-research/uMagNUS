__kernel void
reducemaxabs(         __global real_t* __restrict     src,
             volatile __global real_t* __restrict     dst,
                               real_t             initVal,
                                  int                   n,
             volatile  __local real_t*            scratch) {

    // Calculate indices
    unsigned int    local_idx = get_local_id(0);   // Work-item index within workgroup
    unsigned int       grp_id = get_group_id(0);   // ID of workgroup
    unsigned int       grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    unsigned int        grp_i = grp_id*grp_sz;
    unsigned int       stride = get_global_size(0);
    real_t               mine = initVal;
    unsigned int          itr = 0;

    while (grp_i < (unsigned int)(n)) {
        unsigned int i = grp_i + local_idx;
        if (i < (unsigned int)(n)) {
            mine = fmax(mine, src[i]);
        }

        grp_i += stride;
    }

    // Load workitem value into local buffer and synchronize
    scratch[local_idx] = mine;
    barrier(CLK_LOCAL_MEM_FENCE);

    // Perform reduction in the shared memory.
    if (grp_sz >= 512) {
        if (local_idx < 256)
            scratch[local_idx] = fmax(scratch[local_idx], scratch[local_idx + 256]);
        barrier(CLK_LOCAL_MEM_FENCE);
    }
    if (grp_sz >= 256) {
        if (local_idx < 128)
            scratch[local_idx] = fmax(scratch[local_idx], scratch[local_idx + 128]);
        barrier(CLK_LOCAL_MEM_FENCE);
    }
    if (grp_sz >= 128) {
        if (local_idx < 64)
            scratch[local_idx] = fmax(scratch[local_idx], scratch[local_idx + 64]);
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
        smem[local_idx] = fmax(smem[local_idx], smem[local_idx +  1]);
        mine = fmax(mine, scratch[local_idx]);
    }

    // Store reduction result for each iteration and move to next
    if (local_idx == 0) {
        mine = fmax(scratch[0], scratch[1]);
        atomicMax_r(dst, mine);
//        dst[grp_id] = fmax(scratch[0], scratch[1]);
//        dst[grp_id] = mine;
    }

}
