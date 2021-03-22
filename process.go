package main

type Process struct {
	*Simulation
	sync chan struct{}
}

func (proc Process) Wait(ev *Event) {
	if !ev.addHandler(proc) {
		// event was not pending
		return
	}

	// yield to simulation
	proc.sync <- struct{}{}

	// wait for simulation
	<-proc.sync
}
