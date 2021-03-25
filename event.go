// Copyright © 2021 Felix Schütz
// Licensed under the MIT license. See the LICENSE file for details.

package simgo

import "fmt"

type state int

const (
	pending state = iota
	triggered
	processed
	aborted
)

type Event struct {
	sim      *Simulation
	state    state
	handlers []Process
}

func (ev *Event) Trigger() bool {
	if !ev.Pending() {
		return false
	}

	ev.state = triggered
	ev.sim.schedule(ev, 0)
	return true
}

func (ev *Event) TriggerDelayed(delay float64) bool {
	if delay < 0 {
		panic(fmt.Sprintf("(*Event).TriggerDelayed: delay must not be negative: %f\n", delay))
	}

	if !ev.Pending() {
		return false
	}

	ev.sim.schedule(ev, delay)
	return true
}

func (ev *Event) Abort() bool {
	if !ev.Pending() {
		return false
	}

	ev.state = aborted
	ev.handlers = nil

	return true
}

func (ev *Event) Pending() bool {
	return ev.state == pending
}

func (ev *Event) Triggered() bool {
	return ev.state == triggered || ev.Processed()
}

func (ev *Event) Processed() bool {
	return ev.state == processed
}

func (ev *Event) Aborted() bool {
	return ev.state == aborted
}

func (ev *Event) process() {
	if ev.Processed() || ev.Aborted() {
		return
	}

	ev.state = processed

	for _, proc := range ev.handlers {
		// yield to process
		proc.sync <- true

		// wait for process
		<-proc.sync
	}

	ev.handlers = nil
}

func (ev *Event) addHandler(proc Process) {
	if ev.Processed() || ev.Aborted() {
		// event will not be processed (again), do not store handler
		return
	}

	ev.handlers = append(ev.handlers, proc)
}
