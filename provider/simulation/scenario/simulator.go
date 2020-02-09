package main

import (
	"github.com/synerex/synerex_alpha/api/simulation/clock"
)

// SynerexSimulator :
type Simulator struct {
	Clock *clock.Clock
}

// NewSenerexSimulator:
func NewSimulator(clockInfo *clock.Clock) *Simulator {

	sim := &Simulator{
		Clock: clockInfo,
	}

	return sim
}

// ForwardStep :
func (sim *Simulator) ForwardStep() {
	sim.Clock.Forward()
}

// ForwardStep :
func (sim *Simulator) BackwardStep() {
	sim.Clock.Backward()
}
