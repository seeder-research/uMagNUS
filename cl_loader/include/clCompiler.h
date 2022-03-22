// This file serves as a stub for the actual library. Actual library
// needs to be compiled at the deployment machine to work.
#include <stdlib.h>

extern const char deviceNames[];
extern const size_t deviceNameLen;
extern const int NumDevices;
extern const size_t binIdx[];
extern const size_t binSizes[];
extern const char * hexPtrs[];

extern char * sendStringPtr(size_t idx);

extern size_t sendBinSize(size_t idx);

extern size_t sendBinIdx(size_t idx);
