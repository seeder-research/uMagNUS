__kernel void
reducesum(         __global real_t*    __restrict     src,
          volatile __global real_t*    __restrict     dst,
                            real_t                initVal,
                               int                      n,
                               int      threads_per_group,
          volatile __local  real_t*               scratch){

    ulong const block_size = get_local_size(0);
    ulong const idx_in_block = get_local_id(0);
    ulong idx_global = get_group_id(0) * threads_per_group + idx_in_block;
    ulong const stride =  threads_per_group * get_num_groups(0);
    uint v1idx0 = 0; uint v1idx1 = 0;
    uint v2idx0 = 0; uint v2idx1 = 0;

    real_t v1[8];
    real_t v2[4];

    scratch[idx_in_block] = (real_t)(0.);
    v1[0] = (real_t)(0.); v1[1] = (real_t)(0.); v1[2] = (real_t)(0.); v1[3] = (real_t)(0.);
    v1[4] = (real_t)(0.); v1[5] = (real_t)(0.); v1[6] = (real_t)(0.); v1[7] = (real_t)(0.);
    v2[0] = (real_t)(0.); v2[1] = (real_t)(0.); v2[2] = (real_t)(0.); v2[3] = (real_t)(0.);

    // Read from global and accumulate in work-item registers
    while (idx_global < n) {
        // First 4 datapoints
        v1[v1idx0] += (idx_in_block < threads_per_group) ? src[idx_global] : 0;
        idx_global += stride;
        if (v1idx0 == 8) {
            v1idx0 = 0;
            v1idx1++;
            if (v1idx1 == 5) {
               // accumulate all values in v1 into an entry in v2
               v1[0] += v1[4]; v1[1] += v1[5]; v1[2] += v1[6]; v1[3] += v1[7];
               v1[0] += v1[2]; v1[1] += v1[3];
               v1[0] += v1[1];
               v2[v2idx0] += v1[0];
               v1[0] = (real_t)(0.); v1[1] = (real_t)(0.); v1[2] = (real_t)(0.); v1[3] = (real_t)(0.);
               v1[4] = (real_t)(0.); v1[5] = (real_t)(0.); v1[6] = (real_t)(0.); v1[7] = (real_t)(0.);
               v1idx1 = 0;
               v2idx0++;
               if (v2idx0 == 4) {
                   v2idx0 = 0;
                   v2idx1++;
                   if (v2idx1 == 4) {
                       // accumulate all values in v2 into local memory
                       v2[0] += v2[2]; v2[1] += v2[3];
                       v2[0] += v2[1];
                       scratch[idx_in_block] += v2[0];
                       v2[0] = (real_t)(0.); v2[1] = (real_t)(0.); v2[2] = (real_t)(0.); v2[3] = (real_t)(0.);
                       v2idx1 = 0;
                   }
               }
            }
    }

    v1[0] += v1[4]; v1[1] += v1[5]; v1[2] += v1[6]; v1[3] += v1[7];
    v1[0] += v1[2]; v1[1] += v1[3];
    v1[0] += v1[1];
    v2[3] += v1[0];
    v2[0] += v2[2]; v2[1] += v2[3];
    v2[0] += v2[1];
    scratch[idx_in_block] += v2[0];
    barrier(CLK_LOCAL_MEM_FENCE);

    // Perform reduction in the shared memory.
    if (block_size >= 512) {
        if (idx_in_block < 256)
            scratch[idx_in_block] += scratch[idx_in_block + 256];
        barrier(CLK_LOCAL_MEM_FENCE);
    }
    if (block_size >= 256) {
        if (idx_in_block < 128)
            scratch[idx_in_block] += scratch[idx_in_block + 128];
        barrier(CLK_LOCAL_MEM_FENCE);
    }
    if (block_size >= 128) {
        if (idx_in_block < 64)
            scratch[idx_in_block] += scratch[idx_in_block + 64];
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    if (idx_in_block < 32) {
        if (block_size >= 64)
            scratch[idx_in_block] += scratch[idx_in_block + 32];
        if (block_size >= 32)
            scratch[idx_in_block] += scratch[idx_in_block + 16];
        if (block_size >= 16)
            scratch[idx_in_block] += scratch[idx_in_block + 8];
        if (block_size >= 8)
            scratch[idx_in_block] += scratch[idx_in_block + 4];
        if (block_size >= 4)
            scratch[idx_in_block] += scratch[idx_in_block + 2];
        if (block_size >= 2)
            scratch[idx_in_block] += scratch[idx_in_block + 1];
    }

    // Add atomically to global buffer
    if (idx_in_block == 0) {
        atomicAdd_r(dst, scratch[0]);
//            dst[get_group_id(0)] = scratch[0];
    }

}
