package main

import (
	"github.com/seeder-research/uMagNUS/cl"
)

func compileProgram(ctx *cl.Context, source []string) (*cl.Program, error) {
	return ctx.CreateProgramWithSource(source)
}

func linkProgram(programs []*cl.Program) {
}
