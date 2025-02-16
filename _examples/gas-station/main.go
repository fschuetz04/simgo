package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fschuetz04/simgo"
)

const (
	Target             = 2000
	StationFuelCap     = 200
	NPumps             = 2
	MinArrivalInterval = 30
	MaxArrivalInterval = 300
	MinCarFuelLvl      = 5
	MaxCarFuelLvl      = 25
	CarFuelCap         = 50
	RefuelingSpeed     = 2
	Threshold          = 0.25
	TankTruckTime      = 300
)

func randUniformInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func car(proc simgo.Process, i int, pumps *simgo.Resource, stationFuel *Container) {
	fmt.Printf("[%6.1f] Car %d arrives\n", proc.Now(), i)
	start := proc.Now()

	proc.Wait(pumps.Request())

	fmt.Printf("[%6.1f] Car %d gets to a pump, station fuel lvl = %f\n", proc.Now(), i, stationFuel.lvl)

	lvl := randUniformInt(MinCarFuelLvl, MaxCarFuelLvl)
	amount := float64(CarFuelCap - lvl)
	proc.Wait(stationFuel.Get(amount))

	fmt.Printf("[%6.1f] Car %d starts refueling\n", proc.Now(), i)

	delay := amount / RefuelingSpeed
	proc.Wait(proc.Timeout(delay))

	fmt.Printf("[%6.1f] Car %d finished refueling and leaves after %.1f\n", proc.Now(), i, proc.Now()-start)

	pumps.Release()
}

func carSource(proc simgo.Process, stationFuel *Container) {
	pumps := simgo.NewResource(proc.Simulation, NPumps)

	for i := 1; ; i++ {
		delay := randUniformInt(MinArrivalInterval, MaxArrivalInterval)
		proc.Wait(proc.Timeout(float64(delay)))

		proc.ProcessReflect(car, i, pumps, stationFuel)
	}
}

func tankTruck(proc simgo.Process, stationFuel *Container) {
	for {
		proc.Wait(proc.Timeout(TankTruckTime))

		if stationFuel.lvl <= stationFuel.cap*Threshold {
			fmt.Printf("[%6.1f] Station fuel level below threshold, tank truck is called\n", proc.Now())

			proc.Wait(proc.Timeout(TankTruckTime))

			fmt.Printf("[%6.1f] Tank truck arrives and refills station\n", proc.Now())

			stationFuel.Put(stationFuel.cap - stationFuel.lvl)
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	sim := simgo.NewSimulation()

	stationFuel := NewFilledCappedContainer(sim, StationFuelCap, StationFuelCap)

	sim.ProcessReflect(carSource, stationFuel)
	sim.ProcessReflect(tankTruck, stationFuel)

	sim.RunUntil(Target)
}
