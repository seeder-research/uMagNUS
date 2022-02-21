/*
 * @file mtgp64.cl
 *
 * @brief MTGP Sample Program for openCL 1.1
 * 1 parameter for 1 generator
 * MEXP = 11213
 */

/* ================================ */
/* mtgp64 sample kernel code        */
/* ================================ */
/**
 * This function sets up initial state by seed.
 * kernel function.
 *
 * @param[in] param_tbl recursion parameters
 * @param[in] temper_tbl tempering parameters
 * @param[in] double_temper_tbl tempering parameters for double
 * @param[in] pos_tbl pic-up positions
 * @param[in] sh1_tbl shift parameters
 * @param[in] sh2_tbl shift parameters
 * @param[out] d_status kernel I/O data
 * @param[in] seed initializing seed
 */
__kernel void
mtgp64_seed(
    __constant ulong * param_tbl,
    __constant ulong * temper_tbl,
    __constant ulong * double_temper_tbl,
    __constant uint * pos_tbl,
    __constant uint * sh1_tbl,
    __constant uint * sh2_tbl,
    __global ulong * __restrict d_status,
    __global ulong * __restrict     seed) {
    const int gid = get_group_id(0);
    const int lid = get_local_id(0);
    const int local_size = get_local_size(0);
    __local ulong status[MTGP64_N];
    __local ulong seedVal;
    mtgp64_t mtgp;
    mtgp.status = status;
    mtgp.param_tbl = &param_tbl[MTGP64_TS * gid];
    mtgp.temper_tbl = &temper_tbl[MTGP64_TS * gid];
    mtgp.double_temper_tbl = &double_temper_tbl[MTGP64_TS * gid];
    mtgp.pos = pos_tbl[gid];
    mtgp.sh1 = sh1_tbl[gid];
    mtgp.sh2 = sh2_tbl[gid];

    // initialize
    if (lid == 0) {
        seedVal = seed[gid];
    }
    mtgp64_init_state(&mtgp, seedVal);
    barrier(CLK_LOCAL_MEM_FENCE);

    d_status[gid * MTGP64_N + lid] = status[lid];
    if ((local_size < MTGP64_N) && (lid < MTGP64_N - MTGP64_TN)) {
	d_status[gid * MTGP64_N + MTGP64_TN + lid] = status[MTGP64_TN + lid];
    }
}
