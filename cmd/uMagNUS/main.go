// uMagNUS main source
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	engine "github.com/seeder-research/uMagNUS/engine"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	script "github.com/seeder-research/uMagNUS/script"
	util "github.com/seeder-research/uMagNUS/util"
)

// flags in engine/gofiles.go
var ()

func main() {
	flag.Parse()
	log.SetPrefix("")
	log.SetFlags(0)

	opencl.Synchronous = *engine.Flag_sync
	opencl.Debug = *engine.Flag_debug

	// Check flag and initialize engine
	if len(*engine.Flag_gpulist) > 0 {
		var gpu_arr []int
		gpuList := strings.Split(*engine.Flag_gpulist, ",")
		if len(gpuList) == 0 {
			engine.Flag_gpu = 0
		} else {
			for _, item := range gpuList {
				if id, err := strconv.Atoi(item); err == nil {
					if id < 0 {
						log.Println("Invalid GPU number detected! Must be an integer >= 0!")
					} else {
						gpu_arr = append(gpu_arr, id)
					}
				}
			}
			if len(gpu_arr) == 0 {
				engine.Flag_gpu = 0
			} else {
				engine.Flag_gpu = gpu_arr[0]
			}
		}
	}
	if *engine.Flag_host {
		if engine.Flag_gpu < 0 {
			opencl.Init(engine.Flag_gpu)
		} else {
			log.Fatalln("Cannot disable GPU acceleration while requesting GPU \n")
		}
	} else {
		if engine.Flag_gpu < 0 {
			opencl.Init(0)
		} else {
			opencl.Init(engine.Flag_gpu)
		}
	}

	if *engine.Flag_version {
		printVersion()
	}

	// used by bootstrap launcher to test opencl
	// successful exit means opencl was initialized fine
	if *engine.Flag_test {
		fmt.Println(opencl.GPUInfo)
		os.Exit(0)
	}

	if len(opencl.ClDevices) <= 0 {
		fmt.Print("No OpenCL devices found \n")
		os.Exit(0)
	}

	defer engine.Close() // flushes pending output, if any

	if *engine.Flag_vet {
		vet()
		return
	}

	switch flag.NArg() {
	case 0:
		runInteractive()
	case 1:
		runFileAndServe(flag.Arg(0))
	default:
		RunQueue(flag.Args())
	}
}

func runInteractive() {
	fmt.Println("//no input files: starting interactive session")
	//initEngine()

	// setup outut dir
	now := time.Now()
	outdir := fmt.Sprintf("uMagNUS-%v-%02d-%02d_%02dh%02d.out", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute())
	engine.InitIO(outdir, outdir, *engine.Flag_forceclean)

	engine.Timeout = 365 * 24 * time.Hour // basically forever

	// set up some sensible start configuration
	engine.Eval(`SetGridSize(128, 64, 1)
		SetCellSize(4e-9, 4e-9, 4e-9)
		Msat = 1e6
		Aex = 10e-12
		alpha = 1
		m = RandomMag()`)
	addr := goServeGUI()
	openbrowser("http://127.0.0.1" + addr)
	engine.RunInteractive()
}

func runFileAndServe(fname string) {
	if path.Ext(fname) == ".go" {
		runGoFile(fname)
	} else {
		runScript(fname)
	}
}

func runScript(fname string) {
	outDir := util.NoExt(fname) + ".out"
	if *engine.Flag_od != "" {
		outDir = *engine.Flag_od
	}
	engine.InitIO(fname, outDir, *engine.Flag_forceclean)

	fname = engine.InputFile

	var code *script.BlockStmt
	var err2 error
	if fname != "" {
		// first we compile the entire file into an executable tree
		code, err2 = engine.CompileFile(fname)
		util.FatalErr(err2)
	}

	// now the parser is not used anymore so it can handle web requests
	goServeGUI()

	if *engine.Flag_interactive {
		openbrowser("http://127.0.0.1" + *engine.Flag_port)
	}

	// start executing the tree, possibly injecting commands from web gui
	engine.EvalFile(code)

	if *engine.Flag_interactive {
		engine.RunInteractive()
	}
}

func runGoFile(fname string) {

	// pass through flags
	flags := []string{"run", fname}
	flag.Visit(func(f *flag.Flag) {
		if f.Name != "o" {
			flags = append(flags, fmt.Sprintf("-%v=%v", f.Name, f.Value))
		}
	})

	if *engine.Flag_od != "" {
		flags = append(flags, fmt.Sprintf("-o=%v", *engine.Flag_od))
	}

	cmd := exec.Command("go", flags...)
	log.Println("go", flags)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		engine.Close()
		os.Exit(1)
	}
}

// start Gui server and return server address
func goServeGUI() string {
	if *engine.Flag_port == "" {
		log.Println(`//not starting GUI (-http="")`)
		return ""
	}
	addr := engine.GoServe(*engine.Flag_port)
	fmt.Print("//starting GUI at http://127.0.0.1", addr, "\n")
	return addr
}

// print version to stdout
func printVersion() {
	engine.LogOut("//", engine.UNAME, "\n")
	engine.LogOut("//", opencl.GPUInfo, "\n")
	engine.LogOut("//(c) Xuanyao Fong, SEEDER Research Group", "\n")
	engine.LogOut("//@ National University of Singapore, Singapore", "\n")
	engine.LogOut("//Web site: https://blog.nus.edu.sg/seeder", "\n")
	engine.LogOut("//Email: kelvin.xy.fong@nus.edu.sg", "\n")
	engine.LogOut("//Source code can be downloaded at https://github.com/seeder-research/uMagNUS", "\n")
	engine.LogOut("This is free software without any warranty. See license.txt")
	engine.LogOut("********************************************************************//")
	engine.LogOut("  If you use uMagNUS in any work or publication,                    //")
	engine.LogOut("  we kindly ask you to cite the references in references.bib        //")
	engine.LogOut("********************************************************************//")
	engine.LogOut("//uMagNUS is an OpenCL-based derivative of MuMax 3.10: (c) Arne Vansteenkiste, Dynamat LAB, Ghent University, Belgium", "\n")
}
