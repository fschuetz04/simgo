package simgo_test

import (
	"testing"

	"github.com/fschuetz04/simgo"
)

func TestAnyOfEmpty(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		anyOf := proc.AnyOf()
		proc.Wait(anyOf)
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAnyOfTriggered(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev1 := proc.Event()
		ev2 := proc.Event()
		ev2.Trigger()
		anyOf := proc.AnyOf(ev1, ev2)
		proc.Wait(anyOf)
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAnyOfPending(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev1 := proc.Event()
		ev2 := proc.Timeout(5)
		anyOf := proc.AnyOf(ev1, ev2)
		proc.Wait(anyOf)
		assertf(t, proc.Now() == 5, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAnyOfProcessed(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev1 := proc.Event()
		ev2 := proc.Timeout(5)
		proc.Wait(ev2)
		assertf(t, proc.Now() == 5, "proc.Now() == %f", proc.Now())
		anyOf := proc.AnyOf(ev1, ev2)
		proc.Wait(anyOf)
		assertf(t, proc.Now() == 5, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAnyOfAllAborted(t *testing.T) {
	sim := simgo.NewSimulation()

	ev1 := sim.Event()
	ev2 := sim.Event()
	ev1.Abort()
	ev2.Abort()

	anyOf := sim.AnyOf(ev1, ev2)
	assertf(t, anyOf.Aborted(), "anyOf should be aborted when all events are aborted")
}

func TestAnyOfSomeAborted(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		ev1 := proc.Event()
		ev2 := proc.Timeout(5)
		ev1.Abort()

		anyOf := proc.AnyOf(ev1, ev2)
		assertf(t, !anyOf.Aborted(), "anyOf should not be aborted when some events are pending")

		proc.Wait(anyOf)
		assertf(t, proc.Now() == 5, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished, "process did not finish")
}

func TestAnyOfAllBecomingAborted(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false
	aborted := false

	sim.Process(func(proc simgo.Process) {
		ev1 := proc.Event()
		ev2 := proc.Event()

		anyOf := proc.AnyOf(ev1, ev2)

		// Add an abort handler to detect if anyOf gets aborted
		anyOf.AddAbortHandler(func(ev *simgo.Event) {
			aborted = true
		})

		// In a separate process, abort both events
		proc.Process(func(proc2 simgo.Process) {
			proc2.Wait(proc2.Timeout(2))
			ev1.Abort()
			ev2.Abort()
		})

		// Wait for the anyOf event
		finished = true
		proc.Wait(anyOf)
		t.Error("process continued after waiting for aborted event")
	})

	sim.Run()
	assertf(t, finished, "process did not start waiting")
	assertf(t, aborted, "anyOf event was not aborted when all events became aborted")
}
