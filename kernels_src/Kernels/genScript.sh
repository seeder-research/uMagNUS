#!/bin/bash

WORKDIR=$1

for f in ${WORKDIR}/*.cl
do
    fn=$(basename ${f} .cl)
    if [ ! -f ${fn}.h ]; then
        echo "#include \"cl/"${fn}".cl\"" >> ${fn}.h
    fi
done
