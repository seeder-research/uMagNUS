/*
 * @file mtgp64.cl
 *
 * @brief MTGP Sample Program for openCL 1.1
 * 1 parameter for 1 generator
 * MEXP = 11213
 */

#if defined(__REAL_IS_DOUBLE__)
/**
 * kernel function.
 * This function generates 64-bit unsigned integers in d_data
 *
 * @param[in] param_tbl recursion parameters
 * @param[in] temper_tbl tempering parameters
 * @param[in] double_temper_tbl tempering parameters for double
 * @param[in] pos_tbl pic-up positions
 * @param[in] sh1_tbl shift parameters
 * @param[in] sh2_tbl shift parameters
 * @param[in,out] d_status kernel I/O data
 * @param[out] d_data output
 * @param[in] size number of output data requested.
 */
__kernel void
mtgp64_uniform(
    __constant ulong * param_tbl,
    __constant ulong * temper_tbl,
    __constant ulong * double_temper_tbl,
    __constant uint * pos_tbl,
    __constant uint * sh1_tbl,
    __constant uint * sh2_tbl,
    __global ulong * d_status,
    __global double * d_data,
    int size)
{
    const int gid = get_group_id(0);
    const int lid = get_local_id(0);
    __local ulong status[MTGP64_LS];
    mtgp64_t mtgp;
    ulong r;
    ulong o;

    mtgp.status = status;
    mtgp.param_tbl = &param_tbl[MTGP64_TS * gid];
    mtgp.temper_tbl = &temper_tbl[MTGP64_TS * gid];
    mtgp.double_temper_tbl = &double_temper_tbl[MTGP64_TS * gid];
    mtgp.pos = pos_tbl[gid];
    mtgp.sh1 = sh1_tbl[gid];
    mtgp.sh2 = sh2_tbl[gid];

    int pos = mtgp.pos;

    // copy status data from global memory to shared memory.
    status_read64(status, d_status, gid, lid);

    ulong tmpNum[6];
    // main loop
    for (int i = 0; i < size; i += MTGP64_LS) {
	r = para_rec64(&mtgp,
		       status[MTGP64_LS - MTGP64_N + lid],
		       status[MTGP64_LS - MTGP64_N + lid + 1],
		       status[MTGP64_LS - MTGP64_N + lid + pos]);
	status[lid] = r;
	o = temper64(&mtgp, r, status[MTGP64_LS - MTGP64_N + lid + pos - 1]);
	d_data[size * gid + i + lid] = o;
	tmpNum[0] = o;
	barrier(CLK_LOCAL_MEM_FENCE);
	r = para_rec64(&mtgp,
		       status[(4 * MTGP64_TN - MTGP64_N + lid) % MTGP64_LS],
		       status[(4 * MTGP64_TN - MTGP64_N + lid + 1) % MTGP64_LS],
		       status[(4 * MTGP64_TN - MTGP64_N + lid + pos)
			      % MTGP64_LS]);
	status[lid + MTGP64_TN] = r;
	o = temper64(&mtgp,
		     r,
		     status[(4 * MTGP64_TN - MTGP64_N + lid + pos - 1)
			    % MTGP64_LS]);
	tmpNum[1] = o;
	barrier(CLK_LOCAL_MEM_FENCE);
	r = para_rec64(&mtgp,
		       status[2 * MTGP64_TN - MTGP64_N + lid],
		       status[2 * MTGP64_TN - MTGP64_N + lid + 1],
		       status[2 * MTGP64_TN - MTGP64_N + lid + pos]);
	status[lid + 2 * MTGP64_TN] = r;
	o = temper64(&mtgp, r, status[lid + pos - 1 + 2 * MTGP64_TN - MTGP64_N]);
	tmpNum[2] = o;
	barrier(CLK_LOCAL_MEM_FENCE);
	r = para_rec64(&mtgp,
		       status[MTGP64_LS - MTGP64_N + lid],
		       status[MTGP64_LS - MTGP64_N + lid + 1],
		       status[MTGP64_LS - MTGP64_N + lid + pos]);
	status[lid] = r;
	o = temper64(&mtgp, r, status[MTGP64_LS - MTGP64_N + lid + pos - 1]);
	d_data[size * gid + i + lid] = o;
	tmpNum[3] = o;
	barrier(CLK_LOCAL_MEM_FENCE);
	r = para_rec64(&mtgp,
		       status[(4 * MTGP64_TN - MTGP64_N + lid) % MTGP64_LS],
		       status[(4 * MTGP64_TN - MTGP64_N + lid + 1) % MTGP64_LS],
		       status[(4 * MTGP64_TN - MTGP64_N + lid + pos)
			      % MTGP64_LS]);
	status[lid + MTGP64_TN] = r;
	o = temper64(&mtgp,
		     r,
		     status[(4 * MTGP64_TN - MTGP64_N + lid + pos - 1)
			    % MTGP64_LS]);
	tmpNum[4] = o;
	barrier(CLK_LOCAL_MEM_FENCE);
	r = para_rec64(&mtgp,
		       status[2 * MTGP64_TN - MTGP64_N + lid],
		       status[2 * MTGP64_TN - MTGP64_N + lid + 1],
		       status[2 * MTGP64_TN - MTGP64_N + lid + pos]);
	status[lid + 2 * MTGP64_TN] = r;
	o = temper64(&mtgp, r, status[lid + pos - 1 + 2 * MTGP64_TN - MTGP64_N]);
	tmpNum[5] = o;
	barrier(CLK_LOCAL_MEM_FENCE);

        double tmpRes1 = ulong2double(tmpNum[0], tmpNum[1]); // output value
        d_data[size * gid + i + lid] = tmpRes1;

        tmpRes1 = ulong2double(tmpNum[2], tmpNum[3]); // output value
        d_data[size * gid + MTGP32_TN + i + lid] = tmpRes1;

        tmpRes1 = ulong2double(tmpNum[4], tmpNum[5]); // output value
        d_data[size * gid + 2 * MTGP32_TN + i + lid] = tmpRes1;
    }
    // write back status for next call
    status_write64(d_status, status, gid, lid);
}
#endif // __ REAL_IS_DOUBLE__
