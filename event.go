package main

type Handler func()

type State int

const (
	Pending State = iota
	Triggered
	Processed
)

type Event struct {
	Simulation *Simulation
	State      State
	Waiting    int
	Channel    chan struct{}
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

func (ev *Event) Wait() {
	if ev.State == Triggered || ev.State == Processed {
		return
	}

	ev.Waiting++
	<-ev.Channel
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

	for i := 0; i < ev.Waiting; i++ {
		ev.Channel <- struct{}{}
	}
}
