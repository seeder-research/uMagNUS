// Package opencl provides GPU interaction
package opencl

import (
	"fmt"
	"runtime"

	"github.com/seeder-research/uMagNUS/data"
	"github.com/seeder-research/uMagNUS/opencl/cl"
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
	Synchronous  bool                      // for debug: synchronize stream0 at every kernel launch
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
	ClWGSize     int                       // Get maximum size of work group per compute unit
	ClPrefWGSz   int                       // Get preferred work group size of device
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
	platforms, err := cl.GetPlatforms()
	tmpClPlatforms := []*cl.Platform{}
	tmpGpuList := []GPU{}
	tmpClDevices := []*cl.Device{}

	if err != nil {
		fmt.Printf("Failed to get platforms: %+v \n", err)
		return
	}
	for _, plat := range platforms {
		var pDevices []*cl.Device
		pDevices, err = plat.GetDevices(cl.DeviceTypeGPU)
		if err != nil {
			fmt.Printf("Failed to get devices: %+v \n", err)
		}
		for ii, gpDev := range pDevices {
			if ii == 0 {
				tmpClPlatforms = append(tmpClPlatforms, plat)
			}
			tmpGpuList = append(tmpGpuList, GPU{Platform: plat, Device: gpDev})
			tmpClDevices = append(tmpClDevices, gpDev)
		}
	}
	if len(tmpGpuList) == 0 {
		fmt.Printf("No devices found!\n")
		return
	} else {
		if gpu > len(tmpGpuList)-1 {
			fmt.Printf("Unselectable GPU! Falling back to default selection\n")
		} else {
			selection = gpu
		}
	}

	GPUList = tmpGpuList
	ClDevices = tmpClDevices
	ClPlatforms = tmpClPlatforms
	selectedGPU := GPUList[selection]
	ClPlatform = selectedGPU.getGpuPlatform()
	ClDevice = selectedGPU.getGpuDevice()

	fmt.Printf("// GPU: %d\n", selection)
	PlatformName := ClPlatform.Name()
	PlatformVendor := ClPlatform.Vendor()
	PlatformProfile := ClPlatform.Profile()
	PlatformVersion := ClPlatform.Version()
	PlatformInfo = fmt.Sprint("//   Platform Name: ", PlatformName, "\n//   Vendor: ", PlatformVendor, "\n//   Profile: ", PlatformProfile, "\n//   Version: ", PlatformVersion, "\n")

	DevName = ClDevice.Name()
	TotalMem = ClDevice.GlobalMemSize()
	Version = ClDevice.OpenCLCVersion()
	GPUInfo = fmt.Sprint("OpenCL C Version ", Version, "\n// GPU: ", DevName, "(", (TotalMem)/(1024*1024), "MB) \n")
	context, err := cl.CreateContext([]*cl.Device{ClDevice})
	if err != nil {
		fmt.Printf("CreateContext failed: %+v \n", err)
	}
	queue, err := context.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("CreateCommandQueue failed: %+v \n", err)
	}
	program, err := context.CreateProgramWithSource([]string{GenMergedKernelSource()})
	if err != nil {
		fmt.Printf("CreateProgramWithSource failed: %+v \n", err)
	}
	if err := program.BuildProgram([]*cl.Device{ClDevice}, "-cl-std=CL1.2 -cl-fp32-correctly-rounded-divide-sqrt -cl-kernel-arg-info"); err != nil {
		fmt.Printf("BuildProgram failed: %+v \n", err)
	}

	for _, kernname := range KernelNames {
		KernList[kernname], err = program.CreateKernel(kernname)
		if err != nil {
			fmt.Printf("CreateKernel failed: %+v \n", err)
		}
	}
	ClCtx = context
	ClCmdQueue = queue
	ClProgram = program
	// Set basic configuration for distributing
	// work-items across compute units
	ClCUnits = ClDevice.MaxComputeUnits()
	ClWGSize = ClDevice.MaxWorkGroupSize()
	reducecfg.Grid[0] = ClWGSize
	reducecfg.Block[0] = ClWGSize
	reduceintcfg.Grid[0] = ClWGSize * ClCUnits
	reduceintcfg.Block[0] = ClWGSize
	ClPrefWGSz, err = KernList["madd2"].PreferredWorkGroupSizeMultiple(ClDevice)
	if err != nil {
		fmt.Printf("PreferredWorkGroupSizeMultiple failed: %+v \n", err)
	}

	data.EnableGPU(memFree, memFree, MemCpy, MemCpyDtoH, MemCpyHtoD)

	//	fmt.Printf("Initializing clFFT library \n")
	//	if err := cl.SetupCLFFT(); err != nil {
	//		fmt.Printf("failed to initialize clFFT \n")
	//	}
}

func (s *GPU) getGpuDevice() *cl.Device {
	return s.Device
}

func (s *GPU) getGpuPlatform() *cl.Platform {
	return s.Platform
}

func ReleaseAndClean() {
	//	cl.TeardownCLFFT()
	ClCmdQueue.Release()
	ClProgram.Release()
	ClCtx.Release()
}

// Global stream used for everything
//const stream0 = cu.Stream(0)

// Synchronize the global stream
// This is called before and after all memcopy operations between host and device.
//func Sync() {
//	stream0.Synchronize()
//}
