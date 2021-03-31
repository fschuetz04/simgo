// Copyright © 2021 Felix Schütz
// Licensed under the MIT license. See the LICENSE file for details.

package simgo

import "runtime"

// Process is a process in a discrete-event simulation.
type Process struct {
	// Simulation is used to generate timeouts and other events, and start new
	// processes.
	*Simulation

	// ev is triggered when the process finishes or aborted when the process is
	// aborted.
	ev *Event

	// sync is used to yield to the process / simulation and wait for the
	// process / simulation.
	sync chan bool
}

// Wait yields from the process to the simulation and waits until the given
// awaitable is processed.
//
// If the awaitable is already processed, the process is not paused. If the
// awaitable is aborted, the process is aborted too.
func (proc Process) Wait(ev Awaitable) {
	if ev.Processed() {
		// event was already processed, do not wait
		return
	}

	if ev.Aborted() {
		// event aborted, abort process
		proc.ev.Abort()
		runtime.Goexit()
	}

	// handler called when the event is processed
	ev.AddHandler(func(*Event) {
		// yield to process
		proc.sync <- true

		// wait for process
		<-proc.sync
	})

	// handler called when the event is aborted
	ev.AddAbortHandler(func(*Event) {
		// abort process
		proc.sync <- false

		// wait for process
		<-proc.sync
	})

	// yield to simulation
	proc.sync <- true

	// wait for simulation
	processed := <-proc.sync

	if !processed {
		// event aborted, abort process
		proc.ev.Abort()
		runtime.Goexit()
	}
}

// Pending returns whether the underlying event is pending.
func (proc Process) Pending() bool {
	return proc.ev.Pending()
}

// Triggered retrusn whether the underlying event is triggered.
func (proc Process) Triggered() bool {
	return proc.ev.Triggered()
}

// Pending returns whether the underlying event is processed.
func (proc Process) Processed() bool {
	return proc.ev.Processed()
}

// Pending returns whether the underlying event is aborted.
func (proc Process) Aborted() bool {
	return proc.ev.Aborted()
}

// AddHandler adds the given handler to the underlying event.
func (proc Process) AddHandler(handler Handler) {
	proc.ev.AddHandler(handler)
}

// AddAbortHandler adds the given abort handler to the underlying event.
func (proc Process) AddAbortHandler(handler Handler) {
	proc.ev.AddAbortHandler(handler)
}
