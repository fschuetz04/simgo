package main

import (
	"container/heap"
	"fmt"
)

type Proc func()

type Simulation struct {
	Now        float64
	EventQueue EventQueue
	Channel    chan struct{}
	Processes  int
}

func NewSimulation() *Simulation {
	return &Simulation{Channel: make(chan struct{})}
}

func (sim *Simulation) Event() *Event {
	return &Event{Simulation: sim, Channel: make(chan struct{})}
}

func (sim *Simulation) Timeout(delay float64) *Event {
	ev := sim.Event()
	ev.TriggerDelayed(delay)
	return ev
}

func (sim *Simulation) Process() {
	sim.Processes++
	<-sim.Channel
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
	sim.Start()

	for sim.Step() {
	}
}

func (sim *Simulation) Start() {
	for i := 0; i < sim.Processes; i++ {
		fmt.Println("Starting process")
		sim.Channel <- struct{}{}
	}
}
