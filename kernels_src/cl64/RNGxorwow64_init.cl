/**
@file

Implements a 64-bit xorwow* generator that returns 32-bit values.

// G. Marsaglia, Xorshift RNGs, 2003
// http://www.jstatsoft.org/v08/i14/paper
*/

/**
State buffer stores 6*N uint, where N is the total number of RNGs
The first contiguous N entries are the x[0] of the xorwow states
The second contiguous N entries are the x[1] of the xorwow states
The third contiguous N entries are the x[2] of the xorwow states
The fourth contiguous N entries are the x[3] of the xorwow states
The fifth contiguous N entries are the x[4] of the xorwow states
The sixth contiguous N entries are the d (Weyl sequence number) of the xorwow states
*/

/**
Seeds xorwow RNG.

@param state_buf Variable, that holds state of the generator to be seeded.
@param seed Value used for seeding. Should be randomly generated for each instance of generator (thread).
*/
#if defined(__REAL_IS_DOUBLE__)
__kernel void
xorwow64_seed(__global uint* __restrict state_buf,
              __global uint* __restrict g_jump_matrices,
              ulong seed) {
    // Calculate indices
    int local_idx = get_local_id(0); // Work-item index within workgroup
    int grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int grp_id = get_group_id(0); // Index of workgroup
    int global_idx = grp_id * grp_sz + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * grp_sz; // Offset for memory access

    // Using local registers to compute state from seed	
    uint x[XORWOW_N];
    uint d;

    // Initialize state buffer and Weyl sequence number
    x[0] = 123456789U;
    x[1] = 362436069U;
    x[2] = 521288629U;
    x[3] = 88675123U;
    x[4] = 5783321U;
    d = 6615241U;

    // Update RNG state with seed value
    // Constants are arbitrary prime numbers
    const uint s0 = (uint)(seed) ^ 0x2c7f967fU;
    const uint s1 = (uint)(seed >> 32) ^ 0xa03697cbU;
    const uint t0 = 1228688033U * s0;
    const uint t1 = 2073658381U * s1;
    x[0] += t0;
    x[1] ^= t0;
    x[2] += t1;
    x[3] ^= t1;
    x[4] += t0;
    d += t1 + t0;

    // discarding subsequences to obtain non-overlapping random bit streams via parallelism...
    if (global_idx != 0) {
        xorwow_discard_subsequence(global_idx, x, g_jump_matrices);
    }

    // Write out state to global buffer
    int idx = global_idx;
    state_buf[idx] = x[0];
    idx += grp_offset;
    state_buf[idx] = x[1];
    idx += grp_offset;
    state_buf[idx] = x[2];
    idx += grp_offset;
    state_buf[idx] = x[3];
    idx += grp_offset;
    state_buf[idx] = x[4];
    idx += grp_offset;
    state_buf[idx] = d;
}
#endif // __REAL_IS_DOUBLE__
