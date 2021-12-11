__kernel void
reducemaxabs(__global real_t* __restrict src, __global real_t* __restrict dst, real_t initVal, int n, __local real_t* scratch) {
    // Initialize memory
    int global_idx = get_global_id(0);
    int  local_idx = get_local_id(0);
    real_t currVal = initVal;

    // Loop over input elements in chunks and store max in each chunk into local memory
    while (global_idx < n) {
        real_t element = fabs(src[global_idx]);
        currVal = fmax(currVal, element);
        global_idx += get_global_size(0);
    }

    // At this point, max values on chunks are in local memory. Perform parallel reduction
    scratch[local_idx] = currVal;
    // Add barrier to sync all threads
    barrier(CLK_LOCAL_MEM_FENCE);

    for (int offset = get_local_size(0) / 2; offset > 0; offset = offset / 2) {
        if (local_idx < offset) {
            real_t other = scratch[local_idx + offset];
            real_t  mine = scratch[local_idx];
            scratch[local_idx] = fmax(mine, other);
        }
        // barrier for syncing work group
        barrier(CLK_LOCAL_MEM_FENCE);
    }

    if (local_idx == 0) {
        dst[get_group_id(0)] = scratch[0];
    }
}