// Add magneto-elastic coupling field to B.
// H = - δUmel / δM, 
// where Umel is magneto-elastic energy denstiy given by the eq. (12.18) of Gurevich&Melkov "Magnetization Oscillations and Waves", CRC Press, 1996
__kernel void
addmagnetoelasticfield(__global float* __restrict  Bx, __global float* __restrict  By, __global float* __restrict  Bz,
                       __global float* __restrict  mx, __global float* __restrict  my, __global float* __restrict  mz,
					  __global float* __restrict exx_, float exx_mul,
					  __global float* __restrict eyy_, float eyy_mul,
					  __global float* __restrict ezz_, float ezz_mul,
					  __global float* __restrict exy_, float exy_mul,
					  __global float* __restrict exz_, float exz_mul,
					  __global float* __restrict eyz_, float eyz_mul,
					  __global float* __restrict B1_, float B1_mul, 
					  __global float* __restrict B2_, float B2_mul,
					  __global float* __restrict Ms_, float Ms_mul,
                      int N) {

	int I =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);

	if (I < N) {

	    float Exx = amul(exx_, exx_mul, I);
	    float Eyy = amul(eyy_, eyy_mul, I);
	    float Ezz = amul(ezz_, ezz_mul, I);
	    
	    float Exy = amul(exy_, exy_mul, I);
	    float Eyx = Exy;

	    float Exz = amul(exz_, exz_mul, I);
	    float Ezx = Exz;

	    float Eyz = amul(eyz_, eyz_mul, I);
	    float Ezy = Eyz;

		float invMs = inv_Msat(Ms_, Ms_mul, I);

		float B1 = amul(B1_, B1_mul, I) * invMs;
	    float B2 = amul(B2_, B2_mul, I) * invMs;

	    float3 m = {mx[I], my[I], mz[I]};

	    Bx[I] += -(2.0f*B1*m.x*Exx + B2*(m.y*Exy + m.z*Exz));
	    By[I] += -(2.0f*B1*m.y*Eyy + B2*(m.x*Eyx + m.z*Eyz));
	    Bz[I] += -(2.0f*B1*m.z*Ezz + B2*(m.x*Ezx + m.y*Ezy));
	}
}
