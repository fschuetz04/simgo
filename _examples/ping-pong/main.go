package main

import (
	"fmt"

	"github.com/fschuetz04/simgo"
)

type PingPongEvent struct {
	*simgo.Event
	otherEv *PingPongEvent
}

func party(proc simgo.Process, name string, ev *PingPongEvent, delay float64) {
	for {
		proc.Wait(ev)
		fmt.Printf("[%.0f] %s\n", proc.Now(), name)
		theirEv := ev.otherEv
		ev = &PingPongEvent{Event: proc.Event()}
		theirEv.otherEv = ev
		theirEv.TriggerDelayed(delay)
	}
}

func main() {
	sim := simgo.NewSimulation()

	pongEv := &PingPongEvent{Event: sim.Event(), otherEv: nil}
	pingEv := &PingPongEvent{Event: sim.Timeout(0), otherEv: pongEv}

	sim.ProcessReflect(party, "ping", pingEv, 1)
	sim.ProcessReflect(party, "pong", pongEv, 2)

	sim.RunUntil(8)
}
