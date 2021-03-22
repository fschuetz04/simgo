package main

import (
	"fmt"
)

func process(proc Process, name string, timeout float64) {
	fmt.Printf("[%4.1f] %s 0\n", proc.Now, name)

	for i := 1; i <= 2; i++ {
		proc.Wait(proc.Timeout(timeout))
		fmt.Printf("[%4.1f] %s %d\n", proc.Now, name, i)
	}
}

func main() {
	sim := &Simulation{}

	sim.Start(func(proc Process) { process(proc, "A", 2) })
	sim.Start(func(proc Process) { process(proc, "B", 5) })

	sim.Run()
}
