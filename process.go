package main

type Awaitable interface {
	addHandler(proc Process) bool
}

type Process struct {
	*Simulation
	*Event
	sync chan struct{}
}

func (proc Process) Wait(ev Awaitable) {
	if !ev.addHandler(proc) {
		// event was not pending
		return
	}

	// yield to simulation
	proc.sync <- struct{}{}

	// wait for simulation
	<-proc.sync
}
