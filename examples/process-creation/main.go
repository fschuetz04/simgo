package main

import (
	"fmt"
	"runtime"

	"github.com/fschuetz04/simgo"
)

func badProcess(proc simgo.Process) {
	ev := proc.AllOf(proc.Event(), proc.Timeout(10_000))
	proc.Wait(ev)
}

func processSource(proc simgo.Process) {
	for i := 1; i <= 1_000_000; i++ {
		proc.Wait(proc.Timeout(1))

		proc.Process(badProcess)

		if i%100_000 == 0 {
			fmt.Printf("[%7.0f] %7d => %7d\n", proc.Now(), i, runtime.NumGoroutine())
		}
	}
}

func main() {
	sim := simgo.Simulation{}

	sim.Process(processSource)

	sim.Run()
	fmt.Printf("[    End]         => %7d\n", runtime.NumGoroutine())
}
