package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/seeder-research/uMagNUS/opencl"
	"github.com/seeder-research/uMagNUS/opencl/cl"
	"math/rand"
	"os"
	"time"
)

var (
	d_length      = flag.Int("size", 1024, "Total number of random numbers to generate")
	r_seed        = flag.Uint("seed", 0, "Seed value of RNG")
	d_dump        = flag.Bool("dump", false, "Whether to dump generated values to screen")
	Flag_platform = flag.Int("platform", 0, "Specify OpenCL platform")
	Flag_gpu      = flag.Int("gpu", 0, "Specify GPU")
)

func main() {

	flag.Parse()

	opencl.Init(*Flag_gpu)
	platforms := opencl.ClPlatforms
	fmt.Printf("Discovered platforms: \n")
	for i, p := range platforms {
		fmt.Printf("Platform %d: \n", i)
		fmt.Printf("  Name: %s \n", p.Name())
		fmt.Printf("  Vendor: %s \n", p.Vendor())
		fmt.Printf("  Profile: %s \n", p.Profile())
		fmt.Printf("  Version: %s \n", p.Version())
		fmt.Printf("  Extensions: %s \n", p.Extensions())
	}
	platform := opencl.ClPlatform
	fmt.Printf("In use: \n")
	fmt.Printf("  Vendor: %s \n", platform.Vendor())
	fmt.Printf("  Profile: %s \n", platform.Profile())
	fmt.Printf("  Version: %s \n", platform.Version())
	fmt.Printf("  Extensions: %s \n", platform.Extensions())

	fmt.Printf("Discovered devices: \n")
	devices := opencl.ClDevices
	deviceIndex := -1
	for i, d := range devices {
		if deviceIndex < 0 && d.Type() == cl.DeviceTypeGPU {
			deviceIndex = i
		}
		fmt.Printf("Device %d (%s): %s \n", i, d.Type(), d.Name())
		fmt.Printf("  Address Bits: %d \n", d.AddressBits())
		fmt.Printf("  Available: %+v \n", d.Available())
		fmt.Printf("  Compiler Available: %+v \n", d.CompilerAvailable())
		fmt.Printf("  Double FP Config: %s \n", d.DoubleFPConfig())
		fmt.Printf("  Driver Version: %s \n", d.DriverVersion())
		fmt.Printf("  Error Correction Supported: %+v \n", d.ErrorCorrectionSupport())
		fmt.Printf("  Execution Capabilities: %s \n", d.ExecutionCapabilities())
		fmt.Printf("  Extensions: %s \n", d.Extensions())
		fmt.Printf("  Global Memory Cache Type: %s \n", d.GlobalMemCacheType())
		fmt.Printf("  Global Memory Cacheline Size: %d KB \n", d.GlobalMemCachelineSize()/1024)
		fmt.Printf("  Global Memory Size: %d MB \n", d.GlobalMemSize()/(1024*1024))
		fmt.Printf("  Half FP Config: %s \n", d.HalfFPConfig())
		fmt.Printf("  Host Unified Memory: %+v \n", d.HostUnifiedMemory())
		fmt.Printf("  Image Support: %+v \n", d.ImageSupport())
		fmt.Printf("  Image2D Max Dimensions: %d x %d \n", d.Image2DMaxWidth(), d.Image2DMaxHeight())
		fmt.Printf("  Image3D Max Dimensions: %d x %d x %d \n", d.Image3DMaxWidth(), d.Image3DMaxHeight(), d.Image3DMaxDepth())
		fmt.Printf("  Little Endian: %+v \n", d.EndianLittle())
		fmt.Printf("  Local Mem Size Size: %d KB \n", d.LocalMemSize()/1024)
		fmt.Printf("  Local Mem Type: %s \n", d.LocalMemType())
		fmt.Printf("  Max Clock Frequency: %d \n", d.MaxClockFrequency())
		fmt.Printf("  Max Compute Units: %d \n", d.MaxComputeUnits())
		fmt.Printf("  Max Constant Args: %d \n", d.MaxConstantArgs())
		fmt.Printf("  Max Constant Buffer Size: %d KB \n", d.MaxConstantBufferSize()/1024)
		fmt.Printf("  Max Mem Alloc Size: %d KB \n", d.MaxMemAllocSize()/1024)
		fmt.Printf("  Max Parameter Size: %d \n", d.MaxParameterSize())
		fmt.Printf("  Max Read-Image Args: %d \n", d.MaxReadImageArgs())
		fmt.Printf("  Max Samplers: %d \n", d.MaxSamplers())
		fmt.Printf("  Max Work Group Size: %d \n", d.MaxWorkGroupSize())
		fmt.Printf("  Preferred Work Group Size: %d \n", opencl.ClPrefWGSz)
		fmt.Printf("  Max Work Item Dimensions: %d \n", d.MaxWorkItemDimensions())
		fmt.Printf("  Max Work Item Sizes: %d \n", d.MaxWorkItemSizes())
		fmt.Printf("  Max Write-Image Args: %d \n", d.MaxWriteImageArgs())
		fmt.Printf("  Memory Base Address Alignment: %d \n", d.MemBaseAddrAlign())
		fmt.Printf("  Native Vector Width Char: %d \n", d.NativeVectorWidthChar())
		fmt.Printf("  Native Vector Width Short: %d \n", d.NativeVectorWidthShort())
		fmt.Printf("  Native Vector Width Int: %d \n", d.NativeVectorWidthInt())
		fmt.Printf("  Native Vector Width Long: %d \n", d.NativeVectorWidthLong())
		fmt.Printf("  Native Vector Width Float: %d \n", d.NativeVectorWidthFloat())
		fmt.Printf("  Native Vector Width Double: %d \n", d.NativeVectorWidthDouble())
		fmt.Printf("  Native Vector Width Half: %d \n", d.NativeVectorWidthHalf())
		fmt.Printf("  OpenCL C Version: %s \n", d.OpenCLCVersion())
		fmt.Printf("  Profile: %s \n", d.Profile())
		fmt.Printf("  Profiling Timer Resolution: %d \n", d.ProfilingTimerResolution())
		fmt.Printf("  Vendor: %s \n", d.Vendor())
		fmt.Printf("  Version: %s \n", d.Version())
	}
	kernels := opencl.KernList

	fmt.Printf("Initializing THREEFRY RNG and generate uniformly distributed numbers... \n")

	seed := InitRNG()
	fmt.Println("Seed: ", seed)
	rng := opencl.NewGenerator("threefry")
	rng.Init(seed, nil)

	fmt.Printf("Creating output buffer... \n")
	d_size := int(*d_length)

	output := opencl.Buffer(1, [3]int{d_size, 1, 1})

	resultsSlice := output.HostCopy()
	resultsArr := resultsSlice.Host()

	if *d_dump {
		fmt.Println("Results before execution: ", resultsArr[0])
	}

	event := rng.Uniform(output.DevPtr(0), d_size, []*cl.Event{output.GetEvent(0)})
	err := cl.WaitForEvents([]*cl.Event{event})
	if err != nil {
		fmt.Printf("Uniform RN generation failed for output: %+v \n", err)
		return
	}

	resultsSlice = output.HostCopy()
	resultsArr = resultsSlice.Host()

	if *d_dump {
		fmt.Println("Results after execution: ", resultsArr[0])
	}

	fOut, fErr := os.Create("uniform_bytes.bin")
	if fErr != nil {
		panic(fErr)
	}

	wr := bufio.NewWriter(fOut)
	outBytes := new(bytes.Buffer)
	for _, v := range resultsArr[0] {
		vErr := binary.Write(outBytes, binary.LittleEndian, v)
		if vErr != nil {
			fmt.Println("binary.Write failed: ", vErr)
		}
	}

	nn, wrErr := wr.Write(outBytes.Bytes())
	if wrErr != nil {
		fmt.Println("bufio.Write failed: ", wrErr)
	} else {
		wr.Flush()
		fmt.Println("Wrote ", nn, "bytes to file")
	}
	fOut.Close()

	opencl.Recycle(output)

	fmt.Printf("Re-initializing THREEFRY RNG and generate normally distributed numbers... \n")

	rng.Init(seed, nil)

	output = opencl.Buffer(1, [3]int{d_size, 1, 1})

	resultsSlice = output.HostCopy()
	resultsArr = resultsSlice.Host()

	if *d_dump {
		fmt.Println("Results before execution: ", resultsArr[0])
	}

	event = rng.Normal(output.DevPtr(0), d_size, []*cl.Event{output.GetEvent(0)})
	err = cl.WaitForEvents([]*cl.Event{event})
	if err != nil {
		fmt.Printf("CreateBuffer failed for output: %+v \n", err)
		return
	}

	resultsSlice = output.HostCopy()
	resultsArr = resultsSlice.Host()

	if *d_dump {
		fmt.Println("Results after execution: ", resultsArr[0])
	}

	fOut, fErr = os.Create("norm_bytes.bin")
	if fErr != nil {
		panic(fErr)
	}

	wr = bufio.NewWriter(fOut)
	outBytes = new(bytes.Buffer)
	for _, v := range resultsArr[0] {
		vErr := binary.Write(outBytes, binary.LittleEndian, v)
		if vErr != nil {
			fmt.Println("binary.Write failed: ", vErr)
		}
	}

	nn, wrErr = wr.Write(outBytes.Bytes())
	if wrErr != nil {
		fmt.Println("bufio.Write failed: ", wrErr)
	} else {
		wr.Flush()
		fmt.Println("Wrote ", nn, "bytes to file")
	}

	fOut.Close()

	opencl.Recycle(output)
	fmt.Printf("Finished tests on THREEFRY RNG\n")

	fmt.Printf("freeing resources \n")
	for _, krn := range kernels {
		krn.Release()
	}

	opencl.ReleaseAndClean()
}

func InitRNG() uint64 {
	if *r_seed == (uint)(0) {
		rand.Seed(time.Now().UTC().UnixNano())
		return rand.Uint64()
	}
	return (uint64)(*r_seed)
}
