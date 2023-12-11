//go:build ignore
// +build ignore

// This program generates Go wrappers for opencl sources.
// The opencl file should contain exactly one __kernel void.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/scanner"
	"text/template"

	util "github.com/seeder-research/uMagNUS/util"
)

var (
	Flag_outdir = flag.String("outdir", ".", "directory to output kernel wrappers to")
	Flag_fp64   = flag.Bool("double", false, "compile for fp64 kernels")
	Flag_lib64  = flag.Bool("lib64", false, "compile for opencl64 package")

	RealType     = string("")
	SIZEOF_REALS = int(0)
)

func main() {
	flag.Parse()
	if *Flag_lib64 {
		RealType = "float64"
	} else {
		RealType = "float32"
	}
	tm["float"] = RealType
	tm["double"] = RealType
	tm["real_t"] = RealType
	for _, fname := range flag.Args() {
		ocl2go(fname)
	}
}

// generate opencl wrapper for file.
func ocl2go(fname string) {
	// open opencl file
	f, err := os.Open(fname)
	util.PanicErr(err)
	defer f.Close()

	// read tokens
	var token []string
	var s scanner.Scanner
	s.Init(f)
	tok := s.Scan()
	// Go through the file and filters out specific words, storing
	// the other words in sequence. Preserves code
	for tok != scanner.EOF {
		if !filter(s.TokenText()) {
			token = append(token, s.TokenText())
		}
		tok = s.Scan()
	}

	// find function name and arguments
	funcname := ""
	argstart, argstop := -1, -1
	for i := 0; i < len(token); i++ {
		// Looks for "__kernel" first. At this point, we
		// know which of the recorded words is the name
		// for the kernel and we can start extracting the
		// arguments for the kernel. Hint: if the current
		// word is "__kernel", the next word must be the
		// return type (per C function definition syntax)
		// followed by the function (kernel in this case)
		// name, and then the list of input arguments in
		// parenthese
		if token[i] == "__kernel" {
			funcname = token[i+2]
			argstart = i + 4
		}
		if argstart > 0 && token[i] == ")" {
			argstop = i + 1
			break
		}
	}
	argl := token[argstart:argstop]

	// isolate individual arguments
	var args [][]string
	start := 0
	for i, a := range argl {
		if a == "," || a == ")" {
			args = append(args, argl[start:i])
			start = i + 1
		}
	}

	// separate arg names/types and make pointers Go-style
	argn := make([]string, len(args))
	argt := make([]string, len(args))
	setn := make([]string, len(args))
	for i := range args {
		// Scan through the argument to locate "__global",
		// "__constant", "__local" and "__private"
		var currarg []string
		setn[i] = ""
		for j, txt := range args[i] {
			flag := chkArgMemType(txt)
			if flag == 2 {
				setn[i] = "__local"
			} else if flag == 0 {
				currarg = append(currarg, args[i][j])
			}
		}
		if currarg[1] == "*" {
			currarg = []string{currarg[0] + "*", currarg[2]}
		}
		argt[i] = typemap(currarg[0])
		argn[i] = currarg[1]
	}
	var argtt, argnn []string
	for i, a := range argn {
		if setn[i] == "__local" {
			if *Flag_lib64 {
				setn[i] = "KernList[\"" + funcname + "\"].SetArgUnsafe(" + strconv.Itoa(i) + ", cfg.Block[0]*cfg.Block[1]*cfg.Block[2]*SIZEOF_FLOAT64, nil)"
			} else {
				setn[i] = "KernList[\"" + funcname + "\"].SetArgUnsafe(" + strconv.Itoa(i) + ", cfg.Block[0]*cfg.Block[1]*cfg.Block[2]*SIZEOF_FLOAT32, nil)"
			}
		} else {
			argnn = append(argnn, argn[i])
			argtt = append(argtt, argt[i])
			setn[i] = "SetKernelArgWrapper(\"" + funcname + "\", " + strconv.Itoa(i) + ", " + a + ")"
		}
	}
	wrapgen(fname, funcname, argtt, argnn, setn)
}

// translate C type to Go type.
func typemap(ctype string) string {
	if gotype, ok := tm[ctype]; ok {
		return gotype
	}
	panic(fmt.Errorf("unsupported OpenCL type: %v", ctype))
}

var tm = map[string]string{
	"float*":    "unsafe.Pointer",
	"float2*":   "unsafe.Pointer",
	"float3*":   "unsafe.Pointer",
	"float4*":   "unsafe.Pointer",
	"double*":   "unsafe.Pointer",
	"double2*":  "unsafe.Pointer",
	"double3*":  "unsafe.Pointer",
	"double4*":  "unsafe.Pointer",
	"int":       "int",
	"real_t*":   "unsafe.Pointer",
	"real_t2*":  "unsafe.Pointer",
	"real_t3*":  "unsafe.Pointer",
	"real_t4*":  "unsafe.Pointer",
	"uint*":     "unsafe.Pointer",
	"uint4*":    "unsafe.Pointer",
	"uint8*":    "unsafe.Pointer",
	"uint8_t*":  "unsafe.Pointer",
	"uint16*":   "unsafe.Pointer",
	"uint16_t*": "unsafe.Pointer",
	"uint32*":   "unsafe.Pointer",
	"uint32_t*": "unsafe.Pointer",
	"uint":      "uint32",
	"uint8":     "uint8",
	"uint16":    "uint16",
	"uint32":    "uint32",
	"uint8_t":   "uint8",
	"uint16_t":  "uint16",
	"uint32_t":  "uint32",
	"ulong*":    "unsafe.Pointer",
	"ulong":     "uint64",
}

// template data
type Kernel struct {
	Name  string
	TType string
	ArgT  []string
	ArgN  []string
	SetN  []string
}

var ls []string

// generate wrapper code from template
func wrapgen(filename, funcname string, argt, argn, setn []string) {
	ttype := string("")
	if *Flag_lib64 {
		ttype = string("64")
	}
	kernel := &Kernel{funcname, ttype, argt, argn, setn}
	basename := util.NoExt(filename)
	wrapfname := basename
	if *Flag_fp64 {
		wrapfname = *Flag_outdir + "/" + filepath.Base(basename) + "_fp64_wrapper.go"
	} else {
		wrapfname = *Flag_outdir + "/" + filepath.Base(basename) + "_wrapper.go"
	}
	wrapout, err := os.OpenFile(wrapfname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	util.PanicErr(err)
	defer wrapout.Close()
	util.PanicErr(templ.Execute(wrapout, kernel))
}

// wrapper code template text
const templText = `package opencl{{.TType}}

/*
 THIS FILE IS AUTO-GENERATED BY OCL2GO.
 EDITING IS FUTILE.
*/

import(
	"unsafe"
	cl "github.com/seeder-research/uMagNUS/cl"
	timer "github.com/seeder-research/uMagNUS/timer"
	"sync"
	"fmt"
)


// Stores the arguments for {{.Name}} kernel invocation
type {{.Name}}_args_t struct{
	{{range $i, $_ := .ArgN}} arg_{{.}} {{index $.ArgT $i}}
	{{end}} argptr [{{len .ArgN}}]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for {{.Name}} kernel invocation
var {{.Name}}_args {{.Name}}_args_t

func init(){
	// OpenCL driver kernel call wants pointers to arguments, set them up once.
	{{range $i, $t := .ArgN}} {{$.Name}}_args.argptr[{{$i}}] = unsafe.Pointer(&{{$.Name}}_args.arg_{{.}})
	{{end}} }

// Wrapper for {{.Name}} OpenCL kernel, asynchronous.
func k_{{.Name}}_async ( {{range $i, $t := .ArgT}}{{index $.ArgN $i}} {{$t}}, {{end}} cfg *config, queue *cl.CommandQueue, events []*cl.Event) *cl.Event {
	if Synchronous{ // debug
		if err := queue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish in beginning of {{.Name}}: %+v", err)
		}
		timer.Start("{{.Name}}")
	}

	{{.Name}}_args.Lock()
	defer {{.Name}}_args.Unlock()

	{{range $i, $t := .ArgN}} {{$.Name}}_args.arg_{{.}} = {{.}}
	{{end}}

	{{range $i, $t := .SetN}}{{$t}}
	{{end}}

//	args := {{.Name}}_args.argptr[:]
	event := LaunchKernel("{{.Name}}", cfg.Grid, cfg.Block, queue, events)

	if Synchronous{ // debug
		if err := queue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish at end of {{.Name}}: %+v", err)
		}
		timer.Stop("{{.Name}}")
	}

	return event
}

`

// wrapper code template
var templ = template.Must(template.New("wrap").Parse(templText))

// should token be filtered out of stream?
func filter(token string) bool {
	switch token {
	case "__restrict":
		return true
	case "volatile":
		return true
	case "unsigned":
		return true
	case "signed":
		return true
	case "const":
		return true
	}
	return false
}

func chkArgMemType(token string) int {
	switch token {
	case "__global":
		return 1
	case "__local":
		return 2
	case "__constant":
		return 3
	case "__private":
		return 4
	}
	return 0
}
