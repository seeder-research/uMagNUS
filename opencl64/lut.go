package opencl64

// Look-up tables holding per-region parameter values.
// LUT[regions[cellindex]] gives parameter value for cell.

import "unsafe"

type LUTPtr unsafe.Pointer    // points to 256 float64's
type LUTPtrs []unsafe.Pointer // elements point to 256 float64's
type SymmLUT unsafe.Pointer   // points to 256x256 symmetric matrix, only lower half stored. See exchange.cu
