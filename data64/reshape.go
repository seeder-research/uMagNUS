package data64

// Array reshaping.

import "fmt"

// Re-interpret a contiguous array as a multi-dimensional array of given size.
// Underlying storage is shared.
func reshape(array []float64, size [3]int) [][][]float64 {
	Nx, Ny, Nz := size[0], size[1], size[2]
	if Nx*Ny*Nz != len(array) {
		panic(fmt.Errorf("reshape: size mismatch: %v*%v*%v != %v", Nx, Ny, Nz, len(array)))
	}
	sliced := make([][][]float64, Nz)
	for i := range sliced {
		sliced[i] = make([][]float64, Ny)
	}
	for i := range sliced {
		for j := range sliced[i] {
			sliced[i][j] = array[(i*Ny+j)*Nx+0 : (i*Ny+j)*Nx+Nx]
		}
	}
	return sliced
}
