// Copyright © 2024 Felix Schütz
// Licensed under the MIT license. See the LICENSE file for details.

package simgo

import "fmt"

// state holds the state of an event.
type state int

const (
	// The event has not yet been triggered or processed. This is the initial
	// state of a new event.
	pending state = iota

	// The event has been triggered and will be processed at the current
	// simulation time.
	triggered

	// The event has been processed. All normal handlers are currently being
	// called or have been called.
	processed

	// The event has been aborted. All abort handlers are currently being called
	// or have been called.
	aborted
)

// Handler is either a normal handler or an abort handler, which is called
// when an event is processed / aborted. The handler gets a pointer to the
// relevant pointer, so one function can be used to handle multiple events
// and distinguish between them.
type Handler func(ev *Event)

// Event is an event in a discrete-event simulation. The event does not contain
// information about whether it is scheduled to be processed.
//
// To create a new event, use (*Simulation).Timeout, (*Simulation).Event,
// (*Simulation).AnyOf, or (*Simulation).AllOf:
//
//	ev1 := sim.Event()
//	ev2 := sim.Timeout(5)
//	ev3 := sim.AnyOf(ev1, ev2)
//	ev4 := sim.AllOf(ev1, ev2)
type Event struct {
	// sim is used to schedule the event to be processed.
	sim *Simulation

	// state holds the state of the event.
	state state

	// handlers holds all normal handlers of the event. These handlers will be
	// called when the event is processed.
	handlers []Handler

	// handlers holds all abort handlers of the event. These handlers will be
	// called when the event is aborted.
	abortHandlers []Handler
}

// Trigger schedules the event to be processed immediately. This will call all
// normal handlers of the event.
//
// If the event is not pending, it will not be scheduled.
func (ev *Event) Trigger() bool {
	if !ev.Pending() {
		return false
	}

	ev.state = triggered
	ev.sim.schedule(ev, 0)
	return true
}

// TriggerDelayed schedules the event to be processed after the given delay.
// This will call all normal handlers of the event.
//
// If the event is not pending, it will not be scheduled.
//
// Returns true if the event has been scheduled or false otherwise. Panics if
// the given delay is negative.
//
// Note that an event can be triggered delayed multiple times, or triggered
// immediately after it is already triggered delayed. The event will be
// processed once at the earliest scheduled time.
func (ev *Event) TriggerDelayed(delay float64) bool {
	if delay < 0 {
		panic(fmt.Sprintf("(*Event).TriggerDelayed: delay must not be negative: %f", delay))
	}

	if !ev.Pending() {
		return false
	}

	ev.sim.schedule(ev, delay)

	return true
}

// Abort aborts the event and calls all abort handlers of the event.
//
// If the event is not pending, it will not be aborted.
//
// Returns true if the event has been aborted or false otherwise.
//
// TODO: Abort immediately calls it handlers, which leads to a growing call
// stack. Instead, the processing could be scheduled similar to that of
// (*Event).Trigger.
func (ev *Event) Abort() bool {
	if !ev.Pending() {
		return false
	}

	ev.state = aborted

	for _, handler := range ev.abortHandlers {
		handler(ev)
	}

	// handlers will not be required again
	ev.handlers = nil
	ev.abortHandlers = nil

	return true
}

// Pending returns whether the event is pending. A pending event has not yet
// been triggered or processed. This is the initial state of a new event.
func (ev *Event) Pending() bool {
	return ev.state == pending
}

// Triggered returns whether the event has been triggered. A triggered event
// will be processed at the current simulation time if it has not been processed
// already.
func (ev *Event) Triggered() bool {
	return ev.state == triggered || ev.Processed()
}

// Processed returns whether the event has been processed. All normal handlers
// of a processed event are currently being called or have been called.
func (ev *Event) Processed() bool {
	return ev.state == processed
}

// Aborted returns whether the event has been aborted. All abort handlers of an
// aborted event are currently being called or have been called.
func (ev *Event) Aborted() bool {
	return ev.state == aborted
}

// AddHandler adds the given handler as a normal handler to the event. The
// handler will be called when the event is processed.
//
// If the event is already processed or aborted, the handler is not stored,
// since it will never be called.
func (ev *Event) AddHandler(handler Handler) {
	if ev.Processed() || ev.Aborted() {
		// event will not be processed (again), do not store handler
		return
	}

	ev.handlers = append(ev.handlers, handler)
}

// AddAbortHandler adds the given handler as an abort handler to the event. The
// handler will be called when the event is aborted.
//
// If the event is already processed or aborted, the handler is not stored,
// since it will never be called.
func (ev *Event) AddAbortHandler(handler Handler) {
	if ev.Processed() || ev.Aborted() {
		// event will not be aborted (again), do not store handler
		return
	}

	ev.abortHandlers = append(ev.abortHandlers, handler)
}

// process processes the event and calls all normal handlers.
//
// If the event is already processed or aborted, it will not be processed.
//
// Returns true if the event has been processed or false otherwise.
func (ev *Event) process() bool {
	if ev.Processed() || ev.Aborted() {
		return false
	}

	ev.state = processed

	for _, handler := range ev.handlers {
		handler(ev)
	}

	// handlers will not be required again
	ev.handlers = nil
	ev.abortHandlers = nil

	return true
}
