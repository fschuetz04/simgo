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

func TestTriggerTimeoutEarly(t *testing.T) {
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
