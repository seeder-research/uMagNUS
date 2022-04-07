#! /bin/bash

# builds with -race and runs tests with browser open.

set -e

go install -race github.com/seeder-research/uMagNUS/cmd/uMagNUS

google-chrome http://localhost:35367 &

for f in *.mx3; do
	uMagNUS $f 
done

go install github.com/seeder-research/uMagNUS/cmd/uMagNUS # re-build without race detector

