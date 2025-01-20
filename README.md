# SimGo

SimGo is a discrete event simulation framework for Go.
It is similar to [SimPy](https://simpy.readthedocs.io/en/latest) and aims to be easy
to set up and use.

Processes are defined as simple functions receiving `simgo.Process` as their first
argument.
Each process is executed in a separate goroutine, but it is guaranteed that only
one process is executed at a time.
For examples, look into the `examples` folder.
A short example simulating two clocks ticking in different time intervals looks like
this:

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
    sim := simgo.Simulation{}

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

You can find more examples in the `examples` directory.

## Copyright and License

Copyright © 2024 Felix Schütz.

Licensed under the MIT License.
See the LICENSE file for details.
