package simutil

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/synerex/synerex_alpha/api"
	simapi "github.com/synerex/synerex_alpha/api/simulation"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/provider"
	"github.com/synerex/synerex_alpha/sxutil"
)

var (
	mu                  sync.Mutex
	waitChMap           map[simapi.SupplyType]chan *pb.Supply
	CHANNEL_BUFFER_SIZE int
)

func init() {
	waitChMap = make(map[simapi.SupplyType]chan *pb.Supply)
	CHANNEL_BUFFER_SIZE = 10
}

type Clients struct {
	AgentClient    *sxutil.SMServiceClient
	ClockClient    *sxutil.SMServiceClient
	AreaClient     *sxutil.SMServiceClient
	ProviderClient *sxutil.SMServiceClient
}

type Communicator struct {
	MyClients  *Clients
	MyProvider *provider.Provider
}

func NewCommunicator(providerInfo *provider.Provider) *Communicator {
	c := &Communicator{
		MyProvider: providerInfo,
	}
	return c
}

func (c *Communicator) RegistClients(client pb.SynerexClient, argJson string) {

	agentClient := sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE, argJson)
	clockClient := sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE, argJson)
	providerClient := sxutil.NewSMServiceClient(client, pb.ChannelType_PROVIDER_SERVICE, argJson)

	clients := &Clients{
		AgentClient:    agentClient,
		ClockClient:    clockClient,
		ProviderClient: providerClient,
	}

	c.MyClients = clients
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (c *Communicator) SubscribeAll(demandCallback func(*sxutil.SMServiceClient, *pb.Demand), supplyCallback func(*sxutil.SMServiceClient, *pb.Supply)) error {

	// SubscribeDemand, SubscribeSupply
	go subscribeDemand(c.MyClients.AgentClient, demandCallback)

	go subscribeDemand(c.MyClients.ClockClient, demandCallback)

	go subscribeDemand(c.MyClients.ProviderClient, demandCallback)

	go subscribeSupply(c.MyClients.ClockClient, supplyCallback)

	go subscribeSupply(c.MyClients.ProviderClient, supplyCallback)

	go subscribeSupply(c.MyClients.AgentClient, supplyCallback)

	time.Sleep(3 * time.Second)
	return nil
}

////////////////////////////////////////////////////////////
////////////        Supply Demand Function       ///////////
///////////////////////////////////////////////////////////

func subscribeSupply(client *sxutil.SMServiceClient, supplyCallback func(*sxutil.SMServiceClient, *pb.Supply)) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func subscribeDemand(client *sxutil.SMServiceClient, demandCallback func(*sxutil.SMServiceClient, *pb.Demand)) {

	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func sendDemand(sclient *sxutil.SMServiceClient, simDemand *simapi.SimDemand) uint64 {
	nm := ""
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	mu.Lock()
	id := sclient.RegisterDemand(opts)
	mu.Unlock()
	return id
}

func sendSupply(sclient *sxutil.SMServiceClient, tid uint64, simSupply *simapi.SimSupply) uint64 {
	nm := ""
	js := ""
	opts := &sxutil.SupplyOpts{Target: tid, Name: nm, JSON: js, SimSupply: simSupply}

	mu.Lock()
	id := sclient.ProposeSupply(opts)
	mu.Unlock()
	return id
}

////////////////////////////////////////////////////////////
////////////        Wait Function       ///////////////////
///////////////////////////////////////////////////////////

// SendToSetAgentsResponse : SetAgentsResponseを送る
func (c *Communicator) SendToWaitCh(sp *pb.Supply, supplyType simapi.SupplyType) {
	waitCh := waitChMap[supplyType]
	waitCh <- sp
}

// Wait: 同期が完了するまで待機する関数
func wait(idList []uint64, supplyType simapi.SupplyType) map[uint64]*pb.Supply {

	mu.Lock()
	waitCh := make(chan *pb.Supply, CHANNEL_BUFFER_SIZE)
	waitChMap[supplyType] = waitCh
	mu.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(1)
	pspMap := make(map[uint64]*pb.Supply)
	go func() {
		for {
			select {
			case psp := <-waitCh:
				mu.Lock()
				// spのidがidListに入っているか
				if isSpInIdList(psp, idList) {
					pspMap[psp.SenderId] = psp
					// 同期が終了したかどうか
					if isFinishSync(pspMap, idList) {
						mu.Unlock()
						wg.Done()
						return
					}
				}
				mu.Unlock()
			}
		}
	}()
	wg.Wait()
	return pspMap
}

// isSpInIdList : spのidがidListに入っているか
func isSpInIdList(sp *pb.Supply, idlist []uint64) bool {
	senderId := sp.SenderId
	for _, id := range idlist {
		if senderId == id {
			return true
		}
	}
	return false
}

// isFinishSync : 必要な全てのSupplyを受け取り同期が完了したかどうか
func isFinishSync(pspMap map[uint64]*pb.Supply, idlist []uint64) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			senderId := uint64(sp.SenderId)
			if uint64(id) == senderId {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}

///////////////////////////////////////////
/////////////   Agent API   //////////////
//////////////////////////////////////////

// Agentを取得するDemand
func (c *Communicator) GetAgentsRequest(idList []uint64) (uint64, []*agent.Agent) {
	getAgentsRequest := &agent.GetAgentsRequest{}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_GET_AGENTS_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_GetAgentsRequest{getAgentsRequest},
	}

	id := sendDemand(c.MyClients.AgentClient, simDemand)

	// Wait
	agents := make([]*agent.Agent, 0)
	if idList != nil {
		supplyType := simapi.SupplyType_GET_AGENTS_RESPONSE
		spMap := wait(idList, supplyType)
		for _, sp := range spMap {
			ags := sp.GetSimSupply().GetGetAgentsResponse().GetAgents()
			agents = append(agents, ags...)
		}
	}

	return id, agents
}

// Agentを取得するSupply
func (c *Communicator) GetAgentsResponse(tid uint64, agents []*agent.Agent, agentType agent.AgentType, areaId uint64) uint64 {
	getAgentsResponse := &agent.GetAgentsResponse{
		Agents:    agents,
		AgentType: agentType,
		AreaId:    areaId,
	}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_GET_AGENTS_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_GetAgentsResponse{getAgentsResponse},
	}

	id := sendSupply(c.MyClients.AgentClient, tid, simSupply)

	return id
}

// AgentをセットするDemand
func (c *Communicator) SetAgentsRequest(idList []uint64, agents []*agent.Agent) uint64 {
	setAgentsRequest := &agent.SetAgentsRequest{
		Agents: agents,
	}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_SET_AGENTS_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_SetAgentsRequest{setAgentsRequest},
	}

	id := sendDemand(c.MyClients.AgentClient, simDemand)

	if idList != nil {
		supplyType := simapi.SupplyType_SET_AGENTS_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentのセット完了
func (c *Communicator) SetAgentsResponse(tid uint64) uint64 {
	setAgentsResponse := &agent.SetAgentsResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_SET_AGENTS_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_SetAgentsResponse{setAgentsResponse},
	}

	id := sendSupply(c.MyClients.AgentClient, tid, simSupply)

	return id
}

///////////////////////////////////////////
/////////////   Provider API   //////////////
//////////////////////////////////////////

// Providerを登録するDemand
func (c *Communicator) RegistProviderRequest(idList []uint64, providerInfo *provider.Provider) uint64 {
	registProviderRequest := &provider.RegistProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_REGIST_PROVIDER_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_RegistProviderRequest{registProviderRequest},
	}

	id := sendDemand(c.MyClients.ProviderClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_REGIST_PROVIDER_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator) RegistProviderResponse(tid uint64) uint64 {
	registProviderResponse := &provider.RegistProviderResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_REGIST_PROVIDER_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_RegistProviderResponse{registProviderResponse},
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) KillProviderRequest(idList []uint64, providerInfo *provider.Provider) uint64 {
	killProviderRequest := &provider.KillProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_KILL_PROVIDER_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_KillProviderRequest{killProviderRequest},
	}

	id := sendDemand(c.MyClients.ProviderClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_KILL_PROVIDER_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator) KillProviderResponse(tid uint64) uint64 {
	killProviderResponse := &provider.KillProviderResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_KILL_PROVIDER_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_KillProviderResponse{killProviderResponse},
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) DivideProviderRequest(idList []uint64, providerInfo *provider.Provider) uint64 {
	divideProviderRequest := &provider.DivideProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_DIVIDE_PROVIDER_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_DivideProviderRequest{divideProviderRequest},
	}

	id := sendDemand(c.MyClients.ProviderClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_DIVIDE_PROVIDER_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator) DivideProviderResponse(tid uint64) uint64 {
	divideProviderResponse := &provider.DivideProviderResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_DIVIDE_PROVIDER_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_DivideProviderResponse{divideProviderResponse},
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) UpdateProvidersRequest(idList []uint64, providers []*provider.Provider) uint64 {
	updateProvidersRequest := &provider.UpdateProvidersRequest{
		Providers: providers,
	}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_UPDATE_PROVIDERS_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_UpdateProvidersRequest{updateProvidersRequest},
	}

	id := sendDemand(c.MyClients.ProviderClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_UPDATE_PROVIDERS_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator) UpdateProvidersResponse(tid uint64) uint64 {
	updateProvidersResponse := &provider.UpdateProvidersResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_UPDATE_PROVIDERS_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_UpdateProvidersResponse{updateProvidersResponse},
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) SendProviderStatusRequest(idList []uint64, providerInfo *provider.Provider) uint64 {
	sendProviderStatusRequest := &provider.SendProviderStatusRequest{
		Provider: providerInfo,
	}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_SEND_PROVIDER_STATUS_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_SendProviderStatusRequest{sendProviderStatusRequest},
	}

	id := sendDemand(c.MyClients.ProviderClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_SEND_PROVIDER_STATUS_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator) SendProviderStatusResponse(tid uint64) uint64 {
	sendProviderStatusResponse := &provider.SendProviderStatusResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_SEND_PROVIDER_STATUS_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_SendProviderStatusResponse{sendProviderStatusResponse},
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

///////////////////////////////////////////
/////////////   Clock API   //////////////
//////////////////////////////////////////

func (c *Communicator) UpdateClockRequest(idList []uint64, clockInfo *clock.Clock) uint64 {
	updateClockRequest := &clock.UpdateClockRequest{
		Clock: clockInfo,
	}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_UPDATE_CLOCK_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_UpdateClockRequest{updateClockRequest},
	}

	id := sendDemand(c.MyClients.ClockClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_UPDATE_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator) UpdateClockResponse(tid uint64) uint64 {
	updateClockResponse := &clock.UpdateClockResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_UPDATE_CLOCK_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_UpdateClockResponse{updateClockResponse},
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) SetClockRequest(idList []uint64, clockInfo *clock.Clock) uint64 {
	setClockRequest := &clock.SetClockRequest{
		Clock: clockInfo,
	}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_SET_CLOCK_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_SetClockRequest{setClockRequest},
	}

	id := sendDemand(c.MyClients.ClockClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_SET_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator) SetClockResponse(tid uint64) uint64 {
	setClockResponse := &clock.SetClockResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_SET_CLOCK_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_SetClockResponse{setClockResponse},
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) GetClockRequest(idList []uint64) (uint64, *clock.Clock) {
	getClockRequest := &clock.GetClockRequest{}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_GET_CLOCK_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_GetClockRequest{getClockRequest},
	}

	id := sendDemand(c.MyClients.ClockClient, simDemand)

	// Wait
	var clockInfo *clock.Clock
	if idList != nil {
		supplyType := simapi.SupplyType_GET_CLOCK_RESPONSE
		spMap := wait(idList, supplyType)
		clockInfo = spMap[0].GetSimSupply().GetGetClockResponse().GetClock()
	}

	return id, clockInfo
}

// Agentを取得するSupply
func (c *Communicator) GetClockResponse(tid uint64, clockInfo *clock.Clock) uint64 {
	getClockResponse := &clock.GetClockResponse{
		Clock: clockInfo,
	}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_GET_CLOCK_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_GetClockResponse{getClockResponse},
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) ForwardClockRequest(idList []uint64) uint64 {
	forwardClockRequest := &clock.ForwardClockRequest{}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_FORWARD_CLOCK_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_ForwardClockRequest{forwardClockRequest},
	}

	id := sendDemand(c.MyClients.ClockClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_FORWARD_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator) ForwardClockResponse(tid uint64) uint64 {
	forwardClockResponse := &clock.ForwardClockResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_FORWARD_CLOCK_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_ForwardClockResponse{forwardClockResponse},
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) BackClockRequest(idList []uint64) uint64 {
	backClockRequest := &clock.BackClockRequest{}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_BACK_CLOCK_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_BackClockRequest{backClockRequest},
	}

	id := sendDemand(c.MyClients.ClockClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_BACK_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator) backClockResponse(tid uint64) uint64 {
	backClockResponse := &clock.BackClockResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_BACK_CLOCK_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_BackClockResponse{backClockResponse},
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) StartClockRequest(idList []uint64) uint64 {
	startClockRequest := &clock.StartClockRequest{}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_START_CLOCK_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_StartClockRequest{startClockRequest},
	}

	id := sendDemand(c.MyClients.ClockClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_START_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator) StartClockResponse(tid uint64) uint64 {
	startClockResponse := &clock.StartClockResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_START_CLOCK_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_StartClockResponse{startClockResponse},
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) StopClockRequest(idList []uint64) uint64 {
	stopClockRequest := &clock.StopClockRequest{}

	simDemand := &simapi.SimDemand{
		Type:   simapi.DemandType_STOP_CLOCK_REQUEST,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimDemand_StopClockRequest{stopClockRequest},
	}

	id := sendDemand(c.MyClients.ClockClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_STOP_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator) StopClockResponse(tid uint64) uint64 {
	stopClockResponse := &clock.StopClockResponse{}

	simSupply := &simapi.SimSupply{
		Type:   simapi.SupplyType_STOP_CLOCK_RESPONSE,
		Status: simapi.StatusType_NONE,
		Data:   &simapi.SimSupply_StopClockResponse{stopClockResponse},
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}
