/**
@file

Implements a 64-bit xorwow* generator that returns 32-bit values.

// G. Marsaglia, Xorshift RNGs, 2003
// http://www.jstatsoft.org/v08/i14/paper
*/

/**
Generates a random float using xorwow RNG.

@param state State of the RNG to use.
@param d_data Output.
*/
__kernel void xorwow_normal(
        __global uint* state_buf,
        __global float* d_data,
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
                bool generate = true;
                float z0 = 0.0f;
                float z1 = 0.0f;

                // For each thread that is launched, iterate until the index is out of bounds
                for (uint pos = global_idx; pos < count; pos += grp_offset) {
                        if (generate) {
                                // generate a pair of uint32 (one uint64)
                                // first number...
                                uint t = x[0] ^ (x[0] >> 2);
                                x[0] = x[1];
                                x[1] = x[2];
                                x[2] = x[3];
                                x[3] = x[4];
                                x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

                                d += 362437;

                                uint num1 = d+x[4];

                                // second number...
                                t = x[0] ^ (x[0] >> 2);
                                x[0] = x[1];
                                x[1] = x[2];
                                x[2] = x[3];
                                x[3] = x[4];
                                x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

                                d += 362437;

                                uint num2 = d+x[4];

                                // Find single-precision floating point representation
                                // for the integer formed by the pair of uint32...
                                // the approach is to treat the uint64 as a 32-bit uint
                                // with 32 fractional bits. We divide the number by
                                // 2^32 to get the floating point number. The steps are
                                // to first find the 23-bit mantissa based on the uint64.
                                // We then use the position of the leading 1 to find the
                                // relative exponent. For example, if the MSB in the
                                // first uint32 is 1, then the integer portion is at
                                // least 2^31, which means dividing by 2^32 gives a
                                // result between 0.5f (0x3f000000) and 1.0f (0x3f800000).
                                // Hence, the exponent bits are 01111110 = 126. Note that
                                // the PRNG never returns uint32(0) so we are guaranteed.
                                // that a leading 1 exists in the first uint32.
                                uint finalNum = 0;
                                uint expo = 32;
                                for (;expo > 0; expo--) {
                                        uint flag0 = num1 & 0x80000000;
                                        num1 <<= 1;
                                        if (flag0 != 0) {
                                                break;
                                        }
                                }
                                if (expo < 23) {
                                        uint maskbits = 0xffffffff;
                                        uint shPos = 23 - expo;
                                        maskbits >>= shPos;
                                        maskbits <<= shPos;
                                        maskbits = ~maskbits;
                                        finalNum ^= (num2 & maskbits);
                                }
                                finalNum ^= (num1 >> 9);
                                uint newExpo = 94 + expo;
                                finalNum ^= newExpo << 23;
                                float tmpRes1 = as_float(finalNum); // output value

                                // Repeat for second float...
                                // generate a pair of uint32 (one uint64)
                                // first number...
                                t = x[0] ^ (x[0] >> 2);
                                x[0] = x[1];
                                x[1] = x[2];
                                x[2] = x[3];
                                x[3] = x[4];
                                x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

                                d += 362437;

                                num1 = d+x[4];

                                // second number...
                                t = x[0] ^ (x[0] >> 2);
                                x[0] = x[1];
                                x[1] = x[2];
                                x[2] = x[3];
                                x[3] = x[4];
                                x[4] = (x[4] ^ (x[4] << 4)) ^ (t ^ (t << 1));

                                d += 362437;

                                num2 = d+x[4];

                                // Find single-precision floating point representation
                                // for the integer formed by the pair of uint32...
                                finalNum = 0;
                                expo = 32;
                                for (;expo > 0; expo--) {
                                        uint flag0 = num1 & 0x80000000;
                                        num1 <<= 1;
                                        if (flag0 != 0) {
                                                break;
                                        }
                                }
                                if (expo < 23) {
                                        uint maskbits = 0xffffffff;
                                        uint shPos = 23 - expo;
                                        maskbits >>= shPos;
                                        maskbits <<= shPos;
                                        maskbits = ~maskbits;
                                        finalNum ^= (num2 & maskbits);
                                }
                                finalNum ^= (num1 >> 9);
                                newExpo = 94 + expo;
                                finalNum ^= newExpo << 23;
                                float tmpRes2 = as_float(finalNum); // output value

                                z0 = sqrt( -2.0f * log(tmpRes1)) * cospi(2.0f * tmpRes2);
                                z1 = sqrt( -2.0f * log(tmpRes1)) * sinpi(2.0f * tmpRes2);
                                d_data[pos] = z0; // output normal random value
                                generate = !generate;
                        } else {
                                d_data[pos] = z1; // output normal random value
                        }
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
