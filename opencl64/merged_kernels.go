package opencl64

import (
	"github.com/seeder-research/uMagNUS/kernels"
)

func GenMergedKernelSource() string {
	return kernels.OpenclProgramSource()
}
