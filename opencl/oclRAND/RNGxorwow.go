package oclRAND

import (
        "github.com/seeder-research/uMagNUS/opencl/cl"
)

type XORWOW_status_array_ptr struct {
        Ini         bool
        Seed_Arr    []uint64
        Status_buf  *cl.MemObject
        Status_size int
}

