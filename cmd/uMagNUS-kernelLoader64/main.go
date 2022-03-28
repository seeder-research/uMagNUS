// uMagNUS main source
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/seeder-research/uMagNUS/cl"
	ld "github.com/seeder-research/uMagNUS/loader64"
)

var (
	Flag_gpu     = flag.Int("gpu", -5, "gpu to use")
	Flag_verbose = flag.Int("v", 0, "Verbosity level")
)

func main() {
	flag.Parse()
	// Find available GPUs
	InitGPUs()
	if *Flag_verbose > 0 {
		fmt.Println("Number of GPUs found: ", len(GPUList))
	}

	if *Flag_gpu < 0 {
		fmt.Println("No GPU selected. Exiting...")
		return
	}
	var err error
	gpuNum := int(0)
	if *Flag_gpu < len(GPUList) {
		gpuNum = *Flag_gpu
	}
	ClPlatform = GPUList[gpuNum].Platform
	ClDevice = GPUList[gpuNum].Device
	if ClCtx, err = cl.CreateContext([]*cl.Device{ClDevice}); err != nil {
		fmt.Println("Unable to create context on target device!")
		return
	}
	if *Flag_verbose > 0 {
		fmt.Printf("GPU Name: %+v \n", ClDevice.Name())
	}
	programBytes := ld.GetClDeviceBinary(ClDevice)
	if programBytes == nil {
		fmt.Println("Unable to get program binary!")
		return
	}
	var ClProgram *cl.Program
	if ClProgram, err = ClCtx.CreateProgramWithBinary([]*cl.Device{ClDevice}, []int{len(programBytes)}, [][]byte{programBytes}); err != nil {
		fmt.Println("Unable to create program from binary on context!")
		ClCtx.Release()
		return
	}
	ShowBuildLog(ClProgram, ClDevice)
	var binType cl.ProgramBinaryTypes
	if binType, err = ClProgram.GetProgramBinaryType(ClDevice); err != nil {
		fmt.Println("    Error getting binary type for program on GPU.")
		ClCtx.Release()
		log.Panic(err)
	}
	if *Flag_verbose > 1 {
		switch binType {
		case cl.ProgramBinaryTypeNone:
			fmt.Println("      No compiled binaries available in program.")
		case cl.ProgramBinaryTypeCompiledObject:
			fmt.Println("      Compiled object available in program.")
		case cl.ProgramBinaryTypeLibrary:
			fmt.Println("      Compiled library available in program.")
		case cl.ProgramBinaryTypeExecutable:
			fmt.Println("      Compiled executable available in program.")
		default:
			fmt.Println("      Unknown binary type in program.")
		}
	}
	var kernNum int
	if kernNum, err = ClProgram.GetKernelCounts(); err != nil {
		fmt.Println("    Error getting kernel count for linked program on GPU.")
		ClProgram.Release()
		ClCtx.Release()
		log.Panic(err)
	}
	if *Flag_verbose > 0 {
		fmt.Printf("()  Number of kernels in program: %+v \n", kernNum)
	}
	var kernNames string
	if kernNames, err = ClProgram.GetKernelNames(); err != nil {
		fmt.Println("    Error getting kernel names for linked program on GPU.")
		ClProgram.Release()
		ClCtx.Release()
		log.Panic(err)
	}
	kernNameArray := strings.Split(kernNames, ";")
	if *Flag_verbose > 0 {
		fmt.Println("**  Kernels in program:")
		for _, kn := range kernNameArray {
			fmt.Println("()    kernel: ", kn)
		}
	}
	ClProgram.Release()
	ClCtx.Release()
}
