package main

import (
	//"log"

	"github.com/paulmach/orb/geojson"
	"github.com/synerex/synerex_alpha/api"
	algo "github.com/synerex/synerex_alpha/provider/agent/algorithm"
	"github.com/synerex/synerex_alpha/provider/simutil"
)

var (
	geoInfo *geojson.FeatureCollection
)

// SynerexSimulator :
type Simulator struct {
	Clock     *api.Clock
	Area      *api.Area
	AgentType api.AgentType
	Agents    []*api.Agent
}

// NewSenerexSimulator:
func NewSimulator(clockInfo *api.Clock, areaInfo *api.Area, agentType api.AgentType) *Simulator {

	sim := &Simulator{
		Clock:     clockInfo,
		Area:      areaInfo,
		AgentType: agentType,
		Agents:    make([]*api.Agent, 0),
	}

	return sim
}

// ForwardClock :
func (sim *Simulator) ForwardClock() {
	//log.Printf("-------clock %v", sim.Clock)
	sim.Clock = &api.Clock{
		GlobalTime: sim.Clock.GetGlobalTime() + 1,
	}
}

// ForwardClock :
func (sim *Simulator) BackwardClock() {
	sim.Clock = &api.Clock{
		GlobalTime: sim.Clock.GetGlobalTime() - 1,
	}
}

// SetObstacles :　Obstaclesを追加する関数
func (sim *Simulator) SetGeoInfo(_geoInfo *geojson.FeatureCollection) {
	geoInfo = _geoInfo
}

// SetArea :　Areaを追加する関数
func (sim *Simulator) SetArea(areaInfo *api.Area) {
	sim.Area = areaInfo
}

// GetArea :　Areaを取得する関数
func (sim *Simulator) GetArea() *api.Area {
	return sim.Area
}

// AddAgents :　Agentsを追加する関数
func (sim *Simulator) AddAgents(agentsInfo []*api.Agent) {
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
func (sim *Simulator) SetAgents(agentsInfo []*api.Agent) {
	newAgents := make([]*api.Agent, 0)
	for _, agentInfo := range agentsInfo {
		if agentInfo.Type == sim.AgentType && IsAgentInArea(agentInfo.Route.Position, sim.Area.DuplicateArea) {
			newAgents = append(newAgents, agentInfo)
		}
	}
	sim.Agents = newAgents
}

// ClearAgents :　Agentsを追加する関数
func (sim *Simulator) ClearAgents() {
	sim.Agents = make([]*api.Agent, 0)
}

// GetAgents :　Agentsを取得する関数
func (sim *Simulator) GetAgents() []*api.Agent {
	return sim.Agents
}

// UpdateDuplicateAgents :　重複エリアのエージェントを更新する関数
func (sim *Simulator) UpdateDuplicateAgents(nextControlAgents []*api.Agent, neighborAgents []*api.Agent) []*api.Agent {
	nextAgents := nextControlAgents
	for _, neighborAgent := range neighborAgents {
		//　隣のエージェントが自分のエリアにいてかつ自分のエリアのエージェントと被ってない場合更新
		if len(nextControlAgents) == 0 {
			position := neighborAgent.Route.Position
			if IsAgentInArea(position, sim.Area.DuplicateArea) {
				nextAgents = append(nextAgents, neighborAgent)
			}
		} else {
			isAppendAgent := false
			for _, sameAreaAgent := range nextControlAgents {
				// 自分の管理しているエージェントではなく管理エリアに入っていた場合更新する
				//FIX Duplicateじゃない？
				position := neighborAgent.Route.Position
				if neighborAgent.Id != sameAreaAgent.Id && IsAgentInArea(position, sim.Area.DuplicateArea) {
					isAppendAgent = true
				}
			}
			if isAppendAgent {
				nextAgents = append(nextAgents, neighborAgent)
			}
		}
	}
	return nextAgents
}

// ForwardStep :　次の時刻のエージェントを計算する関数
func (sim *Simulator) ForwardStep(sameAreaAgents []*api.Agent) []*api.Agent {
	IsRVO2 := true
	nextControlAgents := sim.GetAgents()

	if IsRVO2 {
		// RVO2
		timeStep := float64(1.0)
		rvo2route := algo.NewRVO2Route(timeStep, sim.Clock.GlobalTime, sim.Area, sim.Agents, sim.AgentType)
		// Agent計算
		nextControlAgents = rvo2route.CalcNextAgents()

	} else {
		// 干渉なしで目的地へ進む
		timeStep := float64(1.0)
		simpleRoute := algo.NewSimpleRoute(timeStep, sim.Clock.GlobalTime, sim.Area, sim.Agents, sim.AgentType)
		nextControlAgents = simpleRoute.CalcNextAgents()
		/*newAgents := make([]*api.Agent, 0)
		for _, agentInfo := range nextControlAgents {
			if IsAgentInArea(agentInfo.Route.Position, sim.Area.ControlArea) {
				newAgents = append(newAgents, agentInfo)
			}
		}
		nextControlAgents = newAgents*/
		// 干渉なしで目的地へ進む
		//simpleRoute := NewSimpleRoute(sim.TimeStep, sim.GlobalTime, sim.Map, sim.Agents, sim.AgentType)
		//nextControlAgents = simpleRoute.CalcNextAgentsBySimple()

	}
	return nextControlAgents
}

// エージェントがエリアの中にいるかどうか
func IsAgentInArea(position *api.Coord, areaCoords []*api.Coord) bool {
	lat := position.Latitude
	lon := position.Longitude
	maxLat, maxLon, minLat, minLon := simutil.GetCoordRange(areaCoords)
	if minLat < lat && lat < maxLat && minLon < lon && lon < maxLon {
		return true
	} else {
		return false
	}
}
