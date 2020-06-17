package api

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/synerex/synerex_alpha/util"
)

var (
	mu        sync.Mutex
	waitChMap map[SupplyType]chan *Supply
	//spMesMap            map[SupplyType]*Message
	logger              *util.Logger
	CHANNEL_BUFFER_SIZE int
)

func init() {
	waitChMap = make(map[SupplyType]chan *Supply)
	logger = util.NewLogger()
	CHANNEL_BUFFER_SIZE = 10
}

type Clients struct {
	AgentClient    *SMServiceClient
	ClockClient    *SMServiceClient
	ProviderClient *SMServiceClient
}

type SimAPI struct {
	MyClients *Clients
	Waiter    *Waiter
}

func NewSimAPI() *SimAPI {
	s := &SimAPI{
		Waiter: NewWaiter(),
	}
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

	return nil
}

func subscribeSupply(client *SMServiceClient, supplyCallback func(*SMServiceClient, *Supply)) {
	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed? Reconnect...")
	time.Sleep(2 * time.Second)
	subscribeSupply(client, supplyCallback)

}

func subscribeDemand(client *SMServiceClient, demandCallback func(*SMServiceClient, *Demand)) {

	//called as goroutine
	ctx := context.Background() // should check proper context
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
	log.Printf("SMarket Server Closed? Reconnect...")
	time.Sleep(2 * time.Second)
	subscribeDemand(client, demandCallback)
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
func (s *SimAPI) SendSyncDemand(sclient *SMServiceClient, simDemand *SimDemand) ([]*Supply, error) {
	nm := ""
	js := ""
	opts := &DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	mu.Lock()
	msgId := simDemand.GetMsgId()
	CHANNEL_BUFFER_SIZE := 10
	waitCh := make(chan *Supply, CHANNEL_BUFFER_SIZE)
	s.Waiter.WaitSpChMap[msgId] = waitCh
	s.Waiter.SpMap[msgId] = make([]*Supply, 0)
	mu.Unlock()

	mu.Lock()
	sclient.SyncDemand(opts)
	mu.Unlock()

	// waitする
	sps := []*Supply{}
	var err error
	targets := simDemand.GetTargets()
	if len(targets) != 0 {
		msgId := simDemand.GetMsgId()
		sps, err = s.Waiter.WaitSp(msgId, targets, 1000)
		s.Waiter = NewWaiter()
	}

	return sps, err
}

func (s *SimAPI) SendSyncSupply(sclient *SMServiceClient, simSupply *SimSupply) uint64 {
	nm := ""
	js := ""
	opts := &SupplyOpts{Name: nm, JSON: js, SimSupply: simSupply}

	mu.Lock()
	id := sclient.SyncSupply(opts)
	mu.Unlock()
	return id
}

func (s *SimAPI) SendSpToWait(sp *Supply) {
	s.Waiter.SendSpToWait(sp)
}

///////////////////////////////////////////
/////////////    Area API   //////////////
//////////////////////////////////////////

// Areaを送るDemand
func (s *SimAPI) SendAreaInfoRequest(senderId uint64, targets []uint64, areas []*Area) ([]*Supply, error) {

	uid, _ := uuid.NewRandom()
	sendAreaInfoRequest := &SendAreaInfoRequest{
		Areas: areas,
	}

	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_SEND_AREA_INFO_REQUEST,
		Data:     &SimDemand_SendAreaInfoRequest{sendAreaInfoRequest},
		Targets:  targets,
	}

	sps, err := s.SendSyncDemand(s.MyClients.ProviderClient, simDemand)

	return sps, err
}

// Agentのセット完了
func (s *SimAPI) SendAreaInfoResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	sendAreaInfoResponse := &SendAreaInfoResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_SEND_AREA_INFO_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_SendAreaInfoResponse{sendAreaInfoResponse},
		Targets:  targets,
	}

	s.SendSyncSupply(s.MyClients.AgentClient, simSupply)

	return msgId
}

///////////////////////////////////////////
/////////////   Agent API   //////////////
//////////////////////////////////////////

// AgentをセットするDemand
func (s *SimAPI) SetAgentRequest(senderId uint64, targets []uint64, agents []*Agent) ([]*Supply, error) {

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

	sps, err := s.SendSyncDemand(s.MyClients.AgentClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.AgentClient, simSupply)

	return msgId
}

// AgentをセットするDemand
func (s *SimAPI) GetAgentRequest(senderId uint64, targets []uint64) ([]*Supply, error) {

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

	sps, err := s.SendSyncDemand(s.MyClients.AgentClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.AgentClient, simSupply)

	return msgId
}

///////////////////////////////////////////
/////////////   Provider API   //////////////
//////////////////////////////////////////

// Providerを登録するDemand
func (s *SimAPI) ReadyProviderRequest(senderId uint64, targets []uint64, providerInfo *Provider) ([]*Supply, error) {
	readyProviderRequest := &ReadyProviderRequest{
		Provider: providerInfo,
	}

	uid, _ := uuid.NewRandom()
	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_READY_PROVIDER_REQUEST,
		Data:     &SimDemand_ReadyProviderRequest{readyProviderRequest},
		Targets:  targets,
	}

	sps, err := s.SendSyncDemand(s.MyClients.ProviderClient, simDemand)

	return sps, err
}

// Providerを登録するSupply
func (s *SimAPI) ReadyProviderResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	readyProviderResponse := &ReadyProviderResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_READY_PROVIDER_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_ReadyProviderResponse{readyProviderResponse},
		Targets:  targets,
	}

	s.SendSyncSupply(s.MyClients.ProviderClient, simSupply)

	return msgId
}

// Providerを登録するDemand
func (s *SimAPI) RegistProviderRequest(senderId uint64, targets []uint64, providerInfo *Provider) ([]*Supply, error) {
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

	sps, err := s.SendSyncDemand(s.MyClients.ProviderClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.ProviderClient, simSupply)

	return msgId
}

// Providerを登録するDemand
func (s *SimAPI) UpdateProvidersRequest(senderId uint64, targets []uint64, providers []*Provider) ([]*Supply, error) {
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

	sps, err := s.SendSyncDemand(s.MyClients.ProviderClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.ProviderClient, simSupply)

	return msgId
}

///////////////////////////////////////////
/////////////   Clock API   //////////////
//////////////////////////////////////////

func (s *SimAPI) SetClockRequest(senderId uint64, targets []uint64, clockInfo *Clock) ([]*Supply, error) {
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

	sps, err := s.SendSyncDemand(s.MyClients.ClockClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.ClockClient, simSupply)

	return msgId
}

func (s *SimAPI) ForwardClockRequest(senderId uint64, targets []uint64) ([]*Supply, error) {
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

	sps, err := s.SendSyncDemand(s.MyClients.ClockClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.ClockClient, simSupply)

	return msgId
}

func (s *SimAPI) ForwardClockInitRequest(senderId uint64, targets []uint64) ([]*Supply, error) {
	forwardClockInitRequest := &ForwardClockInitRequest{}

	uid, _ := uuid.NewRandom()
	msgId := uint64(uid.ID())
	simDemand := &SimDemand{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     DemandType_FORWARD_CLOCK_INIT_REQUEST,
		Data:     &SimDemand_ForwardClockInitRequest{forwardClockInitRequest},
		Targets:  targets,
	}

	sps, err := s.SendSyncDemand(s.MyClients.ClockClient, simDemand)

	return sps, err
}

// Agentを取得するSupply
func (s *SimAPI) ForwardClockInitResponse(senderId uint64, targets []uint64, msgId uint64) uint64 {
	forwardClockInitResponse := &ForwardClockInitResponse{}

	simSupply := &SimSupply{
		MsgId:    msgId,
		SenderId: senderId,
		Type:     SupplyType_FORWARD_CLOCK_INIT_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_ForwardClockInitResponse{forwardClockInitResponse},
		Targets:  targets,
	}

	s.SendSyncSupply(s.MyClients.ClockClient, simSupply)

	return msgId
}

func (s *SimAPI) StartClockRequest(senderId uint64, targets []uint64) ([]*Supply, error) {
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

	sps, err := s.SendSyncDemand(s.MyClients.ClockClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.ClockClient, simSupply)

	return msgId
}

func (s *SimAPI) StopClockRequest(senderId uint64, targets []uint64) ([]*Supply, error) {
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

	sps, err := s.SendSyncDemand(s.MyClients.ClockClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.ClockClient, simSupply)

	return msgId
}

///////////////////////////////////////////
/////////////   Pod API   //////////////
//////////////////////////////////////////

// AgentをセットするDemand
func (s *SimAPI) CreatePodRequest(senderId uint64, targets []uint64) ([]*Supply, error) {

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

	sps, err := s.SendSyncDemand(s.MyClients.AgentClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.AgentClient, simSupply)

	return msgId
}

// AgentをセットするDemand
func (s *SimAPI) DeletePodRequest(senderId uint64, targets []uint64) ([]*Supply, error) {

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

	sps, err := s.SendSyncDemand(s.MyClients.AgentClient, simDemand)

	return sps, err
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

	s.SendSyncSupply(s.MyClients.AgentClient, simSupply)

	return msgId
}

///////////////////////////////////////////
/////////////      Wait      //////////////
//////////////////////////////////////////

type Waiter struct {
	WaitSpChMap map[uint64]chan *Supply
	SpMap       map[uint64][]*Supply
}

func NewWaiter() *Waiter {
	w := &Waiter{
		WaitSpChMap: make(map[uint64]chan *Supply),
		SpMap:       make(map[uint64][]*Supply),
	}
	return w
}

func (w *Waiter) WaitSp(msgId uint64, targets []uint64, timeout uint64) ([]*Supply, error) {

	waitCh := w.WaitSpChMap[msgId]

	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case sp, _ := <-waitCh:
				mu.Lock()
				//log.Printf("getSP %v, %v\n", sp.GetSimSupply().GetSenderId(), sp.GetSimSupply().GetMsgId())
				// spのidがidListに入っているか
				if sp.GetSimSupply().GetMsgId() == msgId {
					//mu.Lock()
					w.SpMap[sp.GetSimSupply().GetMsgId()] = append(w.SpMap[sp.GetSimSupply().GetMsgId()], sp)
					//mu.Unlock()
					//log.Printf("msgID spId %v, msgId %v targets %v\n", w.SpMap[sp.GetSimSupply().GetMsgId()], msgId, targets)

					// 同期が終了したかどうか
					if w.isFinishSpSync(msgId, targets) {
						logger.Debug("Finish Wait!")
						mu.Unlock()
						wg.Done()
						return
					}
				}
				mu.Unlock()
			case <-time.After(time.Duration(timeout) * time.Millisecond):
				noIds := []uint64{}
				noSps := []*Supply{} // test
				var sp2 *Supply
				for _, tgt := range targets {
					isExist := false
					for _, sp := range w.SpMap[msgId] {
						sp2 = sp
						if tgt == sp.GetSimSupply().GetSenderId() {
							isExist = true
						}
					}
					if isExist == false {
						noIds = append(noIds, tgt)
						noSps = append(noSps, sp2)
					}
				}
				logger.Error("Sync Error... noids %v, msgId %v \n%v\n\n", noIds, msgId, noSps)
				err = fmt.Errorf("Timeout Error")
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	return w.SpMap[msgId], err
}

func (w *Waiter) SendSpToWait(sp *Supply) {
	//log.Printf("getSP2 %v, %v\n", sp.GetSimSupply().GetSenderId(), sp.GetSimSupply().GetMsgId())
	mu.Lock()
	waitCh := w.WaitSpChMap[sp.GetSimSupply().GetMsgId()]
	mu.Unlock()
	waitCh <- sp
}

func (w *Waiter) isFinishSpSync(msgId uint64, targets []uint64) bool {

	for _, pid := range targets {
		isExist := false
		for _, sp := range w.SpMap[msgId] {
			senderId := sp.GetSimSupply().GetSenderId()
			if senderId == pid {
				isExist = true
			}
		}
		if isExist == false {
			return false
		}
	}

	return true
}
