package engine

// Comp is a Derived Quantity pointing to a single component of vector Quantity

import (
	"fmt"

	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

type component struct {
	parent Quantity
	comp   int
}

// Comp returns vector component c of the parent Quantity
func Comp(parent Quantity, c int) ScalarField {
	util.Argument(c >= 0 && c < parent.NComp())
	return AsScalarField(&component{parent, c})
}

func (q *component) NComp() int       { return 1 }
func (q *component) Name() string     { return fmt.Sprint(NameOf(q.parent), "_", compname[q.comp]) }
func (q *component) Unit() string     { return UnitOf(q.parent) }
func (q *component) Mesh() *data.Mesh { return MeshOf(q.parent) }

func (q *component) Slice() (*data.Slice, bool) {
	p := q.parent
	src := ValueOf(p)
	defer opencl.Recycle(src)
	c := opencl.Buffer(1, src.Size())
	return c, true
}

func (q *component) EvalTo(dst *data.Slice) {
	src := ValueOf(q.parent)
	defer opencl.Recycle(src)
	data.Copy(dst, src.Comp(q.comp))

	// sync before returning
	seqQueue := opencl.ClCmdQueue[0]
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish in component.evalto: %+v \n", err)
	}
}

var compname = map[int]string{0: "x", 1: "y", 2: "z"}
