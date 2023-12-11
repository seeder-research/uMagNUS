uMagNUS
======

GPU accelerated micromagnetic simulator based on OpenCL.


Downloads and documentation
---------------------------

The frontend is based on MuMax3 and accepts simulation files written for MuMax3.

Refer to MuMax3 documentation at:
http://mumax.github.io


Paper
-----

- To be updated.


Building from source (for linux and msys2)
--------------------

Consider downloading a pre-compiled binary. If you want to compile nevertheless:

  * install the OpenCL driver, if not yet present.
    - if unsure, it's probably already there
    - requires OpenCL 1.2 support
    - the steps to build from scratch to running the tests were tested on the following setups:
        - Ubuntu 18.04 with Intel HD Graphics 630 (default Intel drivers from ubuntu repo)
        - Ubuntu 20.04 with Nvidia Quadro P100 (nvidia-dkms-510 driver from nvidia repo)
        - Ubuntu 20.04 with Nvidia Quadro P2000 (nvidia-dkms-510 driver from nvidia repo)
        - Ubuntu 22.04 with Nvidia Quadro P100 (nvidia-dkms-515 driver from nvidia repo)
        - Ubuntu 22.04 with Nvidia Quadro P2000 (nvidia-dkms-515 driver from nvidia repo)
        - Ubuntu 20.04 with Nvidia RTX 2080 Super (nvidia-dkms-510 driver from nvidia repo)
        - Ubuntu 22.04 with Nvidia RTX 2080 Super (nvidia-dkms-515 driver from nvidia repo)
        - Ubuntu 20.04 with Nvidia GTX 660 Ti (nvidia-dkms-418 driver from nvidia repo)
        - Ubuntu 20.04 with AMD RX 6500 XT (rocm-dkms driver from rocm repo)
        - Ubuntu 22.04 with AMD RX 6500 XT (rocm-dkms driver from rocm repo)
        - Windows 11 using MSYS2 with Nvidia MX150 (Nvidia 416.xx to 512.xx drivers)
        - Windows 11 using MSYS2 with Intel UHD Graphics 620 (Intel 26.xx to 30.xx drivers)
  * install TDM GCC compiler is gcc from msys2 is not compatible
    - https://jmeubank.github.io/tdm-gcc/download/
  * install Go 
    - https://golang.org/dl/
    - Ensure the directory containing the go binary is in your $PATH
  * if you have git installed:
    - get the source from GitHub.com by running:
        - `git clone https://github.com/seeder-research/uMagNUS`
        - Checkout branch with tags corresponding to the version required
  * if you don't have git:
    - get the source from https://github.com/seeder-research/uMagNUS/releases
    - unzip the source into $GOPATH/src/github.com/seeder-research/uMagNUS
  * `cd $GOPATH/src/github.com/seeder-research/uMagNUS`
    - `make uMagNUS`
    - For double-precision support, run:
        - `make uMagNUS64`
  * optional: install gnuplot if you want pretty graphs
    - Ubuntu: `sudo apt-get install gnuplot`

Your binary is now at `$GOPATH/bin/uMagNUS` (or `$./gopath/bin/uMagNUS` if
GOPATH is not set)
The binary supporting double-precision is now at `GOPATH/bin/uMagNUS64` (or
`./gopath/bin/uMagNUS64` if GOPATH is not set)

To do all at once on Ubuntu:
```
sudo apt-get install git golang-go gcc gnuplot
export GOPATH=$HOME git clone github.com/seeder-research/uMagNUS
cd uMagNUS
make base
```

Adding your own extensions and opencl kernels
------------
  * extensions to the program can be created by looking in engine and opencl
    directories in the source directory. The directories to look in for
    modifying uMagNUS64 end with '64' (i.e., engine64, opencl64)
  * OpenCL kernels are located in kernels_src directory.
    - kernels_src/clh contain header files for the kernels
    - kernels_src/cl contain files in which every file defines one and only
      one kernel that will be used in uMagNUS and uMagNUS64
    - kernels_src/cl64 contain files in which every file defines one and only
      one kernel that will be used in uMagNUS64
  * Data should have real_t type, which are set at compile time as shown in
    clh/typedefs.clh
  * uMagNUS-clCompiler is an offline tool that allows the use of the
    installed OpenCL platform compiler to compiler kernels and output a C
    file which can be compiled into a shared library using GCC.
  * uMagNUS-kernelLoader is an offline tool that checks the shared library
    built using the C file output by uMagNUS-clCompiler
  * Individual kernel files can be tested by running through
    uMagNUS-clCompiler and uMagNUS-kernelLoader to check for compile and link
    errors
  * A kernel file that is free of compile and link errors needs to be added
    to the build steps depending on the deployment method
    - If the kernel is to be compiled every time uMagNUS or uMagNUS64 starts
      up, add the name of the kernel to the list in
      kernels_src/sequence.go
    - If the kernel is to be compiled offline and loaded as a shared library,
      add the names of the kernel files to kernels_src/Kernels/kernels32.h
      (for uMagNUS) and/or kernels_src/Kernels/kernels64.h for uMagNUS64
  * Redo the build steps to create binaries for uMagNUS and uMagNUS64 with
    your extensions

Deployment at the compute machine
------------
The build steps above only build the binaries for uMagNUS and uMagNUS64. There
are additional dependencies that need to be built before uMagNUS can be run.
  * Building the umagnus libraries
    - cd into the the source directory
    - `make loaders`
        - This creates the shared stub library in cl_loader/lib (libumagnus.so
          and libumagnus64.so for linux, umagnus.dll and umagnus64.dll for
          windows) in which there are no OpenCL kernels
        - Make sure the cl_loader/lib directory exists in LD_LIBRARY_PATH for
          linux and PATH for windows
        - Run ldd on the uMagNUS and uMagNUS64 binaries to check if the
          environment is set up to correctly pick up the umagnus libraries
    - `make libs`
        - This step needs to be run on the machine on which the GPU(s) will be
          used to run uMagNUS and uMagNUS64
        - The shared libraries containing the OpenCL kernels are located in
          the libumagnus directory
        - Make sure the libumagnus directory exists in LD_LIBRARY_PATH for
          linux and PATH for windows
        - Run ldd on the uMagNUS and uMagNUS64 binaries to check if the
          environment is set up to correctly pick up the umagnus libraries
        - When these shared libraries are used, the OpenCL kernels that were
          compiled and linked offline by uMagNUS-clCompiler and GCC are
          loaded by uMagNUS and uMagNUS64 when they start up
  * For HPC environments, the `make libumagnus libumagnus64` step should be
    performed on every individual machine that will be running uMagNUS, and
    the generated shared libraries should be saved locally (i.e., in /opt).
    Users running uMagNUS or uMagNUS64 on the machine will need to ensure
    their environment variables are updated so that uMagNUS and uMagNUS64
    see the shared library that is saved locally on the machine

Contributing
------------

Contributions are gratefully accepted. To contribute code, fork the repo on
github and send a pull request.
