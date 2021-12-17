package opencl64

import (
	kernels "github.com/seeder-research/uMagNUS/kernels64"
)

func GenMergedKernelSource() string {
	return kernels.OpenclProgramSource()
}
