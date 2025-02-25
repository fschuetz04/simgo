package simgo_test

import (
	"testing"

	"github.com/fschuetz04/simgo"
)

func TestTriggerDelayedNegative(t *testing.T) {
	defer func() {
		err := recover()
		assertf(t, err != nil, "err == nil")
	}()

	sim := simgo.NewSimulation()

	ev := sim.Event()
	ev.TriggerDelayed(-5)
}

func TestTriggerDelayedTriggered(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev := proc.Timeout(5)
		proc.Wait(ev)
		assertf(t, proc.Now() == 5, "proc.Now() == %f", proc.Now())
		assertf(t, ev.TriggerDelayed(5) == false, "ev.TriggerDelayed(5) == true")
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestTrigger(t *testing.T) {
	sim := simgo.NewSimulation()

	ev := sim.Event()
	ev.Trigger()
	assertf(t, ev.Triggered() == true, "ev.Pending() == false")
}

func TestTriggerTriggered(t *testing.T) {
	sim := simgo.NewSimulation()

	ev := sim.Event()
	assertf(t, ev.Trigger() == true, "ev.Trigger() == false")
	assertf(t, ev.Trigger() == false, "ev.Trigger() == true")
}

func TestAbort(t *testing.T) {
	sim := simgo.NewSimulation()

	ev := sim.Event()
	ev.Abort()
	assertf(t, ev.Aborted() == true, "ev.Aborted() == false")
}

func TestAbortTriggered(t *testing.T) {
	sim := simgo.NewSimulation()

	ev := sim.Event()
	assertf(t, ev.Trigger() == true, "ev.Trigger() == false")
	assertf(t, ev.Triggered() == true, "ev.Triggered() == false")
	assertf(t, ev.Abort() == false, "ev.Abort() == true")
	assertf(t, ev.Aborted() == false, "ev.Aborted() == true")
}
