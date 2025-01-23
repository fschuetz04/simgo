package main

import (
	"fmt"

	"github.com/fschuetz04/simgo"
)

func clock(proc simgo.Process, name string, delay float64) {
	for {
		fmt.Println(name, proc.Now())
		proc.Wait(proc.Timeout(delay))
	}
}

func main() {
	sim := simgo.NewSimulation()

	sim.ProcessReflect(clock, "slow", 2)
	sim.ProcessReflect(clock, "fast", 1)

	sim.RunUntil(5)
}
