package oclRAND

/*
 THIS FILE IS AUTO-GENERATED BY OCL2GO.
 EDITING IS FUTILE.
*/

import (
	"github.com/seeder-research/uMagNUS/cl"
	"github.com/seeder-research/uMagNUS/timer"
	"sync"
	"unsafe"
)

// Stores the arguments for mtgp32_uint32 kernel invocation
type mtgp32_uint32_args_t struct {
	arg_param_tbl         unsafe.Pointer
	arg_temper_tbl        unsafe.Pointer
	arg_single_temper_tbl unsafe.Pointer
	arg_pos_tbl           unsafe.Pointer
	arg_sh1_tbl           unsafe.Pointer
	arg_sh2_tbl           unsafe.Pointer
	arg_d_status          unsafe.Pointer
	arg_d_data            unsafe.Pointer
	arg_size              int
	argptr                [9]unsafe.Pointer
	sync.Mutex
}

// Stores the arguments for mtgp32_uint32 kernel invocation
var mtgp32_uint32_args mtgp32_uint32_args_t

func init() {
	// OpenCL driver kernel call wants pointers to arguments, set them up once.
	mtgp32_uint32_args.argptr[0] = unsafe.Pointer(&mtgp32_uint32_args.arg_param_tbl)
	mtgp32_uint32_args.argptr[1] = unsafe.Pointer(&mtgp32_uint32_args.arg_temper_tbl)
	mtgp32_uint32_args.argptr[2] = unsafe.Pointer(&mtgp32_uint32_args.arg_single_temper_tbl)
	mtgp32_uint32_args.argptr[3] = unsafe.Pointer(&mtgp32_uint32_args.arg_pos_tbl)
	mtgp32_uint32_args.argptr[4] = unsafe.Pointer(&mtgp32_uint32_args.arg_sh1_tbl)
	mtgp32_uint32_args.argptr[5] = unsafe.Pointer(&mtgp32_uint32_args.arg_sh2_tbl)
	mtgp32_uint32_args.argptr[6] = unsafe.Pointer(&mtgp32_uint32_args.arg_d_status)
	mtgp32_uint32_args.argptr[7] = unsafe.Pointer(&mtgp32_uint32_args.arg_d_data)
	mtgp32_uint32_args.argptr[8] = unsafe.Pointer(&mtgp32_uint32_args.arg_size)
}

// Wrapper for mtgp32_uint32 OpenCL kernel, asynchronous.
func k_mtgp32_uint32_async(param_tbl unsafe.Pointer, temper_tbl unsafe.Pointer, single_temper_tbl unsafe.Pointer, pos_tbl unsafe.Pointer, sh1_tbl unsafe.Pointer, sh2_tbl unsafe.Pointer, d_status unsafe.Pointer, d_data unsafe.Pointer, size int, cfg *config, events []*cl.Event) *cl.Event {
	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Start("mtgp32_uint32")
	}

	mtgp32_uint32_args.Lock()
	defer mtgp32_uint32_args.Unlock()

	mtgp32_uint32_args.arg_param_tbl = param_tbl
	mtgp32_uint32_args.arg_temper_tbl = temper_tbl
	mtgp32_uint32_args.arg_single_temper_tbl = single_temper_tbl
	mtgp32_uint32_args.arg_pos_tbl = pos_tbl
	mtgp32_uint32_args.arg_sh1_tbl = sh1_tbl
	mtgp32_uint32_args.arg_sh2_tbl = sh2_tbl
	mtgp32_uint32_args.arg_d_status = d_status
	mtgp32_uint32_args.arg_d_data = d_data
	mtgp32_uint32_args.arg_size = size

	SetKernelArgWrapper("mtgp32_uint32", 0, param_tbl)
	SetKernelArgWrapper("mtgp32_uint32", 1, temper_tbl)
	SetKernelArgWrapper("mtgp32_uint32", 2, single_temper_tbl)
	SetKernelArgWrapper("mtgp32_uint32", 3, pos_tbl)
	SetKernelArgWrapper("mtgp32_uint32", 4, sh1_tbl)
	SetKernelArgWrapper("mtgp32_uint32", 5, sh2_tbl)
	SetKernelArgWrapper("mtgp32_uint32", 6, d_status)
	SetKernelArgWrapper("mtgp32_uint32", 7, d_data)
	SetKernelArgWrapper("mtgp32_uint32", 8, size)

	//	args := mtgp32_uint32_args.argptr[:]
	event := LaunchKernel("mtgp32_uint32", cfg.Grid, cfg.Block, events)

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Stop("mtgp32_uint32")
	}

	return event
}