package main

import (
	"log"
)

type Proc func()

type Simulation struct {
	Now        float64
	EventQueue EventQueue
}

func NewSimulation() *Simulation {
	return &Simulation{}
}

func (sim *Simulation) Process() *Process {
	proc := &Process{Sync: make(chan struct{})}

	ev := sim.Event()
	ev.AddHandler(proc)
	ev.Trigger()

	<-proc.Sync

	return proc
}

func (sim *Simulation) Event() *Event {
	return &Event{Simulation: sim}
}

func (sim *Simulation) Timeout(delay float64) *Event {
	if delay < 0 {
		log.Fatalf("(*Simulation).Timeout: delay must not be negative: %f\n", delay)
	}

	ev := sim.Event()
	ev.TriggerDelayed(delay)
	return ev
}

func (sim *Simulation) Schedule(ev *Event, delay float64) {
	if delay < 0 {
		log.Fatalf("(*Simulation).Schedule: delay must not be negative: %f\n", delay)
	}

	time := sim.Now + delay
	sim.EventQueue.Queue(ev, time)
}

func (sim *Simulation) Step() bool {
	if len(sim.EventQueue) == 0 {
		return false
	}

	qe := sim.EventQueue.Dequeue()
	sim.Now = qe.Time
	qe.Event.Process()

	return true
}

func (sim *Simulation) Run() {
	for sim.Step() {
	}
}
