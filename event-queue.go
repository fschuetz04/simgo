package main

type QueuedEvent struct {
	Time  float64
	Event *Event
}

type EventQueue []QueuedEvent

func (eq EventQueue) Len() int { return len(eq) }

func (eq EventQueue) Less(i, j int) bool { return eq[i].Time < eq[j].Time }

func (eq EventQueue) Swap(i, j int) { eq[i], eq[j] = eq[j], eq[i] }

func (eq *EventQueue) Push(item interface{}) { *eq = append(*eq, item.(QueuedEvent)) }

func (eq *EventQueue) Pop() interface{} {
	n := len(*eq)
	item := (*eq)[n-1]
	*eq = (*eq)[:n-1]
	return item
}
