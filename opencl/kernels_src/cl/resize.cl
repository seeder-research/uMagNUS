// Select and resize one layer for interactive output
__kernel void
resize(__global float* __restrict  dst, int Dx, int Dy, int Dz,
       __global float* __restrict  src, int Sx, int Sy, int Sz,
       int layer, int scalex, int scaley) {

	int ix = get_group_id(0) * get_local_size(0) + get_local_id(0);
	int iy = get_group_id(1) * get_local_size(1) + get_local_id(1);

	if (ix<Dx && iy<Dy) {

		float sum = 0.0f;
		float n = 0.0f;

		for(int J=0; J<scaley; J++) {
			int j2 = iy*scaley+J;

			for(int K=0; K<scalex; K++) {
				int k2 = ix*scalex+K;

				if (j2 < Sy && k2 < Sx) {
					sum += src[(layer*Sy + j2)*Sx + k2];
					n += 1.0f;
				}
			}
		}
		dst[iy*Dx + ix] = sum / n;
	}
}

