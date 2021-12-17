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
State of threefry RNG. We will store in global buffer as a set of ulong
**
counter: ulong[4]
key:     ulong[4]
state:   ulong[4]
index:   ulong
typedef struct{
        ulong counter[4];
        ulong result[4];
        ulong key[4];
        ulong tracker;
} threefry64_state;
**/

/**
Seeds threefry RNG.
@param state Variable, that holds state of the generator to be seeded.
@param seed Value used for seeding. Should be randomly generated for each instance of generator (thread).
**/
#if defined(__REAL_IS_DOUBLE__)
__kernel void
threefry64_seed(__global ulong* __restrict state_key,
                __global ulong* __restrict state_counter,
                __global ulong* __restrict state_result,
                __global ulong* __restrict state_tracker,
                __global ulong* __restrict seed) {
    ulong gid = get_global_id(0);
    ulong rng_count = get_global_size(0);
    ulong idx = gid;
    ulong localJ = seed[gid];
    state_key[idx] = localJ;
    state_counter[idx] = 0x0000000000000000;
    state_result[idx] = 0x0000000000000000;
    state_tracker[idx] = 4;
    idx += rng_count;
    state_key[idx] = 0x0000000000000000;
    state_counter[idx] = 0x0000000000000000;
    state_result[idx] = 0x0000000000000000;
    idx += rng_count;
    state_key[idx] = gid;
    state_counter[idx] = 0x0000000000000000;
    state_result[idx] = 0x0000000000000000;
    idx += rng_count;
    state_key[idx] = 0x0000000000000000;
    state_counter[idx] = 0x0000000000000000;
    state_result[idx] = 0x0000000000000000;
}
#endif // __REAL_IS_DOUBLE__
