package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fschuetz04/simgo"
)

const (
	NCounters           = 1
	NCustomers          = 10
	MeanArrivalInterval = 10
	MaxWaitTime         = 16
	MeanTimeInBank      = 12
)

func customer(proc simgo.Process, id int, counters *simgo.Resource) {
	fmt.Printf("[%5.1f] Customer %d arrives\n", proc.Now(), id)

	request := counters.Request()
	timeout := proc.Timeout(MaxWaitTime)
	proc.Wait(proc.AnyOf(request, timeout))

	if !request.Triggered() {
		request.Abort()
		fmt.Printf("[%5.1f] Customer %d leaves unhappy\n", proc.Now(), id)
		return
	}

	fmt.Printf("[%5.1f] Customer %d is served\n", proc.Now(), id)

	delay := rand.ExpFloat64() * MeanTimeInBank
	proc.Wait(proc.Timeout(delay))

	fmt.Printf("[%5.1f] Customer %d leaves\n", proc.Now(), id)
	counters.Release()
}

func customerSource(proc simgo.Process) {
	counters := simgo.NewResource(proc.Simulation, NCounters)

	for id := 1; id <= NCustomers; id++ {
		proc.ProcessReflect(customer, id, counters)
		delay := rand.ExpFloat64() * MeanArrivalInterval
		proc.Wait(proc.Timeout(delay))
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	sim := simgo.NewSimulation()

	sim.Process(customerSource)

	sim.Run()
}
