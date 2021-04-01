package main

import "github.com/fschuetz04/simgo"

type Resource struct {
	sim  EventGenerator
	reqs []*simgo.Event
	cap  int
}

func NewResource(sim EventGenerator, capacity int) *Resource {
	return &Resource{sim: sim, cap: capacity}
}

func (res *Resource) Request() *simgo.Event {
	req := res.sim.Event()
	res.reqs = append(res.reqs, req)

	res.triggerRequests()

	return req
}

func (res *Resource) Release() {
	res.cap++

	res.triggerRequests()
}

func (res *Resource) triggerRequests() {
	for res.cap > 0 && len(res.reqs) > 0 {
		req := res.reqs[0]
		res.reqs = res.reqs[1:]

		if req.Aborted() {
			continue
		}

		res.cap--
		req.Trigger()
	}
}
