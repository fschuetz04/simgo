# SimGo

SimGo is a discrete event simulation framework for Go.
It is similar to SimPy and aims to be easy to set up and use.

Processes are defined as simple functions receiving `simgo.Process` as their first argument.
Each process is executed in a separate goroutine, but it is guarantueed that only one process is executed at a time.
For examples, look into the `examples` folder.
A short example simulating two clocks ticking in different time intervals looks like this:

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

    sim.ProcessReflect(clock, "fast", 0.5)
    sim.ProcessReflect(clock, "slow", 1)

    sim.RunUntil(2)
}
```

When run, the following output is generated:

```
fast 0
slow 0
fast 0.5
slow 1
fast 1
fast 1.5
slow 2
fast 2
```

You can find more examples in the `examples` directory.

## Copyright and License

Copyright © 2021 Felix Schütz.

Licensed under the MIT License.
See the LICENSE file for details.
