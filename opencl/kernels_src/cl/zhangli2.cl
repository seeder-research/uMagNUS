#define PREFACTOR ((MUB) / (2 * QE * GAMMA0))

// spatial derivatives without dividing by cell size
#define deltax(in) (in[idx(hclampx(ix+1), iy, iz)] - in[idx(lclampx(ix-1), iy, iz)])
#define deltay(in) (in[idx(ix, hclampy(iy+1), iz)] - in[idx(ix, lclampy(iy-1), iz)])
#define deltaz(in) (in[idx(ix, iy, hclampz(iz+1))] - in[idx(ix, iy, lclampz(iz-1))])

__kernel void
addzhanglitorque2(__global float* __restrict tx, __global float* __restrict ty, __global float* __restrict tz,
                  __global float* __restrict mx, __global float* __restrict my, __global float* __restrict mz,
                  __global float* __restrict Ms_, float Ms_mul,
                  __global float* __restrict jx_, float jx_mul,
                  __global float* __restrict jy_, float jy_mul,
                  __global float* __restrict jz_, float jz_mul,
                  __global float* __restrict alpha_, float alpha_mul,
                  __global float* __restrict xi_, float xi_mul,
                  __global float* __restrict pol_, float pol_mul,
                  float cx, float cy, float cz,
                  int Nx, int Ny, int Nz, uint8_t PBC) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if (ix >= Nx || iy >= Ny || iz >= Nz) {
        return;
    }

    int i = idx(ix, iy, iz);

    float alpha = amul(alpha_, alpha_mul, i);
    float xi    = amul(xi_, xi_mul, i);
    float pol   = amul(pol_, pol_mul, i);
	float invMs = inv_Msat(Ms_, Ms_mul, i);
    float b = invMs * PREFACTOR / (1.0f + xi*xi);
	float3 Jvec = vmul(jx_, jy_, jz_, jx_mul, jy_mul, jz_mul, i);
    float3 J = pol*Jvec;

    float3 hspin = make_float3(0.0f, 0.0f, 0.0f); // (u·∇)m
    if (J.x != 0.0f) {
        hspin += (b/cx)*J.x * make_float3(deltax(mx), deltax(my), deltax(mz));
    }
    if (J.y != 0.0f) {
        hspin += (b/cy)*J.y * make_float3(deltay(mx), deltay(my), deltay(mz));
    }
    if (J.z != 0.0f) {
        hspin += (b/cz)*J.z * make_float3(deltaz(mx), deltaz(my), deltaz(mz));
    }

    float3 m      = make_float3(mx[i], my[i], mz[i]);
    float3 torque = (-1.0f/(1.0f + alpha*alpha)) * (
                        (1.0f+xi*alpha) * cross(m, cross(m, hspin))
                        +(  xi-alpha) * cross(m, hspin)           );

    // write back, adding to torque
    tx[i] += torque.x;
    ty[i] += torque.y;
    tz[i] += torque.z;
}

