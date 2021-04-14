// Copyright © 2021 Felix Schütz
// Licensed under the MIT license. See the LICENSE file for details.

package simgo

// queuedEvent is an event which is scheduled to be processed at a particular
// time.
type queuedEvent struct {
	// ev is the scheduled event.
	ev *Event

	// time is the time at which the event will be processed.
	time float64

	// id is an incremental ID to sort events scheduled at the same time by
	// insertion order.
	id uint64
}

// eventQueue holds all scheduled events for a discrete-event simulation.
type eventQueue []queuedEvent

// Len returns the number of scheduled events.
func (eq eventQueue) Len() int {
	return len(eq)
}

// Less returns whether the event at position i is scheduled before the event
// at position j.
func (eq eventQueue) Less(i, j int) bool {
	if eq[i].time != eq[j].time {
		return eq[i].time < eq[j].time
	}

	return eq[i].id < eq[j].id
}

// Swap swaps the scheduled events at position i and j.
func (eq eventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
}

// Push appends the given scheduled event at the back.
func (eq *eventQueue) Push(item interface{}) {
	*eq = append(*eq, item.(queuedEvent))
}

// Pop removes and returns the scheduled event at the front.
func (eq *eventQueue) Pop() interface{} {
	n := len(*eq)
	item := (*eq)[n-1]
	*eq = (*eq)[:n-1]
	return item
}
