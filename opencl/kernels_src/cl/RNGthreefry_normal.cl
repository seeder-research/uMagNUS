/**
@file
Implements threefry RNG.
*******************************************************
 * Modified version of Random123 library:
 * https://www.deshawresearch.com/downloads/download_random123.cgi/
 * The original copyright can be seen here:
 *
 * RANDOM123 LICENSE AGREEMENT
 *
 * Copyright 2010-2011, D. E. Shaw Research. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * * Redistributions of source code must retain the above copyright notice,
 *   this list of conditions, and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright
 *   notice, this list of conditions, and the following disclaimer in the
 *   documentation and/or other materials provided with the distribution.
 *
 * Neither the name of D. E. Shaw Research nor the names of its contributors
 * may be used to endorse or promote products derived from this software
 * without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *********************************************************/

/**
State of threefry RNG. We will store in global buffer as a set of uint
**
typedef struct{
        uint counter[4];
        uint result[4];
        uint key[4];
        uint tracker;
} threefry_state;
**/

/**
Generates a random 32-bit unsigned integer using threefry RNG.
@param state State of the RNG to use.
**/
__kernel void
threefry_normal(__global uint __restrict *state_key,
                __global uint __restrict *state_counter,
                __global uint __restrict *state_result,
                __global uint __restrict *state_tracker,
                __global uint __restrict *output,
                int data_size) {
    uint index = get_group_id(0) * ELEMENTS_PER_BLOCK + get_local_id(0);
    uint totalWorkItems = get_global_size(0);
    uint tmpIdx = index;
    threefry_state state_;
    threefry_state *state = &state_;

    // For first out of four sets...
    // Read in counter
    state->counter[0] = state_counter[tmpIdx];
    // Read in result
    state->result[0] = state_result[tmpIdx];
    // Read in key
    state->key[0] = state_key[tmpIdx];
    // Read in tracker
    state->tracker = state_tracker[tmpIdx];

    // For second out of four sets...
    tmpIdx += totalWorkItems;
    // Read in counter
    state->counter[1] = state_counter[tmpIdx];
    // Read in result
    state->result[1] = state_result[tmpIdx];
    // Read in key
    state->key[1] = state_key[tmpIdx];

    // For third out of four sets...
    tmpIdx += totalWorkItems;
    // Read in counter
    state->counter[2] = state_counter[tmpIdx];
    // Read in result
    state->result[2] = state_result[tmpIdx];
    // Read in key
    state->key[2] = state_key[tmpIdx];

    // For last out of four sets...
    tmpIdx += totalWorkItems;
    // Read in counter
    state->counter[3] = state_counter[tmpIdx];
    // Read in result
    state->result[3] = state_result[tmpIdx];
    // Read in key
    state->key[3] = state_key[tmpIdx];

    for (uint outIndex = index; index < data_size / 2; index += totalWorkItems) {
        uint num1[2];
        float res1[4];
        uint lidx = 0;
        if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx++] = tmp;
        } else {
            num1[lidx++] = state->result[state->tracker++];
        }
        if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx] = tmp;
        } else {
            num1[lidx] = state->result[state->tracker++];
        }
        res1[0] = uint2float(num1[0], num1[1]);
        uint lidx = 0;
        if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx++] = tmp;
        } else {
            num1[lidx++] = state->result[state->tracker++];
        }
        if (state->tracker == 3) {
            uint tmp = state->result[3];
            if (++state->counter[0] == 0) {
                if (++state->counter[1] == 0) {
                    if (++state->counter[2] == 0) {
                        ++state->counter[3];
                    }
                }
            }
            threefry_round(state);
            state->tracker = 0;
            num1[lidx] = tmp;
        } else {
            num1[lidx] = state->result[state->tracker++];
        }
        res1[1] = uint2float(num1[0], num1[1]);
        res1[2] = sqrt( -2.0f * log(res1[0])) * cospi(2.0f * res1[1]);
        res1[3] = sqrt( -2.0f * log(res1[0])) * sinpi(2.0f * res1[1]);
        output[outIndex] = res1[2];
        output[outIndex + (data_size/2)] = res1[3];
    }
    
    // For first out of four sets...
    // Write out counter
    tmpIdx = index;
    state_counter[tmpIdx] = state->counter[0];
    // Write out result
    state_result[tmpIdx] = state->result[0];
    // Write out key
    state_key[tmpIdx] = state->key[0];
    // Write out tracker
    state_tracker[tmpIdx] = state->tracker;

    // For second out of four sets...
    // Write out counter
    tmpIdx += totalWorkItems;
    state_counter[tmpIdx] = state->counter[1];
    // Write out result
    state_result[tmpIdx] = state->result[1];
    // Write out key
    state_key[tmpIdx] = state->key[1];

    // For third out of four sets...
    // Write out counter
    tmpIdx += totalWorkItems;
    state_counter[tmpIdx] = state->counter[2];
    // Write out result
    state_result[tmpIdx] = state->result[2];
    // Write out key
    state_key[tmpIdx] = state->key[2];

    // For last out of four sets...
    // Write out counter
    tmpIdx += totalWorkItems;
    state_counter[tmpIdx] = state->counter[3];
    // Write out result
    state_result[tmpIdx] = state->result[3];
    // Write out key
    state_key[tmpIdx] = state->key[3];

}
