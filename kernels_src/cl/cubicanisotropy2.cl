// add cubic anisotropy field to B.
// B:      effective field in T
// m:      reduced magnetization (unit length)
// Ms:     saturation magnetization in A/m.
// K1:     Kc1 in J/m3
// K2:     Kc2 in T/m3
// C1, C2: anisotropy axes
//
// based on http://www.southampton.ac.uk/~fangohr/software/oxs_cubic8.html
__kernel void
addcubicanisotropy2(__global real_t* __restrict   Bx, __global real_t* __restrict      By, __global real_t* __restrict Bz,
                    __global real_t* __restrict   mx, __global real_t* __restrict      my, __global real_t* __restrict mz,
                    __global real_t* __restrict  Ms_,                      real_t  Ms_mul,
                    __global real_t* __restrict  k1_,                      real_t  k1_mul,
                    __global real_t* __restrict  k2_,                      real_t  k2_mul,
                    __global real_t* __restrict  k3_,                      real_t  k3_mul,
                    __global real_t* __restrict c1x_,                      real_t c1x_mul,
                    __global real_t* __restrict c1y_,                      real_t c1y_mul,
                    __global real_t* __restrict c1z_,                      real_t c1z_mul,
                    __global real_t* __restrict c2x_,                      real_t c2x_mul,
                    __global real_t* __restrict c2y_,                      real_t c2y_mul,
                    __global real_t* __restrict c2z_,                      real_t c2z_mul,
                                            int    N) {

    int   gid = get_global_id(0);
    int gsize = get_global_size(0);

    for (int i = gid; i < N; i += gsize) {

        real_t invMs = inv_Msat(Ms_, Ms_mul, i);
        real_t    k1 = amul(k1_, k1_mul, i);
        real_t    k2 = amul(k2_, k2_mul, i);
        real_t    k3 = amul(k3_, k3_mul, i);

        k1 *= invMs;
        k2 *= invMs;
        k3 *= invMs;

        real_t  u1x = (c1x_ == NULL) ? c1x_mul : (c1x_mul * c1x_[i]);
        real_t3  u1 = normalized(vmul(c1x_, c1y_, c1z_, c1x_mul, c1y_mul, c1z_mul, i));
        real_t3  u2 = normalized(vmul(c2x_, c2y_, c2z_, c2x_mul, c2y_mul, c2z_mul, i));
        real_t3  u3 = cross(u1, u2); // 3rd axis perpendicular to u1,u2
        real_t3   m = make_float3(mx[i], my[i], mz[i]);

        real_t u1m = dot(u1, m);
        real_t u2m = dot(u2, m);
        real_t u3m = dot(u3, m);

        real_t3 B = -2.0f*k1*((pow2(u2m) + pow2(u3m)) * (    (u1m) * u1) +
                              (pow2(u1m) + pow2(u3m)) * (    (u2m) * u2) +
                              (pow2(u1m) + pow2(u2m)) * (    (u3m) * u3))-
                     2.0f*k2*((pow2(u2m) * pow2(u3m)) * (    (u1m) * u1) +
                              (pow2(u1m) * pow2(u3m)) * (    (u2m) * u2) +
                              (pow2(u1m) * pow2(u2m)) * (    (u3m) * u3))-
                     4.0f*k3*((pow4(u2m) + pow4(u3m)) * (pow3(u1m) * u1) +
                              (pow4(u1m) + pow4(u3m)) * (pow3(u2m) * u2) +
                              (pow4(u1m) + pow4(u2m)) * (pow3(u3m) * u3));

        Bx[i] += B.x;
        By[i] += B.y;
        Bz[i] += B.z;
    }
}
