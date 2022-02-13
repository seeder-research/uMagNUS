#!/bin/bash

export UMAGNUS=uMagNUS

## Run OOMMF to obtain results
tclsh ${TCL_OOMMF} boxsi -regression_test true stdprob5.mif && \
cp regression_test*-Spin-00-*.omf \
  oommf-relax.ovf && \
rm -rf /tmp/uMagNUS*.ovf \
  *.out && \
${UMAGNUS} -http "" -f -failfast -paranoid=false -cache /tmp standardproblem5_newell.mx3 && \
mv standardproblem5_newell.out standardproblem5_newell32.res && \
rm -rf /tmp/uMagNUS*.ovf \
  *.out && \
${UMAGNUS}64 -http "" -f -failfast -paranoid=false -cache /tmp standardproblem5_newell.mx3 && \
mv standardproblem5_newell.out standardproblem5_newell64.res
