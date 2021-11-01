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
	n_cycles      = flag.Int("cycles", 5, "Total number of random number generation cycles")
	r_seed        = flag.Uint("seed", 0, "Seed value of RNG")
	d_dump        = flag.Bool("dump", false, "Whether to dump generated values to screen")
	d_norm        = flag.Bool("norm", false, "Whether to generate normally distributed numbers")
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

	kernels := opencl.KernList

	fmt.Printf("Initializing MTGP RNG and generate uniformly distributed numbers... \n")

	seed := InitRNG()
	fmt.Println("Seed: ", seed)
	rng := opencl.NewGenerator("mtgp")
	rng.Init(&seed, nil)

	fmt.Printf("Creating output buffer... \n")
	d_size := int(*d_length)

	output := opencl.Buffer(1, [3]int{d_size, 1, 1})

	resultsSlice := output.HostCopy()
	resultsArr := resultsSlice.Host()

	if *d_dump {
		fmt.Println("Results before execution: ", resultsArr[0])
	}

	for idx := 0; idx < *n_cycles; idx++ {
		if *d_norm {
			rng.Normal(output.DevPtr(0), d_size, []*cl.Event{output.GetEvent(0)})
		} else {
			rng.Uniform(output.DevPtr(0), d_size, []*cl.Event{output.GetEvent(0)})
		}
	}

	resultsSlice = output.HostCopy()
	resultsArr = resultsSlice.Host()

	if *d_dump {
		fmt.Println("Results after execution: ", resultsArr[0])
	}

        var fOut *os.File
        var fErr error

        if *d_norm {
        	fOut, fErr = os.Create("norm_bytes.bin")
        } else {
        	fOut, fErr = os.Create("uniform_bytes.bin")
        }

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
	fmt.Printf("Finished tests on MTGP RNG\n")

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
