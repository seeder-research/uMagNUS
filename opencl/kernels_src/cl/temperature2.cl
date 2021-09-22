// TODO: this could act on x,y,z, so that we need to call it only once.
__kernel void
settemperature2(__global float* __restrict  B,      __global float* __restrict noise, float kB2_VgammaDt,
                __global float* __restrict Ms_, float Ms_mul,
                __global float* __restrict temp_, float temp_mul,
                __global float* __restrict alpha_, float alpha_mul,
                int N) {

    int i =  ( get_group_id(1)*get_num_groups(0) + get_group_id(0) ) * get_local_size(0) + get_local_id(0);
    if (i < N) {
        float invMs = inv_Msat(Ms_, Ms_mul, i);
        float temp = amul(temp_, temp_mul, i);
        float alpha = amul(alpha_, alpha_mul, i);
        B[i] = noise[i] * sqrt((kB2_VgammaDt * alpha * temp * invMs ));
    }
}

