package main

import (
	//"log"

	//"github.com/paulmach/orb/geojson"
	"github.com/synerex/synerex_alpha/api"
	algo "github.com/synerex/synerex_alpha/provider/agent/algorithm"
)

var (
//geoInfo *geojson.FeatureCollection
)

// SynerexSimulator :
type Simulator2 struct {
	Agents []*api.Agent
}

// NewSenerexSimulator:
func NewSimulator2() *Simulator2 {

	sim := &Simulator2{
		Agents: make([]*api.Agent, 0),
	}

	return sim
}

// AddAgents :　Agentsを追加する関数
func (sim *Simulator2) AddAgents(agentsInfo []*api.Agent) {
	newAgents := make([]*api.Agent, 0)
	for _, agentInfo := range agentsInfo {
		newAgents = append(newAgents, agentInfo)
	}
	sim.Agents = append(sim.Agents, newAgents...)
}

// SetAgents :　Agentsをセットする関数
func (sim *Simulator2) SetAgents(agentsInfo []*api.Agent) {
	newAgents := make([]*api.Agent, 0)
	for _, agentInfo := range agentsInfo {
		newAgents = append(newAgents, agentInfo)
	}
	sim.Agents = newAgents
}

// ClearAgents :　Agentsを追加する関数
func (sim *Simulator2) ClearAgents() {
	sim.Agents = make([]*api.Agent, 0)
}

// GetAgents :　Agentsを取得する関数
func (sim *Simulator2) GetAgents() []*api.Agent {
	return sim.Agents
}

// ForwardStep :　次の時刻のエージェントを計算する関数
func (sim *Simulator2) ForwardStep() []*api.Agent {

	nextAgents := sim.GetAgents()
	// Agent計算
	rvo2route := algo.NewRVO2Route2(sim.Agents)
	nextAgents = rvo2route.CalcNextAgents()

	//simpleroute := algo.NewSimpleRoute2(sim.Agents)
	//nextAgents = simpleroute.CalcNextAgents()

	return nextAgents
}
