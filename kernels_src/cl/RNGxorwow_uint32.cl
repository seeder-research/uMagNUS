/**
@file

Implements a 64-bit xorwow* generator that returns 32-bit values.

// G. Marsaglia, Xorshift RNGs, 2003
// http://www.jstatsoft.org/v08/i14/paper
*/

/**
Generates a random 32-bit unsigned integer using xorwow RNG.

@param state_buf State of the RNG to use.
@param d_data Output.
*/
__kernel void
xorwow_uint(
    __global uint* __restrict state_buf,
    __global uint* __restrict   d_data,
    int count){
    // Calculate indices
    int local_idx = get_local_id(0); // Work-item index within workgroup
    int grp_sz = get_local_size(0); // Total number of work-items in each workgroup
    int grp_id = get_group_id(0); // Index of workgroup
    int global_idx = grp_id * grp_sz + local_idx; // Calculate global index of work-item
    int grp_offset = get_num_groups(0) * grp_sz; // Offset for memory access

    // Only threads witin the count bounds will generate the random number
    if (global_idx < count) {
        // Using local registers to compute next state
        uint x[5];
        uint d;

        // Get state from global buffer
        int idx = global_idx;
        x[0] = state_buf[idx];
        idx += grp_offset;
        x[1] = state_buf[idx];
        idx += grp_offset;
        x[2] = state_buf[idx];
        idx += grp_offset;
        x[3] = state_buf[idx];
        idx += grp_offset;
        x[4] = state_buf[idx];
        idx += grp_offset;
        d = state_buf[idx];

        // For each thread that is launched, iterate until the index is out of bounds
        for (uint pos = global_idx; pos < count; pos += grp_offset) {
            const uint t = x[0] ^ (x[0] >> 2);
            x[0] = x[1];
            x[1] = x[2];
            x[2] = x[3];
            x[3] = x[4];
            x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

            d += 362437;

            d_data[pos] = (d + x[4]); // output value
        }

        // update the state buffer with the latest state
        idx = global_idx;
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
}
