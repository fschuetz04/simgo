# SimGo

[![Go Reference](https://pkg.go.dev/badge/github.com/fschuetz04/simgo.svg)](https://pkg.go.dev/github.com/fschuetz04/simgo)

SimGo is a discrete event simulation framework for Go.
It is similar to [SimPy](https://simpy.readthedocs.io/en/latest) and aims to be easy
to set up and use.

Processes are defined as simple functions receiving `simgo.Process` as their first
argument.
Each process is executed in a separate goroutine, but it is guaranteed that only
one process is executed at a time.
For examples, look into the `_examples` folder.

## Basic example: parallel clocks

A short example simulating two clocks ticking in different time intervals:

```go
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
```

When run, the following output is generated:

```text
slow 0
fast 0
fast 1
slow 2
fast 2
fast 3
slow 4
fast 4
```

## Process communication example: ping pong

A more advanced example demonstrating how processes can communicate through events:

```go
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
```

When run, this produces:

```text
[0] ping
[1] pong
[3] ping
[4] pong
[6] ping
[7] pong
```

This example shows two processes communicating through events.
Each process waits for its event to be triggered, prints its name and the current
time, then triggers the event of the other process after a delay.

## License

Licensed under the MIT License.
See the [licence file](LICENSE) for details.
