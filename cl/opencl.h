/*
  This file is used to point the compiler to the actual opencl.h of the system.
  It is also used to check the version of opencl installed
*/
#include <stdlib.h>
#define CL_USE_DEPRECATED_OPENCL_1_2_APIS
#define CL_USE_DEPRECATED_OPENCL_2_0_APIS
#define CL_USE_DEPRECATED_OPENCL_2_1_APIS
#define CL_USE_DEPRECATED_OPENCL_2_2_APIS
#define CL_TARGET_OPENCL_VERSION 120
#ifdef __APPLE__
	#include "OpenCL/OpenCL.h"
#else
	#include "CL/opencl.h"
#endif

#ifndef CL_VERSION_1_2
	#error "This package requires OpenCL 1.2"
#endif

