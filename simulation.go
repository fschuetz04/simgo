package simgo

import (
	"container/heap"
	"fmt"
	"reflect"
	"runtime"
)

// Simulation runs a discrete-event simulation. To create a new simulation, use
// NewSimulation().
type Simulation struct {
	// now holds the current simulation time.
	now float64

	// eq holds the event queue.
	eq eventQueue

	// nextID holds the next ID for scheduling a new event.
	nextID uint64

	// shutdown is used to shutdown all process goroutines of this simulation.
	shutdown chan struct{}
}

// NewSimulation creates a new simulation.
func NewSimulation() *Simulation {
	return &Simulation{shutdown: make(chan struct{})}
}

// Now returns the current simulation time.
func (sim *Simulation) Now() float64 {
	return sim.now
}

// Process starts a new process with the given runner.
//
// Creates and triggers an event. As soon as this event is processed, the
// runner is executed. Whenever the runner waits for a pending event, it is
// paused until the event is processed.
//
// It is ensured that only one process is executed at the same time.
//
// Returns the process. This can be used to wait for the process to finish. As
// soon as the process finishes, the underlying event is triggered.
func (sim *Simulation) Process(runner func(proc Process)) Process {
	proc := Process{
		Simulation: sim,
		ev:         sim.Event(),
		sync:       make(chan bool),
	}

	// schedule an event to be processed immediately and add an handler which
	// is called when the event is processed
	ev := sim.Timeout(0)
	ev.AddHandler(func(*Event) {
		// yield to the process
		proc.sync <- true

		// wait for the process
		<-proc.sync
	})

	go func() {
		// yield to the simulation at the end by closing
		defer close(proc.sync)

		// wait for the simulation
		<-proc.sync

		// execute the runner
		runner(proc)

		// process is finished trigger the underlying event
		proc.ev.Trigger()
	}()

	return proc
}

// ProcessReflect starts a new process with the given runner and the given
// additional argument. This uses reflection.
//
// See (*Simulation).Process for further documentation.
func (sim *Simulation) ProcessReflect(runner any, args ...any) Process {
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

// Event creates and returns a pending event.
func (sim *Simulation) Event() *Event {
	ev := &Event{sim: sim}
	runtime.SetFinalizer(ev, func(ev *Event) {
		ev.Abort()
	})
	return ev
}

// Timeout creates and returns a pending event which is processed after the
// given delay. Panics if the given delay is negative.
func (sim *Simulation) Timeout(delay float64) *Event {
	if delay < 0 {
		panic(fmt.Sprintf("(*Simulation).Timeout: delay must not be negative: %f", delay))
	}

	ev := sim.Event()
	ev.TriggerDelayed(delay)
	return ev
}

// Step sets the current simulation time to the scheduled time of the next event
// in the event queue and processes the next event. Returns false if the event
// queue was empty and no event was processed, true otherwise.
func (sim *Simulation) Step() bool {
	if len(sim.eq) == 0 {
		return false
	}

	qe := heap.Pop(&sim.eq).(queuedEvent)
	sim.now = qe.time
	qe.ev.process()

	return true
}

// Run runs the simulation until the event queue is empty.
func (sim *Simulation) Run() {
	for sim.Step() {
	}
}

// RunUntil runs the simulation until the event queue is empty or the next event
// in the event queue is scheduled at or after the given target time. Sets the
// current simulation time to the target time at the end. Panics if the given
// target time is smaller than the current simulation time.
func (sim *Simulation) RunUntil(target float64) {
	if target < sim.Now() {
		panic(fmt.Sprintf("(*Simulation).RunUntil: target must not be smaller than the current simulation time: %f < %f", target, sim.Now()))
	}

	for len(sim.eq) > 0 && sim.eq[0].time < target {
		sim.Step()
	}

	sim.now = target
}

// Shutdown shuts down all process goroutines of this simulation.
func (sim *Simulation) Shutdown() {
	close(sim.shutdown)
}

// schedule schedules the given event to be processed after the given delay.
// Adds the event to the event queue.
func (sim *Simulation) schedule(ev *Event, delay float64) {
	heap.Push(&sim.eq, queuedEvent{
		ev:   ev,
		time: sim.Now() + delay,
		id:   sim.nextID,
	})
	sim.nextID++
}
