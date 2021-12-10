__kernel void square(__global real_t*  input, __global real_t* output, const unsigned int count) {
   int i = get_global_id(0);
   if (i < count) {
       output[i] = input[i] * input[i];
   }
}
