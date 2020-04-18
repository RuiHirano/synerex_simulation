package main

// main synerex serverからgatewayを介してother synerex serverへ情報を送る
// 基本的に一方通行

import (
	//"flag"
	"fmt"
	"log"
	"os"

	//"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	//"github.com/golang/protobuf/jsonpb"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"google.golang.org/grpc"
)

var (
	workerSynerexAddr1 string
	workerSynerexAddr2 string
	workerNodeIdAddr1  string
	workerNodeIdAddr2  string
	//pm1                *simutil.ProviderManager
	//pm2                *simutil.ProviderManager
	apm        *AgentProviderManager
	waiter     *api.Waiter
	mu         sync.Mutex
	myProvider *api.Provider
	worker1api *api.SimAPI
	worker2api *api.SimAPI
	//scenarioProvider   *provider.Provider
	//com1               *simutil.Communicator
	//com2               *simutil.Communicator
	//providerManager1   *simutil.ProviderManager
	//providerManager2   *simutil.ProviderManager
	logger *simutil.Logger
	//mes1               *Message
	//mes2               *Message
)

func init() {
	//flag.Parse()
	logger = simutil.NewLogger()
	waiter = api.NewWaiter()
	apm = NewAgentProviderManager()
	//myProvider = flagToProviderInfo(*providerJson)
	//scenarioProvider = flagToProviderInfo(*scenarioProviderJson)

	workerSynerexAddr1 = os.Getenv("WORKER_SYNEREX_SERVER1")
	if workerSynerexAddr1 == "" {
		workerSynerexAddr1 = "127.0.0.1:10000"
	}

	workerSynerexAddr2 = os.Getenv("WORKER_SYNEREX_SERVER2")
	if workerSynerexAddr2 == "" {
		workerSynerexAddr2 = "127.0.0.1:10000"
	}

	workerNodeIdAddr1 = os.Getenv("WORKER_NODEID_SERVER1")
	if workerNodeIdAddr1 == "" {
		workerNodeIdAddr1 = "127.0.0.1:9000"
	}

	workerNodeIdAddr2 = os.Getenv("WORKER_NODEID_SERVER2")
	if workerNodeIdAddr2 == "" {
		workerNodeIdAddr2 = "127.0.0.1:9000"
	}
}

////////////////////////////////////////////////////////////
//////////         Agent Provider Manager         /////////
///////////////////////////////////////////////////////////

type AgentProviderManager struct {
	Provider1 *api.Provider
	Provider2 *api.Provider
}

func NewAgentProviderManager() *AgentProviderManager {
	apm := &AgentProviderManager{
		Provider1: nil,
		Provider2: nil,
	}
	return apm
}

func (apm *AgentProviderManager) SetProvider1(ps []*api.Provider) {
	mu.Lock()
	for _, p := range ps {
		if p.GetType() == api.ProviderType_AGENT {
			apm.Provider1 = p
		}
	}
	mu.Unlock()
}

func (apm *AgentProviderManager) SetProvider2(ps []*api.Provider) {
	mu.Lock()
	for _, p := range ps {
		if p.GetType() == api.ProviderType_AGENT {
			apm.Provider2 = p
		}
	}
	mu.Unlock()
}

/*type AgentProviderManager struct {
	Providers1  []*api.Provider
	Providers2  []*api.Provider
	NeighborMap map[uint64][]*api.Provider // 隣接してるProviderマップ
	MsgIdMap    map[uint64]uint64          // msgIdを結びつけるためのマップ
}

func NewAgentProviderManager() *AgentProviderManager {
	apm := &AgentProviderManager{
		Providers1:  []*api.Provider{},
		Providers2:  []*api.Provider{},
		NeighborMap: make(map[uint64][]*api.Provider),
		MsgIdMap:    make(map[uint64]uint64),
	}
	return apm
}

func (apm *AgentProviderManager) SetProviders1(ps []*api.Provider) {
	mu.Lock()
	for _, p := range ps {
		if p.GetProviderType() == api.ProviderType_AGENT {
			apm.Providers1 = append(apm.Providers1, p)
		}
	}
	apm.CreateProvidersMap()
	mu.Unlock()
}

func (apm *AgentProviderManager) SetProviders2(ps []*api.Provider) {
	mu.Lock()
	apm.Providers2 = []*api.Provider{}
	for _, p := range ps {
		if p.GetProviderType() == api.ProviderType_AGENT {
			apm.Providers2 = append(apm.Providers2, p)
		}
	}
	apm.CreateProvidersMap()
	mu.Unlock()
}

func (apm *AgentProviderManager) SetMsgIdMap(msgId1 uint64, msgId2 uint64) {
	mu.Lock()
	apm.MsgIdMap[msgId1] = msgId2
	apm.MsgIdMap[msgId2] = msgId1
	mu.Unlock()
}

func (apm *AgentProviderManager) CreateProvidersMap() {
	neighborMap := make(map[uint64][]*api.Provider)
	for _, p1 := range apm.Providers1 {
		p1Id := p1.GetId()
		for _, p2 := range apm.Providers2 {
			p2Id := p2.GetId()
			//if isNeighborArea(p1, p2) {
			// エリアが隣接していた場合
			neighborMap[p1Id] = append(neighborMap[p1Id], p2)
			neighborMap[p2Id] = append(neighborMap[p2Id], p1)
			//}
		}
	}
	apm.NeighborMap = neighborMap
}*/

/*func isNeighborArea(p1 *api.Provider, p2 *api.Provider) bool {
	myControlArea := pm.MyProvider.GetAgentStatus().GetArea().GetControlArea()
	tControlArea := p.GetAgentStatus().GetArea().GetControlArea()
	maxLat, maxLon, minLat, minLon := GetCoordRange(myControlArea)
	tMaxLat, tMaxLon, tMinLat, tMinLon := GetCoordRange(tControlArea)
	if maxLat == tMinLat && (minLon <= tMaxLon && tMaxLon <= maxLon || minLon <= tMinLon && tMinLon <= maxLon) {
		return true
	}
	if minLat == tMaxLat && (minLon <= tMaxLon && tMaxLon <= maxLon || minLon <= tMinLon && tMinLon <= maxLon) {
		return true
	}
	if maxLon == tMinLon && (minLat <= tMaxLat && tMaxLat <= maxLat || minLat <= tMinLat && tMinLat <= maxLat) {
		return true
	}
	if minLon == tMaxLon && (minLat <= tMaxLat && tMaxLat <= maxLat || minLat <= tMinLat && tMinLat <= maxLat) {
		return true
	}
	return false
}*/

////////////////////////////////////////////////////////////
//////////     Worker1 Demand Supply Callback     /////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func supplyCallback1(clt *api.SMServiceClient, sp *api.Supply) {
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_READY_PROVIDER_RESPONSE:
		//time.Sleep(10 * time.Millisecond)
		worker1api.SendSpToWait(sp)
		fmt.Printf("ready provider response")

	case api.SupplyType_GET_AGENT_RESPONSE:
		//fmt.Printf("Get Sp from Worker1%v\n", sp)

		//time.Sleep(10 * time.Millisecond)
		worker1api.SendSpToWait(sp)
		/*msgId2 := sp.GetSimSupply().GetMsgId()
		// send to worker1
		msgId1 := apm.MsgIdMap[msgId2]
		senderId := myProvider.Id
		agents := sp.GetSimSupply().GetGetAgentResponse()
		worker1api.GetAgentResponse(senderId, targets, msgId1, agents)*/

	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		//masterProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		fmt.Printf("regist provider to Worler1 Provider!\n")
	}
}

// Demandのコールバック関数
func demandCallback1(clt *api.SMServiceClient, dm *api.Demand) {

	switch dm.GetSimDemand().GetType() {
	case api.DemandType_READY_PROVIDER_REQUEST:
		provider := dm.GetSimDemand().GetReadyProviderRequest().GetProvider()
		//pm.SetProviders(providers)

		// workerへ登録
		senderId := myProvider.Id
		targets := []uint64{provider.GetId()}
		worker1api.RegistProviderRequest(senderId, targets, myProvider)
		//waiter.WaitSp(msgId, targets, 1000)

		// response
		targets = []uint64{dm.GetSimDemand().GetSenderId()}
		senderId = myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		worker1api.ReadyProviderResponse(senderId, targets, msgId)
		logger.Info("Finish: Regist Provider from ready ")

	case api.DemandType_GET_AGENT_REQUEST:
		// 隣接エリアがない場合はそのまま返す
		t1 := time.Now()

		/*targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		agents := []*api.Agent{}
		worker1api.GetAgentResponse(senderId, targets, msgId, agents)
		logger.Debug("Finish: Get Agent Request Worker1 %v %v\n", targets, msgId)*/
		/*agents := []*api.Agent{}
		senderId := myProvider.Id
		// worker2のagent-providerから取得
		targets2 := []uint64{apm.Provider2.GetId()}
		msgId2 := worker2api.GetAgentRequest(senderId, targets2)
		//logger.Debug("Get Agent Request to Worker2 %v %v %v\n", targets2, msgId2, dm)
		sps, _ := waiter.WaitSp(msgId2, targets2, 1000)
		for _, sp := range sps {
			ags := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
			agents = append(agents, ags...)
		}

		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		worker1api.GetAgentResponse(senderId, targets, msgId, agents)
		logger.Debug("Finish: Get Agent Response to Worker1 %v %v %v\n", targets, msgId)*/

		/*pid := dm.GetSimDemand().GetSenderId()
		if len(apm.ProvidersMap[pid]) == 0{
			// 隣接エリアがない場合はそのまま返す
			targets := []uint64{dm.GetSimDemand().GetSenderId()}
			senderId := myProvider.Id
			msgId := dm.GetSimDemand().GetMsgId()
			agents := []*api.Agent{}
			worker1api.GetAgentResponse(senderId, targets, msgId, agents)
		}else{
			//隣接エリアが存在していたらそのAgentProviderへ送る
			// senderIDの取り扱いに注意　workerからはgatewayのみが見えているようになっている
			targets := dm.GetSimDemand().GetTargets()
			senderId := myProvider.Id
			msgId1 := dm.GetSimDemand().GetMsgId()
			msgId2 := worker2api.GetAgentRequest(senderId, targets)
			apm.SetMsgIdMap(msgId1, msgId2) // msgIdを紐づける
		}*/
		agents := []*api.Agent{}
		senderId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		worker1api.GetAgentResponse(senderId, targets, msgId, agents)
		t2 := time.Now()
		duration := t2.Sub(t1).Milliseconds()
		logger.Info("Duration: %v", duration)

	case api.DemandType_UPDATE_PROVIDERS_REQUEST:
		ps1 := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
		//apm.SetProviders1(ps1)
		apm.SetProvider1(ps1)
		//pm.SetProviders(providers)

		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		worker1api.UpdateProvidersResponse(senderId, targets, msgId)
		logger.Info("Finish: Update Providers1 num: %v\n", len(ps1))
	}
}

////////////////////////////////////////////////////////////
//////////     Worker2 Demand Supply Callback     /////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func supplyCallback2(clt *api.SMServiceClient, sp *api.Supply) {
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_GET_AGENT_RESPONSE:
		//fmt.Printf("Get Sp from Worker2%v\n", sp)
		//time.Sleep(10 * time.Millisecond)
		worker2api.SendSpToWait(sp)
		// send to worker1 agent-provider

	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		//masterProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		fmt.Printf("regist provider to Worler2 Provider!\n")

	case api.SupplyType_READY_PROVIDER_RESPONSE:
		//time.Sleep(10 * time.Millisecond)
		worker2api.SendSpToWait(sp)
		fmt.Printf("ready provider response")
	}
}

// Demandのコールバック関数
func demandCallback2(clt *api.SMServiceClient, dm *api.Demand) {
	switch dm.GetSimDemand().GetType() {
	case api.DemandType_READY_PROVIDER_REQUEST:
		provider := dm.GetSimDemand().GetReadyProviderRequest().GetProvider()
		//pm.SetProviders(providers)

		// workerへ登録
		senderId := myProvider.Id
		targets := []uint64{provider.GetId()}
		worker2api.RegistProviderRequest(senderId, targets, myProvider)
		//waiter.WaitSp(msgId, targets, 1000)

		// response
		targets = []uint64{dm.GetSimDemand().GetSenderId()}
		senderId = myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		worker2api.ReadyProviderResponse(senderId, targets, msgId)
		logger.Info("Finish: Regist Provider from ready ")

	case api.DemandType_GET_AGENT_REQUEST:

		t1 := time.Now()
		/*// 隣接エリアがない場合はそのまま返す
		agents := []*api.Agent{}
		senderId := myProvider.Id
		// worker2のagent-providerから取得
		targets1 := []uint64{apm.Provider1.GetId()}
		msgId1 := worker1api.GetAgentRequest(senderId, targets1)
		//logger.Debug("Get Agent Request to Worker1 %v %v %v\n", targets1, msgId1, dm)
		sps, _ := waiter.WaitSp(msgId1, targets1, 1000)
		for _, sp := range sps {
			ags := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
			agents = append(agents, ags...)
		}

		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		worker2api.GetAgentResponse(senderId, targets, msgId, agents)
		logger.Debug("Finish: Get Agent Request Worker2 %v %v\n", targets, msgId)
		//隣接エリアが存在していたらそのAgentProviderへ送る*/
		agents := []*api.Agent{}
		senderId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		worker2api.GetAgentResponse(senderId, targets, msgId, agents)
		t2 := time.Now()
		duration := t2.Sub(t1).Milliseconds()
		logger.Info("Duration: %v", duration)
		// ない場合はそのまま返す

	case api.DemandType_UPDATE_PROVIDERS_REQUEST:
		ps2 := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
		//apm.SetProviders2(ps2)
		apm.SetProvider2(ps2)
		//pm.SetProviders(providers)

		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		worker2api.UpdateProvidersResponse(senderId, targets, msgId)
		logger.Info("Finish: Update Providers2 num: %v\n", len(ps2))
	}

}

func main() {
	logger.Info("StartUp Provider")

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: "GatewayProvider",
		Type: api.ProviderType_GATEWAY,
	}
	//pm1 = simutil.NewProviderManager(myProvider)
	//pm2 = simutil.NewProviderManager(myProvider)

	//////////////////////////////////////////////////
	//////////           worker1             ////////
	////////////////////////////////////////////////

	// Connect to Worker1 Node Server
	for {
		err := api.RegisterNodeName(workerNodeIdAddr1, "GatewayProvider", false)
		if err == nil {
			logger.Info("connected NodeID server!")
			go api.HandleSigInt()
			api.RegisterDeferFunction(api.UnRegisterNode)
			break
		} else {
			logger.Warn("NodeID Error... reconnecting...")
			time.Sleep(2 * time.Second)
		}
	}

	// Connect to Worker1 Synerex Server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(workerSynerexAddr1, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	api.RegisterDeferFunction(func() { conn.Close() })
	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Gateway}")

	// Communicator
	worker1api = api.NewSimAPI()
	worker1api.RegistClients(client, myProvider.Id, argJson)  // channelごとのClientを作成
	worker1api.SubscribeAll(demandCallback1, supplyCallback1) // ChannelにSubscribe

	// workerへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	worker1api.RegistProviderRequest(senderId, targets, myProvider)

	//////////////////////////////////////////////////
	//////////           worker2             ////////
	////////////////////////////////////////////////

	// Connect to Worker2 Node Server
	for {
		err := api.RegisterNodeName(workerNodeIdAddr2, "GatewayProvider", false)
		if err == nil {
			logger.Info("connected NodeID server!")
			go api.HandleSigInt()
			api.RegisterDeferFunction(api.UnRegisterNode)
			break
		} else {
			logger.Warn("NodeID Error... reconnecting...")
			time.Sleep(2 * time.Second)
		}
	}

	// Connect to Worker2 Synerex Server
	var opts2 []grpc.DialOption
	opts2 = append(opts, grpc.WithInsecure())
	conn2, err2 := grpc.Dial(workerSynerexAddr2, opts2...)
	if err2 != nil {
		log.Fatalf("fail to dial: %v", err2)
	}
	api.RegisterDeferFunction(func() { conn2.Close() })
	client2 := api.NewSynerexClient(conn2)
	argJson2 := fmt.Sprintf("{Client:Gateway}")

	// Communicator
	worker2api = api.NewSimAPI()
	worker2api.RegistClients(client2, myProvider.Id, argJson2) // channelごとのClientを作成
	worker2api.SubscribeAll(demandCallback2, supplyCallback2)  // ChannelにSubscribe

	// workerへ登録
	senderId = myProvider.Id
	targets = make([]uint64, 0)
	worker2api.RegistProviderRequest(senderId, targets, myProvider)

	wg := sync.WaitGroup{}
	wg.Add(1)

	wg.Wait()
	api.CallDeferFunctions() // cleanup!

}
