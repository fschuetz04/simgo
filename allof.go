package simgo

// AllOf creates and returns a pending event which is triggered when all of the
// given events are processed.
func (sim *Simulation) AllOf(evs ...Awaitable) *Event {
	n := len(evs)

	// Check if any events are already aborted - if so, abort immediately
	for _, ev := range evs {
		if ev.Aborted() {
			ev := sim.Event()
			ev.Abort()
			return ev
		}
	}

	// check how many events are already processed
	for _, ev := range evs {
		if ev.Processed() {
			n--
		}
	}

	// if no events are given or all events are already processed, the returned
	// event is immediately triggered
	if n == 0 {
		return sim.Timeout(0)
	}

	allOf := sim.Event()

	for _, ev := range evs {
		// when the event is processed, check whether the condition is
		// fulfilled, and trigger the returned event if so
		ev.AddHandler(func(ev *Event) {
			n--
			if n == 0 {
				allOf.Trigger()
			}
		})

		// if the event is aborted, the condition cannot be fulfilled, so abort
		// the returned event
		ev.AddAbortHandler(func(ev *Event) {
			allOf.Abort()
		})
	}

	return allOf
}
