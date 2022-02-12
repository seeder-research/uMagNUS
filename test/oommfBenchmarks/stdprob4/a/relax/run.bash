#!/bin/bash

export UMAGNUS=uMagNUS

## Run OOMMF to obtain results
tclsh ${TCL_OOMMF} boxsi stdprob4.mif && \
cp stdprob4-cs2.50x2.50x3.00-Oxs_MinDriver-Spin-00-0001388.omf \
  stdprob4a.ovf && \
rm -rf /tmp/uMagNUS*.ovf \
  *.out && \
${UMAGNUS} -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4_newell-relax.mx3 && \
mv standardproblem4_newell-relax.out standardproblem4_newell-relax32.res && \
${UMAGNUS}64 -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4_newell-relax.mx3 && \
mv standardproblem4_newell-relax.out standardproblem4_newell-relax64.res
