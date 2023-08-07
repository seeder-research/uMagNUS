// Uses 10 registers fo produce reduction sum in one thread
// 5 nested for loops
// 1) Outmost loop adds 8 datapoints into two v1 registers
//   (4 each). If 8 datapoints have been accessed, will go
//   into next inner loop.
//     2) Total of 8 results from v1 registers merged into
//        two v2 registers (4 each). If 64 datapoints have
//        been accessed, will go into next inner loop.
//         3) Total of 8 results from v2 registers merged
//            into v3 registers (4 each). If 512 datapoints
//            have been accessed, will go into next inner
//            loop.
//             4) Total of 8 results from v3 registers
//                merged into v4 registers (4 each). If
//                4096 datapoints have been accessed, will
//                go into next inner loop.
//                 5) Total of 8 results from v4 registers
//                    merged into v5 registers (4 each). If
//                    32768 datapoints have been accessed,
//                    will go into next inner loop.
//                     6) Total of 8 results from v5
//                        registers merged into scratch
//                        (4 into workitem position). At
//                        this point, should have accessed
//                        131072 datapoints
// 7) Results in scratch is reduced to one in scratch[0]
// 8) fac is used to determine how the merged result is
//    stored into global buffer. If fac is 0, all
//    workgroups output to the same location in output
//    buffer atomically. If fac is 1, workgroups output
//    to location in output buffer corresponding to the
//    workgroup ID. Will need one more reduce stage to
//    merge.
// 9) Number of workitems per block should be a power of 2.
//    However, threads_pre_group will be used to determine
//    which workitems in the workgroup will be active.
__kernel void
reducesum(         __global real_t*    __restrict     src,
          volatile __global real_t*    __restrict     dst,
                            real_t                initVal,
                               int                    fac,
                               int                      n,
                               int      threads_per_group,
          volatile __local  real_t*               scratch){

    ulong const block_size = get_local_size(0);
    ulong const idx_in_block = get_local_id(0);
    ulong idx_global = get_group_id(0) * threads_per_group + idx_in_block;
    ulong const stride =  threads_per_group * get_num_groups(0);
    uint v1idx0 = 0;
    uint flag = 0;

    real_t v1[2];
    real_t v2[2];
    real_t v3[2];
    real_t v4[2];
    real_t v5[2];
    real_t tmp;

    scratch[idx_in_block] = (real_t)(0.);
    v1[0] = (real_t)(0.); v1[1] = (real_t)(0.);
    v2[0] = (real_t)(0.); v2[1] = (real_t)(0.);
    v3[0] = (real_t)(0.); v3[1] = (real_t)(0.);
    v4[0] = (real_t)(0.); v4[1] = (real_t)(0.);
    v5[0] = (real_t)(0.); v5[1] = (real_t)(0.);

    // Read from global and accumulate in work-item registers
    while (idx_global < n) {
        // First stage adds 4 datapoints in 2 registers
        flag = v1idx0 & 0x00000001;
        v1[flag] += (idx_in_block < threads_per_group) ? src[idx_global] : 0;
        idx_global += stride;
        v1idx0++;
        flag = v1idx0 & 0x00000007;
        if (flag == 0) {
            // accumulate all values in v1 into an entry in v2
            tmp = v1[0] + v1[1];
            v1[0] = (real_t)(0.0);
            v1[1] = (real_t)(0.0);
            flag = (v1idx0 >> 3);
            v2[(flag & 0x00000001)] += tmp;
            flag = v1idx0 & 0x0000003f;
            if (flag == 0) {
                tmp = v2[0] + v2[1];
                v2[0] = (real_t)(0.0);
                v2[1] = (real_t)(0.0);
                flag = (v1idx0 >> 6);
                v3[(flag & 0x00000001)] += tmp;
                flag = v1idx0 & 0x000001FF;
                if (flag == 0) {
                    tmp = v3[0] + v3[1];
                    v3[0] = (real_t)(0.0);
                    v3[1] = (real_t)(0.0);
                    flag = (v1idx0 >> 9);
                    v4[(flag & 0x00000001)] += tmp;
                    flag = v1idx0 & 0x00000FFF;
                    if (flag == 0) {
                        tmp = v4[0] + v4[1];
                        v4[0] = (real_t)(0.0);
                        v4[1] = (real_t)(0.0);
                        flag = (v1idx0 >> 12);
                        v5[(flag & 0x00000001)] += tmp;
                        flag = v1idx0 & 0x00007FFF;
                        if (flag == 0) {
                            tmp = v5[0] + v5[1];
                            v5[0] = (real_t)(0.0);
                            v5[1] = (real_t)(0.0);
                            scratch[idx_in_block] += tmp;
                            flag = v1idx0 & 0x0003FFFF;
                            if (flag == 0) {
                                break;
                            }
                        }
                    }
                }
            }
        }
    }

    // Merge results for access into workitem storage in
    // shared memory.
    v1[0] += v1[1];
    v2[1] += v1[0];
    v2[0] += v2[1];
    v3[1] += v2[0];
    v3[0] += v3[1];
    v4[1] += v3[0];
    v4[0] += v4[1];
    v5[1] += v4[0];
    v5[0] += v5[1];

    scratch[idx_in_block] += v5[0];
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
        if (fac == 0) {
            atomicAdd_r(dst, scratch[0]);
        } else {
            atomicAdd_r(dst[get_group_id(0)], scratch[0]);
        }
    }

}
