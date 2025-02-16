package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fschuetz04/simgo"
)

const (
	NInitialCars    = 4
	WashTime        = 5
	NMachines       = 2
	MeanArrivalTime = 5
)

func wash(proc simgo.Process, id int) {
	proc.Wait(proc.Timeout(WashTime))
	fmt.Printf("[%5.1f] Car %d washed\n", proc.Now(), id)
}

func car(proc simgo.Process, id int, machines *simgo.Resource) {
	fmt.Printf("[%5.1f] Car %d arrives\n", proc.Now(), id)

	proc.Wait(machines.Request())

	fmt.Printf("[%5.1f] Car %d enters\n", proc.Now(), id)

	proc.Wait(proc.ProcessReflect(wash, id))

	fmt.Printf("[%5.1f] Car %d leaves\n", proc.Now(), id)
	machines.Release()
}

func carSource(proc simgo.Process, machines *simgo.Resource) {
	for id := 1; ; id++ {
		if id > NInitialCars {
			proc.Wait(proc.Timeout(rand.ExpFloat64() * MeanArrivalTime))
		}

		proc.ProcessReflect(car, id, machines)
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	sim := simgo.NewSimulation()
	machines := simgo.NewResource(sim, NMachines)

	sim.ProcessReflect(carSource, machines)

	sim.RunUntil(20)
}
