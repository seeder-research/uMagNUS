// See maxangle.go for more details.
__kernel void
setmaxangle(__global float* __restrict dst,
            __global float* __restrict mx, __global float* __restrict my, __global float* __restrict mz,
            __global float* __restrict aLUT2d, __global uint8_t* __restrict regions,
            int Nx, int Ny, int Nz, uint8_t PBC) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
	int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

	if (ix >= Nx || iy >= Ny || iz >= Nz) {
		return;
	}

	// central cell
	int I = idx(ix, iy, iz);
	float3 m0 = make_float3(mx[I], my[I], mz[I]);

	if (is0(m0)) {
		return;
	}

	uint8_t r0 = regions[I];
	float angle  = 0.0f;

	int i_;    // neighbor index
	float3 m_; // neighbor mag
	float a__; // inter-cell exchange stiffness

	// left neighbor
	i_  = idx(lclampx(ix-1), iy, iz);           // clamps or wraps index according to PBC
	m_  = make_float3(mx[i_], my[i_], mz[i_]);  // load m
	m_  = ( is0(m_)? m0: m_ );                  // replace missing non-boundary neighbor
	a__ = aLUT2d[symidx(r0, regions[i_])];
	if (a__ != 0) {
		angle = max(angle, acos(dot(m_,m0)));
	}

	// right neighbor
	i_  = idx(hclampx(ix+1), iy, iz);
	m_  = make_float3(mx[i_], my[i_], mz[i_]);
	m_  = ( is0(m_)? m0: m_ );
	a__ = aLUT2d[symidx(r0, regions[i_])];
	if (a__ != 0) {
		angle = max(angle, acos(dot(m_,m0)));
	}

	// back neighbor
	i_  = idx(ix, lclampy(iy-1), iz);
	m_  = make_float3(mx[i_], my[i_], mz[i_]);
	m_  = ( is0(m_)? m0: m_ );
	a__ = aLUT2d[symidx(r0, regions[i_])];
	if (a__ != 0) {
		angle = max(angle, acos(dot(m_,m0)));
	}

	// front neighbor
	i_  = idx(ix, hclampy(iy+1), iz);
	m_  = make_float3(mx[i_], my[i_], mz[i_]);
	m_  = ( is0(m_)? m0: m_ );
	a__ = aLUT2d[symidx(r0, regions[i_])];
	if (a__ != 0) {
		angle = max(angle, acos(dot(m_,m0)));
	}

	// only take vertical derivative for 3D sim
	if (Nz != 1) {
		// bottom neighbor
		i_  = idx(ix, iy, lclampz(iz-1));
		m_  = make_float3(mx[i_], my[i_], mz[i_]);
		m_  = ( is0(m_)? m0: m_ );
		a__ = aLUT2d[symidx(r0, regions[i_])];
		if (a__ != 0) {
			angle = max(angle, acos(dot(m_,m0)));
		}

		// top neighbor
		i_  = idx(ix, iy, hclampz(iz+1));
		m_  = make_float3(mx[i_], my[i_], mz[i_]);
		m_  = ( is0(m_)? m0: m_ );
		a__ = aLUT2d[symidx(r0, regions[i_])];
		if (a__ != 0) {
			angle = max(angle, acos(dot(m_,m0)));
		}
	}

	dst[I] = angle;
}

