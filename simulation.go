package simgo

import (
	"log"
	"sync"
)

type Simulation struct {
	Now   float64
	eq    eventQueue
	mutex sync.Mutex
}

type Runner func(proc Process)

func (sim *Simulation) Start(runner Runner) Process {
	proc := Process{
		Simulation: sim,
		Event:      sim.Event(),
		sync:       make(chan struct{}),
	}

	ev := sim.Event()
	ev.addHandlerProcess(proc)

	sim.mutex.Lock()
	ev.Trigger()
	sim.mutex.Unlock()

	go func() {
		// wait for simulation
		<-proc.sync

		runner(proc)

		proc.Trigger()

		// yield to simulation
		proc.sync <- struct{}{}
	}()

	return proc
}

func (sim *Simulation) Event() *Event {
	return &Event{sim: sim}
}

func (sim *Simulation) Timeout(delay float64) *Event {
	if delay < 0 {
		log.Fatalf("(*Simulation).Timeout: delay must not be negative: %f\n", delay)
	}

	ev := sim.Event()
	ev.TriggerDelayed(delay)
	return ev
}

func (sim *Simulation) Step() bool {
	if len(sim.eq) == 0 {
		return false
	}

	qe := sim.eq.dequeue()
	sim.Now = qe.time
	qe.event.process()

	return true
}

func (sim *Simulation) Run() {
	for sim.Step() {
	}
}

func (sim *Simulation) RunUntil(target float64) {
	if target < 0 {
		log.Fatalf("(*Simulation).RunUntil: target must not be negative: %f\n", target)
	}

	for len(sim.eq) > 0 && sim.eq[0].time <= target {
		sim.Step()
	}

	sim.Now = target
}

func (sim *Simulation) schedule(ev *Event, delay float64) {
	time := sim.Now + delay
	sim.eq.queue(ev, time)
}
