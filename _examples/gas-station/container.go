package main

import (
	"math"

	"github.com/fschuetz04/simgo"
)

type Container struct {
	sim  EventGenerator
	lvl  float64
	cap  float64
	gets []AmountEvent
	puts []AmountEvent
}

type AmountEvent struct {
	*simgo.Event
	amount float64
}

type EventGenerator interface {
	Event() *simgo.Event
}

func NewContainer(sim EventGenerator) *Container {
	return NewFilledCappedContainer(sim, 0, math.Inf(1))
}

func NewCappedContainer(sim EventGenerator, cap float64) *Container {
	return NewFilledCappedContainer(sim, 0, cap)
}

func NewFilledContainer(sim EventGenerator, lvl float64) *Container {
	return NewFilledCappedContainer(sim, lvl, math.Inf(1))
}

func NewFilledCappedContainer(sim EventGenerator, lvl float64, cap float64) *Container {
	return &Container{sim: sim, lvl: lvl, cap: cap}
}

func (con *Container) Get(amount float64) AmountEvent {
	if amount < 0 {
		panic("(*Container).Get: amount must not be negative")
	}

	ev := con.newAmountEvent(amount)
	con.gets = append(con.gets, ev)

	con.triggerGets(true)

	return ev
}

func (con *Container) Put(amount float64) AmountEvent {
	if amount < 0 {
		panic("(*Container).Put: amount must not be negative")
	}

	ev := con.newAmountEvent(amount)
	con.puts = append(con.puts, ev)

	con.triggerPuts(true)

	return ev
}

func (con *Container) newAmountEvent(amount float64) AmountEvent {
	return AmountEvent{Event: con.sim.Event(), amount: amount}
}

func (con *Container) triggerGets(triggerPuts bool) {
	for {
		triggered := false

		for len(con.gets) > 0 && con.gets[0].amount <= con.lvl {
			get := con.gets[0]
			con.gets = con.gets[1:]

			if get.Aborted() {
				continue
			}

			con.lvl -= get.amount
			get.Trigger()
			triggered = true
		}

		if triggered && triggerPuts {
			con.triggerPuts(false)
		} else {
			break
		}
	}
}

func (con *Container) triggerPuts(triggerGets bool) {
	for {
		triggered := false

		for len(con.puts) > 0 && con.puts[0].amount <= con.cap-con.lvl {
			put := con.puts[0]
			con.puts = con.puts[1:]

			if put.Aborted() {
				continue
			}

			con.lvl += put.amount
			put.Trigger()
			triggered = true
		}

		if triggered && triggerGets {
			con.triggerGets(false)
		} else {
			break
		}
	}
}
