package simgo

import (
	"testing"
)

func assertf(t *testing.T, condition bool, format string, args ...interface{}) {
	t.Helper()
	if !condition {
		t.Errorf(format, args...)
	}
}

func TestAnyOfEmpty(t *testing.T) {
	sim := Simulation{}

	anyOf := sim.AnyOf()

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(anyOf)
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAnyOfTriggered(t *testing.T) {
	sim := Simulation{}

	ev1 := sim.Event()
	ev2 := sim.Event()
	ev2.Trigger()
	anyOf := sim.AnyOf(ev1, ev2)

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(anyOf)
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAnyOfPending(t *testing.T) {
	sim := Simulation{}

	ev1 := sim.Event()
	ev2 := sim.Timeout(5)
	anyOf := sim.AnyOf(ev1, ev2)

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(anyOf)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAllOfEmpty(t *testing.T) {
	sim := Simulation{}

	allOf := sim.AllOf()

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(allOf)
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAllOfTriggered(t *testing.T) {
	sim := Simulation{}

	ev1 := sim.Event()
	ev1.Trigger()
	ev2 := sim.Event()
	ev2.Trigger()
	allOf := sim.AllOf(ev1, ev2)

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(allOf)
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAllOfPending(t *testing.T) {
	sim := Simulation{}

	ev1 := sim.Timeout(10)
	ev2 := sim.Timeout(5)
	allOf := sim.AllOf(ev1, ev2)

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(allOf)
		assertf(t, proc.Now == 10, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}
