#! /bin/bash

go build ocl2go.go || exit 1

for f in *.cl; do
	g=$(echo $f | sed 's/\.cl$//') # file basename
	if [[ $f -nt '../../'$g'_wrapper.go' ]]; then
		./ocl2go $f || exit 1
	fi
done

