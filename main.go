package main

import (
	"fmt"
)

func process(sim *Simulation, name string, timeout float64) {
	proc := sim.Process()

	fmt.Printf("[%4.1f] %s 0\n", sim.Now, name)

	for i := 1; i <= 2; i++ {
		proc.Wait(sim.Timeout(timeout))
		fmt.Printf("[%4.1f] %s %d\n", sim.Now, name, i)
	}

	proc.Exit()
}

func main() {
	sim := &Simulation{}

	sim.Add(2)
	go process(sim, "A", 2)
	go process(sim, "B", 10)

	sim.Run()
}
