package main

import "container/heap"

type QueuedEvent struct {
	Event *Event
	Time  float64
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

func (eq *EventQueue) Queue(ev *Event, time float64) {
	heap.Push(eq, QueuedEvent{
		Event: ev,
		Time:  time,
	})
}

func (eq *EventQueue) Dequeue() QueuedEvent {
	return heap.Pop(eq).(QueuedEvent)
}
