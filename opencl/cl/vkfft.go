package cl

/*
#cgo CFLAGS: -I./
#cgo LDFLAGS: -lm
#include "fft_interface.h"
*/
import "C"

import (
	"errors"
	"fmt"
	//        "unsafe"
)

//////////////// Basic Errors ////////////////

var ErrUnsupportedFFT = errors.New("vkFFT: unsupported")

var (
	ErrUnknownFFT = errors.New("vkFFT: unknown error") // Generally an unexpected result from a vkFFT function (e.g. CL_SUCCESS but null pointer)
)

type ErrOtherFFT int

func (e ErrOtherFFT) Error() string {
	return fmt.Sprintf("vkFFT: error %d", int(e))
}

var (
	ErrVkFFTFailMalloc                     = errors.New("vkFFT: Failed to malloc")
	ErrVkFFTFailInsufficientCodeBuffer     = errors.New("vkFFT: Insufficient code buffer")
	ErrVkFFTFailInsufficientTempBuffer     = errors.New("vkFFT: Insufficient temp buffer")
	ErrVkFFTFailPlanNotInitialized         = errors.New("vkFFT: Plan not initialized")
	ErrVkFFTFailNullTemp                   = errors.New("vkFFT: Null temp")
	ErrVkFFTFailInvalidPhysicalDevice      = errors.New("vkFFT: Invalid physical device")
	ErrVkFFTFailInvalidDevice              = errors.New("vkFFT: Invalid device")
	ErrVkFFTFailInvalidQueue               = errors.New("vkFFT: Invalid queue")
	ErrVkFFTFailInvalidCommandPool         = errors.New("vkFFT: Invalid command pool")
	ErrVkFFTFailInvalidFence               = errors.New("vkFFT: Invalid fence")
	ErrVkFFTFailOnlyForwardFFTInitialized  = errors.New("vkFFT: Only forward FFT initialized")
	ErrVkFFTFailOnlyBackwardFFTInitialized = errors.New("vkFFT: Only backward FFT initialized")
	ErrVkFFTFailInvalidContext             = errors.New("vkFFT: Invalid context")
	ErrVkFFTFailInvalidPlatform            = errors.New("vkFFT: Invalid platform")
	ErrVkFFTFailEmptyFFTdim                = errors.New("vkFFT: Empty FFTdim")
	ErrVkFFTFailEmptySize                  = errors.New("vkFFT: Empty bufferSize")
	ErrVkFFTFailEmptyBufferSize            = errors.New("vkFFT: Empty buffer")
	ErrVkFFTFailEmptyBuffer                = errors.New("vkFFT: Empty buffer")
	ErrVkFFTFailEmptyTempBufferSize        = errors.New("vkFFT: Empty tempBufferSize")
	ErrVkFFTFailEmptyTempBuffer            = errors.New("vkFFT: Empty tempBuffer")
	ErrVkFFTFailEmptyInputBufferSize       = errors.New("vkFFT: Empty inputBufferSize")
	ErrVkFFTFailEmptyInputBuffer           = errors.New("vkFFT: Empty inputBuffer")
	ErrVkFFTFailEmptyOutputBufferSize      = errors.New("vkFFT: Empty outputBufferSize")
	ErrVkFFTFailEmptyOutputBuffer          = errors.New("vkFFT: Empty outputBuffer")
	ErrVkFFTFailEmptyKernelSize            = errors.New("vkFFT: Empty kernelSize")
	ErrVkFFTFailEmptyKernel                = errors.New("vkFFT: Empty kernel")
	ErrVkFFTFailUnsupportedRadix           = errors.New("vkFFT: Unsupported radix")
	ErrVkFFTFailUnsupportedFFTLength       = errors.New("vkFFT: Unsupported FFT length")
	ErrVkFFTFailUnsupportedFFTLengthR2C    = errors.New("vkFFT: Unsupported r2c FFT length")
	ErrVkFFTFailUnsupportedFFTLengthDCT    = errors.New("vkFFT: Unsupported DCT length")
	ErrVkFFTFailUnsupportedFFTOmit         = errors.New("vkFFT: Unsupported FFT omit")
	ErrVkFFTFailAllocate                   = errors.New("vkFFT: Failed to allocate")
	ErrVkFFTFailMapMemory                  = errors.New("vkFFT: Failed to map memory")
	ErrVkFFTFailAllocateCommandBuffers     = errors.New("vkFFT: Failed to allocate command buffers")
	ErrVkFFTFailBeginCommandBuffer         = errors.New("vkFFT: Failed to begin command buffer")
	ErrVkFFTFailEndCommandBuffer           = errors.New("vkFFT: Failed to end command buffer")
	ErrVkFFTFailSubmitQueue                = errors.New("vkFFT: Failed to submit queue")
	ErrVkFFTFailWaitForFences              = errors.New("vkFFT: Failed to wait for fences")
	ErrVkFFTFailResetFences                = errors.New("vkFFT: Failed to reset fences")
	ErrVkFFTFailCreateDescriptorPool       = errors.New("vkFFT: Failed to create descriptor pool")
	ErrVkFFTFailCreateDescriptorSetLayout  = errors.New("vkFFT: Failed to create descriptor set layout")
	ErrVkFFTFailAllocateDescriptorSets     = errors.New("vkFFT: Failed to allocate descriptor sets")
	ErrVkFFTFailCreatePipelineLayout       = errors.New("vkFFT: Failed to create pipeline layout")
	ErrVkFFTFailShaderPreprocess           = errors.New("vkFFT: Failed to preprocess shaer")
	ErrVkFFTFailShaderParse                = errors.New("vkFFT: Failed to parse shader")
	ErrVkFFTFailShaderLink                 = errors.New("vkFFT: Failed to link shader")
	ErrVkFFTFailSPIRVGenerate              = errors.New("vkFFT: Failed to generate SPIRV")
	ErrVkFFTFailCreateShaderModule         = errors.New("vkFFT: Failed to create shader module")
	ErrVkFFTFailCreateInstance             = errors.New("vkFFT: Failed to create instance")
	ErrVkFFTFailSetupDebugMessenger        = errors.New("vkFFT: Failed to setup debug messenger")
	ErrVkFFTFailFindPhysicalDevice         = errors.New("vkFFT: Failed to find physical device")
	ErrVkFFTFailCreateDevice               = errors.New("vkFFT: Failed to create device")
	ErrVkFFTFailCreateFence                = errors.New("vkFFT: Failed to create fence")
	ErrVkFFTFailCreateCommandPool          = errors.New("vkFFT: Failed to create command pool")
	ErrVkFFTFailCreateBuffer               = errors.New("vkFFT: Failed to create buffer")
	ErrVkFFTFailAllocateMemory             = errors.New("vkFFT: Failed to allocate memory")
	ErrVkFFTFailBindBufferMemory           = errors.New("vkFFT: Failed to bind buffer memory")
	ErrVkFFTFailFindMemory                 = errors.New("vkFFT: Failed to find memory")
	ErrVkFFTFailSynchronize                = errors.New("vkFFT: Failed to synchronize")
	ErrVkFFTFailCopy                       = errors.New("vkFFT: Failed to copy")
	ErrVkFFTFailCreateProgram              = errors.New("vkFFT: Failed to create program")
	ErrVkFFTFailCompileProgram             = errors.New("vkFFT: Failed to compile program")
	ErrVkFFTFailGetCodeSize                = errors.New("vkFFT: Failed to get code size")
	ErrVkFFTFailGetCode                    = errors.New("vkFFT: Failed to get code")
	ErrVkFFTFailDestroyProgram             = errors.New("vkFFT: Failed to destroy program")
	ErrVkFFTFailLoadModule                 = errors.New("vkFFT: Failed to load module")
	ErrVkFFTFailGetFunction                = errors.New("vkFFT: Failed to get function")
	ErrVkFFTFailSetDynamicSharedMemory     = errors.New("vkFFT: Failed to set dynamic shared memory")
	ErrVkFFTFailModuleGetGlobal            = errors.New("vkFFT: Module failed to get global")
	ErrVkFFTFailLaunchKernel               = errors.New("vkFFT: Failed to launch kernel")
	ErrVkFFTFailEventRecord                = errors.New("vkFFT: Failed to record event")
	ErrVkFFTFailAddNameExpression          = errors.New("vkFFT: Failed to add name expression")
	ErrVkFFTFailInitialize                 = errors.New("vkFFT: Failed to initialize")
	ErrVkFFTFailSetDeviceID                = errors.New("vkFFT: Failed to set device id")
	ErrVkFFTFailGetDevice                  = errors.New("vkFFT: Failed to get device id")
	ErrVkFFTFailCreateContext              = errors.New("vkFFT: Failed to create context")
	ErrVkFFTFailCreatePipeline             = errors.New("vkFFT: Failed to crete pipeline")
	ErrVkFFTFailSetKernelArg               = errors.New("vkFFT: Failed to set kernel argument")
	ErrVkFFTFailCreateCommandQueue         = errors.New("vkFFT: Failed to create command queue")
	ErrVkFFTFailReleaseCommandQueue        = errors.New("vkFFT: Failed to release command queue")
	ErrVkFFTFailEnumerateDevices           = errors.New("vkFFT: Failed to enumerate devices")
	ErrVkFFTFailGetAttribute               = errors.New("vkFFT: Fail to get attribute")
	ErrVkFFTFailCreateEvent                = errors.New("vkFFT: Failed to create event")
)

var errorMapVkFFT = map[C.VkFFTResult]error{
	C.VKFFT_SUCCESS:                                      nil,
	C.VKFFT_ERROR_MALLOC_FAILED:                          ErrVkFFTFailMalloc,
	C.VKFFT_ERROR_INSUFFICIENT_CODE_BUFFER:               ErrVkFFTFailInsufficientCodeBuffer,
	C.VKFFT_ERROR_INSUFFICIENT_TEMP_BUFFER:               ErrVkFFTFailInsufficientTempBuffer,
	C.VKFFT_ERROR_PLAN_NOT_INITIALIZED:                   ErrVkFFTFailPlanNotInitialized,
	C.VKFFT_ERROR_NULL_TEMP_PASSED:                       ErrVkFFTFailNullTemp,
	C.VKFFT_ERROR_INVALID_PHYSICAL_DEVICE:                ErrVkFFTFailInvalidPhysicalDevice,
	C.VKFFT_ERROR_INVALID_DEVICE:                         ErrVkFFTFailInvalidDevice,
	C.VKFFT_ERROR_INVALID_QUEUE:                          ErrVkFFTFailInvalidQueue,
	C.VKFFT_ERROR_INVALID_COMMAND_POOL:                   ErrVkFFTFailInvalidCommandPool,
	C.VKFFT_ERROR_INVALID_FENCE:                          ErrVkFFTFailInvalidFence,
	C.VKFFT_ERROR_ONLY_FORWARD_FFT_INITIALIZED:           ErrVkFFTFailOnlyForwardFFTInitialized,
	C.VKFFT_ERROR_ONLY_INVERSE_FFT_INITIALIZED:           ErrVkFFTFailOnlyBackwardFFTInitialized,
	C.VKFFT_ERROR_INVALID_CONTEXT:                        ErrVkFFTFailInvalidContext,
	C.VKFFT_ERROR_INVALID_PLATFORM:                       ErrVkFFTFailInvalidPlatform,
	C.VKFFT_ERROR_EMPTY_FFTdim:                           ErrVkFFTFailEmptyFFTdim,
	C.VKFFT_ERROR_EMPTY_size:                             ErrVkFFTFailEmptySize,
	C.VKFFT_ERROR_EMPTY_bufferSize:                       ErrVkFFTFailEmptyBufferSize,
	C.VKFFT_ERROR_EMPTY_buffer:                           ErrVkFFTFailEmptyBuffer,
	C.VKFFT_ERROR_EMPTY_tempBufferSize:                   ErrVkFFTFailEmptyTempBufferSize,
	C.VKFFT_ERROR_EMPTY_tempBuffer:                       ErrVkFFTFailEmptyTempBuffer,
	C.VKFFT_ERROR_EMPTY_inputBufferSize:                  ErrVkFFTFailEmptyInputBufferSize,
	C.VKFFT_ERROR_EMPTY_inputBuffer:                      ErrVkFFTFailEmptyInputBuffer,
	C.VKFFT_ERROR_EMPTY_outputBufferSize:                 ErrVkFFTFailEmptyOutputBufferSize,
	C.VKFFT_ERROR_EMPTY_outputBuffer:                     ErrVkFFTFailEmptyOutputBuffer,
	C.VKFFT_ERROR_EMPTY_kernelSize:                       ErrVkFFTFailEmptyKernelSize,
	C.VKFFT_ERROR_EMPTY_kernel:                           ErrVkFFTFailEmptyKernel,
	C.VKFFT_ERROR_UNSUPPORTED_RADIX:                      ErrVkFFTFailUnsupportedRadix,
	C.VKFFT_ERROR_UNSUPPORTED_FFT_LENGTH:                 ErrVkFFTFailUnsupportedFFTLength,
	C.VKFFT_ERROR_UNSUPPORTED_FFT_LENGTH_R2C:             ErrVkFFTFailUnsupportedFFTLengthR2C,
	C.VKFFT_ERROR_UNSUPPORTED_FFT_LENGTH_DCT:             ErrVkFFTFailUnsupportedFFTLengthDCT,
	C.VKFFT_ERROR_UNSUPPORTED_FFT_OMIT:                   ErrVkFFTFailUnsupportedFFTOmit,
	C.VKFFT_ERROR_FAILED_TO_ALLOCATE:                     ErrVkFFTFailAllocate,
	C.VKFFT_ERROR_FAILED_TO_MAP_MEMORY:                   ErrVkFFTFailMapMemory,
	C.VKFFT_ERROR_FAILED_TO_ALLOCATE_COMMAND_BUFFERS:     ErrVkFFTFailAllocateCommandBuffers,
	C.VKFFT_ERROR_FAILED_TO_BEGIN_COMMAND_BUFFER:         ErrVkFFTFailBeginCommandBuffer,
	C.VKFFT_ERROR_FAILED_TO_END_COMMAND_BUFFER:           ErrVkFFTFailEndCommandBuffer,
	C.VKFFT_ERROR_FAILED_TO_SUBMIT_QUEUE:                 ErrVkFFTFailSubmitQueue,
	C.VKFFT_ERROR_FAILED_TO_WAIT_FOR_FENCES:              ErrVkFFTFailWaitForFences,
	C.VKFFT_ERROR_FAILED_TO_RESET_FENCES:                 ErrVkFFTFailResetFences,
	C.VKFFT_ERROR_FAILED_TO_CREATE_DESCRIPTOR_POOL:       ErrVkFFTFailCreateDescriptorPool,
	C.VKFFT_ERROR_FAILED_TO_CREATE_DESCRIPTOR_SET_LAYOUT: ErrVkFFTFailCreateDescriptorSetLayout,
	C.VKFFT_ERROR_FAILED_TO_ALLOCATE_DESCRIPTOR_SETS:     ErrVkFFTFailAllocateDescriptorSets,
	C.VKFFT_ERROR_FAILED_TO_CREATE_PIPELINE_LAYOUT:       ErrVkFFTFailCreatePipelineLayout,
	C.VKFFT_ERROR_FAILED_SHADER_PREPROCESS:               ErrVkFFTFailShaderPreprocess,
	C.VKFFT_ERROR_FAILED_SHADER_PARSE:                    ErrVkFFTFailShaderParse,
	C.VKFFT_ERROR_FAILED_SHADER_LINK:                     ErrVkFFTFailShaderLink,
	C.VKFFT_ERROR_FAILED_SPIRV_GENERATE:                  ErrVkFFTFailSPIRVGenerate,
	C.VKFFT_ERROR_FAILED_TO_CREATE_SHADER_MODULE:         ErrVkFFTFailCreateShaderModule,
	C.VKFFT_ERROR_FAILED_TO_CREATE_INSTANCE:              ErrVkFFTFailCreateInstance,
	C.VKFFT_ERROR_FAILED_TO_SETUP_DEBUG_MESSENGER:        ErrVkFFTFailSetupDebugMessenger,
	C.VKFFT_ERROR_FAILED_TO_FIND_PHYSICAL_DEVICE:         ErrVkFFTFailFindPhysicalDevice,
	C.VKFFT_ERROR_FAILED_TO_CREATE_DEVICE:                ErrVkFFTFailCreateDevice,
	C.VKFFT_ERROR_FAILED_TO_CREATE_FENCE:                 ErrVkFFTFailCreateFence,
	C.VKFFT_ERROR_FAILED_TO_CREATE_COMMAND_POOL:          ErrVkFFTFailCreateCommandPool,
	C.VKFFT_ERROR_FAILED_TO_CREATE_BUFFER:                ErrVkFFTFailCreateBuffer,
	C.VKFFT_ERROR_FAILED_TO_ALLOCATE_MEMORY:              ErrVkFFTFailAllocateMemory,
	C.VKFFT_ERROR_FAILED_TO_BIND_BUFFER_MEMORY:           ErrVkFFTFailBindBufferMemory,
	C.VKFFT_ERROR_FAILED_TO_FIND_MEMORY:                  ErrVkFFTFailFindMemory,
	C.VKFFT_ERROR_FAILED_TO_SYNCHRONIZE:                  ErrVkFFTFailSynchronize,
	C.VKFFT_ERROR_FAILED_TO_COPY:                         ErrVkFFTFailCopy,
	C.VKFFT_ERROR_FAILED_TO_CREATE_PROGRAM:               ErrVkFFTFailCreateProgram,
	C.VKFFT_ERROR_FAILED_TO_COMPILE_PROGRAM:              ErrVkFFTFailCompileProgram,
	C.VKFFT_ERROR_FAILED_TO_GET_CODE_SIZE:                ErrVkFFTFailGetCodeSize,
	C.VKFFT_ERROR_FAILED_TO_GET_CODE:                     ErrVkFFTFailGetCode,
	C.VKFFT_ERROR_FAILED_TO_DESTROY_PROGRAM:              ErrVkFFTFailDestroyProgram,
	C.VKFFT_ERROR_FAILED_TO_LOAD_MODULE:                  ErrVkFFTFailLoadModule,
	C.VKFFT_ERROR_FAILED_TO_GET_FUNCTION:                 ErrVkFFTFailGetFunction,
	C.VKFFT_ERROR_FAILED_TO_SET_DYNAMIC_SHARED_MEMORY:    ErrVkFFTFailSetDynamicSharedMemory,
	C.VKFFT_ERROR_FAILED_TO_MODULE_GET_GLOBAL:            ErrVkFFTFailModuleGetGlobal,
	C.VKFFT_ERROR_FAILED_TO_LAUNCH_KERNEL:                ErrVkFFTFailLaunchKernel,
	C.VKFFT_ERROR_FAILED_TO_EVENT_RECORD:                 ErrVkFFTFailEventRecord,
	C.VKFFT_ERROR_FAILED_TO_ADD_NAME_EXPRESSION:          ErrVkFFTFailAddNameExpression,
	C.VKFFT_ERROR_FAILED_TO_INITIALIZE:                   ErrVkFFTFailInitialize,
	C.VKFFT_ERROR_FAILED_TO_SET_DEVICE_ID:                ErrVkFFTFailSetDeviceID,
	C.VKFFT_ERROR_FAILED_TO_GET_DEVICE:                   ErrVkFFTFailGetDevice,
	C.VKFFT_ERROR_FAILED_TO_CREATE_CONTEXT:               ErrVkFFTFailCreateContext,
	C.VKFFT_ERROR_FAILED_TO_CREATE_PIPELINE:              ErrVkFFTFailCreatePipeline,
	C.VKFFT_ERROR_FAILED_TO_SET_KERNEL_ARG:               ErrVkFFTFailSetKernelArg,
	C.VKFFT_ERROR_FAILED_TO_CREATE_COMMAND_QUEUE:         ErrVkFFTFailCreateCommandQueue,
	C.VKFFT_ERROR_FAILED_TO_RELEASE_COMMAND_QUEUE:        ErrVkFFTFailReleaseCommandQueue,
	C.VKFFT_ERROR_FAILED_TO_ENUMERATE_DEVICES:            ErrVkFFTFailEnumerateDevices,
	C.VKFFT_ERROR_FAILED_TO_GET_ATTRIBUTE:                ErrVkFFTFailGetAttribute,
	C.VKFFT_ERROR_FAILED_TO_CREATE_EVENT:                 ErrVkFFTFailCreateEvent,
}

type VkfftDirection int

const (
	VkfftForwardDirection  VkfftDirection = C.VKFFT_FORWARD_TRANSFORM
	VkfftBackwardDirection VkfftDirection = C.VKFFT_BACKWARD_TRANSFORM
)

type VkfftPlan struct {
	vkfftPlanStruct C.interfaceFFTPlan
}

func (vPlan *VkfftPlan) GetPlanPointer() *C.interfaceFFTPlan {
	return &vPlan.vkfftPlanStruct
}

func NewVkFFTPlan(ctx *Context) *VkfftPlan {
	var outPlan *C.interfaceFFTPlan
	outPlan = C.vkfftCreateR2CFFTPlan(ctx.clContext)
	return &VkfftPlan{*outPlan}
}

func (p *VkfftPlan) VkFFTSetFFTPlanSize(lengths []int) {
	var cLengths [3]C.size_t
	dim := len(lengths)
	if dim > 3 {
		fmt.Printf("lengths is longer than expected. Will only use the first 3 entries!\n")
	}
	cLengths[0] = (C.size_t)(lengths[0])
	cLengths[1] = 1
	cLengths[2] = 1
	if dim > 1 {
		cLengths[1] = (C.size_t)(lengths[1])
	}
	if dim > 2 {
		cLengths[2] = (C.size_t)(lengths[2])
	}
	C.vkfftSetFFTPlanSize(p.GetPlanPointer(), &cLengths[0])
}

func (p *VkfftPlan) VkFFTEnqueueTransformUnsafe(dir VkfftDirection, input []*MemObject, output []*MemObject) error {
	return toError(C.vkfftEnqueueTransform(p.GetPlanPointer(), (C.vkfft_transform_dir)(dir), &(input[0].clMem), &(output[0].clMem)))
}

func (p *VkfftPlan) EnqueueForwardTransform(input []*MemObject, output []*MemObject) error {
	return p.VkFFTEnqueueTransformUnsafe(VkfftForwardDirection, input, output)
}

func (p *VkfftPlan) EnqueueBackwardTransform(input []*MemObject, output []*MemObject) error {
	return p.VkFFTEnqueueTransformUnsafe(VkfftBackwardDirection, input, output)
}

func (plan *VkfftPlan) Destroy() {
	C.vkfftDestroyFFTPlan(plan.GetPlanPointer())
}
