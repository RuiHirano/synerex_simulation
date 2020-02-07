
import (
	common "github.com/synerex/synerex_alpha/api/simulation/common"
	agent "github.com/synerex/synerex_alpha/api/simulation/agent"
)

func NewAgent(agentType common.AgentType)*Agent{
	a := &Agent{
		Id: 0,
		Type: agentType,
	}
	return a
}

func NewPedestrian(agentType common.AgentType, ped *agent.Pedestrian)*Agent{
	a := &Agent{
		Id: 0,
		Type: agentType,
	}
	a.WithPedestrian(ped)
	return a
}

func (a *Agent) WithPedestrian(p *agent.Pedestrian) *Agent {
	a.Data = &Agent_Pedestrian{a}
	return a
}

func NewCar(agentType common.AgentType, car *agent.Car)*Agent{
	a := &Agent{
		Id: 0,
		Type: agentType,
	}
	a.WithCar(car)
	return a
}

func (a *Agent) WithCar(c *agent.Car) *Agent {
	a.Data = &Agent_Pedestrian{c}
	return a
}

// Agentを取得するDemand
func GetAgents(sclient *sxutil.SMServiceClient) uint64{
	getAgentsRequest := &agent.GetAgentsRequest{}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_GET_AGENTS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_GetAgentsRequest{getAgentsRequest},
	}

	nm := "GetAgentsRequest"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	mu.Lock()
	id := sclient.RegisterDemand(opts)
	mu.Unlock()
	return id
}

// Agentを取得するSupply
func SendAgents(sclient *sxutil.SMServiceClient, agents []*agent.Agent, agentType agent.AgentType, areaId uint64) uint64{
	getAgentsResponse := &agent.GetAgentsResponse{
		Agents: agents,
		AgentType: agentType,
		AreaId: areaId,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_GET_AGENTS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_GetAgentsResponse{getAgentsResponse},
	}

	nm := "GetAgentsResponse"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	mu.Lock()
	id := sclient.ProposeSupply(opts)
	mu.Unlock()
	return id
}

// AgentをセットするDemand
func SetAgentsOrder(sclient *sxutil.SMServiceClient) uint64{
	getAgentsRequest := &agent.GetAgentsRequest{
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_GET_AGENTS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_GetAgentsRequest{getAgentsRequest},
	}

	nm := "GetAgentsRequest"
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	mu.Lock()
	id := sclient.RegisterDemand(opts)
	mu.Unlock()
	return id
}

// Agentのセット完了
func FinishSetAgents(sclient *sxutil.SMServiceClient, agents []*agent.Agent, agentType agent.AgentType, areaId uint64) uint64{
	getAgentsResponse := &agent.GetAgentsResponse{
		Agents: agents,
		AgentType: agentType,
		AreaId: areaId,
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_GET_AGENTS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_GetAgentsResponse{getAgentsResponse},
	}

	nm := "GetAgentsResponse"
	js := ""
	opts := &sxutil.SupplyOpts{Name: nm, JSON: js, SimSupply: simSupply}

	mu.Lock()
	id := sclient.ProposeSupply(opts)
	mu.Unlock()
	return id
}
