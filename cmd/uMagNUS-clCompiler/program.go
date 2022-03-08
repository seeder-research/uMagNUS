package main

import (
	"fmt"
	"github.com/seeder-research/uMagNUS/cl"
)

func compileProgram(ctx *cl.Context, devices []*cl.Device, source []string) (*cl.Program, error) {
	var program *cl.Program
	program, err = ctx.CreateProgramWithSource(source)
	if err != nil {
		fmt.Println("compileProgram: Unable to get create program in context!")
		return nil, err
	}
	err = program.CompileProgram(devices, generateCompilerOpts(), nil)
	return program, err
}

func linkProgram(programs []*cl.Program) {
}
