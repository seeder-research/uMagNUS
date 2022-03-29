package opencl

import (
	kernels "github.com/seeder-research/uMagNUS/kernels"
)

func GenMergedKernelSource() string {
	return kernels.OpenclProgramSource()
}
