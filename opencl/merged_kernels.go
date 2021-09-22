package opencl

import (
	"github.com/seeder-research/uMagNUS/opencl/kernels"
)

func GenMergedKernelSource() string {
	return kernels.OpenclProgramSource()
}
