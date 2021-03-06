package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
)

// m = 1 / (4 + τ²(m x H)²) [{4 - τ²(m x H)²} m - 4τ(m x m x H)]
// note: torque from LLNoPrecess has negative sign
func Minimize(m, m0, torque *data.Slice, dt float64) {
	N := m.Len()
	cfg := make1DConf(N)

	eventList := [](*cl.Event){}
	tmpEvt := torque.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = torque.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = torque.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m0.GetEvent(X)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m0.GetEvent(Y)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvt = m0.GetEvent(Z)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	if len(eventList) == 0 {
		eventList = nil
	}

	event := k_minimize_async(m.DevPtr(X), m.DevPtr(Y), m.DevPtr(Z),
		m0.DevPtr(X), m0.DevPtr(Y), m0.DevPtr(Z),
		torque.DevPtr(X), torque.DevPtr(Y), torque.DevPtr(Z),
		dt, N, cfg, eventList)

	m.SetEvent(X, event)
	m.SetEvent(Y, event)
	m.SetEvent(Z, event)
	m0.SetEvent(X, event)
	m0.SetEvent(Y, event)
	m0.SetEvent(Z, event)
	torque.SetEvent(X, event)
	torque.SetEvent(Y, event)
	torque.SetEvent(Z, event)

	if Debug {
		if err := cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in minimize: %+v \n", err)
		}
	}
}
