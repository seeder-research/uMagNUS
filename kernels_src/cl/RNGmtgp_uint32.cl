/* ================================ */
/* mtgp32 sample kernel code        */
/* ================================ */
/**
 * This kernel function generates 32-bit unsigned integers in d_data.
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
__kernel void mtgp32_uint32(
    __constant uint* param_tbl,
    __constant uint* temper_tbl,
    __constant uint* single_temper_tbl,
    __constant uint* pos_tbl,
    __constant uint* sh1_tbl,
    __constant uint* sh2_tbl,
    __global uint* d_status,
    __global uint* d_data,
    int size)
{
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

    // main loop
    for (int i = 0; i < size; i += MTGP32_LS) {
	r = para_rec(&mtgp,
		     status[MTGP32_LS - MTGP32_N + lid],
		     status[MTGP32_LS - MTGP32_N + lid + 1],
		     status[MTGP32_LS - MTGP32_N + lid + pos]);
	status[lid] = r;
	o = temper(&mtgp,
			  r,
			  status[MTGP32_LS - MTGP32_N + lid + pos - 1]);
	d_data[size * gid + i + lid] = o;
	barrier(CLK_LOCAL_MEM_FENCE);
	r = para_rec(&mtgp,
		     status[(4 * MTGP32_TN - MTGP32_N + lid) % MTGP32_LS],
		     status[(4 * MTGP32_TN - MTGP32_N + lid + 1) % MTGP32_LS],
		     status[(4 * MTGP32_TN - MTGP32_N + lid + pos)
			    % MTGP32_LS]);
	status[lid + MTGP32_TN] = r;
	o = temper(
	    &mtgp,
	    r,
	    status[(4 * MTGP32_TN - MTGP32_N + lid + pos - 1) % MTGP32_LS]);
	d_data[size * gid + MTGP32_TN + i + lid] = o;
	barrier(CLK_LOCAL_MEM_FENCE);
	r = para_rec(&mtgp,
		     status[2 * MTGP32_TN - MTGP32_N + lid],
		     status[2 * MTGP32_TN - MTGP32_N + lid + 1],
		     status[2 * MTGP32_TN - MTGP32_N + lid + pos]);
	status[lid + 2 * MTGP32_TN] = r;
	o = temper(&mtgp,
			  r,
			  status[lid + pos - 1 + 2 * MTGP32_TN - MTGP32_N]);
	d_data[size * gid + 2 * MTGP32_TN + i + lid] = o;
	barrier(CLK_LOCAL_MEM_FENCE);
    }
    // write back status for next call
    status_write(d_status, status, gid, lid);
}
