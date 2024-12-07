// Copyright © 2024 Felix Schütz
// Licensed under the MIT license. See the LICENSE file for details.

package simgo

import "testing"

func TestStoreCapacityZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewStoreWithCapacity did not panic with a capacity of 0")
		}
	}()

	sim := &Simulation{}
	NewStoreWithCapacity[int](sim, 0)
}

func TestStoreCapacityPut(t *testing.T) {
	sim := &Simulation{}
	capacity := 2
	store := NewStoreWithCapacity[int](sim, capacity)

	sim.Process(func(proc Process) {
		for i := 0; i <= capacity; i++ {
			store.Put(i)
		}
	})

	sim.Run()
	assertf(t, store.Capacity() == capacity, "store.Capacity() == %d", store.Capacity())
	assertf(t, store.Available() == capacity, "store.Available() == %d", store.Available())
}

func TestStoreImmediatePut(t *testing.T) {
	sim := &Simulation{}
	store := NewStoreWithCapacity[int](sim, 1)
	finished := false

	sim.Process(func(proc Process) {
		// store is not full, immediate put request
		put_ev := store.Put(0)

		assertf(t, len(store.gets) == 0, "len(store.gets) == %d", len(store.puts))
		assertf(t, len(store.puts) == 0, "len(store.puts) == %d", len(store.puts))
		assertf(t, store.Available() == 1, "store.Available() == %d", store.Available())
		assertf(t, put_ev.Triggered(), "put_ev.Triggered() == false")

		// store is not empty, immediate get request
		get_ev := store.Get()

		assertf(t, len(store.gets) == 0, "len(store.gets) == %d", len(store.puts))
		assertf(t, len(store.puts) == 0, "len(store.puts) == %d", len(store.puts))
		assertf(t, store.Available() == 0, "store.Available() == %d", store.Available())
		assertf(t, get_ev.Triggered(), "get_ev.Triggered() == false")

		// store is empty, get request queued
		get_ev = store.Get()

		assertf(t, len(store.gets) == 1, "len(store.gets) == %d", len(store.puts))
		assertf(t, len(store.puts) == 0, "len(store.puts) == %d", len(store.puts))
		assertf(t, store.Available() == 0, "store.Available() == %d", store.Available())
		assertf(t, !get_ev.Triggered(), "get_ev.Triggered() == true")

		// store is not full, immediate put request
		put_ev = store.Put(1)

		assertf(t, len(store.gets) == 1, "len(store.gets) == %d", len(store.puts))
		assertf(t, len(store.puts) == 0, "len(store.puts) == %d", len(store.puts))
		assertf(t, store.Available() == 1, "store.Available() == %d", store.Available())
		assertf(t, put_ev.Triggered(), "put_ev.Triggered() == false")

		// get request will now be triggered
		proc.Wait(get_ev)

		assertf(t, len(store.gets) == 0, "len(store.gets) == %d", len(store.puts))
		assertf(t, len(store.puts) == 0, "len(store.puts) == %d", len(store.puts))
		assertf(t, store.Available() == 0, "store.Available() == %d", store.Available())

		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}

func TestStoreImmediateGet(t *testing.T) {
	sim := &Simulation{}
	store := NewStoreWithCapacity[int](sim, 1)
	finished := false

	sim.Process(func(proc Process) {
		// store is not full, immediate put request
		put_ev := store.Put(0)

		assertf(t, len(store.gets) == 0, "len(store.gets) == %d", len(store.puts))
		assertf(t, len(store.puts) == 0, "len(store.puts) == %d", len(store.puts))
		assertf(t, store.Available() == 1, "store.Available() == %d", store.Available())
		assertf(t, put_ev.Triggered(), "put_ev.Triggered() == false")

		// store is full, put request queued
		put_ev = store.Put(1)

		assertf(t, len(store.gets) == 0, "len(store.gets) == %d", len(store.puts))
		assertf(t, len(store.puts) == 1, "len(store.puts) == %d", len(store.puts))
		assertf(t, store.Available() == 1, "store.Available() == %d", store.Available())
		assertf(t, !put_ev.Triggered(), "put_ev.Triggered() == true")

		// store is not empty, immediate get request
		get_ev := store.Get()

		assertf(t, len(store.gets) == 0, "len(store.gets) == %d", len(store.puts))
		assertf(t, len(store.puts) == 1, "len(store.puts) == %d", len(store.puts))
		assertf(t, store.Available() == 0, "store.Available() == %d", store.Available())
		assertf(t, get_ev.Triggered(), "get_ev.Triggered() == false")

		// put request will now be triggered
		proc.Wait(put_ev)

		assertf(t, len(store.gets) == 0, "len(store.gets) == %d", len(store.puts))
		assertf(t, len(store.puts) == 0, "len(store.puts) == %d", len(store.puts))
		assertf(t, store.Available() == 1, "store.Available() == %d", store.Available())

		finished = true
	})

	sim.Run()
	assertf(t, finished == true, "finished == false")
}
