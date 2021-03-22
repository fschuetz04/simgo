package main

type Process chan struct{}

func (proc Process) Wait(ev *Event) {
	if !ev.AddHandler(proc) {
		return
	}

	proc <- struct{}{}
	<-proc
}

func (proc Process) Exit() {
	proc <- struct{}{}
}
