package api

import (
	"context"
	"log"
	"sync"
	"time"
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

func (s *SimAPI) RegistClients(client SynerexClient, argJson string) {

	agentClient := NewSMServiceClient(client, ChannelType_AGENT_SERVICE, argJson)
	clockClient := NewSMServiceClient(client, ChannelType_CLOCK_SERVICE, argJson)
	providerClient := NewSMServiceClient(client, ChannelType_PROVIDER_SERVICE, argJson)

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
	log.Printf("SMarket Server Closed?")
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
	id := sclient.ProposeSupply(opts)
	mu.Unlock()
	return id
}

///////////////////////////////////////////
/////////////   Agent API   //////////////
//////////////////////////////////////////

// AgentをセットするDemand
func (s *SimAPI) SetAgentRequest(senderId uint64, agents []*Agent) uint64 {
	setAgentRequest := &SetAgentRequest{
		Agents: agents,
	}

	simDemand := &SimDemand{
		SenderId: senderId,
		Type:     DemandType_SET_AGENT_REQUEST,
		Data:     &SimDemand_SetAgentRequest{setAgentRequest},
	}

	id := sendDemand(s.MyClients.AgentClient, simDemand)

	return id
}

// Agentのセット完了
func (s *SimAPI) SetAgentResponse(senderId uint64) uint64 {
	setAgentResponse := &SetAgentResponse{}

	simSupply := &SimSupply{
		SenderId: senderId,
		Type:     SupplyType_SET_AGENT_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_SetAgentResponse{setAgentResponse},
	}

	id := sendSupply(s.MyClients.AgentClient, simSupply)

	return id
}

///////////////////////////////////////////
/////////////   Provider API   //////////////
//////////////////////////////////////////

// Providerを登録するDemand
func (s *SimAPI) RegistProviderRequest(senderId uint64, providerInfo *Provider) uint64 {
	registProviderRequest := &RegistProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &SimDemand{
		SenderId: senderId,
		Type:     DemandType_REGIST_PROVIDER_REQUEST,
		Data:     &SimDemand_RegistProviderRequest{registProviderRequest},
	}

	id := sendDemand(s.MyClients.ProviderClient, simDemand)

	return id
}

// Providerを登録するSupply
func (s *SimAPI) RegistProviderResponse(senderId uint64) uint64 {
	registProviderResponse := &RegistProviderResponse{}

	simSupply := &SimSupply{
		SenderId: senderId,
		Type:     SupplyType_REGIST_PROVIDER_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_RegistProviderResponse{registProviderResponse},
	}

	id := sendSupply(s.MyClients.ProviderClient, simSupply)

	return id
}

///////////////////////////////////////////
/////////////   Clock API   //////////////
//////////////////////////////////////////

func (s *SimAPI) SetClockRequest(senderId uint64, clockInfo *Clock) uint64 {
	setClockRequest := &SetClockRequest{
		Clock: clockInfo,
	}

	simDemand := &SimDemand{
		SenderId: senderId,
		Type:     DemandType_SET_CLOCK_REQUEST,
		Data:     &SimDemand_SetClockRequest{setClockRequest},
	}

	id := sendDemand(s.MyClients.ClockClient, simDemand)

	return id
}

// Agentを取得するSupply
func (s *SimAPI) SetClockResponse(senderId uint64) uint64 {
	setClockResponse := &SetClockResponse{}

	simSupply := &SimSupply{
		SenderId: senderId,
		Type:     SupplyType_SET_CLOCK_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_SetClockResponse{setClockResponse},
	}

	id := sendSupply(s.MyClients.ClockClient, simSupply)

	return id
}

func (s *SimAPI) ForwardClockRequest(senderId uint64) uint64 {
	forwardClockRequest := &ForwardClockRequest{}

	simDemand := &SimDemand{
		SenderId: senderId,
		Type:     DemandType_FORWARD_CLOCK_REQUEST,
		Data:     &SimDemand_ForwardClockRequest{forwardClockRequest},
	}

	id := sendDemand(s.MyClients.ClockClient, simDemand)

	return id
}

// Agentを取得するSupply
func (s *SimAPI) ForwardClockResponse(senderId uint64) uint64 {
	forwardClockResponse := &ForwardClockResponse{}

	simSupply := &SimSupply{
		SenderId: senderId,
		Type:     SupplyType_FORWARD_CLOCK_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_ForwardClockResponse{forwardClockResponse},
	}

	id := sendSupply(s.MyClients.ClockClient, simSupply)

	return id
}

func (s *SimAPI) StartClockRequest(senderId uint64) uint64 {
	startClockRequest := &StartClockRequest{}

	simDemand := &SimDemand{
		SenderId: senderId,
		Type:     DemandType_START_CLOCK_REQUEST,
		Data:     &SimDemand_StartClockRequest{startClockRequest},
	}

	id := sendDemand(s.MyClients.ClockClient, simDemand)

	return id
}

// Agentを取得するSupply
func (s *SimAPI) StartClockResponse(senderId uint64) uint64 {
	startClockResponse := &StartClockResponse{}

	simSupply := &SimSupply{
		SenderId: senderId,
		Type:     SupplyType_START_CLOCK_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_StartClockResponse{startClockResponse},
	}

	id := sendSupply(s.MyClients.ClockClient, simSupply)

	return id
}

func (s *SimAPI) StopClockRequest(senderId uint64) uint64 {
	stopClockRequest := &StopClockRequest{}

	simDemand := &SimDemand{
		SenderId: senderId,
		Type:     DemandType_STOP_CLOCK_REQUEST,
		Data:     &SimDemand_StopClockRequest{stopClockRequest},
	}

	id := sendDemand(s.MyClients.ClockClient, simDemand)

	return id
}

// Agentを取得するSupply
func (s *SimAPI) StopClockResponse(senderId uint64) uint64 {
	stopClockResponse := &StopClockResponse{}

	simSupply := &SimSupply{
		SenderId: senderId,
		Type:     SupplyType_STOP_CLOCK_RESPONSE,
		Status:   StatusType_OK,
		Data:     &SimSupply_StopClockResponse{stopClockResponse},
	}

	id := sendSupply(s.MyClients.ClockClient, simSupply)

	return id
}

///////////////////////////////////////////
/////////////      Wait      //////////////
//////////////////////////////////////////

// SendToSetAgentsResponse : SetAgentsResponseを送る
func (s *SimAPI) SendToWaitCh(sp *Supply, supplyType SupplyType) {
	mu.Lock()
	waitCh := waitChMap[supplyType]
	mu.Unlock()
	waitCh <- sp
}

// Wait: 同期が完了するまで待機する関数
func wait(idList []uint64, supplyType SupplyType) map[uint64]*Supply {

	mu.Lock()
	waitCh := make(chan *Supply, CHANNEL_BUFFER_SIZE)
	waitChMap[supplyType] = waitCh
	mu.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(1)
	pspMap := make(map[uint64]*Supply)
	go func() {
		for {
			select {
			case psp, _ := <-waitCh:
				mu.Lock()
				// spのidがidListに入っているか
				if isPidInIdList(psp, idList) {
					//logger.Debug("isPidInIDList %v, %v", psp.GetSimSupply().GetSenderId(), idList)
					pspMap[psp.GetSimSupply().GetSenderId()] = psp
					//logger.Debug("isFinishSync %v, %v", isFinishSync(pspMap, idList), idList)
					//for _, sp := range pspMap {
					//	logger.Debug("pspMap %v", sp.GetSimSupply().GetSenderId(), idList)
					//}
					// 同期が終了したかどうか
					if isFinishSync(pspMap, idList) {
						//logger.Debug("isFinishSync")
						mu.Unlock()
						wg.Done()
						return
					}
				}
				mu.Unlock()
			case <-time.After(1500 * time.Millisecond):
				noIds := make([]uint64, 0)
				for _, id := range idList {
					noFlag := true
					for _, sp := range pspMap {
						if sp.GetSimSupply().GetSenderId() == id {
							noFlag = false
						}

					}
					if noFlag {
						noIds = append(noIds, id)
					}
				}

				//logger.Error("Sync Error: NoIds: %v", noIds)
			}
		}
	}()
	wg.Wait()
	return pspMap
}

// isSpInIdList : spのidがidListに入っているか
func isPidInIdList(sp *Supply, idlist []uint64) bool {
	pid := sp.GetSimSupply().GetSenderId()
	for _, id := range idlist {
		if pid == id {
			return true
		}
	}
	return false
}

/*// isFinishSync : 必要な全てのSupplyを受け取り同期が完了したかどうか
func isFinishSync2(spList []*Supply, idlist []uint64) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range spList {
			pid := sp.GetSimSupply().GetSenderId()
			if id == pid {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}*/

// isFinishSync : 必要な全てのSupplyを受け取り同期が完了したかどうか
func isFinishSync(pspMap map[uint64]*Supply, idlist []uint64) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			pid := sp.GetSimSupply().GetSenderId()
			if id == pid {
				isMatch = true
			}
		}
		if isMatch == false {
			return false
		}
	}
	return true
}
