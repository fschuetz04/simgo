package simgo

import "log"

type state int

const (
	pending state = iota
	triggered
	processed
)

type Event struct {
	sim      *Simulation
	state    state
	handlers []func()
}

func (ev *Event) Trigger() bool {
	if ev.state != pending {
		return false
	}

	ev.state = triggered
	ev.sim.schedule(ev, 0)
	return true
}

func (ev *Event) TriggerDelayed(delay float64) bool {
	if delay < 0 {
		log.Fatalf("(*Event).TriggerDelayed: delay must not be negative: %f\n", delay)
	}

	if ev.state != pending {
		return false
	}

	ev.sim.schedule(ev, delay)
	return true
}

func (ev *Event) process() {
	if ev.state == processed {
		return
	}

	ev.state = processed

	for _, handler := range ev.handlers {
		handler()
	}
}

func (ev *Event) addHandler(handler func()) bool {
	if ev.state != pending {
		return false
	}

	ev.handlers = append(ev.handlers, handler)
	return true
}

func (ev *Event) addHandlerProcess(proc Process) bool {
	return ev.addHandler(func() {
		// yield to process
		proc.sync <- struct{}{}

		// wait for process
		<-proc.sync
	})
}
