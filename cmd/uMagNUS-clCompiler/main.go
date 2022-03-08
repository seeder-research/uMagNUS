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
)

func main() {
	flag.Parse()
	fmt.Println(generateCompilerOpts())
}

// print version to stdout
func printVersion() {
}
