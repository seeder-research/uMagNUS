package main

import (
	"encoding/json"
	"io"

	data "github.com/seeder-research/uMagNUS/data"
)

func dumpJSON(f *data.Slice, info data.Meta, out io.Writer) {
	w := json.NewEncoder(out)
	w.Encode(f.Tensors())
}
