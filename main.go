package main

import (
	"fmt"
	"time"
)

func process(sim *Simulation, name string, timeout float64) {
	fmt.Printf("[%s] [%f] Initializing process\n", name, sim.Now)

	proc := sim.Process()

	fmt.Printf("[%s] [%f] 0\n", name, sim.Now)

	for i := 1; i <= 2; i++ {
		sim.Timeout(timeout)
		fmt.Printf("[%s] [%f] %d\n", name, sim.Now, i)
	}

	proc.Exit()
	fmt.Printf("[%s] [%f] End\n", name, sim.Now)
}

func main() {
	sim := NewSimulation()

	go process(sim, "A", 5)
	go process(sim, "B", 10)

	time.Sleep(time.Second)

	fmt.Println("[M] Running simulation")
	sim.Run()

	time.Sleep(time.Second)

	fmt.Println("[M] End")
}
