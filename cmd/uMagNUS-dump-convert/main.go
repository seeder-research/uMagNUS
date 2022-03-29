/*
dump-convert converts .dump output files between single and double precision.

Usage

Command-line flags must always preceed the input files:
	dump-convert -i input_file -o output_file

Input file is read in to discover the precision of data stored
If data is single precision, output is double precision
If data is double precision, output is single precision
*/
package main

import (
	"flag"
	"log"

	dump "github.com/seeder-research/uMagNUS/dump64"
)

var (
	flag_ifile = flag.String("i", "", "Input file")
	flag_ofile = flag.String("o", "", "Output file")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *flag_ifile == "" {
		log.Fatal("no input file")
	}

	precision, err := dump.ReadFilePrecision(*flag_ifile)
	if err != nil {
		log.Fatal("Unable to read input file")
	} else {
		if precision == 4 {
			d, m, e := dump.ReadFile32(*flag_ifile)
			e = dump.WriteFile(*flag_ofile, d, m)
			if e != nil {
				log.Fatal("Unable to write output file")
			}
		} else if precision == 8 {
			d, m, e := dump.ReadFile(*flag_ifile)
			e = dump.WriteFile32(*flag_ofile, d, m)
			if e != nil {
				log.Fatal("Unable to write output file")
			}
		} else {
			log.Fatal("Unknown precision of input file!")
		}
	}
}
