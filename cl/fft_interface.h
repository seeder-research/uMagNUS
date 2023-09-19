#ifndef __FFT_INTERFACE__
#define __FFT_INTERFACE__
// Interface to vkFFT with convenice functions
#define VKFFT_BACKEND     3
#define __SIZEOF_HALF__   2
#define __SIZEOF_FLOAT__  4
#define __SIZEOF_DOUBLE__ 8
#include <stdbool.h>
#include "vkFFT.h"

typedef enum vkfft_transform_dir {
    VKFFT_FORWARD_TRANSFORM    = -1,
    VKFFT_BACKWARD_TRANSFORM   =  1
} vkfft_transform_dir;

// Use a plan structure
struct interfaceFFTPlan {
    VkFFTConfiguration*   config;
    VkFFTApplication*     app;
    bool                  isBaked;
    bool                  notInit;
    VkFFTLaunchParams*    lParams;
    cl_platform_id        platform;
    cl_device_id          device;
    cl_context            context;
    cl_command_queue      commandQueue;
    int                   dataType;
    uint64_t              inputBufferSize;
    uint64_t              outputBufferSize;
};

typedef struct interfaceFFTPlan interfaceFFTPlan;

// Interface functions for plan creation
interfaceFFTPlan* vkfftCreateDefaultFFTPlan(cl_context ctx);
interfaceFFTPlan* vkfftCreateR2CFFTPlan(cl_context ctx);

// Interface function for modifying the FFT plan details
void vkfftSetFFTPlanBufferSizes(interfaceFFTPlan* plan);
void vkfftSetFFTPlanDataType(interfaceFFTPlan* plan, int dataType);
void vkfftSetFFTPlanSize(interfaceFFTPlan* plan, size_t lengths[3]);

// Interface functions to make the library compatible with other conventional FFT libraries
VkFFTResult vkfftBakeFFTPlan(interfaceFFTPlan* plan);
VkFFTResult vkfftEnqueueTransform(interfaceFFTPlan* plan, vkfft_transform_dir dir, cl_mem* input, cl_mem* dst);
void vkfftDestroyFFTPlan(interfaceFFTPlan* plan);

// Basic function to return a FFT plan.
// This flow is similar to other FFT libraries such as FFTW, cuFFT, clFFT, rocFFT.
interfaceFFTPlan* vkfftCreateDefaultFFTPlan(cl_context ctx) {
    interfaceFFTPlan* plan = (interfaceFFTPlan*)calloc(1, sizeof(interfaceFFTPlan));
    // Empty plan
    plan->config  = (VkFFTConfiguration*)calloc(1, sizeof(VkFFTConfiguration));
    plan->app     = (VkFFTApplication*)calloc(1, sizeof(VkFFTApplication));
    plan->lParams = (VkFFTLaunchParams*)calloc(1, sizeof(VkFFTLaunchParams));

    cl_int res;
    // Grab required information from context given...
    plan->context = ctx;

    // Get device ID from context
    size_t numCount = 0;
    res = clGetContextInfo(plan->context, CL_CONTEXT_DEVICES, sizeof(cl_device_id), &plan->device, &numCount);
    if (res != CL_SUCCESS) {
        free(plan);
        return NULL;
    }

    res = clGetDeviceInfo(plan->device, CL_DEVICE_PLATFORM, sizeof(cl_platform_id), &plan->platform, &numCount);
    if (res != CL_SUCCESS) {
        free(plan);
        return NULL;
    }

    // Create a command queue for the plan
    plan->commandQueue = clCreateCommandQueue(plan->context, plan->device, 0, &res);
    if (res != CL_SUCCESS) {
        free(plan);
        return NULL;
    }

    // Update internal pointers
    plan->config->platform      = &plan->platform;
    plan->config->context       = &plan->context;
    plan->config->device        = &plan->device;
    plan->lParams->commandQueue = &plan->commandQueue;

    // Default to 3D plan but with all dimensions to be 1
    plan->config->FFTdim  = 3;
    plan->config->size[0] = 1;
    plan->config->size[1] = 1;
    plan->config->size[2] = 1;

    // Default to in-place C2C transform where inverse
    // is not normalized
    plan->config->performR2C                 = 0;
    plan->config->inverseReturnToInputBuffer = 0;
    plan->config->normalize                  = 0;

    // Default to out-of-place transform
    plan->config->isInputFormatted = 1;
    plan->config->inputBufferNum = 1;
    plan->config->bufferNum = 1;

    // Default to use on-chip units for sin and cos
    plan->config->useLUT = 0;

    // Default to float
    // half   = -1
    // float  =  0
    // double =  1
    plan->dataType                = 0;
    plan->config->halfPrecision   = false;
    plan->config->doublePrecision = false;

    // Initialize flag to ensure plan is baked before execution can happen
    plan->isBaked = false;
    plan->notInit = false;
    return plan;
}

// A specialized function to return a FFT plan that computes
// R2C (forward) and C2R (backward) transforms
interfaceFFTPlan* vkfftCreateR2CFFTPlan(cl_context ctx) {
    interfaceFFTPlan* plan = vkfftCreateDefaultFFTPlan(ctx);

    plan->config->performR2C                 = 1;
    plan->config->inverseReturnToInputBuffer = 1;
    return plan;
}

// Interface function to set up the data type for the FFT
void vkfftSetFFTPlanDataType(interfaceFFTPlan* plan, int dataType) {
    // Default to float
    // half   = -1
    // float  =  0
    // double =  1
    plan->dataType = dataType;
    if (dataType < 0) {
        plan->config->halfPrecision   = true;
        plan->config->doublePrecision = false;
    } else if (dataType > 0) {
        plan->config->halfPrecision   = false;
        plan->config->doublePrecision = true;
        // Uncomment next line to force use of LUT to compute sin/cos for double
        // plan->config->useLUT = 1;
    } else {
        plan->config->halfPrecision   = false;
        plan->config->doublePrecision = false;
    }
    vkfftSetFFTPlanBufferSizes(plan);
}

// Interface function to set up the FFT sizes
void vkfftSetFFTPlanSize(interfaceFFTPlan* plan, size_t lengths[3]) {
    // If the plan was previously baked, we need to clean up the plan
    if (plan->isBaked) {
        deleteVkFFT(plan->app);
        plan->app = (VkFFTApplication*)calloc(1, sizeof(VkFFTApplication));
    }
    plan->isBaked = false;

    ////// Order the lengths of the FFT so it can be "fast"

    // Default to 3D first but set all sizes to 1
    plan->config->FFTdim  = 3;
    plan->config->size[0] = 1;
    plan->config->size[1] = 1;
    plan->config->size[2] = 1;

    // Find out the desired dimensionality of the FFT
    if (lengths[0] == 1) { plan->config->FFTdim--; }
    if (lengths[1] == 1) { plan->config->FFTdim--; }
    if (lengths[2] == 1) { plan->config->FFTdim--; }

    // Catch when all entries of lengths[] is 1
    if (plan->config->FFTdim == 0) {
        plan->config->FFTdim = 1; // the FFT has all lengths to be 1
    } else if (plan->config->FFTdim == 1) { // Case where FFT is 1D
        // Find the entry of lengths[] that is not 1 and assign to
        // config.size[0] (the other entries default to 1 from before)
        if (lengths[0] != 1) {
            plan->config->size[0] = lengths[0];
        } else if (lengths[1] != 1) {
            plan->config->size[0] = lengths[1];
        } else {
            plan->config->size[0] = lengths[2];
        }
    } else if (plan->config->FFTdim == 2) { // Case where FFT is 2D
        // Find the entry of lengths[] that is 1 and assign to remaining
        // to config.size[0] and config.size[1] (the remaining entry
        // default to 1 from before)
        if (lengths[0] == 1) {
            plan->config->size[0] = lengths[1];
            plan->config->size[1] = lengths[2];
        } else if (lengths[1] == 1) {
            plan->config->size[0] = lengths[0];
            plan->config->size[1] = lengths[2];
        } else {
            plan->config->size[0] = lengths[0];
            plan->config->size[1] = lengths[1];
        }
    } else { // Case where FFT is 3D
        plan->config->size[0] = lengths[0];
        plan->config->size[1] = lengths[1];
        plan->config->size[2] = lengths[2];
    }
    vkfftSetFFTPlanBufferSizes(plan);
}

// Function to determine the input and output buffer sizes
void vkfftSetFFTPlanBufferSizes(interfaceFFTPlan* plan) {
    // Input and output buffer sizes if transform is C2C
    plan->inputBufferSize  = plan->config->size[1] * plan->config->size[2];

    if (plan->dataType < 0) {
        plan->inputBufferSize  *= __SIZEOF_HALF__;
    } else if (plan->dataType > 0) {
        plan->inputBufferSize  *= __SIZEOF_DOUBLE__;
    } else {
        plan->inputBufferSize  *= __SIZEOF_FLOAT__;
    }

    plan->outputBufferSize = plan->inputBufferSize;

    // If plan is already defined as R2C, then we can set the input and output buffer sizes
    // as well as the strides
    if (plan->config->performR2C == 1) {
        plan->inputBufferSize  *= plan->config->size[0];
        plan->outputBufferSize *= 2 * (plan->config->size[0] / 2 + 1);
    } else { // Otherwise, plan is C2C
        plan->inputBufferSize  *= 2 * plan->config->size[0];
        plan->outputBufferSize  = plan->inputBufferSize;
    }

    // Update plan
    plan->config->inputBufferSize       = &plan->inputBufferSize;
    plan->config->inputBufferStride[0]  = plan->config->size[0];
    plan->config->inputBufferStride[1]  = plan->config->size[0]*plan->config->size[1];
    plan->config->inputBufferStride[2]  = plan->config->size[0]*plan->config->size[1]*plan->config->size[2];
    plan->config->bufferSize            = &plan->outputBufferSize;
}

// Interface to initializeVkFFT()
// Provide this function so that initialization can be checked prior to
// any execution
VkFFTResult vkfftBakeFFTPlan(interfaceFFTPlan* plan) {
    VkFFTResult res;
#if(__DEBUG__>0)
    printf("Begin initialization...\n");
#endif
    // If the plan was baked previously, the previous plan needs to be deleted
    if ((plan->app != NULL) && (plan->isBaked)) {
        deleteVkFFT(plan->app);
        plan->app = (VkFFTApplication*)calloc(1,sizeof(VkFFTApplication));
    }
    VkFFTConfiguration tmpConfig = *plan->config;
    res = initializeVkFFT(plan->app, tmpConfig);
#if(__DEBUG__>0)
    printf("    Done with initialization...\n");
#endif
    if (res == VKFFT_SUCCESS) {
        plan->isBaked = true;
    } else {
        plan->isBaked = false;
    }
    plan->notInit = true;
    return res;
}

// Interface function to perform a FFT.
// This function will ensure the plan is initialized prior to execution.
VkFFTResult vkfftEnqueueTransform(interfaceFFTPlan* plan, vkfft_transform_dir dir, cl_mem* input, cl_mem* dst) {
    // Set up buffers for input and output so that vkFFT can recognize them
    if (dir < 0) {
        plan->lParams->inputBuffer = input;
        plan->lParams->buffer = dst;
    } else {
        plan->lParams->inputBuffer = dst;
        plan->lParams->buffer = input;
    }

    VkFFTResult res;
    // Initialize the plan if it is not already initialized
    if (!plan->isBaked) {
        res = vkfftBakeFFTPlan(plan);
        if (res != VKFFT_SUCCESS) {
            return res;
        }
    }

    // Plan is guaranteed to be initialized so we launch the execution
    return VkFFTAppend(plan->app, dir, plan->lParams);
}

// Interface function to clean up
void vkfftDestroyFFTPlan(interfaceFFTPlan* plan) {
    if (plan->notInit) {
        deleteVkFFT(plan->app);
    }
    free(plan->config);
    free(plan->lParams);
}

cl_event vkfftGetPlanEvent(interfaceFFTPlan* plan) {
    return plan->app->configuration.queueEvent;
}

cl_command_queue vkfftPlanGetCommandQueue(interfaceFFTPlan* plan) {
    return plan->app->configuration.commandQueue[0];
}

cl_device_id vkfftPlanGetDevice(interfaceFFTPlan* plan) {
    return *(plan->app->configuration.device);
}

cl_int vkfftPlanQueueFinish(interfaceFFTPlan* plan) {
    return clFinish(plan->app->configuration.commandQueue[0]);
}

cl_int vkfftPlanQueueFlush(interfaceFFTPlan* plan) {
    return clFlush(plan->app->configuration.commandQueue[0]);
}

#endif // __FFT_INTERFACE__
