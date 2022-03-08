// uMagNUS main source
package main

import (
	//	"bufio"
	//	"bytes"
	"flag"
	"fmt"
	//	"github.com/seeder-research/uMagNUS/cl"
	//	"io"
	//	"log"
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
	if len(flag.Args()) == 0 {
		fmt.Println("No files given!")
	} else {
		for _, fname := range flag.Args() {
			fmt.Println("Processing file: ", fname)
			fcode := readFile(fname)
			if *Flag_verbose > 0 {
				fmt.Printf("%+v \n", fcode)
			}
		}
	}
}

// print version to stdout
func printVersion() {
}
