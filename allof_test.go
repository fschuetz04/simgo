package simgo_test

import (
	"testing"

	"github.com/fschuetz04/simgo"
)

func TestAllOfEmpty(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		allOf := proc.AllOf()
		proc.Wait(allOf)
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAllOfTriggered(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev1 := proc.Event()
		ev1.Trigger()
		ev2 := proc.Event()
		ev2.Trigger()
		allOf := proc.AllOf(ev1, ev2)
		proc.Wait(allOf)
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAllOfPending(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev1 := proc.Timeout(10)
		ev2 := proc.Timeout(5)
		allOf := proc.AllOf(ev1, ev2)
		proc.Wait(allOf)
		assertf(t, proc.Now() == 10, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAllOfProcessed(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev1 := proc.Event()
		ev1.Trigger()
		ev2 := proc.Timeout(5)
		proc.Wait(ev1)
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		allOf := proc.AllOf(ev1, ev2)
		proc.Wait(allOf)
		assertf(t, proc.Now() == 5, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAllOfInitiallyAborted(t *testing.T) {
	sim := simgo.NewSimulation()

	ev1 := sim.Event()
	ev2 := sim.Event()
	ev1.Abort()

	allOf := sim.AllOf(ev1, ev2)
	assertf(t, allOf.Aborted(), "allOf should be aborted when any event is aborted")
}

func TestAllOfBecomingAborted(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false
	aborted := false

	ev1 := sim.Event()
	ev2 := sim.Event()

	sim.Process(func(proc simgo.Process) {
		allOf := proc.AllOf(ev1, ev2)

		allOf.AddAbortHandler(func(ev *simgo.Event) {
			aborted = true
		})

		finished = true
		proc.Wait(allOf)
		t.Error("Process continued after waiting for aborted event")
	})

	sim.Process(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(2))
		ev1.Abort()
	})

	sim.Run()
	assertf(t, finished, "Process did not start waiting")
	assertf(t, aborted, "allOf was not aborted when one event became aborted")
}
