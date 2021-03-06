//go:build ignore
// +build ignore

// This program generates Go wrappers for opencl sources.
// The opencl file should contain exactly one __kernel void.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"os"
	"regexp"
	"text/scanner"

	kernels_src "github.com/seeder-research/uMagNUS/kernels_src"
	util "github.com/seeder-research/uMagNUS/util"
)

// template data
type Kernel_stuff struct {
	OCL  map[string]string
	Code map[string]string
}

var ls_dirclh []string
var ls_dircl []string
var ls_dircl64 []string
var (
	Flag_indir  = flag.String("indir", ".", "directory containing clh, cl and cl64 directories to search for kernel files")
	Flag_outdir = flag.String("outdir", "..", "directory containing kernels and kernels64 directories to output program_wrapper.go")
	Flag_fp64   = flag.Bool("double", false, "generating for double-precision")
)

func main() {
	flag.Parse()
	// find .clh files
	if ls_dirclh == nil {
		dirclh, errd := os.Open(*Flag_indir + "/clh")
		defer dirclh.Close()
		util.PanicErr(errd)
		var errls error
		ls_dirclh, errls = dirclh.Readdirnames(-1)
		util.PanicErr(errls)
	}

	// find .cl files
	if ls_dircl == nil {
		dircl, errd := os.Open(*Flag_indir + "/cl")
		defer dircl.Close()
		util.PanicErr(errd)
		var errls error
		ls_dircl, errls = dircl.Readdirnames(-1)
		util.PanicErr(errls)
	}

	// get header codes in .clh files
	opencl_codes := &Kernel_stuff{make(map[string]string), make(map[string]string)}
	for _, f := range ls_dirclh {
		match, e := regexp.MatchString("..clh$", f)
		util.PanicErr(e)
		if match {
			fkey := f[:len(f)-len(".clh")]
			opencl_codes.OCL[fkey] = getFile(*Flag_indir + "/clh/" + f)
		}
	}

	// get names of kernels available in .cl files
	for _, f := range ls_dircl {
		match, e := regexp.MatchString("..cl$", f)
		util.PanicErr(e)
		if match {
			kname := getKernelName(*Flag_indir + "/cl/" + f)
			opencl_codes.Code[kname] = getFile(*Flag_indir + "/cl/" + f)
		}
	}

	if *Flag_fp64 {
		// find .cl files
		if ls_dircl64 == nil {
			dircl64, errd := os.Open(*Flag_indir + "/cl64")
			defer dircl64.Close()
			util.PanicErr(errd)
			var errls error
			ls_dircl64, errls = dircl64.Readdirnames(-1)
			util.PanicErr(errls)
		}
		// get names of kernels available in .cl files
		for _, f := range ls_dircl64 {
			match, e := regexp.MatchString("..cl$", f)
			util.PanicErr(e)
			if match {
				kname := getKernelName(*Flag_indir + "/cl64/" + f)
				opencl_codes.Code[kname] = getFile(*Flag_indir + "/cl64/" + f)
			}
		}
	}

	tmpBuffer := new(bytes.Buffer)
	if *Flag_fp64 {
		tmpBuffer.WriteString("package kernels64\n")
	} else {
		tmpBuffer.WriteString("package kernels\n")
	}
	tmpBuffer.WriteString("\n\n// THIS FILE WAS CREATED BY OPENCL2GO\n")
	tmpBuffer.WriteString("// MODIFYING THIS FILE IS FUTILE!!!!!\n\n")
	tmpBuffer.WriteString("func OpenclProgramSource() string {\n")
	tmpBuffer.WriteString("	opencl_codes := `\n")
	for _, keynames := range kernels_src.OCLHeadersList {
		tmpBuffer.WriteString(opencl_codes.OCL[keynames])
	}
	for _, keynames := range kernels_src.OCLKernelsList {
		tmpBuffer.WriteString(opencl_codes.Code[keynames])
	}
	if *Flag_fp64 {
		for _, keynames := range kernels_src.OCL64KernelsList {
			tmpBuffer.WriteString(opencl_codes.Code[keynames])
		}
	}
	tmpBuffer.WriteString("\n`\n\n	return opencl_codes\n}\n")

	wrapfname := string("")
	if *Flag_fp64 {
		wrapfname = *Flag_outdir + "/kernels64/program_wrapper.go"
	} else {
		wrapfname = *Flag_outdir + "/kernels/program_wrapper.go"
	}
	wrapout, err := os.OpenFile(wrapfname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	util.PanicErr(err)
	wrapout.WriteString(tmpBuffer.String())
	wrapout.Close()
}

func getKernelName(fname string) string {
	// open opencl file
	f, err := os.Open(fname)
	util.PanicErr(err)
	defer f.Close()

	// read tokens
	var token []string
	var s scanner.Scanner
	s.Init(f)
	tok := s.Scan()
	for tok != scanner.EOF {
		if !filter(s.TokenText()) {
			token = append(token, s.TokenText())
		}
		tok = s.Scan()
	}

	// find function name and arguments
	funcname := ""
	for i := 0; i < len(token); i++ {
		if token[i] == "__kernel" {
			funcname = token[i+2]
			break
		}
	}
	return funcname
}

func getFile(fname string) string {
	f, err := os.Open(fname)
	util.PanicErr(err)
	defer f.Close()
	in := bufio.NewReader(f)
	var out bytes.Buffer
	line, err := in.ReadBytes('\n')
	for err != io.EOF {
		util.PanicErr(err)
		out.Write(line)
		line, err = in.ReadBytes('\n')
	}
	return out.String()
}

// should token be filtered out of stream?
func filter(token string) bool {
	switch token {
	case "__restrict":
		return true
	case "__global":
		return true
	case "__constant":
		return true
	case "__local":
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
