package engine

// support for running Go files as if they were mx3 files.

import (
	"flag"
	"github.com/seeder-research/uMagNUS/opencl"
	"github.com/seeder-research/uMagNUS/util"
	"os"
	"path"
)

var (
	// These flags are shared between cmd/uMagNUS and Go input files.
	Flag_cachedir    = flag.String("cache", os.TempDir(), "Kernel cache directory (empty disables caching)")
	Flag_gpu         = flag.Int("gpu", 0, "Specify GPU")
	Flag_interactive = flag.Bool("i", false, "Open interactive browser session")
	Flag_od          = flag.String("o", "", "Override output directory")
	Flag_port        = flag.String("http", ":35367", "Port to serve web gui")
	Flag_selftest    = flag.Bool("paranoid", false, "Enable convolution self-test for cuFFT sanity.")
	Flag_silent      = flag.Bool("s", false, "Silent") // provided for backwards compatibility
	Flag_sync        = flag.Bool("sync", false, "Synchronize all CUDA calls (debug)")
	Flag_forceclean  = flag.Bool("f", false, "Force start, clean existing output directory")
	Flag_failfast    = flag.Bool("failfast", false, "If one simulation fails, stop entire batch immediately")
	Flag_test        = flag.Bool("test", false, "OpenCL test (internal)")
	Flag_version     = flag.Bool("v", true, "Print version")
	Flag_vet         = flag.Bool("vet", false, "Check input files for errors, but don't run them")
)

// Usage: in every Go input file, write:
//
// 	func main(){
// 		defer InitAndClose()()
// 		// ...
// 	}
//
// This initialises the GPU, output directory, etc,
// and makes sure pending output will get flushed.
func InitAndClose() func() {

	flag.Parse()

	opencl.Init(*Flag_gpu)
	opencl.Synchronous = *Flag_sync

	od := *Flag_od
	if od == "" {
		od = path.Base(os.Args[0]) + ".out"
	}
	inFile := util.NoExt(od)
	InitIO(inFile, od, *Flag_forceclean)

	GoServe(*Flag_port)

	return func() {
		Close()
	}
}
