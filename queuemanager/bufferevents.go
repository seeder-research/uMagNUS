package queuemanager

import (
	"context"
	"log"

	cl "github.com/seeder-research/uMagNUS/cl"
)

type BufferType interface {
	NComp() int
	SetEvent(index int, event *cl.Event)
	GetEvent(index int) *cl.Event
	SetReadEvents(index int, eventList []*cl.Event)
	InsertReadEvent(index int, event *cl.Event)
	RemoveReadEvent(index int, event *cl.Event)
	GetReadEvents(index int) []*cl.Event
	GetAllEvents(index int) []*cl.Event
}

// Get all event markers (producers and consumers)
// Use for producers to know when data is no longer required
// from consumers, and ensure data is written in sequence of
// producer queue
func GetAllEventsOfBuffers(list []BufferType) []*cl.Event {
	var outList []*cl.Event
	for _, buf := range list {
		if buf != nil {
			for i := 0; i < buf.NComp(); i++ {
				bufEventList := buf.GetAllEvents(i)
				outList = append(outList, bufEventList...)
			}
		}
	}
	return outList
}

// Get all event markers (producers only).
// Use for consumers to know when data from producer is ready
func GetEventsOfBuffers(list []BufferType) []*cl.Event {
	var outList []*cl.Event
	for _, buf := range list {
		if buf != nil {
			for i := 0; i < buf.NComp(); i++ {
				bufEventList := buf.GetEvent(i)
				outList = append(outList, bufEventList)
			}
		}
	}
	return outList
}

// Remove producer event from list of buffers
// Use after producer command completes
func RemoveEventFromBuffers(list []BufferType, ev *cl.Event) {
	for _, buf := range list {
		if buf != nil {
			for i := 0; i < buf.NComp(); i++ {
				if buf.GetEvent(i) == ev {
					buf.SetEvent(i, nil)
				}
			}
		}
	}
}

// Remove consumer event from list of buffers
// Use after consumer command completes
func RemoveReadEventFromBuffers(list []BufferType, ev *cl.Event) {
	for _, buf := range list {
		if buf != nil {
			for i := 0; i < buf.NComp(); i++ {
				buf.RemoveReadEvent(i, ev)
			}
		}
	}
}

// Insert producer event into list of buffers
// Use after producer command is launched but before it completes
func SetWriteEventToBuffers(list []BufferType, ev *cl.Event) {
	for _, buf := range list {
		if buf != nil {
			for i := 0; i < buf.NComp(); i++ {
				// All concumers and producers launched after this producer
				// needs to sync to this producer. So we do not need to track
				// consumers launched before this producer
				buf.SetEvent(i, ev)
				buf.SetReadEvents(i, []*cl.Event{})
			}
		}
	}
}

// Insert consumer event into list of buffers
// Use after consumer command is launched but before it completes
func AddReadEventsToBuffers(list []BufferType, ev *cl.Event) {
	for _, buf := range list {
		if buf != nil {
			for i := 0; i < buf.NComp(); i++ {
				buf.InsertReadEvent(i, ev)
			}
		}
	}
}

// /////////////////////////////////////////////////////////////////////////////////////
//
//	Single goroutine to support updates o events tracked in buffers
//
// /////////////////////////////////////////////////////////////////////////////////////
type EventAndBuffers struct {
	ev    *cl.Event
	cList []BufferType
	pList []BufferType
}

// Wait for list of events to complete and update
// the events in the buffer lists
func WaitAndUpdateEventsInBuffers(cList, pList []BufferType, ev *cl.Event) {
	// Wait for event to complete
	err := cl.WaitForEvents([]*cl.Event{ev})
	if err != nil {
		log.Println("ERROR: WaitAndUpdateEventsInBuffers failed to wait for event!")
	}
	// Regardless of whether there was error, remove the event from the buffers by
	// signaling goroutine
	EventInChan <- EventAndBuffers{ev: ev, cList: cList, pList: pList}
}

var (
	EventOutChan  chan []*cl.Event
	EventInChan   chan EventAndBuffers
	bufUpdaterCtx context.Context
	bufCtxCancFcn context.CancelFunc
)

func initEventRoutine() {
	EventOutChan = make(chan []*cl.Event)
	EventInChan = make(chan EventAndBuffers)

	bufUpdaterCtx, bufCtxCancFcn = context.WithCancel(context.Background())

	go EventUpdateRoutine(EventInChan, EventOutChan)
}

func killEventRoutine() {
	bufCtxCancFcn()
}

func EventUpdateRoutine(in <-chan EventAndBuffers, out chan<- []*cl.Event) {
	for {
		select {
		case data := <-in:
			if data.ev == nil {
				// This signals that routine should return list of all events
				// for buffers in cList and pList
				outEventList := GetAllEventsOfBuffers(data.cList)
				outEventList = append(outEventList, GetAllEventsOfBuffers(data.pList)...)
				out <- outEventList
			} else {
				// This signals that routine should remove event
				// from buffers in cList and pList
				// Event has successfully completed...
				// Purge event from buffers consumed...
				RemoveReadEventFromBuffers(data.cList, data.ev)
				RemoveEventFromBuffers(data.pList, data.ev)
			}
		case <-bufUpdaterCtx.Done():
			return
		}
	}
}
