package engine

import (
	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

func fixedPtIterations(hFac float32, Y, ks *data.Slice) (float64, float64, int) {
	// For implicit solvers, there is a need to solve for
	// k_{s} = f( (t_{n} + c_{s} h),
	//            (y_{n} + \sum_{i=1}^{s}( a_{s,i} h k_{i}) ),
	// where c_{s} and a_{s, i} are coefficients of the solver
	// method typically listed in a Butcher tableau format.
	// This function evaluates k_{s} = f( T, Y + k_{s}) to find
	// the solution to the implicit step using the fixed point
	// method.

	// Expectation of function:
	//   - When called, M is at the initial guess for the fixed point
	//     iteration.
	//   - hFac is h multiplied by the coefficient of stage s
	//   - G() is g() excluding the term in hFac * k_{s}
	//   - The function solves for k_{s}, which is returned in dy1
	//   - ypred = y_{n} + G(h, k_{i}), i running from 1 to s - 1

	kPrev := opencl.Buffer(VECTOR, ks.Size())
	errVector := opencl.Buffer(VECTOR, ks.Size())
	defer opencl.Recycle(kPrev)
	defer opencl.Recycle(errVector)
	torqueFn(kPrev)

	yPred := M.Buffer()

	// Initialize loop state so at least one iteration
	// of for loop gets executed
	ErrIter := float64(0.0)
	relErr := float64(0.0)
	iterate := true
	Niters := 0

	// checkout queues for executing kernel
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}
	seqQueue := opencl.ClCmdQueue[0]
	// sync in the beginning
	opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})

	// fixed point iterations until converence criterion reached
	for ; iterate && (Niters < NConv); Niters++ {
		// Update guess
		opencl.Madd2(yPred, Y, kPrev, 1.0, hFac, queues, nil) // y = y0 + dt * dy
		// sync queues
		opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
		M.normalize()
		torqueFn(ks)

		// Calculate error as the difference in calculated predictions
		// in consecutive fixed point iterations
		ErrIter = float64(opencl.MaxVecDiff(ks, kPrev, seqQueue, nil))
		ksNorm := float64(opencl.MaxVecNorm(ks, seqQueue, nil))
		relErr = RelErrConv * ksNorm
		iterate = (ErrIter > AbsErrConv) && (ErrIter > relErr)

		if iterate && (Niters < NConv-1) {
			// Record fixed point result for next iteration
			data.Copy(kPrev, ks)

			// sync before next iteration
			opencl.SyncQueues(queues, []*cl.CommandQueue{seqQueue})
		}
	}
	if Niters == NConv {
		util.Log("fixed point iterations exceeded limit!")
	}
	return ErrIter, relErr, Niters
}
