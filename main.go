package main

import (
	"fmt"
)

func clock(proc Process, name string, timeout float64) {
	fmt.Printf("[%4.1f] %s 0\n", proc.Now, name)

	for i := 1; i <= 2; i++ {
		proc.Wait(proc.Timeout(timeout))
		fmt.Printf("[%4.1f] %s %d\n", proc.Now, name, i)
	}

	ev := proc.Timeout(5)
	proc2 := proc.Start(func(proc Process) { dummy(proc, ev) })
	proc.Wait(proc2)
	fmt.Printf("[%4.1f] %s end\n", proc.Now, name)
}

func dummy(proc Process, ev *Event) {
	fmt.Printf("[%4.1f] dummy start\n", proc.Now)
	proc.Wait(ev)
	fmt.Printf("[%4.1f] dummy end\n", proc.Now)
}

func main() {
	sim := Simulation{}

	sim.Start(func(proc Process) { clock(proc, "A", 2) })
	sim.Start(func(proc Process) { clock(proc, "B", 5) })

	sim.Run()
}
