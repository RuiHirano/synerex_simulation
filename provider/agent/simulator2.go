package main

import (
	//"log"

	//"github.com/paulmach/orb/geojson"
	"math"

	"github.com/synerex/synerex_alpha/api"
	algo "github.com/synerex/synerex_alpha/provider/agent/algorithm"
)

var (
//geoInfo *geojson.FeatureCollection
)

// SynerexSimulator :
type Simulator2 struct {
	Agents    []*api.Agent
	Area      *api.Area
	AgentType api.AgentType
}

// NewSenerexSimulator:
func NewSimulator2(areaInfo *api.Area, agentType api.AgentType) *Simulator2 {

	sim := &Simulator2{
		Agents:    make([]*api.Agent, 0),
		Area:      areaInfo,
		AgentType: agentType,
	}

	return sim
}

// AddAgents :　Agentsを追加する関数
func (sim *Simulator2) AddAgents(agentsInfo []*api.Agent) {
	newAgents := make([]*api.Agent, 0)
	for _, agentInfo := range agentsInfo {
		if agentInfo.Type == sim.AgentType {
			position := agentInfo.Route.Position
			//("Debug %v, %v", position, sim.Area.DuplicateArea)
			if IsAgentInArea(position, sim.Area.DuplicateArea) {
				newAgents = append(newAgents, agentInfo)
			}
		}
	}
	sim.Agents = append(sim.Agents, newAgents...)
}

// SetAgents :　Agentsをセットする関数
func (sim *Simulator2) SetAgents(agentsInfo []*api.Agent) {
	newAgents := make([]*api.Agent, 0)
	for _, agentInfo := range agentsInfo {
		if agentInfo.Type == sim.AgentType {
			position := agentInfo.Route.Position
			//("Debug %v, %v", position, sim.Area.DuplicateArea)
			if IsAgentInArea(position, sim.Area.DuplicateArea) {
				newAgents = append(newAgents, agentInfo)
			}
		}
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

// UpdateDuplicateAgents :　重複エリアのエージェントを更新する関数
func (sim *Simulator2) UpdateDuplicateAgents(nextControlAgents []*api.Agent, neighborAgents []*api.Agent) []*api.Agent {
	nextAgents := nextControlAgents
	for _, neighborAgent := range neighborAgents {
		isAppendAgent := true
		position := neighborAgent.Route.Position
		for _, sameAreaAgent := range nextControlAgents {
			// 自分の管理しているエージェントではなく重複エリアに入っていた場合更新する
			//FIX Duplicateじゃない？
			if neighborAgent.Id == sameAreaAgent.Id {
				isAppendAgent = false
			}
		}
		if isAppendAgent && IsAgentInArea(position, sim.Area.DuplicateArea) {
			nextAgents = append(nextAgents, neighborAgent)
		}
	}
	return nextAgents
}

// ForwardStep :　次の時刻のエージェントを計算する関数
func (sim *Simulator2) ForwardStep(sameAgents []*api.Agent) []*api.Agent {

	nextAgents := sim.GetAgents()
	// Agent計算
	rvo2route := algo.NewRVO2Route2(sim.Agents, sim.Area)
	nextAgents = rvo2route.CalcNextAgents()

	//simpleroute := algo.NewSimpleRoute2(sim.Agents)
	//nextAgents = simpleroute.CalcNextAgents()

	return nextAgents
}

// エージェントがエリアの中にいるかどうか
func IsAgentInArea(position *api.Coord, areaCoords []*api.Coord) bool {
	lat := position.Latitude
	lon := position.Longitude
	maxLat, maxLon, minLat, minLon := GetCoordRange(areaCoords)
	if minLat < lat && lat < maxLat && minLon < lon && lon < maxLon {
		return true
	} else {
		return false
	}
}

func GetCoordRange(coords []*api.Coord) (float64, float64, float64, float64) {
	maxLon, maxLat := math.Inf(-1), math.Inf(-1)
	minLon, minLat := math.Inf(0), math.Inf(0)
	for _, coord := range coords {
		if coord.Latitude > maxLat {
			maxLat = coord.Latitude
		}
		if coord.Longitude > maxLon {
			maxLon = coord.Longitude
		}
		if coord.Latitude < minLat {
			minLat = coord.Latitude
		}
		if coord.Longitude < minLon {
			minLon = coord.Longitude
		}
	}
	return maxLat, maxLon, minLat, minLon
}
