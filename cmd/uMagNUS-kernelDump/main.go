package main

import (
	"fmt"
	"flag"
	"os"

	k32 "github.com/seeder-research/uMagNUS/kernels"
	k64 "github.com/seeder-research/uMagNUS/kernels64"
)

var (
        Flag_k64    = flag.Bool("k64", false, "Output kernels for 64-bit.")
	Flag_dump   = flag.Bool("dump", false, "Dump to screen.")
	Flag_file   = flag.String("o", "", "Filename to output string to.")
)

func main() {
	flag.Parse()

	var outString string

	if *Flag_k64 {
		outString = k64.OpenclProgramSource()
	} else {
		outString = k32.OpenclProgramSource()
	}

	if *Flag_dump {
		fmt.Printf("%+v \n", outString)
	}

	if len(*Flag_file) > 0 {
		fmt.Printf("Outputting to file %+v \n", *Flag_file)
		file, err := os.Create(*Flag_file)
		if err != nil {
			fmt.Printf("ERROR: unable to open file for output! \n")
			panic(err)
		}
		defer file.Close()
		file.WriteString(outString)
		fmt.Printf("Done.\n")
	}

	fmt.Printf("Exiting...\n")
}
