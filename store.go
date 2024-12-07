package simgo

import (
	"math"
)

// Store is a resource for storing objects. The objects are put and retrieved
// from the store in a first-in first-out order.
type Store[T any] struct {
	// sim is the reference to the simulation.
	sim *Simulation

	// gets holds the list of pending get events.
	gets []*GetEvent[T]

	// puts holds the list of pending put events.
	puts []*PutEvent[T]

	// items holds the items currently in the store.
	items []T

	// capacity is the maximum number of items in the store.
	capacity int
}

// GetEvent is the event returned from (*Store).Get.
type GetEvent[T any] struct {
	// Event is the underlying event.
	*Event

	// Item holds the item retrieved from the store after the underlying event is
	// triggered.
	Item T
}

// PutEvent is the returned from (*Store).Put.
type PutEvent[T any] struct {
	// Event is the underlying event.
	*Event

	// item holds the item to be returned to the store.
	item T
}

// NewStore creates a store for the given simulation with an unlimited capacity.
func NewStore[T any](sim *Simulation) *Store[T] {
	return NewStoreWithCapacity[T](sim, math.MaxInt)
}

// NewStoreWithCapacity crates a store for the given simulation with the given
// capacity.
func NewStoreWithCapacity[T any](sim *Simulation, capacity int) *Store[T] {
	if capacity <= 0 {
		panic("NewStoreWithCapacity: capacity must be > 0")
	}

	return &Store[T]{sim: sim, capacity: capacity}
}

// Capacity returns the capacity of the store.
func (store *Store[T]) Capacity() int {
	return store.capacity
}

// Available returns the number of items currently in the store.
func (store *Store[T]) Available() int {
	return len(store.items)
}

// Get returns an event that is triggered when an item is retrieved from the
// store, which may be immediately.
func (store *Store[T]) Get() *GetEvent[T] {
	ev := &GetEvent[T]{Event: store.sim.Event()}
	ev.AddHandler(func(*Event) {
		// the store has one less item, so check whether any pending puts can be
		// triggered.
		store.triggerPuts()
	})

	store.gets = append(store.gets, ev)
	store.triggerGets()

	return ev
}

// Put returns an event that is triggered when the given item is returned to the
// store, which may be immediately.
func (store *Store[T]) Put(item T) *PutEvent[T] {
	ev := &PutEvent[T]{Event: store.sim.Event(), item: item}
	ev.AddHandler(func(*Event) {
		// the store has one more item, so check whether any pending gets can be
		// triggered
		store.triggerGets()
	})

	store.puts = append(store.puts, ev)
	store.triggerPuts()

	return ev
}

// triggerGets triggers pending get events until the store is empty.
// left in the store.
func (store *Store[T]) triggerGets() {
	for len(store.gets) > 0 && len(store.items) > 0 {
		get := store.gets[0]
		store.gets = store.gets[1:]

		if !get.Trigger() {
			continue
		}

		item := store.items[0]
		store.items = store.items[1:]

		get.Item = item
	}
}

// triggerPuts triggers pending put events until the store is full.
func (store *Store[T]) triggerPuts() {
	for len(store.puts) > 0 && len(store.items) < store.Capacity() {
		put := store.puts[0]
		store.puts = store.puts[1:]

		if !put.Trigger() {
			continue
		}

		store.items = append(store.items, put.item)
	}
}
