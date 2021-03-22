package simgo

import "container/heap"

type queuedEvent struct {
	event *Event
	time  float64
}

type eventQueue []queuedEvent

func (eq eventQueue) Len() int { return len(eq) }

func (eq eventQueue) Less(i, j int) bool { return eq[i].time < eq[j].time }

func (eq eventQueue) Swap(i, j int) { eq[i], eq[j] = eq[j], eq[i] }

func (eq *eventQueue) Push(item interface{}) { *eq = append(*eq, item.(queuedEvent)) }

func (eq *eventQueue) Pop() interface{} {
	n := len(*eq)
	item := (*eq)[n-1]
	*eq = (*eq)[:n-1]
	return item
}

func (eq *eventQueue) queue(ev *Event, time float64) {
	heap.Push(eq, queuedEvent{
		event: ev,
		time:  time,
	})
}

func (eq *eventQueue) dequeue() queuedEvent {
	return heap.Pop(eq).(queuedEvent)
}
