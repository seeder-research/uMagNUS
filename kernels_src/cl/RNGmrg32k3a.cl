/******************************************************************************
 * @file     prngcl_mrg32k3a.cl
 * @author   Vadim Demchik <vadimdi@yahoo.com>
 * @version  1.1.2
 *
 * @brief    [PRNGCL library]
 *           contains OpenCL implementation of MRG32k3a pseudo-random number generator
 *
 *
 * @section  CREDITS
 *
 *   Pierre L'Ecuyer,
 *   "Good Parameter Sets for Combined Multiple Recursive Random Number Generators",
 *   Operations Research, 47, 1 (1999), 159--164.
 *
 *
 * @section  LICENSE
 *
 * Copyright (c) 2013-2015 Vadim Demchik
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted provided that the following conditions are met:
 *
 *    Redistributions of source code must retain the above copyright notice,
 *      this list of conditions and the following disclaimer.
 *
 *    Redistributions in binary form must reproduce the above copyright notice,
 *      this list of conditions and the following disclaimer in the documentation
 *      and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
 * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 * 
 *****************************************************************************/

__kernel void
mrg32k3a(__global uint4* seed_table, 
         __global float4* randoms,
            const uint N)
{
    uint giddst = GID;
    float rnd;
    float4 result;

    uint4 seed1 = seed_table[GID];
    uint4 seed2 = seed_table[GID + GID_SIZE];

    for (uint i = 0; i < N; i++) {
        mrg32k3a_step(&seed1,&seed2,&rnd);
            result.x = rnd;
        mrg32k3a_step(&seed1,&seed2,&rnd);
            result.y = rnd;
        mrg32k3a_step(&seed1,&seed2,&rnd);
            result.z = rnd;
        mrg32k3a_step(&seed1,&seed2,&rnd);
            result.w = rnd;
        randoms[giddst] = result;
        giddst += GID_SIZE;
    }
    seed_table[GID] = seed1;
    seed_table[GID + GID_SIZE] = seed2;
}
