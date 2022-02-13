#!/bin/bash

export UMAGNUS=uMagNUS

## Run OOMMF to obtain results
tclsh ${TCL_OOMMF} boxsi stdprob4.mif && \
cp stdprob4a-cs2.50x2.50x3.00-Oxs_TimeDriver-Spin-00-*.omf \
  oommf-end.ovf && \
rm -rf /tmp/uMagNUS*.ovf \
  *.out && \
sed -i "s/precision := 0/precision := 0/" standardproblem4a_newell-oommf.mx3 && \
${UMAGNUS} -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4a_newell-oommf.mx3 && \
mv standardproblem4a_newell-oommf.out standardproblem4a_newell-oommf32.res && \
sed -i "s/precision := 0/precision := 1/" standardproblem4_newell-oommf.mx3 && \
${UMAGNUS}64 -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4a_newell-oommf.mx3 && \
mv standardproblem4_newell-oommf.out standardproblem4a_newell-oommf64.res && \
sed -i "s/precision := 1/precision := 0/" standardproblem4a_newell-oommf.mx3 && \
sed -i "s/precision := 0/precision := 0/" standardproblem4a_newell.mx3 && \
${UMAGNUS} -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4a_newell.mx3 && \
mv standardproblem4a_newell.out standardproblem4a_newell32.res && \
sed -i "s/precision := 0/precision := 1/" standardproblem4a_newell.mx3 && \
${UMAGNUS}64 -http "" -f -failfast -paranoid=false -cache /tmp standardproblem4a_newell.mx3 && \
mv standardproblem4a_newell.out standardproblem4a_newell64.res && \
sed -i "s/precision := 1/precision := 0/" standardproblem4a_newell.mx3
