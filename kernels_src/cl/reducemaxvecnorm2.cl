__kernel void
reducemaxvecnorm2(__global real_t* __restrict       x,
                  __global real_t* __restrict       y,
                  __global real_t* __restrict       z,
                  __global real_t* __restrict     dst,
                           real_t             initVal,
                              int                   n,
                  __local  real_t*            scratch) {

    // Initialize memory
    int global_idx = get_global_id(0);
    int  local_idx = get_local_id(0);
    real_t currVal = initVal;

    // Loop over input elements in chunks and accumulate each chunk into local memory
    while (global_idx < n) {
        real_t element = (x[global_idx]*x[global_idx]) + (y[global_idx]*y[global_idx]) + (z[global_idx]*z[global_idx]);
        currVal = fmax(currVal, element);
        global_idx += get_global_size(0);
    }

    // At this point, accumulated values on chunks are in local memory. Perform parallel reduction
    scratch[local_idx] = currVal;
    // Add barrier to sync all threads
    barrier(CLK_LOCAL_MEM_FENCE);

    for (int offset = (get_local_size(0) >> 1); offset > 1; offset >>= 1) {
        if (local_idx < offset) {
            real_t other = scratch[local_idx + offset];
            real_t  mine = scratch[local_idx];
            scratch[local_idx] = fmax(mine, other);
        }
        // barrier for syncing work group
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    if (local_idx == 0) {
        dst[get_group_id(0)] = fmax(scratch[0], scratch[1]);
    }
}
