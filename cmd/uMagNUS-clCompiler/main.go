// uMagNUS main source
package main

import (
	//	"bufio"
	//	"bytes"
	"flag"
	"fmt"
	"github.com/seeder-research/uMagNUS/cl"
	//	"io"
	"log"
	//	"os"
	//	"path"
	//	"strconv"
	"strings"
	//	"text/scanner"
	//	"time"
)

var (
	Flag_defines  = flag.String("dopts", "", "-D arguments to pass to compiler")
	Flag_includes = flag.String("iopts", "", "-I arguments to pass to compiler")
	Flag_libpaths = flag.String("lopts", "", "-L arguments to pass to compiler")
	Flag_libs     = flag.String("libs", "", "-l arguments to pass to compiler")
	Flag_linkopts = flag.String("link", "", "arguments to pass to linker")
	Flag_ClStd    = flag.String("std", "", "-std argument to pass to compiler")
	Flag_ComArgs  = flag.String("args", "", "Other arguments to pass to compiler")
	Flag_verbose  = flag.Int("v", 0, "Verbosity level")
)

func main() {
	flag.Parse()
	fmt.Println("Compiler options: ", generateCompilerOpts())
	fmt.Println("Linker options: ", generateLinkerOpts())
	var fcode string
	fcode = ""
	if len(flag.Args()) == 0 {
		fmt.Println("No files given!")
		return
	} else {
		for _, fname := range flag.Args() {
			fmt.Println("Processing file: ", fname)
			fcode = readFile(fname)
			if *Flag_verbose > 6 {
				fmt.Printf("%+v \n", fcode)
			}
		}
	}
	// Find available GPUs
	InitGPUs()
	if *Flag_verbose > 0 {
		fmt.Println("Number of GPUs found: ", len(GPUList))
	}
	for gpuId, _ := range GPUList {
		var gpuArg []*cl.Device
		gpuArg = append(gpuArg, GPUList[gpuId].Device)
		if *Flag_verbose > 2 {
			fmt.Println("    Creating context on GPU: ", gpuId)
		}
		tmpContext, err := cl.CreateContext(gpuArg)
		if err != nil {
			fmt.Println("    Error creating context on GPU.")
			log.Panic(err)
		}

		if *Flag_verbose > 2 {
			fmt.Println("      Create and compile program on GPU: ", gpuId)
		}
		tmpProgram, err := tmpContext.CreateProgramWithSource([]string{fcode})
		if err != nil {
			fmt.Println("    Error creatinging program on GPU.")
			tmpContext.Release()
			log.Panic(err)
		}
		ShowBuildLog(tmpProgram, GPUList[gpuId].Device)
		var binType cl.ProgramBinaryTypes
		binType, err = tmpProgram.GetProgramBinaryType(GPUList[gpuId].Device)
		if err != nil {
			fmt.Println("    Error getting binary type for program on GPU.")
			tmpContext.Release()
			log.Panic(err)
		}
		if *Flag_verbose > 2 {
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

		fmt.Println("    Building program...")
		buildOpts := generateCompilerOpts()
		if buildOpts != "" {
			buildOpts = buildOpts + " "
		}
		buildOpts = buildOpts + generateLinkerOpts()
		fmt.Println("        using options: ", buildOpts)
		err = tmpProgram.BuildProgram([]*cl.Device{GPUList[gpuId].Device}, buildOpts)
                if err != nil {
                        fmt.Println("    Error building binary for program on GPU.")
                        tmpProgram.Release()
                        tmpContext.Release()
                        log.Panic(err)
                }
                ShowBuildLog(tmpProgram, GPUList[gpuId].Device)
                binType, err = tmpProgram.GetProgramBinaryType(GPUList[gpuId].Device)
                if err != nil {
                        fmt.Println("    Error getting binary type for program on GPU.")
                        tmpContext.Release()
                        log.Panic(err)
                }
                if *Flag_verbose > 2 {
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
		kernNum, err = tmpProgram.GetKernelCounts()
		if err != nil {
			fmt.Println("    Error getting kernel count for linked program on GPU.")
			tmpProgram.Release()
			tmpContext.Release()
			log.Panic(err)
		}
		if *Flag_verbose > 2 {
			fmt.Printf("    Program has %+v number of kernels\n", kernNum)
		}
		var kernNames string
		kernNames, err = tmpProgram.GetKernelNames()
                if err != nil {
                        fmt.Println("    Error getting kernel names for linked program on GPU.")
                        tmpProgram.Release()
                        tmpContext.Release()
                        log.Panic(err)
                }
		kernNameArray := strings.Split(kernNames, ";")
		fmt.Println("  Kernels in program:")
		for _, kn := range kernNameArray {
			fmt.Println("    kernel: ", kn)
		}

		var binSizes []int
		binSizes, err = tmpProgram.GetBinarySizes()
                if err != nil {
                        fmt.Println("    Error getting binary sizes program on GPU.")
                        tmpProgram.Release()
                        tmpContext.Release()
                        log.Panic(err)
                }

		fmt.Println("  Number of program binaries: ", len(binSizes))
		for idx, binSz := range binSizes {
			fmt.Printf("    Size of program binary number %+v: %+v bytes\n", idx+1, binSz)
		}

		var bins []*uint8
		bins, err = tmpProgram.GetBinaries()
                if err != nil {
                        fmt.Println("    Error getting binaries for program on GPU.")
                        tmpProgram.Release()
                        tmpContext.Release()
                        log.Panic(err)
                }

		fmt.Printf("      Check of number of binaries: %+v \n", len(bins))

		if *Flag_verbose > 2 {
			fmt.Println("    Releasing program on GPU: ", gpuId)
		}
		tmpProgram.Release()
		if *Flag_verbose > 2 {
			fmt.Println("    Releasing context on GPU: ", gpuId)
		}
		tmpContext.Release()
	}
}

// print version to stdout
func printVersion() {
}
