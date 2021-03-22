package main

import (
	"container/heap"
)

type Proc func()

type Simulation struct {
	Now        float64
	EventQueue EventQueue
}

func NewSimulation() *Simulation {
	return &Simulation{}
}

func (sim *Simulation) Event() *Event {
	return &Event{Simulation: sim}
}

func (sim *Simulation) Timeout(delay float64) *Event {
	ev := sim.Event()
	ev.TriggerDelayed(delay)
	return ev
}

func (sim *Simulation) Schedule(ev *Event, delay float64) {
	qe := QueuedEvent{Time: sim.Now, Event: ev}
	heap.Push(&sim.EventQueue, qe)
}

func (sim *Simulation) Step() bool {
	if len(sim.EventQueue) == 0 {
		return false
	}

	qe := heap.Pop(&sim.EventQueue).(QueuedEvent)
	sim.Now = qe.Time
	qe.Event.Process()

	return true
}

func (sim *Simulation) Run() {
	for sim.Step() {
	}
}
