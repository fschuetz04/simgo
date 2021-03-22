package main

import (
	"fmt"
)

func process(sim *Simulation, name string, timeout float64) {
	fmt.Printf("[%.0f] %s 0\n", sim.Now, name)
	for i := 1; i <= 2; i++ {
		sim.Timeout(timeout).Wait()
		fmt.Printf("[%.0f] %s %d\n", sim.Now, name, i)
	}
}

func main() {
	sim := NewSimulation()

	go process(sim, "A", 5)
	go process(sim, "B", 10)

	fmt.Println("Running simulation")
	sim.Run()
}
