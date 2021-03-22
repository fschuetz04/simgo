package main

import "log"

type State int

const (
	Pending State = iota
	Triggered
	Processed
)

type Event struct {
	Simulation *Simulation
	State      State
	Handlers   []Process
}

func (state State) String() string {
	switch state {
	case Pending:
		return "Pending"
	case Triggered:
		return "Triggered"
	case Processed:
		return "Processed"
	default:
		return "Invalid"
	}
}

func (ev *Event) Trigger() bool {
	if ev.State != Pending {
		return false
	}

	ev.State = Triggered
	ev.Simulation.Schedule(ev, 0)
	return true
}

func (ev *Event) TriggerDelayed(delay float64) bool {
	if delay < 0 {
		log.Fatalf("(*Event).TriggerDelayed: delay must not be negative: %f\n", delay)
	}

	if ev.State != Pending {
		return false
	}

	ev.Simulation.Schedule(ev, delay)
	return true
}

func (ev *Event) Process() {
	if ev.State == Processed {
		return
	}

	ev.State = Processed

	for _, process := range ev.Handlers {
		process <- struct{}{}
		<-process
	}
}

func (ev *Event) AddHandler(handler Process) bool {
	if ev.State != Pending {
		return false
	}

	ev.Handlers = append(ev.Handlers, handler)
	return true
}
