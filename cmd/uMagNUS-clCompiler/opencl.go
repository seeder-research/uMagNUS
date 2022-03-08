package main

import (
	"fmt"
	"github.com/seeder-research/uMagNUS/cl"
)

type GPU struct {
	Platform *cl.Platform
	Device   *cl.Device
}

var (
	Version      string         // OpenCL version
	DevName      string         // GPU name
	TotalMem     int64          // total GPU memory
	PlatformInfo string         // Human-readable OpenCL platform description
	GPUInfo      string         // Human-readable GPU description
	GPUList      []GPU          // List of GPUs available
	ClPlatforms  []*cl.Platform // list of platforms available
	ClPlatform   *cl.Platform   // platform the global OpenCL context is attached to
	ClDevices    []*cl.Device   // list of devices global OpenCL context may be associated with
	ClDevice     *cl.Device     // device associated with global OpenCL context
	ClCtx        *cl.Context    // global OpenCL context
)

func InitGPUs() {
	if *Flag_verbose > 0 {
		fmt.Println("Initializing OpenCL...")
	}
	platforms, err := cl.GetPlatforms()
	tmpClPlatforms := []*cl.Platform{}
	tmpGpuList := []GPU{}
	tmpClDevices := []*cl.Device{}

	if err != nil {
		fmt.Printf("Failed to get platforms: %+v \n", err)
		return
	}
	if *Flag_verbose > 0 {
		fmt.Println("  Scanning platforms for GPUs...")
	}
	for _, plat := range platforms {
		var pDevices []*cl.Device
		pDevices, err = plat.GetDevices(cl.DeviceTypeGPU)
		if err != nil {
			fmt.Printf("Failed to get devices: %+v \n", err)
		}
		if *Flag_verbose > 0 {
			fmt.Println("    Scanning found GPUs...")
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
	}

	GPUList = tmpGpuList
}

func InitCPUs() {
	if *Flag_verbose > 0 {
		fmt.Println("Initializing OpenCL...")
	}
	platforms, err := cl.GetPlatforms()
	tmpClPlatforms := []*cl.Platform{}
	tmpGpuList := []GPU{}
	tmpClDevices := []*cl.Device{}

	if err != nil {
		fmt.Printf("Failed to get platforms: %+v \n", err)
		return
	}
	if *Flag_verbose > 0 {
		fmt.Println("  Scanning platforms for CPUs...")
	}
	for _, plat := range platforms {
		var pDevices []*cl.Device
		pDevices, err = plat.GetDevices(cl.DeviceTypeCPU)
		if err != nil {
			fmt.Printf("Failed to get devices: %+v \n", err)
		}
		if *Flag_verbose > 0 {
			fmt.Println("    Scanning found CPUs...")
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
	}

	GPUList = tmpGpuList
}
