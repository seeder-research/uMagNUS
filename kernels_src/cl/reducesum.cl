// Uses 16 registers (if single-precision, then 64B per
// workitem. Twice if double-precision, i.e., 128B per
// workitem) to produce reduction sum in one thread. Allow
// maximum of 64 workitems per workgroup to run.
// 8 nested for loops
// Single workgroup of 64 workitems can handle up to
// 1024x1024x1024 items
#ifdef(REDUCE_SUM_THREADS_PER_GROUP)
#define REDUCE_SUM_THREADS_PER_GROUP_OLD REDUCE_SUM_THREADS_PER_GROUP
#endif
#define REDUCE_SUM_THREADS_PER_GROUP 64
#define REDUCE_SUM_NUM_PER_STAGE 8
__kernel void
reducesum(         __global real_t*    __restrict     src,
          volatile __global real_t*    __restrict     dst,
                            real_t                initVal,
                               int                      n,
          volatile __local  real_t*               scratch) {

    ulong const block_size = get_local_size(0);
    ulong const idx_in_block = get_local_id(0);

    if (idx_in_block < REDUCE_SUM_THREADS_PER_GROUP) {
        // Calculate index for workitem
        ulong idx_global = get_group_id(0) * REDUCE_SUM_THREADS_PER_GROUP + idx_in_block;
        ulong const stride =  REDUCE_SUM_THREADS_PER_GROUP * get_num_groups(0);

        // Initialize registers for reduction
        real_t v0[2];
        v0[0] = (idx_global == 0) ? initVal : real_t(0.0); v0[1] = real_t(0.0);
        real_t v1[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v2[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v3[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v4[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v5[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v6[2] = {(real_t)(0.0), (real_t)(0.0)};
        real_t v7[2] = {(real_t)(0.0), (real_t)(0.0)};

        // Initialize scratch spaces
        scratch[idx_in_block] = (real_t)(0.0);
        scratch[idx_in_block + REDUCE_SUM_THREADS_PER_GROUP] = (real_t)(0.0);

        // Begin loop
        for (uint i0 = 0; i0 < REDUCE_SUM_NUM_PER_STAGE; i0++) {
            for (uint i1 = 0; i1 < REDUCE_SUM_NUM_PER_STAGE; i1++) {
                for (uint i2 = 0; i2 < REDUCE_SUM_NUM_PER_STAGE; i2++) {
                    for (uint i3 = 0; i3 < REDUCE_SUM_NUM_PER_STAGE; i3++) {
                        for (uint i4 = 0; i4 < REDUCE_SUM_NUM_PER_STAGE; i4++) {
                            for (uint i5 = 0; i5 < REDUCE_SUM_NUM_PER_STAGE; i5++) {
                                for (uint i6 = 0; i6 < REDUCE_SUM_NUM_PER_STAGE; i6++) {
                                    for (uint i7 = 0; i7 < REDUCE_SUM_NUM_PER_STAGE; i7++) {
                                        v0[i7 & 0x00000001] += (idx_global < n) ? src[idx_global] : real_t(0.0);
                                        idx_global += stride;
                                        if (idx_global >= n) {
                                            break;
                                        }
                                    }
                                    v1[i6 & 0x00000001] += (v0[0] + v0[1]);
                                    if (idx_global >= n) {
                                        break;
                                    }
                                    v0[0] = real_t(0.0);
                                    v0[1] = real_t(0.0);
                                }
                                v2[i5 & 0x00000001] += (v1[0] + v1[1]);
                                if (idx_global >= n) {
                                    break;
                                }
                                v1[0] = real_t(0.0);
                                v1[1] = real_t(0.0);
                            }
                            v3[i4 & 0x00000001] += (v2[0] + v2[1]);
                            if (idx_global >= n) {
                                break;
                            }
                            v2[0] = real_t(0.0);
                            v2[1] = real_t(0.0);
                        }
                        v4[i3 & 0x00000001] += (v3[0] + v3[1]);
                        if (idx_global >= n) {
                            break;
                        }
                        v3[0] = real_t(0.0);
                        v3[1] = real_t(0.0);
                    }
                    v5[i2 & 0x00000001] += (v4[0] + v4[1]);
                    if (idx_global >= n) {
                        break;
                    }
                    v4[0] = real_t(0.0);
                    v4[1] = real_t(0.0);
                }
                v6[i1 & 0x00000001] += (v5[0] + v5[1]);
                if (idx_global >= n) {
                    break;
                }
                v5[0] = real_t(0.0);
                v5[1] = real_t(0.0);
            }
            v7[i0 & 0x00000001] += (v6[0] + v6[1]);
            if (idx_global >= n) {
                break;
            }
            v6[0] = real_t(0.0);
            v6[1] = real_t(0.0);
        }

        // Accumulate in scratch
        scratch[idx_in_block] += (v7[0] + v7[1]);
        barrier(CLK_LOCAL_MEM_FENCE);

        // Reduce in scratch
        if (idx_in_block < 32) {
            scratch[idx_in_block] += scratch[idx_in_block + 32];
            scratch[idx_in_block] += scratch[idx_in_block + 16];
            scratch[idx_in_block] += scratch[idx_in_block + 8];
            scratch[idx_in_block] += scratch[idx_in_block + 4];
            scratch[idx_in_block] += scratch[idx_in_block + 2];
            scratch[idx_in_block] += scratch[idx_in_block + 1];
        }

        // Add atomically to global buffer
        if (idx_in_block == 0) {
            atomicAdd_r(dst, scratch[0]);
        }
        
    }
}
