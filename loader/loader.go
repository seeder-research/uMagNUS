package loader
/*
#include clCompiler.h

static size_t returnStringLength() {
	return deviceNameLen;
}

static int returnNumDevices() {
	return NumDevices;
}

static char * returnStringPointer() {
	return &deviceNames[0];
}
*/
import "C"

import (
	"strings"

	"github.com/seeder-research/uMagNUS/cl"
)

func checkDevice(d *cl.Device) int {
	// Get the string names in the library
	stringPtr := C.returnStringPointer()
	stringLen := (C.int)(C.returnStringLength())

	// Error check if library is empty
	if (stringPtr == nil) || (stringLen <= 0) {
		return -1
	}

	// Library is not empty so try to get the semi-colon
	// separated string consisting of all device names
	fullNameString := C.GoStringN(stringPtr, stringLen)

	// Get the name of the target device
	targDeviceName := d.Name()

	// Compare device name to found OpenCL device names
	// to get the index
	deviceNameArray := strings.Split(fullNameString, ";")
	for idx, devName := range deviceNameArray {
		if devName == targDeviceName {
			return idx
		}
	}
	return -2
}

func GetClDeviceBinary(d cl*Device) []byte {
	// Check library for index of OpenCL device
	// if it is available
	idx := checkDevice(d)

	// Error check for invalid device
	if idx < 0 {
		return []byte{}
	}

	// OpenCL device is available
	// Proceeding to get binary...
	binIdx := C.sendBinIdx((C.size_t)(idx))
	binSize := C.sendBinSize(binIdx)
	ProgramBytes, err := e.DecodeString(C.GoStringN(C.sendStringPtr(binIdx), (C.int)(2*binSize)))
	if err != nil {
		fmt.Println("Unable to get opencl program from library")
		return []byte{}
	}
	return ProgramBytes
}
