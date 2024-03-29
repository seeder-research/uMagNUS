// Returns the topological charge contribution on an elementary triangle ijk
// Order of arguments is important here to preserve the same measure of chirality
// Note: the result is zero if an argument is zero, or when two arguments are the same
static inline real_t triangleCharge(real_t3 mi, real_t3 mj, real_t3 mk) {
    real_t numer   = dot(mi, cross(mj, mk));
    real_t denom   = 1.0f + dot(mi, mj) + dot(mi, mk) + dot(mj, mk);
    return 2.0f * atan2(numer, denom);
}

// Set s to the toplogogical charge density for lattices based on the solid angle 
// subtended by triangle associated with three spins: a,b,c
//
//       s = 2 atan[(a . b x c /(1 + a.b + a.c + b.c)] / (dx dy)
//
// After M Boettcher et al, New J Phys 20, 103014 (2018), adapted from
// B. Berg and M. Luescher, Nucl. Phys. B 190, 412 (1981), and implemented by
// Joo-Von Kim.
//
// A unit cell comprises two triangles, but s is a site-dependent quantity so we
// double-count and average over four triangles.
__kernel void
settopologicalchargelattice(__global real_t* __restrict     s,
                            __global real_t* __restrict    mx, __global real_t* __restrict my, __global real_t* __restrict mz,
                                                 real_t icxcy,
                                                    int    Nx,                         int Ny,                         int Nz,
                                                uint8_t   PBC) {

    int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
    int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);
    int iz = get_group_id(2) * get_local_size(2) + get_local_id(2);

    if ((ix >= Nx) || (iy >= Ny) || (iz >= Nz)) {
        return;
    }

    int     i0 = idx(ix, iy, iz);                     // central cell index
    real_t3 m0 = make_float3(mx[i0], my[i0], mz[i0]); // central cell magnetization

    if (is0(m0)) {
        s[i0] = 0.0f;
        return;
    }

    // indices of the 4 neighbors (counter clockwise)
    int i1 = idx(hclampx(ix+1), iy, iz); // (i+1,j)
    int i2 = idx(ix, hclampy(iy+1), iz); // (i,j+1)
    int i3 = idx(lclampx(ix-1), iy, iz); // (i-1,j)
    int i4 = idx(ix, lclampy(iy-1), iz); // (i,j-1)

    // magnetization of the 4 neighbors
    real_t3 m1 = make_float3(mx[i1], my[i1], mz[i1]);
    real_t3 m2 = make_float3(mx[i2], my[i2], mz[i2]);
    real_t3 m3 = make_float3(mx[i3], my[i3], mz[i3]);
    real_t3 m4 = make_float3(mx[i4], my[i4], mz[i4]);

    // local topological charge (accumulator)
    real_t topcharge = 0.0f;

    // charge contribution from the upper right triangle
    // if diagonally opposite neighbor is not zero, use a weight of 1/2 to avoid counting charges twice
    if (((ix+1<Nx) || PBCx) && ((iy+1<Ny) || PBCy)) { 
        int         i_ = idx(hclampx(ix+1), hclampy(iy+1), iz); // diagonal opposite neighbor in upper right quadrant
        real_t3     m_ = make_float3(mx[i_], my[i_], mz[i_]);
        real_t  weight = is0(m_) ? (real_t)1.0 : (real_t)0.5;
        topcharge     += weight * triangleCharge(m0, m1, m2);
    }

    // upper left
    if (((ix-1>=0) || PBCx) && ((iy+1<Ny) || PBCy)) { 
        int         i_ = idx(lclampx(ix-1), hclampy(iy+1), iz); 
        real_t3     m_ = make_float3(mx[i_], my[i_], mz[i_]);
        real_t  weight = is0(m_) ? (real_t)1.0 : (real_t)0.5;
        topcharge     += weight * triangleCharge(m0, m2, m3);
    }

    // bottom left
    if (((ix-1>=0) || PBCx) && ((iy-1>=0) || PBCy)) { 
        int         i_ = idx(lclampx(ix-1), lclampy(iy-1), iz); 
        real_t3     m_ = make_float3(mx[i_], my[i_], mz[i_]);
        real_t  weight = is0(m_) ? (real_t)1.0 : (real_t)0.5;
        topcharge     += weight * triangleCharge(m0, m3, m4);
    }

    // bottom right
    if (((ix+1<Nx) || PBCx) && ((iy-1>=0) || PBCy)) { 
        int         i_ = idx(hclampx(ix+1), lclampy(iy-1), iz); 
        real_t3     m_ = make_float3(mx[i_], my[i_], mz[i_]);
        real_t  weight = is0(m_) ? (real_t)1.0 : (real_t)0.5;
        topcharge     += weight * triangleCharge(m0, m4, m1);
    }

    s[i0] = icxcy * topcharge;
}
