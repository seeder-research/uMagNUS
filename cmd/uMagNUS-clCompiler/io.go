package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readFile(fname string) string {
	if *Flag_verbose > 5 {
		fmt.Println("Attempting to open file: ", fname)
	}
	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Unable to open file: %v \n", fname)
		return ""
	}
	defer f.Close()

	in := bufio.NewReader(f)
	var nline string
	line, eof := readLine(in)
	for !eof {
		line = line + nline + "\n"
		nline, eof = readLine(in)
	}
	if nline != "" {
		line += nline
	}
	return line
}
