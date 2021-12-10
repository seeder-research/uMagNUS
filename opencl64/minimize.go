package opencl64

import (
	"fmt"

	data "github.com/seeder-research/uMagNUS/data64"
	"github.com/seeder-research/uMagNUS/cl"
)

// m = 1 / (4 + τ²(m x H)²) [{4 - τ²(m x H)²} m - 4τ(m x m x H)]
// note: torque from LLNoPrecess has negative sign
func Minimize(m, m0, torque *data.Slice, dt float64) {
	N := m.Len()
	cfg := make1DConf(N)

	event := k_minimize_async(m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		m0.DevPtr(X), m0.DevPtr(Y), m0.DevPtr(Z),
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		dt, N, cfg,
		[](*cl.Event){m.GetEvent(X), m.GetEvent(Y), m.GetEvent(Z),
			m0.GetEvent(X), m0.GetEvent(Y), m0.GetEvent(Z),
			torque.GetEvent(X), torque.GetEvent(Y), torque.GetEvent(Z)})
	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	m0.SetEvent(X, event)
	m0.SetEvent(Y, event)
	m0.SetEvent(Z, event)
	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)
	err := cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in minimize: %+v \n", err)
	}
}
