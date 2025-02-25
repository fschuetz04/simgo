package simgo

// AnyOf creates and returns a pending event which is triggered when any of the
// given events is processed.
func (sim *Simulation) AnyOf(evs ...Awaitable) *Event {
	// if no events are given, the returned event is immediately triggered
	if len(evs) == 0 {
		return sim.Timeout(0)
	}

	// if any event is already processed, the returned event is immediately
	// triggered
	for _, ev := range evs {
		if ev.Processed() {
			return sim.Timeout(0)
		}
	}

	// check how many events are aborted
	n := len(evs)
	for _, ev := range evs {
		if ev.Aborted() {
			n--
		}
	}

	// if all events are aborted, the returned event is aborted
	if n == 0 {
		ev := sim.Event()
		ev.Abort()
		return ev
	}

	anyOf := sim.Event()

	for _, ev := range evs {
		// when the event is processed, the condition is fulfilled, so trigger
		// the returned event
		ev.AddHandler(func(ev *Event) { anyOf.Trigger() })

		// if the event gets aborted, check whether this was the last non-aborted
		// event non-aborted event, and aborted the returned event if so
		ev.AddAbortHandler(func(ev *Event) {
			n--
			if n == 0 {
				anyOf.Abort()
			}
		})
	}

	return anyOf
}
