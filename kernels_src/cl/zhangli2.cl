#define PREFACTOR ((MUB) / (2 * QE * GAMMA0))

// spatial derivatives without dividing by cell size
#define deltax(in) (in[idx(hclampx(ix+1), iy, iz)] - in[idx(lclampx(ix-1), iy, iz)])
#define deltay(in) (in[idx(ix, hclampy(iy+1), iz)] - in[idx(ix, lclampy(iy-1), iz)])
#define deltaz(in) (in[idx(ix, iy, hclampz(iz+1))] - in[idx(ix, iy, lclampz(iz-1))])

__kernel void
addzhanglitorque2(__global real_t* __restrict     tx, __global real_t* __restrict        ty, __global real_t* __restrict tz,
                  __global real_t* __restrict     mx, __global real_t* __restrict        my, __global real_t* __restrict mz,
                  __global real_t* __restrict    Ms_,                      real_t    Ms_mul,
                  __global real_t* __restrict    jx_,                      real_t    jx_mul,
                  __global real_t* __restrict    jy_,                      real_t    jy_mul,
                  __global real_t* __restrict    jz_,                      real_t    jz_mul,
                  __global real_t* __restrict alpha_,                      real_t alpha_mul,
                  __global real_t* __restrict    xi_,                      real_t    xi_mul,
                  __global real_t* __restrict   pol_,                      real_t   pol_mul,
                                       real_t     cx,                      real_t        cy,                      real_t cz,
                                          int     Nx,                        int        Ny,                        int Nz,
                                      uint8_t    PBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    int i = idx(ix, iy, iz);

    real_t  alpha = amul(alpha_, alpha_mul, i);
    real_t     xi = amul(xi_, xi_mul, i);
    real_t    pol = amul(pol_, pol_mul, i);
    real_t  invMs = inv_Msat(Ms_, Ms_mul, i);
    real_t      b = invMs * PREFACTOR / (1.0f + xi*xi);
    real_t3  Jvec = vmul(jx_, jy_, jz_, jx_mul, jy_mul, jz_mul, i);
    real_t3     J = pol*Jvec;
    real_t3 hspin = make_float3(0.0f, 0.0f, 0.0f); // (u·∇)m

    if (J.x != 0.0f) {
        hspin += (b/cx)*J.x * make_float3(deltax(mx), deltax(my), deltax(mz));
    }
    if (J.y != 0.0f) {
        hspin += (b/cy)*J.y * make_float3(deltay(mx), deltay(my), deltay(mz));
    }
    if (J.z != 0.0f) {
        hspin += (b/cz)*J.z * make_float3(deltaz(mx), deltaz(my), deltaz(mz));
    }

    real_t3      m = make_float3(mx[i], my[i], mz[i]);
    real_t3 torque = (-1.0f/(1.0f + alpha*alpha)) * (
                         (1.0f+xi*alpha) * cross(m, cross(m, hspin))
                         +(  xi-alpha) * cross(m, hspin)           );

    // write back, adding to torque
    tx[i] += torque.x;
    ty[i] += torque.y;
    tz[i] += torque.z;
}
