/**
@file
Implements threefry RNG.
/*******************************************************
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
        uint2 counter;
        uint2 result;
        uint2 key;
        uint tracker;
} threefry_state;
**/


__kernel void threefry_uint32(__global uint *state_counter, __global uint *state_result, __global uint *state_key, __global uint *state_tracker, uint* output, uint rng_count){
    uint index = get_group_id(0) * ELEMENTS_PER_BLOCK + get_local_id(0);
    threefry_state state;

    // Read in counter
    state.counter.x = state_counter[index];
    state.counter.y = state_counter[index + rng_count];

    // Read in result
    state.result.x = state_result[index];
    state.result.y = state_result[index + rng_count];

    // Read in key
    state.key.x = state_key[index];
    state.key.y = state_key[index + rng_count];

    // Read in tracker
    state.tracker.x = state_tracker[index];

    if (state.tracker == 1) {
        uint tmp = state.result.y;
        state.counter.x += index;
        state.counter.y += (state.counter.y < index);
        threefry_round(state);
        state.tracker = 0;
        return tmp;
    } else {
        state->tracker++;
        return state.result.x;
    }
}
