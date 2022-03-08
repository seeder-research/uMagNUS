package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func readFile(fname string) string {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Unable to open file: %v \n", fname)
		return ""
	}
	defer f.Close()

	in := bufio.NewReader(f)
	var out bytes.Buffer
	line, err := in.ReadBytes('\n')
	for err != io.EOF {
		log.Panic(err)
		out.Write(line)
		line, err = in.ReadBytes('\n')
	}
	return out.String()
}
