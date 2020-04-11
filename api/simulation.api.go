package api

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	mu        sync.Mutex
	waitChMap map[SupplyType]chan *Supply
	//spMesMap            map[SupplyType]*Message
	CHANNEL_BUFFER_SIZE int
)

func init() {
	waitChMap = make(map[SupplyType]chan *Supply)
	CHANNEL_BUFFER_SIZE = 10
}

type Clients struct {
	AgentClient    *SMServiceClient
	ClockClient    *SMServiceClient
	ProviderClient *SMServiceClient
}

type SimAPI struct {
	MyClients *Clients
}

func NewSimAPI() *SimAPI {
	s := &SimAPI{}
	return s
}

////////////////////////////////////////////////////////////
////////////        Supply Demand Function       ///////////
///////////////////////////////////////////////////////////

func (s *SimAPI) RegistClients(client SynerexClient, providerId uint64, argJson string) {

	agentClient := NewSMServiceClient(client, ChannelType_AGENT_SERVICE, providerId, argJson)
	clockClient := NewSMServiceClient(client, ChannelType_CLOCK_SERVICE, providerId, argJson)
	providerClient := NewSMServiceClient(client, ChannelType_PROVIDER_SERVICE, providerId, argJson)

	clients := &Clients{
		AgentClient:    agentClient,
		ClockClient:    clockClient,
		ProviderClient: providerClient,
	}

	s.MyClients = clients
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (s *SimAPI) SubscribeAll(demandCallback func(*SMServiceClient, *Demand), supplyCallback func(*SMServiceClient, *Supply)) error {

	// SubscribeDemand, SubscribeSupply
	go subscribeDemand(s.MyClients.AgentClient, demandCallback)

	go subscribeDemand(s.MyClients.ClockClient, demandCallback)

	go subscribeDemand(s.MyClients.ProviderClient, demandCallback)

	go subscribeSupply(s.MyClients.ClockClient, supplyCallback)

	go subscribeSupply(s.MyClients.ProviderClient, supplyCallback)

	go subscribeSupply(s.MyClients.AgentClient, supplyCallback)

	time.Sleep(1 * time.Second)
	return nil
}

func subscribeSupply(client *SMServiceClient, supplyCallback func(*SMServiceClient, *Supply)) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed? Reconnect...")
	time.Sleep(1 * time.Second)
	subscribeSupply(client, supplyCallback)

}

func subscribeDemand(client *SMServiceClient, demandCallback func(*SMServiceClient, *Demand)) {

	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed?")
}

func sendDemand(sclient *SMServiceClient, simDemand *SimDemand) uint64 {
	nm := ""
	js := ""
	opts := &DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	mu.Lock()
	id := sclient.RegisterDemand(opts)
	mu.Unlock()
	return id
}

func sendSupply(sclient *SMServiceClient, simSupply *SimSupply) uint64 {
	nm := ""
	js := ""
	opts := &SupplyOpts{Name: nm, JSON: js, SimSupply: simSupply}

	mu.Lock()
	id := sclient.RegisterSupply(opts)
	mu.Unlock()
	return id
}

//////////////////////////
// add new function////////
/////////////////////////
func sendSyncDemand(sclient *SMServiceClient, simDemand *SimDemand) uint64 {
	nm := ""
	js := ""
	opts := &DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	mu.Lock()
	id := sclient.SyncDemand(opts)
	mu.Unlock()
	return id
}

func sendSyncSupply(sclient *SMServiceClient, simSupply *SimSupply) uint64 {
	nm := ""
	js := ""
	opts := &SupplyOpts{Name: nm, JSON: js, SimSupply: simSupply}

	mu.Lock()
	id := sclient.SyncSupply(opts)
	mu.Unlock()
	return id
}

///////////////////////////////////////////
/////////////   Agent API   //////////////
//////////////////////////////////////////

// AgentをセットするDemand
func (s *SimAPI) SetAgentRequest(senderId uint64, targets []uint64, agents []*Agent) uint64 {

	uid, _ := uuid.NewRandom()
	setAgentRequest := &SetAgentRequest{
		Agents: agents,
	}

	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_SET_AGENT_REQUEST,
		Data:     &SimDemand_SetAgentRequest{setAgentRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.AgentClient, simDemand)

	return msgId
}

// Agentのセット完了
func (s *SimAPI) SetAgentResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	setAgentResponse := &SetAgentResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_SET_AGENT_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_SetAgentResponse{setAgentResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.AgentClient, simSupply)

	return msgId
}

// AgentをセットするDemand
func (s *SimAPI) GetAgentRequest(senderId uint64, targets []uint64) uint64 {

	uid, _ := uuid.NewRandom()
	getAgentRequest := &GetAgentRequest{}

	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_GET_AGENT_REQUEST,
		Data:     &SimDemand_GetAgentRequest{getAgentRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.AgentClient, simDemand)

	return msgId
}

// Agentのセット完了
func (s *SimAPI) GetAgentResponse(senderId uint64, targets []uint64, msgId uint64, agents []*Agent) uint64 {
	getAgentResponse := &GetAgentResponse{
		Agents: agents,
	}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_GET_AGENT_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_GetAgentResponse{getAgentResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.AgentClient, simSupply)

	return msgId
}

///////////////////////////////////////////
/////////////   Provider API   //////////////
//////////////////////////////////////////

// Providerを登録するDemand
func (s *SimAPI) RegistProviderRequest(senderId uint64, targets []uint64, providerInfo *Provider) uint64 {
	registProviderRequest := &RegistProviderRequest{
		Provider: providerInfo,
	}

	uid, _ := uuid.NewRandom()
	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_REGIST_PROVIDER_REQUEST,
		Data:     &SimDemand_RegistProviderRequest{registProviderRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.ProviderClient, simDemand)

	return msgId
}

// Providerを登録するSupply
func (s *SimAPI) RegistProviderResponse(senderId uint64, targets []uint64, msgId uint64, providerInfo *Provider) uint64 {
	registProviderResponse := &RegistProviderResponse{
		Provider: providerInfo,
	}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_REGIST_PROVIDER_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_RegistProviderResponse{registProviderResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.ProviderClient, simSupply)

	return msgId
}

// Providerを登録するDemand
func (s *SimAPI) UpdateProvidersRequest(senderId uint64, targets []uint64, providers []*Provider) uint64 {
	updateProvidersRequest := &UpdateProvidersRequest{
		Providers: providers,
	}

	uid, _ := uuid.NewRandom()
	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_UPDATE_PROVIDERS_REQUEST,
		Data:     &SimDemand_UpdateProvidersRequest{updateProvidersRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.ProviderClient, simDemand)

	return msgId
}

// Providerを登録するSupply
func (s *SimAPI) UpdateProvidersResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	updateProvidersResponse := &UpdateProvidersResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_UPDATE_PROVIDERS_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_UpdateProvidersResponse{updateProvidersResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.ProviderClient, simSupply)

	return msgId
}

///////////////////////////////////////////
/////////////   Clock API   //////////////
//////////////////////////////////////////

func (s *SimAPI) SetClockRequest(senderId uint64, targets []uint64, clockInfo *Clock) uint64 {
	setClockRequest := &SetClockRequest{
		Clock: clockInfo,
	}

	uid, _ := uuid.NewRandom()
	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_SET_CLOCK_REQUEST,
		Data:     &SimDemand_SetClockRequest{setClockRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.ClockClient, simDemand)

	return msgId
}

// Agentを取得するSupply
func (s *SimAPI) SetClockResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	setClockResponse := &SetClockResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_SET_CLOCK_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_SetClockResponse{setClockResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.ClockClient, simSupply)

	return msgId
}

func (s *SimAPI) ForwardClockRequest(senderId uint64, targets []uint64) uint64 {
	forwardClockRequest := &ForwardClockRequest{}

	uid, _ := uuid.NewRandom()
	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_FORWARD_CLOCK_REQUEST,
		Data:     &SimDemand_ForwardClockRequest{forwardClockRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.ClockClient, simDemand)

	return msgId
}

// Agentを取得するSupply
func (s *SimAPI) ForwardClockResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	forwardClockResponse := &ForwardClockResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_FORWARD_CLOCK_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_ForwardClockResponse{forwardClockResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.ClockClient, simSupply)

	return msgId
}

func (s *SimAPI) StartClockRequest(senderId uint64, targets []uint64) uint64 {
	startClockRequest := &StartClockRequest{}

	uid, _ := uuid.NewRandom()
	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_START_CLOCK_REQUEST,
		Data:     &SimDemand_StartClockRequest{startClockRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.ClockClient, simDemand)

	return msgId
}

// Agentを取得するSupply
func (s *SimAPI) StartClockResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	startClockResponse := &StartClockResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_START_CLOCK_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_StartClockResponse{startClockResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.ClockClient, simSupply)

	return msgId
}

func (s *SimAPI) StopClockRequest(senderId uint64, targets []uint64) uint64 {
	stopClockRequest := &StopClockRequest{}

	uid, _ := uuid.NewRandom()
	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_STOP_CLOCK_REQUEST,
		Data:     &SimDemand_StopClockRequest{stopClockRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.ClockClient, simDemand)

	return msgId
}

// Agentを取得するSupply
func (s *SimAPI) StopClockResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	stopClockResponse := &StopClockResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_STOP_CLOCK_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_StopClockResponse{stopClockResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.ClockClient, simSupply)

	return msgId
}

///////////////////////////////////////////
/////////////   Pod API   //////////////
//////////////////////////////////////////

// AgentをセットするDemand
func (s *SimAPI) CreatePodRequest(senderId uint64, targets []uint64) uint64 {

	uid, _ := uuid.NewRandom()
	createPodRequest := &CreatePodRequest{}

	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_SET_AGENT_REQUEST,
		Data:     &SimDemand_CreatePodRequest{createPodRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.AgentClient, simDemand)

	return msgId
}

// Agentのセット完了
func (s *SimAPI) CreatePodResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	createPodResponse := &CreatePodResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_SET_AGENT_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_CreatePodResponse{createPodResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.AgentClient, simSupply)

	return msgId
}

// AgentをセットするDemand
func (s *SimAPI) DeletePodRequest(senderId uint64, targets []uint64) uint64 {

	uid, _ := uuid.NewRandom()
	deletePodRequest := &DeletePodRequest{}

	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_GET_AGENT_REQUEST,
		Data:     &SimDemand_DeletePodRequest{deletePodRequest},
		Targets:  targets,
	}

	sendSyncDemand(s.MyClients.AgentClient, simDemand)

	return msgId
}

// Agentのセット完了
func (s *SimAPI) DeletePodResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	deletePodResponse := &DeletePodResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_GET_AGENT_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_DeletePodResponse{deletePodResponse},
		Targets:  targets,
	}

	sendSyncSupply(s.MyClients.AgentClient, simSupply)

	return msgId
}

///////////////////////////////////////////
/////////////      Wait      //////////////
//////////////////////////////////////////

type Waiter struct {
	WaitSpChMap map[uint64]chan *Supply
	SpMap       map[uint64][]*Supply
	WaitDmChMap map[uint64]chan *Demand
	DmMap       map[uint64][]*Demand
}

func NewWaiter() *Waiter {
	w := &Waiter{
		WaitSpChMap: make(map[uint64]chan *Supply),
		SpMap:       make(map[uint64][]*Supply),
		WaitDmChMap: make(map[uint64]chan *Demand),
		DmMap:       make(map[uint64][]*Demand),
	}
	return w
}

func (w *Waiter) WaitSp(msgId uint64, targets []uint64) []*Supply {
	if len(targets) == 0 {
		return []*Supply{}
	}
	mu.Lock()
	CHANNEL_BUFFER_SIZE := 10
	waitCh := make(chan *Supply, CHANNEL_BUFFER_SIZE)
	w.WaitSpChMap[msgId] = waitCh
	w.SpMap[msgId] = make([]*Supply, 0)
	mu.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case sp, _ := <-waitCh:
				mu.Lock()
				// spのidがidListに入っているか
				if sp.GetSimSupply().GetMsgId() == msgId {
					w.SpMap[sp.GetSimSupply().GetMsgId()] = append(w.SpMap[sp.GetSimSupply().GetMsgId()], sp)

					// 同期が終了したかどうか
					if w.isFinishSpSync(msgId, targets) {
						log.Printf("Finish Wait!")
						mu.Unlock()
						wg.Done()
						return
					}
				}
				mu.Unlock()
			case <-time.After(1000 * time.Millisecond):
				log.Printf("Sync Error... \n")
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	return w.SpMap[msgId]
}

func (w *Waiter) SendSpToWait(sp *Supply) {
	mu.Lock()
	waitCh := w.WaitSpChMap[sp.GetSimSupply().GetMsgId()]
	mu.Unlock()
	waitCh <- sp
}

func (w *Waiter) isFinishSpSync(msgId uint64, targets []uint64) bool {
	for _, sp := range w.SpMap[msgId] {
		senderId := sp.GetSimSupply().GetSenderId()
		isMatch := false
		for _, pid := range targets {
			if senderId == pid {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}

func (w *Waiter) WaitDm(msgId uint64, targets []uint64) []*Demand {
	if len(targets) == 0 {
		return []*Demand{}
	}
	mu.Lock()
	CHANNEL_BUFFER_SIZE := 10
	waitCh := make(chan *Demand, CHANNEL_BUFFER_SIZE)
	w.WaitDmChMap[msgId] = waitCh
	w.DmMap[msgId] = make([]*Demand, 0)
	mu.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case dm, _ := <-waitCh:
				mu.Lock()
				// dmのidがidListに入っているか
				if dm.GetSimDemand().GetMsgId() == msgId {
					w.DmMap[dm.GetSimDemand().GetMsgId()] = append(w.DmMap[dm.GetSimDemand().GetMsgId()], dm)

					// 同期が終了したかどうか
					if w.isFinishDmSync(msgId, targets) {
						log.Printf("Finish Wait!")
						mu.Unlock()
						wg.Done()
						return
					}
				}
				mu.Unlock()
			case <-time.After(1000 * time.Millisecond):
				log.Printf("Sync Error... \n")
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	return w.DmMap[msgId]
}

func (w *Waiter) SendDmToWait(dm *Demand) {
	mu.Lock()
	waitCh := w.WaitDmChMap[dm.GetSimDemand().GetMsgId()]
	mu.Unlock()
	waitCh <- dm
}

func (w *Waiter) isFinishDmSync(msgId uint64, targets []uint64) bool {
	for _, dm := range w.DmMap[msgId] {
		senderId := dm.GetSimDemand().GetSenderId()
		isMatch := false
		for _, pid := range targets {
			if senderId == pid {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}
