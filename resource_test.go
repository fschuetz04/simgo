package simgo

import "testing"

func TestResourceRequestRelease(t *testing.T) {
	sim := NewSimulation()
	res := NewResource(sim, 1)

	sim.Process(func(proc Process) {
		// resource is not empty, immediate request
		req_ev := res.Request()

		assertf(t, len(res.reqs) == 0, "len(res.reqs) == %d", len(res.reqs))
		assertf(t, res.Available() == 0, "res.Available() == %d", res.Available())
		assertf(t, req_ev.Triggered(), "req_ev.Triggered() == false")

		// resource is empty, request queued
		req_ev = res.Request()

		assertf(t, len(res.reqs) == 1, "len(res.reqs) == %d", len(res.reqs))
		assertf(t, res.Available() == 0, "res.Available() == %d", res.Available())
		assertf(t, !req_ev.Triggered(), "req_ev.Triggered() == true")

		res.Release()

		assertf(t, len(res.reqs) == 0, "len(res.reqs) == %d", len(res.reqs))
		assertf(t, res.Available() == 0, "res.Available() == %d", res.Available())
		assertf(t, req_ev.Triggered(), "req_ev.Triggered() == false")
	})
}
