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
	sim := simgo.Simulation{}

	sim.ProcessReflect(clock, "fast", 0.5)
	sim.ProcessReflect(clock, "slow", 1)

	sim.RunUntil(2)
}
