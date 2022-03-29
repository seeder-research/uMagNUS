package loader64

/*
#include "clCompiler.h"

const char * returnDeviceNamePtr() {
        return &deviceNames[0];
}

size_t returnDeviceNameLen() {
        return deviceNameLen;
}

int returnNumDevices() {
        return NumDevices;
}

const char * sendStringPtr(size_t idx) {
        return hexPtrs[idx];
}

size_t sendBinSize(size_t idx) {
        return binSizes[idx];
}

size_t sendBinIdx(size_t idx) {
        return binIdx[idx];
}
*/
import "C"

import (
	e "encoding/hex"
	"fmt"
	"log"
	"strings"

	cl "github.com/seeder-research/uMagNUS/cl"
)

func checkDevice(d *cl.Device) int {

	// First error check if library is empty
	numDevices := C.returnNumDevices()
	if numDevices <= 0 {
		return -1
	}

	// Get the string names in the library
	stringPtr := C.returnDeviceNamePtr()
	stringLen := (C.int)(C.returnDeviceNameLen())

	// Second error check if library is empty
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

func GetClDeviceBinary(d *cl.Device) []byte {
	// Check library for index of OpenCL device
	// if it is available
	idx := checkDevice(d)

	// Error check for invalid device
	if idx < 0 {
		return nil
	}

	// OpenCL device is available
	// Proceeding to get binary...
	binIdx := C.sendBinIdx((C.size_t)(idx))
	binSize := (int)(C.sendBinSize(binIdx))
	//byteString := C.GoStringN(C.sendStringPtr(binIdx), (C.int)(2*binSize))
	byteString := C.GoString(C.sendStringPtr(binIdx))
	ProgramBytes, err := e.DecodeString(byteString)
	if err != nil {
		fmt.Println("Unable to get opencl program from library")
		log.Fatal(err)
		return nil
	}
	if len(ProgramBytes) != binSize {
		fmt.Printf("Decoded program bytes (%+v bytes) do not match stored program bytes (%+v bytes)! \n", len(ProgramBytes), binSize)
	}
	return ProgramBytes
}
