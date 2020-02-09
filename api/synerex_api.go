package api

import (
	simapi "github.com/synerex/synerex_alpha/api/simulation"
)

// Demand
// NewDemand returns empty Demand.
func NewDemand() *Demand {
	return &Demand{}
}

// NewSupply returns empty Supply.
func NewSupply() *Supply {
	return &Supply{}
}

func (dm *Demand) WithSimDemand(r *simapi.SimDemand) *Demand {
	dm.ArgOneof = &Demand_SimDemand{r}
	return dm
}

func (sp *Supply) WithSimSupply(c *simapi.SimSupply) *Supply {
	sp.ArgOneof = &Supply_SimSupply{c}
	return sp
}
