package simgo

// Resource can be used by a limited number of processes at a time.
type Resource struct {
	// sim is the reference to the simulation.
	sim *Simulation

	// reqs holds the list of pending request events.
	reqs []*Event

	// available is the number of available instances.
	available int
}

// NewResource creates a resource for the given simulation with the given
// number of available instances.
func NewResource(sim *Simulation, available int) *Resource {
	return &Resource{sim: sim, available: available}
}

// Available returns the number of available instances of the resource.
func (res *Resource) Available() int {
	return res.available
}

// Request requests an instance of the resource.
func (res *Resource) Request() *Event {
	req := res.sim.Event()
	res.reqs = append(res.reqs, req)

	res.triggerRequests()

	return req
}

// Release releases an instance of the resource.
func (res *Resource) Release() {
	res.available++

	res.triggerRequests()
}

// triggerRequests triggers pending request events until no more instances are
// available.
func (res *Resource) triggerRequests() {
	for len(res.reqs) > 0 && res.available > 0 {
		req := res.reqs[0]
		res.reqs = res.reqs[1:]

		if !req.Trigger() {
			continue
		}

		res.available--
	}
}
