package main

import (
	"io"
)

func generateCompilerOpts() string {
	var opts string
	opts = ""
	if *Flag_ComArgs != "" {
		opts = opts + *Flag_ComArgs
	}
	if *Flag_ClStd != "" {
		if opts != "" {
			opts = opts + " "
		}
		opts = opts + "-std=" + *Flag_ClStd
	}
	if *Flag_includes != "" {
		if opts != "" {
			opts = opts + " "
		}
		opts = opts + *Flag_includes
	}
	if *Flag_defines != "" {
		if opts != "" {
			opts = opts + " "
		}
		opts = opts + *Flag_defines
	}
	return opts
}

func generateLinkerOpts() string {
	var opts string
	opts = ""
	if *Flag_libpaths != "" {
		opts = opts + *Flag_libpaths
	}
	if *Flag_libs != "" {
		if opts != "" {
			opts = opts + " "
		}
		opts = opts + *Flag_libs
	}
	return opts
}

func readLine(in io.Reader) (line string, eof bool) {
	char := readChar(in)
	eof = isEOF(char)

	for !isEndline(char) {
		line += string(byte(char))
		char = readChar(in)
	}
	return line, eof
}

func isEOF(char int) bool {
	return char == -1
}

func isEndline(char int) bool {
	return isEOF(char) || char == int('\n')
}

//// Blocks until all requested bytes are read.
//type fullReader struct{ io.Reader }
//
//func (r fullReader) Read(p []byte) (n int, err error) {
//      return io.ReadFull(r.Reader, p)
//}

// Reads one character from the Reader.
// -1 means EOF.
// Errors are cought and cause panic
func readChar(in io.Reader) int {
	buffer := [1]byte{}
	switch nr, err := in.Read(buffer[0:]); true {
	case nr < 0: // error
		panic(err)
	case nr == 0: // eof
		return -1
	case nr > 0: // ok
		return int(buffer[0])
	}
	panic("unreachable")
}
