#! /bin/bash

# builds with -race and runs tests with browser open.

set -e

go install -race github.com/mumax/3cl/cmd/mumax3cl

google-chrome http://localhost:35367 &

for f in *.mx3; do
	mumax3cl $f 
done

go install github.com/mumax/3cl/cmd/mumax3cl # re-build without race detector

