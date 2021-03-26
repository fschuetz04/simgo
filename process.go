// Copyright © 2021 Felix Schütz
// Licensed under the MIT license. See the LICENSE file for details.

package simgo

import "runtime"

type awaitable interface {
	Pending() bool
	Triggered() bool
	Processed() bool
	Aborted() bool
	addHandler(proc Process)
}

type Process struct {
	*Simulation
	ev   *Event
	sync chan bool
}

func (proc Process) Wait(ev awaitable) {
	if ev.Processed() {
		// event was already processed, do not wait
		return
	}

	if ev.Aborted() {
		// event aborted, exit process
		proc.ev.Abort()
		runtime.Goexit()
	}

	// resume this process when the event is processed
	ev.addHandler(proc)

	// yield to simulation
	proc.sync <- true

	// wait for simulation
	processed := <-proc.sync

	if !processed {
		// event aborted, exit process
		proc.ev.Abort()
		runtime.Goexit()
	}
}

func (proc Process) Pending() bool {
	return proc.ev.Pending()
}

func (proc Process) Triggered() bool {
	return proc.ev.Triggered()
}

func (proc Process) Processed() bool {
	return proc.ev.Processed()
}

func (proc Process) Aborted() bool {
	return proc.ev.Aborted()
}

func (proc Process) addHandler(handlerProc Process) {
	proc.ev.addHandler(handlerProc)
}
