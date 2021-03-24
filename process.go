package simgo

type awaitable interface {
	addHandlerProcess(proc Process) bool
}

type Process struct {
	*Simulation
	ev   *Event
	sync chan struct{}
}

func (proc Process) Wait(ev awaitable) {
	if !ev.addHandlerProcess(proc) {
		// event was not pending
		return
	}

	// yield to simulation
	proc.sync <- struct{}{}

	// wait for simulation
	<-proc.sync
}

func (proc Process) addHandlerProcess(handlerProc Process) bool {
	return proc.ev.addHandlerProcess(handlerProc)
}
