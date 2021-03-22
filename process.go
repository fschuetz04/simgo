package simgo

type awaitable interface {
	addHandler(proc Process) bool
}

type Process struct {
	*Simulation
	*Event
	sync chan struct{}
}

func (proc Process) Wait(ev awaitable) {
	if !ev.addHandler(proc) {
		// event was not pending
		return
	}

	// yield to simulation
	proc.sync <- struct{}{}

	// wait for simulation
	<-proc.sync
}
