// uMagNUS main source
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	cl "github.com/seeder-research/uMagNUS/cl"
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
	Flag_dump     = flag.Bool("dump", false, "dump C code to screen")
	Flag_outfile  = flag.String("outfile", "", "output file to print C code to")
)

func main() {
	flag.Parse()
	if *Flag_verbose > 0 {
		fmt.Println("Compiler options: ", generateCompilerOpts())
		fmt.Println("Linker options: ", generateLinkerOpts())
	}
	var fcode string
	fcode = ""
	if len(flag.Args()) == 0 {
		fmt.Println("No files given!")
		return
	} else {
		for _, fname := range flag.Args() {
			if *Flag_verbose > 1 {
				fmt.Println("Processing file: ", fname)
			}
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

	// Generate header of output C code as a string
	outcode := printHeader()

	if len(gpuIdMap) > 0 {
		binariesMap = make(map[string]*cl.ProgramBinaries)
	} else {
		fmt.Println("No GPUs available...exiting")
		return
	}

	gpuNameList := string("")
	binSizeList := string("")
	outBinIdx := string("")
	mergedHexString := string("")
	hexPtrsString := string("")
	binIdx := 0

	for gpuName, gpuId := range gpuIdMap {
		if len(gpuNameList) > 0 {
			gpuNameList += ";"
		}
		gpuNameList += gpuName
		if len(outBinIdx) > 0 {
			outBinIdx += ", "
		}
		outBinIdx += strconv.Itoa(binIdx)
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
			fmt.Println("    Error creating program on GPU.")
			tmpContext.Release()
			log.Panic(err)
		}
		if *Flag_verbose > 0 {
			ShowBuildLog(tmpProgram, GPUList[gpuId].Device)
		}
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

		if *Flag_verbose > 1 {
			fmt.Println("    Building program...")
		}
		buildOpts := generateCompilerOpts()
		if buildOpts != "" {
			buildOpts = buildOpts + " "
		}
		currPlatform := GPUList[gpuId].Device.Platform()
		currPlatformInfo := fmt.Sprint(currPlatform.Name(), " ", currPlatform.Vendor(), " ", currPlatform.Profile(), " ", currPlatform.Version())
		if strings.Contains(strings.ToUpper(currPlatformInfo), "NVIDIA") {
			buildOpts += string(" -D__NVCODE__ ")
		} else {
			if strings.EqualFold(DevName, "gfx908") {
				buildOpts += fmt.Sprint(buildOpts, " -D__AMDGPU_FP32ATOMICS_1__ ")
			}
			if strings.EqualFold(DevName, "gfx90a") {
				buildOpts += fmt.Sprint(buildOpts, " -D__AMDGPU_FP32ATOMICS_1__ -D__AMDGPU_FP64ATOMICS_0__ ")
			}
			if strings.EqualFold(DevName, "gfx940") {
				buildOpts += fmt.Sprint(buildOpts, " -D__AMDGPU_FP32ATOMICS_0__ -D__AMDGPU_FP64ATOMICS_0__ ")
			}
		}
		buildOpts = buildOpts + generateLinkerOpts()
		if *Flag_verbose > 1 {
			fmt.Println("        using options: ", buildOpts)
		}
		err = tmpProgram.BuildProgram([]*cl.Device{GPUList[gpuId].Device}, buildOpts)
		if err != nil {
			fmt.Println("    Error building binary for program on GPU.")
			tmpProgram.Release()
			tmpContext.Release()
			log.Panic(err)
		}
		if *Flag_verbose > 0 {
			ShowBuildLog(tmpProgram, GPUList[gpuId].Device)
		}
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
		if *Flag_verbose > 1 {
			fmt.Println("()  Kernels in program:")
			for _, kn := range kernNameArray {
				fmt.Println("()    kernel: ", kn)
			}
		}

		var binSizes []int
		binSizes, err = tmpProgram.GetBinarySizes()
		if err != nil {
			fmt.Println("    Error getting binary sizes program on GPU.")
			tmpProgram.Release()
			tmpContext.Release()
			log.Panic(err)
		}

		if *Flag_verbose > 1 {
			fmt.Println("  Number of program binaries: ", len(binSizes))
			for idx, binSz := range binSizes {
				fmt.Printf("    Size of program binary number %+v: %+v bytes\n", idx+1, binSz)
			}
		}

		var bins *cl.ProgramBinaries
		bins, err = tmpProgram.GetBinaries()
		if err != nil {
			fmt.Println("    Error getting binaries for program on GPU.")
			tmpProgram.Release()
			tmpContext.Release()
			log.Panic(err)
		}

		binariesMap[gpuName] = bins
		binsArrays := bins.GetBinaryArray()
		binsArraysPtrs := bins.GetBinaryArrayPointers()
		if *Flag_verbose > 1 {
			fmt.Printf("      Check of number of binaries: %+v ; Expected %+v \n", len(binsArraysPtrs), len(binSizes))
		}
		if len(binsArraysPtrs) == len(binSizes) {
			for idx, binSz := range binSizes {
				if *Flag_verbose > 1 {
					fmt.Printf("        Check binary size: %+v \n", len(binsArrays[idx]))
					fmt.Printf("              Expect size: %+v \n", binSz)
				}
			}
		} else {
			fmt.Println("   Expected binary sizes do not match!")
		}

		if *Flag_verbose > 2 {
			fmt.Println("****  Attempting to create program by loading binary...")
		}

		tmpProgram.Release()
		tmpProgram, err = tmpContext.CreateProgramWithBinary([]*cl.Device{GPUList[gpuId].Device}, bins.GetBinarySizes(), binsArrays)
		if err != nil {
			if *Flag_verbose > 2 {
				fmt.Printf("**** CreateProgramWithBinary(): failed to create program with binary: %+v \n", err)
			}
			tmpContext.Release()
		} else {
			if *Flag_verbose > 3 {
				fmt.Println("****    Successfully created program with binary")
			}
		}

		if *Flag_verbose > 0 {
			ShowBuildLog(tmpProgram, GPUList[gpuId].Device)
		}
		binType, err = tmpProgram.GetProgramBinaryType(GPUList[gpuId].Device)
		if err != nil {
			fmt.Println("    Error getting binary type for program on GPU.")
			tmpContext.Release()
			log.Panic(err)
		} else {
			if *Flag_verbose > 3 {
				fmt.Println("****    Attempted to show build log")
			}
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

		kernNum = int(0)
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
		kernNames = string("")
		kernNames, err = tmpProgram.GetKernelNames()
		if err != nil {
			fmt.Println("    Error getting kernel names for linked program on GPU.")
			tmpProgram.Release()
			tmpContext.Release()
			log.Panic(err)
		}
		kernNameArray = strings.Split(kernNames, ";")
		if *Flag_verbose > 1 {
			fmt.Println("**  Kernels in program:")
			for _, kn := range kernNameArray {
				fmt.Println("**    kernel: ", kn)
			}
		}

		if *Flag_verbose > 2 {
			fmt.Println("    Releasing program on GPU: ", gpuId)
		}
		tmpProgram.Release()
		if *Flag_verbose > 2 {
			fmt.Println("    Releasing context on GPU: ", gpuId)
		}
		tmpContext.Release()
		for _, val := range bins.GetBinarySizes() {
			if len(binSizeList) > 0 {
				binSizeList += ", "
			}
			binSizeList += strconv.Itoa(val)
		}
		hexString := printHex(binsArrays[0])
		mergedHexString += "const char hexString" + strconv.Itoa(binIdx) + "[] = " + hexString
		if len(hexPtrsString) > 0 {
			hexPtrsString += ", "
		}
		hexPtrsString += "&hexString" + strconv.Itoa(binIdx) + "[0]"
		binIdx++
	}
	outcode += "extern const char deviceNames[] = \"" + gpuNameList + "\";\n"
	outcode += "extern const size_t deviceNameLen = " + strconv.Itoa(len(gpuNameList)) + ";\n"
	outcode += "extern const int NumDevices = " + strconv.Itoa(len(gpuIdMap)) + ";\n"
	outcode += "extern const size_t binIdx[] = { " + outBinIdx + "};\n"
	outcode += "extern const size_t binSizes[] = { " + binSizeList + "};\n"
	outcode += "\n" + mergedHexString + "\n"
	outcode += "extern const char * hexPtrs[] = { " + hexPtrsString + "};\n"
	if *Flag_dump {
		fmt.Printf("%+v\n", outcode)
	}
}
