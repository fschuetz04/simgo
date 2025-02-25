package simgo_test

import (
	"testing"

	"github.com/fschuetz04/simgo"
)

func assertf(t *testing.T, condition bool, format string, args ...any) {
	t.Helper()
	if !condition {
		t.Errorf(format, args...)
	}
}

func TestWaitForProc(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := 0

	proc1 := sim.Process(func(proc simgo.Process) {
		proc.Wait(proc.Timeout(5))
		finished++
	})

	sim.Process(func(proc simgo.Process) {
		proc.Wait(proc1)
		assertf(t, proc.Now() == 5, "proc.Now() == %f", proc.Now())
		finished++
	})

	sim.Run()
	assertf(t, finished == 2, "finished == %d", finished)
}

func TestWaitForProcessed(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev := proc.Timeout(5)
		proc.Wait(ev)
		assertf(t, proc.Now() == 5, "proc.Now() == %f", proc.Now())
		proc.Wait(ev)
		assertf(t, proc.Now() == 5, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestWaitForAborted(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev := proc.Event()
		ev.Abort()
		finished = true
		proc.Wait(ev)
		t.Error("Process was executed too far")
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestRunUntil(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		proc.Wait(proc.Timeout(4))
		assertf(t, proc.Now() == 4, "proc.Now() == %f", proc.Now())
		finished = true
		proc.Wait(proc.Timeout(1))
		t.Error("Simulation was executed too far")
	})

	sim.RunUntil(5)
	assertf(t, finished == true, "finished == false")
}

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

func TestRunUntilNegative(t *testing.T) {
	defer func() {
		err := recover()
		assertf(t, err != nil, "err == nil")
	}()

	sim := simgo.NewSimulation()

	sim.RunUntil(-5)
}

func TestTimeoutNegative(t *testing.T) {
	defer func() {
		err := recover()
		assertf(t, err != nil, "err == nil")
	}()

	sim := simgo.NewSimulation()

	sim.Timeout(-5)
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

func TestTriggerEarly(t *testing.T) {
	sim := simgo.NewSimulation()
	finished := false

	sim.Process(func(proc simgo.Process) {
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		ev := proc.Timeout(5)
		ev.Trigger()
		proc.Wait(ev)
		assertf(t, proc.Now() == 0, "proc.Now() == %f", proc.Now())
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
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
