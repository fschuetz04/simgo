package main

import (
	"fmt"
	"math/rand"

	"github.com/fschuetz04/simgo"
)

const (
	NMachines  = 10
	RepairTime = 30
	NRepairMen = 1
	NWeeks     = 4

	// normal distribution
	TimeForPartMean   = 10
	TimeForPartStdDev = 2

	// exponential distribution
	TimeToFailureMean = 300
)

func machineProduction(proc simgo.Process, nPartsMade *int, repairMen *simgo.Resource, failure **simgo.Event) {
	for {
		timeForPart := rand.NormFloat64()*TimeForPartStdDev + TimeForPartMean

		for {
			start := proc.Now()
			timeout := proc.Timeout(timeForPart)
			proc.Wait(proc.AnyOf(timeout, *failure))

			if timeout.Triggered() {
				// part is finished
				*nPartsMade++
				break
			}

			// machine failred, calculate remaining time for part and wait for repair
			timeForPart -= proc.Now() - start
			proc.Wait(repairMen.Request())
			proc.Wait(proc.Timeout(RepairTime))
			repairMen.Release()
		}
	}
}

func machineFailure(proc simgo.Process, failure **simgo.Event) {
	for {
		proc.Wait(proc.Timeout(rand.ExpFloat64() * TimeToFailureMean))
		(*failure).Trigger()
		*failure = proc.Event()
	}
}

func machine(proc simgo.Process, nPartsMade *int, repairMen *simgo.Resource) {
	failure := proc.Event()
	proc.ProcessReflect(machineProduction, nPartsMade, repairMen, &failure)
	proc.ProcessReflect(machineFailure, &failure)
}

func main() {
	sim := simgo.NewSimulation()

	repairMen := simgo.NewResource(sim, NRepairMen)
	nPartsMade := make([]int, NMachines)

	for i := 0; i < NMachines; i++ {
		sim.ProcessReflect(machine, &nPartsMade[i], repairMen)
	}

	sim.RunUntil(NWeeks * 7 * 24 * 60)

	fmt.Printf("Machine shop results after %d weeks:\n", NWeeks)
	total := 0
	for i := 0; i < NMachines; i++ {
		fmt.Printf("- Machine %d made %d parts\n", i, nPartsMade[i])
		total += nPartsMade[i]
	}
	fmt.Printf("Total: %d parts\n", total)
}
