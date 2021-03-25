// Copyright © 2021 Felix Schütz
// Licensed under the MIT license. See the LICENSE file for details.

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
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(proc.AnyOf())
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAnyOfTriggered(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		ev1 := proc.Event()
		ev2 := proc.Event()
		ev2.Trigger()
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(proc.AnyOf(ev1, ev2))
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAnyOfPending(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := proc.Event()
		ev2 := proc.Timeout(5)
		proc.Wait(proc.AnyOf(ev1, ev2))
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAnyOfProcessed(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := proc.Event()
		ev2 := proc.Timeout(5)
		proc.Wait(ev2)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		proc.Wait(proc.AnyOf(ev1, ev2))
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAllOfEmpty(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(proc.AllOf())
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAllOfTriggered(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := proc.Event()
		ev1.Trigger()
		ev2 := proc.Event()
		ev2.Trigger()
		proc.Wait(proc.AllOf(ev1, ev2))
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAllOfPending(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := proc.Timeout(10)
		ev2 := proc.Timeout(5)
		proc.Wait(proc.AllOf(ev1, ev2))
		assertf(t, proc.Now == 10, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAllOfProcessed(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := proc.Event()
		ev1.Trigger()
		ev2 := proc.Timeout(5)
		proc.Wait(ev1)
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(proc.AllOf(ev1, ev2))
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestWaitForProc(t *testing.T) {
	sim := Simulation{}
	finished := 0

	proc1 := sim.Process(func(proc Process) {
		proc.Wait(proc.Timeout(5))
		finished++
	})

	sim.Process(func(proc Process) {
		proc.Wait(proc1)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		finished++
	})

	sim.Run()
	assertf(t, finished == 2, "finished == %d", finished)
}

func TestWaitForProcessed(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev := proc.Timeout(5)
		proc.Wait(ev)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		proc.Wait(ev)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestWaitForAborted(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
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
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(proc.Timeout(5))
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		finished = true
		proc.Wait(proc.Timeout(5))
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

	sim := Simulation{}

	ev := sim.Event()
	ev.TriggerDelayed(-5)
}

func TestTriggerDelayedTriggered(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev := proc.Timeout(5)
		proc.Wait(ev)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
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

	sim := Simulation{}

	sim.RunUntil(-5)
}

func TestTimeoutNegative(t *testing.T) {
	defer func() {
		err := recover()
		assertf(t, err != nil, "err == nil")
	}()

	sim := Simulation{}

	sim.Timeout(-5)
}

func TestTrigger(t *testing.T) {
	sim := Simulation{}

	ev := sim.Event()
	ev.Trigger()
	assertf(t, ev.Triggered() == true, "ev.Pending() == false")
}

func TestTriggerTriggered(t *testing.T) {
	sim := Simulation{}

	ev := sim.Event()
	assertf(t, ev.Trigger() == true, "ev.Trigger() == false")
	assertf(t, ev.Trigger() == false, "ev.Trigger() == true")
}

func TestTriggerEarly(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Process(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev := proc.Timeout(5)
		ev.Trigger()
		proc.Wait(ev)
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestAbort(t *testing.T) {
	sim := Simulation{}

	ev := sim.Event()
	ev.Abort()
	assertf(t, ev.Aborted() == true, "ev.Aborted() == false")
}

func TestAbortTriggered(t *testing.T) {
	sim := Simulation{}

	ev := sim.Event()
	assertf(t, ev.Trigger() == true, "ev.Trigger() == false")
	assertf(t, ev.Triggered() == true, "ev.Triggered() == false")
	assertf(t, ev.Abort() == false, "ev.Abort() == true")
	assertf(t, ev.Aborted() == false, "ev.Aborted() == true")
}
