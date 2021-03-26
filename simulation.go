// Copyright © 2021 Felix Schütz
// Licensed under the MIT license. See the LICENSE file for details.

package simgo

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

type Simulation struct {
	now   float64
	eq    eventQueue
	mutex sync.Mutex
}

type Runner func(proc Process)

func (sim *Simulation) Now() float64 {
	return sim.now
}

func (sim *Simulation) Process(runner Runner) Process {
	proc := Process{
		Simulation: sim,
		ev:         sim.Event(),
		sync:       make(chan bool),
	}

	// yield to the process when the event is triggered
	ev := sim.Event()
	ev.addProcess(proc)

	sim.mutex.Lock()
	ev.Trigger()
	sim.mutex.Unlock()

	go func() {
		// yield to the simulation at the end by closing
		defer close(proc.sync)

		// wait for simulation
		<-proc.sync

		runner(proc)

		proc.ev.Trigger()
	}()

	return proc
}

func (sim *Simulation) ProcessReflect(runner interface{}, args ...interface{}) Process {
	return sim.Process(func(proc Process) {
		reflectF := reflect.ValueOf(runner)
		reflectArgs := make([]reflect.Value, len(args)+1)
		reflectArgs[0] = reflect.ValueOf(proc)
		for i, arg := range args {
			expected := reflectF.Type().In(i + 1)
			reflectArgs[i+1] = reflect.ValueOf(arg).Convert(expected)
		}
		reflectF.Call(reflectArgs)
	})
}

func (sim *Simulation) Event() *Event {
	ev := &Event{sim: sim}
	runtime.SetFinalizer(ev, func(ev *Event) {
		ev.Abort()
	})
	return ev
}

func (sim *Simulation) Timeout(delay float64) *Event {
	if delay < 0 {
		panic(fmt.Sprintf("(*Simulation).Timeout: delay must not be negative: %f\n", delay))
	}

	ev := sim.Event()
	ev.TriggerDelayed(delay)
	return ev
}

func (sim *Simulation) AnyOf(evs ...*Event) *Event {
	if len(evs) == 0 {
		return sim.Timeout(0)
	}

	for _, ev := range evs {
		if ev.Processed() {
			return sim.Timeout(0)
		}
	}

	anyOf := sim.Event()

	for _, ev := range evs {
		ev.AddHandler(func(ev *Event) { anyOf.Trigger() })
	}

	return anyOf
}

func (sim *Simulation) AllOf(evs ...*Event) *Event {
	n := len(evs)

	for _, ev := range evs {
		if ev.Processed() {
			n--
		}
	}

	if n == 0 {
		return sim.Timeout(0)
	}

	allOf := sim.Event()

	for _, ev := range evs {
		ev.AddHandler(func(ev *Event) {
			n--
			if n == 0 {
				allOf.Trigger()
			}
		})

		ev.AddAbortHandler(func(ev *Event) {
			allOf.Abort()
		})
	}

	return allOf
}

func (sim *Simulation) Step() bool {
	if len(sim.eq) == 0 {
		return false
	}

	qe := sim.eq.dequeue()
	sim.now = qe.time
	qe.event.process()

	return true
}

func (sim *Simulation) Run() {
	for sim.Step() {
	}
}

func (sim *Simulation) RunUntil(target float64) {
	if target < 0 {
		panic(fmt.Sprintf("(*Simulation).RunUntil: target must not be negative: %f\n", target))
	}

	for len(sim.eq) > 0 && sim.eq[0].time <= target {
		sim.Step()
	}

	sim.now = target
}

func (sim *Simulation) schedule(ev *Event, delay float64) {
	time := sim.Now() + delay
	sim.eq.queue(ev, time)
}
