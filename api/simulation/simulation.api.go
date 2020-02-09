package simulation

/*import (
	"github.com/synerex/synerex_alpha/sxutil"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/provider"
)

var (
	mu                  sync.Mutex
	waitChMap map[simapi.SupplyType]chan *pb.Supply
)

func init(){
	waitChMap = make(map[simapi.SupplyType]chan *pb.Supply)
}

type Clients struct {
	AgentClient       *sxutil.SMServiceClient
	ClockClient       *sxutil.SMServiceClient
	AreaClient        *sxutil.SMServiceClient
	ProviderClient       *sxutil.SMServiceClient
}

type Communicator struct{
	MyClients *Clients
	MyProvider *provider.Provider
}

func NewCommunicator(providerInfo *provider.Provider)*Communicator{
	c := &Communicator{
		MyProvider: providerInfo,
	}
}

func (c *Communicator) RegistClients(client pb.SynerexClient, argJson string) {

	agentClient := sxutil.NewSMServiceClient(client, pb.ChannelType_AGENT_SERVICE, argJson)
	clockClient := sxutil.NewSMServiceClient(client, pb.ChannelType_CLOCK_SERVICE, argJson)
	providerClient := sxutil.NewSMServiceClient(client, pb.ChannelType_PROVIDER_SERVICE, argJson)

	clients := &Clients{
		AgentClient:       agentClient,
		ClockClient:       clockClient,
		ProviderClient: providerClient,
	}

	c.MyClients = clients
}

// SubscribeAll: 全てのチャネルに登録、SubscribeSupply, SubscribeDemandする
func (c *Communicator) SubscribeAll(demandCallback func(*sxutil.SMServiceClient, *pb.Demand), supplyCallback func(*sxutil.SMServiceClient, *pb.Supply)) error{

	// SubscribeDemand, SubscribeSupply
	go subscribeDemand(c.MyClients.AgentClient, demandCallback)

	go subscribeDemand(c.MyClients.ClockClient, demandCallback)

	go subscribeDemand(c.MyClients.AreaClient, demandCallback)

	go subscribeDemand(c.MyClients.ParticipantClient, demandCallback)

	go subscribeDemand(c.MyClients.RouteClient, demandCallback)

	go subscribeSupply(c.MyClients.ClockClient, supplyCallback)

	go subscribeSupply(c.MyClients.AreaClient, supplyCallback)

	go subscribeSupply(c.MyClients.AgentClient, supplyCallback)

	go subscribeSupply(c.MyClients.ParticipantClient, supplyCallback)

	go subscribeSupply(c.MyClients.RouteClient, supplyCallback))

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

func sendDemand(sclient *sxutil.SMServiceClient, simDemand *simapi.SimDemand) uint64{
	nm := ""
	js := ""
	opts := &sxutil.DemandOpts{Name: nm, JSON: js, SimDemand: simDemand}

	mu.Lock()
	id := sclient.RegisterDemand(opts)
	mu.Unlock()
	return id
}

func sendSupply(sclient *sxutil.SMServiceClient, tid uint64, simSupply *simapi.SimSupply) uint64{
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
				if isSpInIdList(psp, idList){
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
		if senderId == id{
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
func (c *Communicator)GetAgentsRequest(idList []uint64) (uint64, []*agent.Agent){
	getAgentsRequest := &agent.GetAgentsRequest{}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_GET_AGENTS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_GetAgentsRequest{getAgentsRequest},
	}

	id := sendDemand(c.Clients.AgentClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_GET_AGENTS_RESPONSE
		spMap := wait(idList, supplyType)
		agents := make([]*agent.Agent, 0)
		for _, sp := range spMap {
			ags = sp.GetSimSupply().GetGetNeighborAreaAgentsResponse().GetAgents()
			agents = append(agents, ags...)
		}
	}

	return id, agents
}

// Agentを取得するSupply
func (c *Communicator)GetAgentsResponse(tid uint64, agents []*agent.Agent, agentType agent.AgentType, areaId uint64) uint64{
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

	id := sendSupply(c.Clients.AgentClient, tid, simSupply)

	return id
}

// AgentをセットするDemand
func (c *Communicator)SetAgentsRequest(idList []uint64, agents []*agent.Agent) uint64{
	setAgentsRequest := &agent.SetAgentsRequest{
		Agents: agents,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_SET_AGENTS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_SetAgentsRequest{setAgentsRequest},
	}

	id := sendDemand(c.Clients.AgentClient, simDemand)

	if idList != nil{
		supplyType := synerex.SupplyType_SET_AGENTS_RESPONSE
		wait(idList, supplyType)
	}

	return id,
}

// Agentのセット完了
func (c *Communicator)SetAgentsResponse(tid uint64, agents []*agent.Agent, agentType agent.AgentType, areaId uint64) uint64{
	setAgentsResponse := &agent.SetAgentsResponse{
		Agents: agents,
		AgentType: agentType,
		AreaId: areaId,
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_SET_AGENTS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_SetAgentsResponse{setAgentsResponse},
	}

	id := sendSupply(c.Clients.AgentClient, tid, simSupply)

	return id
}


///////////////////////////////////////////
/////////////   Provider API   //////////////
//////////////////////////////////////////

// Providerを登録するDemand
func (c *Communicator)RegistProviderRequest(idList []uint64, providerInfo *provider.Provider) uint64{
	registProviderRequest := &provider.RegistProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_REGIST_PROVIDER_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_RegistProviderRequest{registProviderRequest},
	}

	id := sendDemand(c.Clients.ProviderClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_REGIST_PROVIDER_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator)RegistProviderResponse(tid uint64) uint64{
	registProviderResponse := &provider.RegistProviderResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_REGIST_PROVIDER_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_RegistProviderResponse{registProviderResponse},
	}

	id := sendSupply(c.Clients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator)KillProviderRequest(idList []uint64, providerInfo *provider.Provider) uint64{
	killProviderRequest := &provider.KillProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_KILL_PROVIDER_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_KillProviderRequest{killProviderRequest},
	}

	id := sendDemand(c.Clients.ProviderClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_KILL_PROVIDER_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator)KillProviderResponse(tid uint64) uint64{
	killProviderResponse := &provider.KillProviderResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_KILL_PROVIDER_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_KillProviderResponse{killProviderResponse},
	}

	id := sendSupply(c.Clients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator)DivideProviderRequest(idList []uint64, providerInfo *provider.Provider) uint64{
	divideProviderRequest := &provider.DivideProviderRequest{
		Provider: providerInfo,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_DIVIDE_PROVIDER_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_DivideProviderRequest{divideProviderRequest},
	}

	id := sendDemand(c.Clients.ProviderClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_DIVIDE_PROVIDER_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator)DivideProviderResponse(tid uint64) uint64{
	divideProviderResponse := &provider.DivideProviderResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_DIVIDE_PROVIDER_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_DivideProviderResponse{divideProviderResponse},
	}

	id := sendSupply(c.Clients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator)UpdateProvidersRequest(idList []uint64, providers []*Provider) uint64{
	updateProvidersRequest := &provider.UpdateProviderRequest{
		Providers: providers,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_UPDATE_PROVIDERS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_UpdateProvidersRequest{updateProvidersRequest},
	}

	id := sendDemand(c.Clients.ProviderClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_UPDATE_PROVIDERS_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator)UpdateProvidersResponse(tid uint64) uint64{
	updateProvidersResponse := &provider.UpdateProvidersResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_UPDATE_PROVIDERS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_UpdateProvidersResponse{updateProvidersResponse},
	}

	id := sendSupply(c.Clients.ProviderClient, tid, simSupply)

	return id
}

// Providerを登録するDemand
func (c *Communicator)SendProviderStatusRequest(idList []uint64, providerInfo *provider.Provider) uint64{
	sendProviderStatusRequest := &provider.SendProviderStatusRequest{
		Provider: providerInfo,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_SEND_PROVIDER_STATUS_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_SendProviderStatusRequest{sendProviderStatusRequest},
	}

	id := sendDemand(c.Clients.ProviderClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_SEND_PROVIDER_STATUS_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Providerを登録するSupply
func (c *Communicator)SendProviderStatusResponse(tid uint64) uint64{
	sendProviderStatusResponse := &provider.SendProviderStatusResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_SEND_PROVIDER_STATUS_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_SendProviderStatusResponse{sendProviderStatusResponse},
	}

	id := sendSupply(c.Clients.ProviderClient, tid, simSupply)

	return id
}


///////////////////////////////////////////
/////////////   Clock API   //////////////
//////////////////////////////////////////

func (c *Communicator)UpdateClockRequest(idList []uint64, clockInfo *clock.Clock) uint64{
	updateClockRequest := &clock.UpdateClockRequest{
		Clock: clockInfo,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_UPDATE_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_UpdateClockRequest{updateClockRequest},
	}

	id := sendDemand(c.Clients.ClockClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_UPDATE_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator)UpdateClockResponse(tid uint64) uint64{
	updateClockResponse := &clock.UpdateClockResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_UPDATE_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_UpdateClockResponse{UpdateClockResponse},
	}

	id := sendSupply(c.Clients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator)SetClockRequest(idList []uint64, clockInfo *clock.Clock) uint64{
	setClockRequest := &clock.SetClockRequest{
		Clock: clockInfo,
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_SET_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_SetClockRequest{setClockRequest},
	}

	id := sendDemand(c.Clients.ClockClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_SET_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator)SetClockResponse(tid uint64) uint64{
	setClockResponse := &clock.SetClockResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_SET_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_SetClockResponse{setClockResponse},
	}

	id := sendSupply(c.Clients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator)GetClockRequest(idList []uint64) (uint64, *clock.Clock){
	getClockRequest := &clock.GetClockRequest{
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_GET_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_GetClockRequest{getClockRequest},
	}

	id := sendDemand(c.Clients.ClockClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_GET_CLOCK_RESPONSE
		spMap := wait(idList, supplyType)
		clockInfo := sp.GetSimSupply().GetGetClockResponse().GetClock()
	}

	return id, clockInfo
}

// Agentを取得するSupply
func (c *Communicator)GetClockResponse(tid uint64, clockInfo *clock.Clock) uint64{
	getClockResponse := &clock.GetClockResponse{
		Clock: clockInfo,
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_GET_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_GetClockResponse{getClockResponse},
	}

	id := sendSupply(c.Clients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator)ForwardClockRequest(idList []uint64) uint64{
	forwardClockRequest := &clock.ForwardClockRequest{
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_FORWARD_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_ForwardClockRequest{forwardClockRequest},
	}

	id := sendDemand(c.Clients.ClockClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_FORWARD_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator)ForwardClockResponse(tid uint64) uint64{
	forwardClockResponse := &clock.ForwardClockResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_FORWARD_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_ForwardClockResponse{forwardClockResponse},
	}

	id := sendSupply(c.Clients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator)BackClockRequest(idList []uint64) uint64{
	backClockRequest := &clock.BackClockRequest{
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_BACK_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_BackClockRequest{backClockRequest},
	}

	id := sendDemand(c.Clients.ClockClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_BACK_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator)backClockResponse(tid uint64) uint64{
	BackClockResponse := &clock.BackClockResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_BACK_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_BackClockResponse{backClockResponse},
	}

	id := sendSupply(c.Clients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator)StartClockRequest(idList []uint64) uint64{
	startClockRequest := &clock.StartClockRequest{
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_START_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_StartClockRequest{startClockRequest},
	}

	id := sendDemand(c.Clients.ClockClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_START_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator)StartClockResponse(tid uint64) uint64{
	startClockResponse := &clock.StartClockResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_START_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_StartClockResponse{startClockResponse},
	}

	id := sendSupply(c.Clients.ClockClient, tid, simSupply)

	return id
}

func (c *Communicator)StopClockRequest(idList []uint64) uint64{
	stopClockRequest := &clock.StopClockRequest{
	}

	simDemand := &synerex.SimDemand{
		DemandType: synerex.DemandType_STOP_CLOCK_REQUEST,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimDemand_StopClockRequest{stopClockRequest},
	}

	id := sendDemand(c.Clients.ClockClient, simDemand)

	// Wait
	if idList != nil{
		supplyType := synerex.SupplyType_STOP_CLOCK_RESPONSE
		wait(idList, supplyType)
	}

	return id
}

// Agentを取得するSupply
func (c *Communicator)StopClockResponse(tid uint64) uint64{
	stopClockResponse := &clock.StopClockResponse{
	}

	simSupply := &synerex.SimSupply{
		SupplyType: synerex.SupplyType_STOP_CLOCK_RESPONSE,
		StatusType: synerex.StatusType_NONE,
		Data:       &synerex.SimSupply_StopClockResponse{stopClockResponse},
	}

	id := sendSupply(c.Clients.ClockClient, tid, simSupply)

	return id
}*/
