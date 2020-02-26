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
	mu        sync.Mutex
	waitChMap map[simapi.SupplyType]chan *pb.Supply
	//spMesMap            map[simapi.SupplyType]*Message
	CHANNEL_BUFFER_SIZE int
	logger              *Logger
)

func init() {
	waitChMap = make(map[simapi.SupplyType]chan *pb.Supply)
	//spMesMap = make(map[simapi.SupplyType]*Message)
	CHANNEL_BUFFER_SIZE = 10
	logger = NewLogger()
}

type Clients struct {
	AgentClient    *sxutil.SMServiceClient
	ClockClient    *sxutil.SMServiceClient
	ProviderClient *sxutil.SMServiceClient
}

type Communicator struct {
	MyClients *Clients
}

func NewCommunicator() *Communicator {
	c := &Communicator{}
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

	time.Sleep(1 * time.Second)
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

/*////////////////////////////////////////////////////////////
////////////            Message Class           ///////////
///////////////////////////////////////////////////////////

type Message struct {
	Ready  chan struct{}
	SpCh   chan *pb.Supply
	SpList []*pb.Supply
}

func NewMessage() *Message {
	return &Message{Ready: make(chan struct{}), SpCh: make(chan *pb.Supply)}
}

func (m *Message) Add(sp *pb.Supply) {
	m.SpList = append(m.SpList, sp)
}

func (m *Message) Get(idList []uint64) []*pb.Supply {
	go func() {
		for {
			sp := <-m.SpCh
			if isPidInIdList(sp, idList) {
				m.SpList = append(m.SpList, sp)
				if isFinishSync2(m.SpList, idList) {
					close(m.Ready)
					return
				}
			}
		}
	}()
	<-m.Ready
	return m.SpList
}

////////////////////////////////////////////////////////////
////////////        Wait Function       ///////////////////
///////////////////////////////////////////////////////////

// SendToSetAgentsResponse : SetAgentsResponseを送る
func (c *Communicator) SendToWaitCh(sp *pb.Supply, supplyType simapi.SupplyType) {
	mu.Lock()
	//waitCh := waitChMap[supplyType]
	spMes := spMesMap[supplyType]
	if spMes == nil {
		spMes = NewMessage()
	}
	spMes.Add(sp)
	spMesMap[supplyType] = spMes
	mu.Unlock()
}

// Wait: 同期が完了するまで待機する関数
func wait(idList []uint64, supplyType simapi.SupplyType) []*pb.Supply {
	mu.Lock()
	spMes := spMesMap[supplyType]
	if spMes == nil {
		spMes = NewMessage()
	}
	mu.Unlock()
	var wg sync.WaitGroup
	wg.Add(1)
	var spList []*pb.Supply
	go func() {
		spList = spMes.Get(idList)
		wg.Done()
	}()
	wg.Wait()
	return spList
}*/

// SendToSetAgentsResponse : SetAgentsResponseを送る
func (c *Communicator) SendToWaitCh(sp *pb.Supply, supplyType simapi.SupplyType) {
	mu.Lock()
	waitCh := waitChMap[supplyType]
	mu.Unlock()
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
			case psp, _ := <-waitCh:
				mu.Lock()
				// spのidがidListに入っているか
				if isPidInIdList(psp, idList) {
					//logger.Debug("isPidInIDList %v, %v", psp.GetSimSupply().GetPid(), idList)
					pspMap[psp.GetSimSupply().GetPid()] = psp
					//logger.Debug("isFinishSync %v, %v", isFinishSync(pspMap, idList), idList)
					//for _, sp := range pspMap {
					//	logger.Debug("pspMap %v", sp.GetSimSupply().GetPid(), idList)
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
						if sp.GetSimSupply().GetPid() == id {
							noFlag = false
						}

					}
					if noFlag {
						noIds = append(noIds, id)
					}
				}
				for _, sp := range pspMap {

					logger.Error("pspMap %v", sp.GetSimSupply().GetPid(), idList)
				}
				logger.Error("Sync Error: NoIds: %v", noIds)
			}
		}
	}()
	wg.Wait()
	return pspMap
}

// isSpInIdList : spのidがidListに入っているか
func isPidInIdList(sp *pb.Supply, idlist []uint64) bool {
	pid := sp.GetSimSupply().GetPid()
	for _, id := range idlist {
		if pid == id {
			return true
		}
	}
	return false
}

/*// isFinishSync : 必要な全てのSupplyを受け取り同期が完了したかどうか
func isFinishSync2(spList []*pb.Supply, idlist []uint64) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range spList {
			pid := sp.GetSimSupply().GetPid()
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
func isFinishSync(pspMap map[uint64]*pb.Supply, idlist []uint64) bool {
	for _, id := range idlist {
		isMatch := false
		for _, sp := range pspMap {
			pid := sp.GetSimSupply().GetPid()
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

///////////////////////////////////////////
/////////////   Agent API   //////////////
//////////////////////////////////////////

// Agentを取得するDemand
func (c *Communicator) GetAgentsRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64) (uint64, []*agent.Agent) {
	getAgentsRequest := &agent.GetAgentsRequest{}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_GET_AGENTS_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_GetAgentsRequest{getAgentsRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) GetAgentsResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64, agents []*agent.Agent, agentType agent.AgentType, areaId uint64) uint64 {
	getAgentsResponse := &agent.GetAgentsResponse{
		Agents:    agents,
		AgentType: agentType,
		AreaId:    areaId,
	}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_GET_AGENTS_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_GetAgentsResponse{getAgentsResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.AgentClient, tid, simSupply)

	return id
}

// AgentをセットするDemand
func (c *Communicator) SetAgentsRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64, agents []*agent.Agent) uint64 {
	setAgentsRequest := &agent.SetAgentsRequest{
		Agents: agents,
	}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_SET_AGENTS_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_SetAgentsRequest{setAgentsRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendDemand(c.MyClients.AgentClient, simDemand)

	if idList != nil {
		supplyType := simapi.SupplyType_SET_AGENTS_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentのセット完了
func (c *Communicator) SetAgentsResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	setAgentsResponse := &agent.SetAgentsResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_SET_AGENTS_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_SetAgentsResponse{setAgentsResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.AgentClient, tid, simSupply)

	return id
}

///////////////////////////////////////////
/////////////   Provider API   //////////////
//////////////////////////////////////////

// Providerを登録するDemand
func (c *Communicator) RegistProviderRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64, providerInfo *provider.Provider) uint64 {
	registProviderRequest := &provider.RegistProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_REGIST_PROVIDER_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_RegistProviderRequest{registProviderRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) RegistProviderResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	registProviderResponse := &provider.RegistProviderResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_REGIST_PROVIDER_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_RegistProviderResponse{registProviderResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) KillProviderRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64, providerInfo *provider.Provider) uint64 {
	killProviderRequest := &provider.KillProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_KILL_PROVIDER_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_KillProviderRequest{killProviderRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) KillProviderResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	killProviderResponse := &provider.KillProviderResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_KILL_PROVIDER_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_KillProviderResponse{killProviderResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) DivideProviderRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64, providerInfo *provider.Provider) uint64 {
	divideProviderRequest := &provider.DivideProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_DIVIDE_PROVIDER_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_DivideProviderRequest{divideProviderRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) DivideProviderResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	divideProviderResponse := &provider.DivideProviderResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_DIVIDE_PROVIDER_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_DivideProviderResponse{divideProviderResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) UpdateProvidersRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64, providers []*provider.Provider) uint64 {
	updateProvidersRequest := &provider.UpdateProvidersRequest{
		Providers: providers,
	}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_UPDATE_PROVIDERS_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_UpdateProvidersRequest{updateProvidersRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) UpdateProvidersResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	updateProvidersResponse := &provider.UpdateProvidersResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_UPDATE_PROVIDERS_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_UpdateProvidersResponse{updateProvidersResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) SendProviderStatusRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64, providerInfo *provider.Provider) uint64 {
	sendProviderStatusRequest := &provider.SendProviderStatusRequest{
		Provider: providerInfo,
	}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_SEND_PROVIDER_STATUS_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_SendProviderStatusRequest{sendProviderStatusRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) SendProviderStatusResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	sendProviderStatusResponse := &provider.SendProviderStatusResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_SEND_PROVIDER_STATUS_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_SendProviderStatusResponse{sendProviderStatusResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) SetProvidersRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64, providers []*provider.Provider) uint64 {
	setProvidersRequest := &provider.SetProvidersRequest{
		Providers: providers,
	}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_SET_PROVIDERS_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_SetProvidersRequest{setProvidersRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendDemand(c.MyClients.ProviderClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_SET_PROVIDERS_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator) SetProvidersResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	setProvidersResponse := &provider.SetProvidersResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_SET_PROVIDERS_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_SetProvidersResponse{setProvidersResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator) GetProvidersRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64) uint64 {
	getProvidersRequest := &provider.GetProvidersRequest{}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_GET_PROVIDERS_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_GetProvidersRequest{getProvidersRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendDemand(c.MyClients.ProviderClient, simDemand)

	// Wait
	if idList != nil {
		supplyType := simapi.SupplyType_GET_PROVIDERS_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator) GetProvidersResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64, providers []*provider.Provider) uint64 {
	getProvidersResponse := &provider.GetProvidersResponse{
		Providers: providers,
	}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_GET_PROVIDERS_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_GetProvidersResponse{getProvidersResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ProviderClient, tid, simSupply)

	return id
}

///////////////////////////////////////////
/////////////   Clock API   //////////////
//////////////////////////////////////////

func (c *Communicator) UpdateClockRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64, clockInfo *clock.Clock) uint64 {
	updateClockRequest := &clock.UpdateClockRequest{
		Clock: clockInfo,
	}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_UPDATE_CLOCK_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_UpdateClockRequest{updateClockRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) UpdateClockResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	updateClockResponse := &clock.UpdateClockResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_UPDATE_CLOCK_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_UpdateClockResponse{updateClockResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) SetClockRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64, clockInfo *clock.Clock) uint64 {
	setClockRequest := &clock.SetClockRequest{
		Clock: clockInfo,
	}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_SET_CLOCK_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_SetClockRequest{setClockRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) SetClockResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	setClockResponse := &clock.SetClockResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_SET_CLOCK_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_SetClockResponse{setClockResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) GetClockRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64) (uint64, *clock.Clock) {
	getClockRequest := &clock.GetClockRequest{}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_GET_CLOCK_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_GetClockRequest{getClockRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendDemand(c.MyClients.ClockClient, simDemand)

	// Wait
	var clockInfo *clock.Clock
	if idList != nil {
		supplyType := simapi.SupplyType_GET_CLOCK_RESPONSE
		spMap := wait(idList, supplyType)
		//var clockInfo *clock.Clock
		for _, sp := range spMap {
			clockInfo = sp.GetSimSupply().GetGetClockResponse().GetClock()
		}
		logger.Info("Clock Info: %v", clockInfo)
	}

	return id, clockInfo
}

// Agentを取得するSupply
func (c *Communicator) GetClockResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64, clockInfo *clock.Clock) uint64 {
	getClockResponse := &clock.GetClockResponse{
		Clock: clockInfo,
	}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_GET_CLOCK_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_GetClockResponse{getClockResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) ForwardClockRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64) uint64 {
	forwardClockRequest := &clock.ForwardClockRequest{}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_FORWARD_CLOCK_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_ForwardClockRequest{forwardClockRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) ForwardClockResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	forwardClockResponse := &clock.ForwardClockResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_FORWARD_CLOCK_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_ForwardClockResponse{forwardClockResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) BackClockRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64) uint64 {
	backClockRequest := &clock.BackClockRequest{}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_BACK_CLOCK_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_BackClockRequest{backClockRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) backClockResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	backClockResponse := &clock.BackClockResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_BACK_CLOCK_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_BackClockResponse{backClockResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) StartClockRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64) uint64 {
	startClockRequest := &clock.StartClockRequest{}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_START_CLOCK_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_StartClockRequest{startClockRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) StartClockResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	startClockResponse := &clock.StartClockResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_START_CLOCK_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_StartClockResponse{startClockResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator) StopClockRequest(senderInfo *provider.Provider, targets []uint64, pid uint64, idList []uint64) uint64 {
	stopClockRequest := &clock.StopClockRequest{}

	simDemand := &simapi.SimDemand{
		Pid:        pid,
		Type:       simapi.DemandType_STOP_CLOCK_REQUEST,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimDemand_StopClockRequest{stopClockRequest},
		SenderInfo: senderInfo,
		Targets:    targets,
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
func (c *Communicator) StopClockResponse(senderInfo *provider.Provider, targets []uint64, pid uint64, tid uint64) uint64 {
	stopClockResponse := &clock.StopClockResponse{}

	simSupply := &simapi.SimSupply{
		Pid:        pid,
		Type:       simapi.SupplyType_STOP_CLOCK_RESPONSE,
		Status:     simapi.StatusType_NONE,
		Data:       &simapi.SimSupply_StopClockResponse{stopClockResponse},
		SenderInfo: senderInfo,
		Targets:    targets,
	}

	id := sendSupply(c.MyClients.ClockClient, tid, simSupply)

	return id
}
