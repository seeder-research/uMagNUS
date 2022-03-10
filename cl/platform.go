package cl

/*
#include "./opencl.h"

static cl_int CLGetPlatformInfoParamSize(cl_platform_id                  platform,
                                         cl_platform_info              param_name,
                                         size_t             *param_value_size_ret) {
    return clGetPlatformInfo(platform, param_name, NULL, NULL, param_value_size_ret);
}

static cl_int CLGetPlatformInfoParamUnsafe(cl_platform_id            platform,
                                           cl_platform_info        param_name,
                                           size_t            param_value_size,
                                           void                  *param_value) {
    return clGetPlatformInfo(platform, param_name, param_value_size, param_value, NULL);
}
*/
import "C"

import "unsafe"

//////////////// Constants ////////////////
const maxPlatforms = 32

//////////////// Abstract Types ////////////////
type Platform struct {
	id C.cl_platform_id
}

//////////////// Basic Functions ////////////////

// Obtain the list of platforms available.
func GetPlatforms() ([]*Platform, error) {
	var platformIds [maxPlatforms]C.cl_platform_id
	var nPlatforms C.cl_uint
	if err := C.clGetPlatformIDs(C.cl_uint(maxPlatforms), &platformIds[0], &nPlatforms); err != C.CL_SUCCESS {
		return nil, toError(err)
	}
	platforms := make([]*Platform, nPlatforms)
	for i := 0; i < int(nPlatforms); i++ {
		platforms[i] = &Platform{id: platformIds[i]}
	}
	return platforms, nil
}

//////////////// Abstract Functions ////////////////
func (p *Platform) GetDevices(deviceType DeviceType) ([]*Device, error) {
	return GetDevices(p, deviceType)
}

func (p *Platform) getInfoString(param C.cl_platform_info) (string, error) {
	var strN C.size_t
	if err := C.CLGetPlatformInfoParamSize(p.id, param, &strN); err != C.CL_SUCCESS {
		return "", toError(err)
	}
	strC := (*C.char)(C.calloc(strN, 1))
	defer C.free(unsafe.Pointer(strC))
	if err := C.CLGetPlatformInfoParamUnsafe(p.id, param, strN, unsafe.Pointer(strC)); err != C.CL_SUCCESS {
		return "", toError(err)
	}
	retString := C.GoStringN(strC, C.int(strN-1))
	return retString, nil
}

func (p *Platform) Name() string {
	if str, err := p.getInfoString(C.CL_PLATFORM_NAME); err != nil {
		panic("Platform.Name() should never fail")
	} else {
		return str
	}
}

func (p *Platform) Vendor() string {
	if str, err := p.getInfoString(C.CL_PLATFORM_VENDOR); err != nil {
		panic("Platform.Vendor() should never fail")
	} else {
		return str
	}
}

func (p *Platform) Profile() string {
	if str, err := p.getInfoString(C.CL_PLATFORM_PROFILE); err != nil {
		panic("Platform.Profile() should never fail")
	} else {
		return str
	}
}

func (p *Platform) Version() string {
	if str, err := p.getInfoString(C.CL_PLATFORM_VERSION); err != nil {
		panic("Platform.Version() should never fail")
	} else {
		return str
	}
}

func (p *Platform) Extensions() string {
	if str, err := p.getInfoString(C.CL_PLATFORM_EXTENSIONS); err != nil {
		panic("Platform.Extensions() should never fail")
	} else {
		return str
	}
}
