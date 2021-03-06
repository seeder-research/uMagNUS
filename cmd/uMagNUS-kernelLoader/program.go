package main

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
)

var binariesMap map[string]*cl.ProgramBinaries

func ShowBuildLog(p *cl.Program, d *cl.Device) {
	status, err := p.GetBuildStatus(d)
	if err != nil {
		fmt.Println("  ERROR: unable to get build status of program!")
		return
	}
	switch status {
	case cl.BuildStatusSuccess:
		if *Flag_verbose > 2 {
			fmt.Println("  Successfully built program")
		}
	case cl.BuildStatusNone:
		if *Flag_verbose > 2 {
			fmt.Println("  Program was not built/compiled/linked")
			fmt.Println("    Please run clBuildProgram, clCompileProgram or clLinkProgram")
		}
		return
	case cl.BuildStatusError:
		if *Flag_verbose > 2 {
			fmt.Println("  Program is built with errors!")
		}
	case cl.BuildStatusInProgress:
		if *Flag_verbose > 2 {
			fmt.Println("  Program build is in progress")
		}
		return
	default:
		if *Flag_verbose > 2 {
			fmt.Println("  ERROR: Unknown status returned")
		}
		return
	}

	var logOutput string
	logOutput, err = p.GetBuildLog(d)
	if err != nil {
		fmt.Println("  ERROR: unable to get build log of program!")
		return
	}
	if logOutput == "" {
		fmt.Println("Empty log!")
	} else {
		fmt.Printf("%+v \n", logOutput)
	}
}
