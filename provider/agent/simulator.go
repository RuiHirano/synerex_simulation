package main

import (
	//"log"

	"github.com/paulmach/orb/geojson"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/common"
	algo "github.com/synerex/synerex_alpha/provider/agent/algorithm"
	"github.com/synerex/synerex_alpha/provider/simutil"
)

var (
	geoInfo *geojson.FeatureCollection
)

// SynerexSimulator :
type Simulator struct {
	Clock     *clock.Clock
	Area      *area.Area
	AgentType agent.AgentType
	Agents    []*agent.Agent
}

// NewSenerexSimulator:
func NewSimulator(clockInfo *clock.Clock, areaInfo *area.Area, agentType agent.AgentType) *Simulator {

	sim := &Simulator{
		Clock:     clockInfo,
		Area:      areaInfo,
		AgentType: agentType,
		Agents:    make([]*agent.Agent, 0),
	}

	return sim
}

// ForwardClock :
func (sim *Simulator) ForwardClock() {
	//log.Printf("-------clock %v", sim.Clock)
	sim.Clock.Forward()
}

// ForwardClock :
func (sim *Simulator) BackwardClock() {
	sim.Clock.Backward()
}

// SetObstacles :　Obstaclesを追加する関数
func (sim *Simulator) SetGeoInfo(_geoInfo *geojson.FeatureCollection) {
	geoInfo = _geoInfo
}

// SetArea :　Areaを追加する関数
func (sim *Simulator) SetArea(areaInfo *area.Area) {
	sim.Area = areaInfo
}

// GetArea :　Areaを取得する関数
func (sim *Simulator) GetArea() *area.Area {
	return sim.Area
}

// AddAgents :　Agentsを追加する関数
func (sim *Simulator) AddAgents(agentsInfo []*agent.Agent) {
	newAgents := make([]*agent.Agent, 0)
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
func (sim *Simulator) SetAgents(agentsInfo []*agent.Agent) {
	newAgents := make([]*agent.Agent, 0)
	for _, agentInfo := range agentsInfo {
		if agentInfo.Type == sim.AgentType && IsAgentInArea(agentInfo.Route.Position, sim.Area.DuplicateArea) {
			newAgents = append(newAgents, agentInfo)
		}
	}
	sim.Agents = newAgents
}

// ClearAgents :　Agentsを追加する関数
func (sim *Simulator) ClearAgents() {
	sim.Agents = make([]*agent.Agent, 0)
}

// GetAgents :　Agentsを取得する関数
func (sim *Simulator) GetAgents() []*agent.Agent {
	return sim.Agents
}

// UpdateDuplicateAgents :　重複エリアのエージェントを更新する関数
func (sim *Simulator) UpdateDuplicateAgents(nextControlAgents []*agent.Agent, neighborAgents []*agent.Agent) []*agent.Agent {
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
func (sim *Simulator) ForwardStep(sameAreaAgents []*agent.Agent) []*agent.Agent {
	IsRVO2 := true
	nextControlAgents := sim.GetAgents()

	if IsRVO2 {
		// RVO2
		rvo2route := algo.NewRVO2Route(sim.Clock.TimeStep, sim.Clock.GlobalTime, sim.Area, sim.Agents, sim.AgentType)
		// Agent計算
		nextControlAgents = rvo2route.CalcNextAgents()

	} else {
		// 干渉なしで目的地へ進む
		simpleRoute := algo.NewSimpleRoute(sim.Clock.TimeStep, sim.Clock.GlobalTime, sim.Area, sim.Agents, sim.AgentType)
		nextControlAgents = simpleRoute.CalcNextAgents()
		/*newAgents := make([]*agent.Agent, 0)
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
func IsAgentInArea(position *common.Coord, areaCoords []*common.Coord) bool {
	lat := position.Latitude
	lon := position.Longitude
	maxLat, maxLon, minLat, minLon := simutil.GetCoordRange(areaCoords)
	if minLat < lat && lat < maxLat && minLon < lon && lon < maxLon {
		return true
	} else {
		return false
	}
}
