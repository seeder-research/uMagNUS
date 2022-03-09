#!/bin/bash

WORKDIR=$1

outFile=merged_kernels.h
if [ ! -f ${outFile} ]; then
    for f in ${WORKDIR}/clh/*.clh
    do
        fn=$(basename ${f} .clh)
        echo "#include \"clh/"${fn}".clh\"" >> ${outFile}
    done
    for f in ${WORKDIR}/cl/*.cl
    do
        fn=$(basename ${f} .cl)
        echo "#include \"cl/"${fn}".cl\"" >> ${outFile}
    done
fi
