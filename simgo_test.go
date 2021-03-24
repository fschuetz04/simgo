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

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(proc.AnyOf())
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == %t", finished)
}

func TestAnyOfTriggered(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Start(func(proc Process) {
		ev1 := proc.Event()
		ev2 := proc.Event()
		ev2.Trigger()
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(proc.AnyOf(ev1, ev2))
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == %t", finished)
}

func TestAnyOfPending(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := proc.Event()
		ev2 := proc.Timeout(5)
		proc.Wait(proc.AnyOf(ev1, ev2))
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == %t", finished)
}

func TestAnyOfProcessed(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Start(func(proc Process) {
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
	assertf(t, finished == true, "finished == %t", finished)
}

func TestAllOfEmpty(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(proc.AllOf())
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == %t", finished)
}

func TestAllOfTriggered(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Start(func(proc Process) {
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
	assertf(t, finished == true, "finished == %t", finished)
}

func TestAllOfPending(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := proc.Timeout(10)
		ev2 := proc.Timeout(5)
		proc.Wait(proc.AllOf(ev1, ev2))
		assertf(t, proc.Now == 10, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == %t", finished)
}

func TestAllOfProcessed(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Start(func(proc Process) {
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
	assertf(t, finished == true, "finished == %t", finished)
}

func TestWaitForProc(t *testing.T) {
	sim := Simulation{}
	finished := 0

	proc1 := sim.Start(func(proc Process) {
		proc.Wait(proc.Timeout(5))
		finished++
	})

	sim.Start(func(proc Process) {
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

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev := proc.Timeout(5)
		proc.Wait(ev)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		proc.Wait(ev)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == %t", finished)
}

func TestRunUntil(t *testing.T) {
	sim := Simulation{}
	finished := false

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(proc.Timeout(5))
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		finished = true
		proc.Wait(proc.Timeout(5))
		t.Error("Simulation was executed too far")
	})

	sim.RunUntil(5)
	assertf(t, finished == true, "finished == %t", finished)
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
