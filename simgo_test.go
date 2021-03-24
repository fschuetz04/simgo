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

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(sim.AnyOf())
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAnyOfTriggered(t *testing.T) {
	sim := Simulation{}

	sim.Start(func(proc Process) {
		ev1 := sim.Event()
		ev2 := sim.Event()
		ev2.Trigger()
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(sim.AnyOf(ev1, ev2))
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAnyOfPending(t *testing.T) {
	sim := Simulation{}

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := sim.Event()
		ev2 := sim.Timeout(5)
		proc.Wait(sim.AnyOf(ev1, ev2))
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAnyOfProcessed(t *testing.T) {
	sim := Simulation{}

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := sim.Event()
		ev2 := sim.Timeout(5)
		proc.Wait(ev2)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
		proc.Wait(sim.AnyOf(ev1, ev2))
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAllOfEmpty(t *testing.T) {
	sim := Simulation{}

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(sim.AllOf())
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAllOfTriggered(t *testing.T) {
	sim := Simulation{}

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := sim.Event()
		ev1.Trigger()
		ev2 := sim.Event()
		ev2.Trigger()
		proc.Wait(sim.AllOf(ev1, ev2))
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAllOfPending(t *testing.T) {
	sim := Simulation{}

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := sim.Timeout(10)
		ev2 := sim.Timeout(5)
		proc.Wait(sim.AllOf(ev1, ev2))
		assertf(t, proc.Now == 10, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestAllOfProcessed(t *testing.T) {
	sim := Simulation{}

	sim.Start(func(proc Process) {
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		ev1 := sim.Event()
		ev1.Trigger()
		ev2 := sim.Timeout(5)
		proc.Wait(ev1)
		assertf(t, proc.Now == 0, "proc.Now == %f", proc.Now)
		proc.Wait(sim.AllOf(ev1, ev2))
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
	})

	sim.Run()
}

func TestWaitForProc(t *testing.T) {
	sim := Simulation{}

	proc1 := sim.Start(func(proc Process) {
		proc.Wait(proc.Timeout(5))
	})

	sim.Start(func(proc Process) {
		proc.Wait(proc1)
		assertf(t, proc.Now == 5, "proc.Now == %f", proc.Now)
	})

	sim.Run()
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

func TestRunUntil(t *testing.T) {
	sim := Simulation{}

	// TODO

	sim.RunUntil(10)
}
