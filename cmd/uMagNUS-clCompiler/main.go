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
	//	"strings"
	//	"text/scanner"
	//	"time"
)

var (
	Flag_defines  = flag.String("dopts", "", "-D arguments to pass to compiler")
	Flag_includes = flag.String("iopts", "", "-I arguments to pass to compiler")
	Flag_libpaths = flag.String("lopts", "", "-L arguments to pass to compiler")
	Flag_libs     = flag.String("libs", "", "-l arguments to pass to compiler")
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
	} else {
		for _, fname := range flag.Args() {
			fmt.Println("Processing file: ", fname)
			fcode := readFile(fname)
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
		tmpProgram, err := compileProgram(tmpContext, gpuArg, []string{fcode})
		if err != nil {
			fmt.Println("    Error creating and compiling program on GPU.")
			tmpContext.Release()
			log.Panic(err)
		}

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
