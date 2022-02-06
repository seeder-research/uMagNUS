/* ================================ */
/* mtgp32 sample kernel code        */
/* ================================ */
/**
 * This kernel function generates single precision floating point numbers
 * in the range [0, 1) in d_data.
 *
 * @param[in] param_tbl recursion parameters
 * @param[in] temper_tbl tempering parameters
 * @param[in] single_temper_tbl tempering parameters for float
 * @param[in] pos_tbl pic-up positions
 * @param[in] sh1_tbl shift parameters
 * @param[in] sh2_tbl shift parameters
 * @param[in,out] d_status kernel I/O data
 * @param[out] d_data output. IEEE single precision format.
 * @param[in] size number of output data requested.
 */
__kernel void
mtgp32_normal(__constant uint*         param_tbl,
              __constant uint*        temper_tbl,
              __constant uint* single_temper_tbl,
              __constant uint*           pos_tbl,
              __constant uint*           sh1_tbl,
              __constant uint*           sh2_tbl,
              __global   uint*          d_status,
              __global  float*            d_data,
              int size) {
    const int gid = get_group_id(0);
    const int lid = get_local_id(0);
    __local uint status[MTGP32_LS];
    mtgp32_t mtgp;
    uint r;
    uint o;

    mtgp.status = status;
    mtgp.param_tbl = &param_tbl[MTGP32_TS * gid];
    mtgp.temper_tbl = &temper_tbl[MTGP32_TS * gid];
    mtgp.single_temper_tbl = &single_temper_tbl[MTGP32_TS * gid];
    mtgp.pos = pos_tbl[gid];
    mtgp.sh1 = sh1_tbl[gid];
    mtgp.sh2 = sh2_tbl[gid];

    int pos = mtgp.pos;

    // copy status data from global memory to shared memory.
    status_read(status, d_status, gid, lid);

    uint tmpNum[12];
    float unif[6];
    float normf_[6];
    // main loop
    for (int i = 0; i < size/2; i += MTGP32_LS) {
        r = para_rec(&mtgp,
        status[MTGP32_LS - MTGP32_N + lid],
        status[MTGP32_LS - MTGP32_N + lid + 1],
        status[MTGP32_LS - MTGP32_N + lid + pos]);
        status[lid] = r;
        o = temper(&mtgp,
                   r,
                   status[MTGP32_LS - MTGP32_N + lid + pos - 1]);
        tmpNum[0] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[(4 * MTGP32_TN - MTGP32_N + lid) % MTGP32_LS],
                     status[(4 * MTGP32_TN - MTGP32_N + lid + 1) % MTGP32_LS],
                     status[(4 * MTGP32_TN - MTGP32_N + lid + pos) % MTGP32_LS]);
        status[lid + MTGP32_TN] = r;
        o = temper(&mtgp,
                   r,
                   status[(4 * MTGP32_TN - MTGP32_N + lid + pos - 1) % MTGP32_LS]);
        tmpNum[1] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[2 * MTGP32_TN - MTGP32_N + lid],
                     status[2 * MTGP32_TN - MTGP32_N + lid + 1],
                     status[2 * MTGP32_TN - MTGP32_N + lid + pos]);
        status[lid + 2 * MTGP32_TN] = r;
        o = temper(&mtgp,
                   r,
                   status[lid + pos - 1 + 2 * MTGP32_TN - MTGP32_N]);
        tmpNum[2] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[MTGP32_LS - MTGP32_N + lid],
                     status[MTGP32_LS - MTGP32_N + lid + 1],
                     status[MTGP32_LS - MTGP32_N + lid + pos]);
        status[lid] = r;
        o = temper(&mtgp,
                   r,
                   status[MTGP32_LS - MTGP32_N + lid + pos - 1]);
        tmpNum[3] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[(4 * MTGP32_TN - MTGP32_N + lid) % MTGP32_LS],
                     status[(4 * MTGP32_TN - MTGP32_N + lid + 1) % MTGP32_LS],
                     status[(4 * MTGP32_TN - MTGP32_N + lid + pos)
                            % MTGP32_LS]);
        status[lid + MTGP32_TN] = r;
        o = temper(&mtgp,
                   r,
                   status[(4 * MTGP32_TN - MTGP32_N + lid + pos - 1) % MTGP32_LS]);
        tmpNum[4] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[2 * MTGP32_TN - MTGP32_N + lid],
                     status[2 * MTGP32_TN - MTGP32_N + lid + 1],
                     status[2 * MTGP32_TN - MTGP32_N + lid + pos]);
        status[lid + 2 * MTGP32_TN] = r;
        o = temper(&mtgp,
                   r,
                   status[lid + pos - 1 + 2 * MTGP32_TN - MTGP32_N]);
        tmpNum[5] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[MTGP32_LS - MTGP32_N + lid],
                     status[MTGP32_LS - MTGP32_N + lid + 1],
                     status[MTGP32_LS - MTGP32_N + lid + pos]);
        status[lid] = r;
        o = temper(&mtgp,
                   r,
                   status[MTGP32_LS - MTGP32_N + lid + pos - 1]);
        tmpNum[6] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[(4 * MTGP32_TN - MTGP32_N + lid) % MTGP32_LS],
                     status[(4 * MTGP32_TN - MTGP32_N + lid + 1) % MTGP32_LS],
                     status[(4 * MTGP32_TN - MTGP32_N + lid + pos)
                            % MTGP32_LS]);
        status[lid + MTGP32_TN] = r;
        o = temper(&mtgp,
                   r,
                   status[(4 * MTGP32_TN - MTGP32_N + lid + pos - 1) % MTGP32_LS]);
        tmpNum[7] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[2 * MTGP32_TN - MTGP32_N + lid],
                     status[2 * MTGP32_TN - MTGP32_N + lid + 1],
                     status[2 * MTGP32_TN - MTGP32_N + lid + pos]);
        status[lid + 2 * MTGP32_TN] = r;
        o = temper(&mtgp,
                   r,
                   status[lid + pos - 1 + 2 * MTGP32_TN - MTGP32_N]);
        tmpNum[8] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[MTGP32_LS - MTGP32_N + lid],
                     status[MTGP32_LS - MTGP32_N + lid + 1],
                     status[MTGP32_LS - MTGP32_N + lid + pos]);
        status[lid] = r;
        o = temper(&mtgp,
                   r,
                   status[MTGP32_LS - MTGP32_N + lid + pos - 1]);
        tmpNum[9] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[(4 * MTGP32_TN - MTGP32_N + lid) % MTGP32_LS],
                     status[(4 * MTGP32_TN - MTGP32_N + lid + 1) % MTGP32_LS],
                     status[(4 * MTGP32_TN - MTGP32_N + lid + pos)
                            % MTGP32_LS]);
        status[lid + MTGP32_TN] = r;
        o = temper(&mtgp,
                   r,
                   status[(4 * MTGP32_TN - MTGP32_N + lid + pos - 1) % MTGP32_LS]);
        tmpNum[10] = o;
        barrier(CLK_LOCAL_MEM_FENCE);
        r = para_rec(&mtgp,
                     status[2 * MTGP32_TN - MTGP32_N + lid],
                     status[2 * MTGP32_TN - MTGP32_N + lid + 1],
                     status[2 * MTGP32_TN - MTGP32_N + lid + pos]);
        status[lid + 2 * MTGP32_TN] = r;
        o = temper(&mtgp,
                   r,
                   status[lid + pos - 1 + 2 * MTGP32_TN - MTGP32_N]);
        tmpNum[11] = o;
        barrier(CLK_LOCAL_MEM_FENCE);

        unif[0] = uint2float(tmpNum[0], tmpNum[1]);
        unif[1] = uint2float(tmpNum[2], tmpNum[3]);
        unif[2] = uint2float(tmpNum[4], tmpNum[5]);
        unif[3] = uint2float(tmpNum[6], tmpNum[7]);
        unif[4] = uint2float(tmpNum[8], tmpNum[9]);
        unif[5] = uint2float(tmpNum[10], tmpNum[11]);

        boxMuller(unif, normf_, 0);
        boxMuller(unif, normf_, 1);
        boxMuller(unif, normf_, 2);

	d_data[size * gid + i + lid] = normf_[0];
	d_data[size * gid + MTGP32_TN + i + lid] = normf_[1];
	d_data[size * gid + 2 * MTGP32_TN + i + lid] = normf_[2];
	d_data[size * gid + (size / 2) + i + lid] = normf_[3];
	d_data[size * gid + (size / 2) + MTGP32_TN + i + lid] = normf_[4];
	d_data[size * gid + (size / 2) + 2 * MTGP32_TN + i + lid] = normf_[5];
    }
    // write back status for next call
    status_write(d_status, status, gid, lid);
}
