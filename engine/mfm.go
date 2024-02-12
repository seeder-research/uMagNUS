package engine

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
)

var (
	MFM        = NewScalarField("MFM", "arb.", "MFM image", SetMFM)
	MFMLift    inputValue
	MFMTipSize inputValue
	mfmconv_   *opencl.MFMConvolution
)

func init() {
	MFMLift = numParam(50e-9, "MFMLift", "m", reinitmfmconv)
	MFMTipSize = numParam(1e-3, "MFMDipole", "m", reinitmfmconv)
	DeclLValue("MFMLift", &MFMLift, "MFM lift height")
	DeclLValue("MFMDipole", &MFMTipSize, "Height of vertically magnetized part of MFM tip")
}

func SetMFM(dst *data.Slice) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting from all queues to finish in setmfm: %+v \n", err)
	}

	buf := opencl.Buffer(3, Mesh().Size())
	defer opencl.Recycle(buf)
	if mfmconv_ == nil {
		reinitmfmconv()
	}

	msat := Msat.MSlice()
	defer msat.Recycle()

	// checkout queues and execute kernel
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
	seqQueue := opencl.ClCmdQueue[0]

	mfmconv_.Exec(buf, M.Buffer(), geometry.Gpu(), msat, seqQueue, nil)
	// sync queues to seqQueue
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})

	opencl.Madd3(dst, buf.Comp(0), buf.Comp(1), buf.Comp(2), 1, 1, 1, queues, nil)

	// sync in before returning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting from all queues to finish after setmfm: %+v \n", err)
	}
}

func reinitmfmconv() {
	SetBusy(true)
	defer SetBusy(false)
	seqQueue := opencl.ClCmdQueue[0]
	if mfmconv_ == nil {
		mfmconv_ = opencl.NewMFM(Mesh(), MFMLift.v, MFMTipSize.v, *Flag_cachedir, seqQueue, nil)
	} else {
		mfmconv_.Reinit(MFMLift.v, MFMTipSize.v, *Flag_cachedir, seqQueue, nil)
	}
	// sync before returning
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish in reinitmfmconv: %+v \n", err)
	}
}
