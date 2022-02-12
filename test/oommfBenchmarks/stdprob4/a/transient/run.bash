#!/bin/bash

export UMAGNUS=uMagNUS

## Run OOMMF to obtain results
tclsh ${TCL_OOMMF} boxsi stdprob4.mif && \
cp stdprob4a-cs2.50x2.50x3.00-Oxs_TimeDriver-Spin-00-0015168.omf \
  oommf-end.ovf && \
rm -rf /tmp/uMagNUS*.ovf \
  *.out && \
${UMAGNUS} -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4_newell-oommf.mx3 && \
mv standardproblem4_newell-oommf.out standardproblem4_newell-oommf32.res && \
${UMAGNUS}64 -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4_newell-oommf.mx3 && \
mv standardproblem4_newell-oommf.out standardproblem4_newell-oommf64.res && \
${UMAGNUS} -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4_newell.mx3 && \
mv standardproblem4_newell.out standardproblem4_newell32.res && \
${UMAGNUS}64 -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4_newell.mx3 && \
mv standardproblem4_newell.out standardproblem4_newell64.res
