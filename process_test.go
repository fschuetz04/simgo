package simgo_test

import (
	"testing"

	"github.com/fschuetz04/simgo"
)

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