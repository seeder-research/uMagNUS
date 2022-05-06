// Package opencl provides GPU interaction
package opencl64

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	ld "github.com/seeder-research/uMagNUS/loader64"
)

type GPU struct {
	Platform *cl.Platform
	Device   *cl.Device
}

var (
	Version      string                    // OpenCL version
	DevName      string                    // GPU name
	TotalMem     int64                     // total GPU memory
	PlatformInfo string                    // Human-readable OpenCL platform description
	GPUInfo      string                    // Human-readable GPU description
	GPUList      []GPU                     // List of GPUs available
	Synchronous  bool                      // for debug: synchronize command queue at every kernel launch
	Debug        bool                      // for debug: synchronize command queue after every kernel launch
	ClPlatforms  []*cl.Platform            // list of platforms available
	ClPlatform   *cl.Platform              // platform the global OpenCL context is attached to
	ClDevices    []*cl.Device              // list of devices global OpenCL context may be associated with
	ClDevice     *cl.Device                // device associated with global OpenCL context
	ClCtx        *cl.Context               // global OpenCL context
	ClCmdQueue   *cl.CommandQueue          // command queue attached to global OpenCL context
	ClProgram    *cl.Program               // handle to program in the global OpenCL context
	KernList     = map[string]*cl.Kernel{} // Store pointers to all compiled kernels
	initialized  = false                   // Initial state defaults to false
	ClCUnits     int                       // Get number of compute units available
	ClWGSize     []int                     // Get maximum size of work group in each dimension
	ClPrefWGSz   int                       // Get preferred work group size of device
	ClMaxWGSize  int                       // Get maximum number of concurrent work-items that can execute simultaneously
	ClMaxWGNum   int                       // Get maximum number of max-sized work groups that can execute simultaneously
	ClTotalPE    int                       // Get total number of processing elements available
	GPUVend      int                       // 1: nvidia, 2: intel, 3: amd, 4: unknown
)

// Locks to an OS thread and initializes CUDA for that thread.
func Init(gpu int) {
	defer func() {
		initialized = true
	}()

	selection := int(0)

	if initialized {
		fmt.Printf("Already initialized \n")
		return // needed for tests
	}

	runtime.LockOSThread()

	// Attempt to get list of opencl platforms. Return if failed.
	var platforms []*cl.Platform
	var err error
	platforms, err = cl.GetPlatforms()
	if err != nil {
		fmt.Printf("Failed to get platforms: %+v \n", err)
		return
	}

	// Build list of opencl devices
	tmpClPlatforms := []*cl.Platform{}
	tmpGpuList := []GPU{}
	tmpClDevices := []*cl.Device{}
	for _, plat := range platforms {
		var pDevices []*cl.Device
		if gpu < 0 {
			pDevices, err = plat.GetDevices(cl.DeviceTypeCPU)
		} else {
			pDevices, err = plat.GetDevices(cl.DeviceTypeGPU)
		}
		if err != nil {
			fmt.Printf("Failed to get devices: %+v \n", err)
		}
		for idx, gpDev := range pDevices {
			if idx == 0 {
				tmpClPlatforms = append(tmpClPlatforms, plat)
			}
			// Only add devices that can support FP64 calculations
			if strings.Contains(gpDev.Extensions(), "cl_khr_fp64") || strings.Contains(gpDev.Extensions(), "cl_amd_fp64") {
				tmpGpuList = append(tmpGpuList, GPU{Platform: plat, Device: gpDev})
				tmpClDevices = append(tmpClDevices, gpDev)
			}
		}
	}

	// Check number of opencl devices detected.
	// Return if none found.
	// Otherwise, attempt to select desired opencl device.
	if len(tmpGpuList) == 0 {
		fmt.Printf("No devices found!\n")
		return
	} else {
		if gpu > len(tmpGpuList)-1 {
			fmt.Printf("Requested GPU: %+v ...\n    Unselectable GPU! Falling back to default selection\n", gpu)
		} else {
			if gpu >= 0 {
				selection = gpu
			}
		}
	}

	// Initialize the library with the selected opencl device
	GPUList = tmpGpuList
	ClDevices = tmpClDevices
	ClPlatforms = tmpClPlatforms
	selectedGPU := GPUList[selection]
	ClPlatform = selectedGPU.getGpuPlatform()
	ClDevice = selectedGPU.getGpuDevice()

	// Output information about platform of selected opencl device
	fmt.Printf("// GPU: %d\n", selection)
	PlatformName := ClPlatform.Name()
	PlatformVendor := ClPlatform.Vendor()
	PlatformProfile := ClPlatform.Profile()
	PlatformVersion := ClPlatform.Version()
	PlatformInfo = fmt.Sprint("//   Platform Name: ", PlatformName, "\n//   Vendor: ", PlatformVendor, "\n//   Profile: ", PlatformProfile, "\n//   Version: ", PlatformVersion, "\n")

	// Output information about selected opencl device
	DevName = ClDevice.Name()
	TotalMem = ClDevice.GlobalMemSize()
	Version = ClDevice.OpenCLCVersion()
	GPUInfo = fmt.Sprint("OpenCL C Version ", Version, "\n// GPU: ", DevName, "(", (TotalMem)/(1024*1024), "MB) \n")

	// Create opencl context on selected device
	var context *cl.Context
	context, err = cl.CreateContext([]*cl.Device{ClDevice})
	if err != nil {
		fmt.Printf("CreateContext failed: %+v \n", err)
		return
	}

	// Create opencl command queue on selected device
	var queue *cl.CommandQueue
	queue, err = context.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("CreateCommandQueue failed: %+v \n", err)
		return
	}

	// Create opencl program on selected opencl device
	var program *cl.Program
	nobinary := bool(false)

	// Attempt to obtain binary from library. Compile from source if unable to...
	programBytes := ld.GetClDeviceBinary(ClDevice)
	if programBytes == nil {
		fmt.Println("Unable to get program binary!")
		nobinary = true
	} else {
		if program, err = context.CreateProgramWithBinary([]*cl.Device{ClDevice}, []int{len(programBytes)}, [][]byte{programBytes}); err != nil {
			fmt.Printf("Unable to load binary from library...continuing to compile code \n")
			nobinary = true
		}
	}

	// Unable to load kernel from binary. Compile kernel source code instead
	if nobinary {
		if program, err = context.CreateProgramWithSource([]string{GenMergedKernelSource()}); err != nil {
			fmt.Printf("CreateProgramWithSource failed: %+v \n", err)
			return
		}

		// Attempt to build binary from opencl program
		if err = program.BuildProgram([]*cl.Device{ClDevice}, "-cl-std=CL1.2 -cl-finite-math-only -cl-no-signed-zeros -cl-fp32-correctly-rounded-divide-sqrt -cl-kernel-arg-info -D__REAL_IS_DOUBLE__"); err != nil {
			fmt.Printf("BuildProgram failed: %+v \n", err)
			return
		}
	}

	// Attempt to build list of kernels in opencl program
	completed := bool(true)
	if kernelsString, errK := program.GetKernelNames(); errK == nil {
		kernelNamesArray := strings.Split(kernelsString, ";")
		for _, kernname := range kernelNamesArray {
			KernList[kernname], err = program.CreateKernel(kernname)
			if err != nil {
				fmt.Printf("CreateKernel failed: %+v \n", err)
				completed = false
			}
		}
	} else {
		fmt.Printf("Unable to get list of kernels in program: %+v \n", errK)
		return
	}
	if completed != true {
		fmt.Println("Unable to completely build map of kernels!")
		return
	}

	ClCtx = context
	ClCmdQueue = queue
	ClProgram = program

	// Set basic configuration for distributing
	// work-items across compute units
	ClCUnits = ClDevice.MaxComputeUnits()
	ClWGSize = ClDevice.MaxWorkItemSizes()
	ClMaxWGSize = ClDevice.MaxWorkGroupSize()

	nvRegExp := regexp.MustCompile("(?i)nvidia")
	inRegExp := regexp.MustCompile("(?i)intel")
	adRegExp0 := regexp.MustCompile("(?i)amd")
	adRegExp1 := regexp.MustCompile("(?i)micro device")
	if chk0 := nvRegExp.Match([]byte(GPUInfo)); chk0 {
		GPUVend = 1
	} else {
		if chk1 := inRegExp.Match([]byte(GPUInfo)); chk1 {
			GPUVend = 2
		} else {
			chk2, chk3 := adRegExp0.Match([]byte(GPUInfo)), adRegExp1.Match([]byte(GPUInfo))
			if (chk2 == true) || (chk3 == true) {
				GPUVend = 3
			} else {
				GPUVend = 4
			}
		}
	}
	ClMaxWGNum = ClCUnits
	if GPUVend == 1 { // Nvidia
		ClTotalPE = ClWGSize[2] * ClCUnits
		if ClMaxWGSize > ClTotalPE {
			ClMaxWGNum = ClTotalPE / ClMaxWGSize
		} else {
			ClMaxWGNum = 1
			ClMaxWGSize = ClTotalPE
		}
	}
	if GPUVend == 2 { // Intel
		ClMaxWGSize = 7 * 32
		ClTotalPE = ClMaxWGNum * ClMaxWGSize
	}

	ClPrefWGSz, err = KernList["madd2"].PreferredWorkGroupSizeMultiple(ClDevice)
	if err != nil {
		fmt.Printf("PreferredWorkGroupSizeMultiple failed: %+v \n", err)
	}

	config1DSize = ClTotalPE

	// Reduce kernel launch parameters are updated on update to mesh size
	reduceSingleSize = 16 * 2 * ClPrefWGSz
	reducecfg.Grid[0] = reduceSingleSize
	reducecfg.Block[0] = reduceSingleSize
	reduceintcfg.Grid[0] = config1DSize
	reduceintcfg.Block[0] = reduceSingleSize

	data.EnableGPU(memFree, memFree, MemCpy, MemCpyDtoH, MemCpyHtoD)

}

func (s *GPU) getGpuDevice() *cl.Device {
	return s.Device
}

func (s *GPU) getGpuPlatform() *cl.Platform {
	return s.Platform
}

func ReleaseAndClean() {
	ClCmdQueue.Release()
	ClProgram.Release()
	ClCtx.Release()
}
