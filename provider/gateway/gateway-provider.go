package main

// main synerex serverからgatewayを介してother synerex serverへ情報を送る
// 基本的に一方通行

import (
	//"flag"
	"fmt"
	"log"
	"os"

	//"strings"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	//"github.com/golang/protobuf/jsonpb"
	"github.com/synerex/synerex_alpha/api"
	napi "github.com/synerex/synerex_alpha/nodeapi"
	"github.com/synerex/synerex_alpha/util"
	"google.golang.org/grpc"
)

var (
	workerSynerexAddr1 string
	workerSynerexAddr2 string
	workerNodeIdAddr1  string
	workerNodeIdAddr2  string
	providerName       string
	//pm1                *simutil.ProviderManager
	//pm2                *simutil.ProviderManager
	apm            *AgentProviderManager
	waiter         *api.Waiter
	mu             sync.Mutex
	myProvider     *api.Provider
	worker1        *api.Provider
	worker2        *api.Provider
	agentProvider1 *api.Provider
	agentProvider2 *api.Provider
	worker1api     *api.SimAPI
	worker2api     *api.SimAPI
	//scenarioProvider   *provider.Provider
	//com1               *simutil.Communicator
	//com2               *simutil.Communicator
	//providerManager1   *simutil.ProviderManager
	//providerManager2   *simutil.ProviderManager
	logger *util.Logger
	//mes1               *Message
	//mes2               *Message
)

func init() {
	//flag.Parse()
	logger = util.NewLogger()
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

	providerName = os.Getenv("PROVIDER_NAME")
	if providerName == "" {
		providerName = "GatewayProvider"
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
		logger.Debug("Provider1: PID: %v, Name: %v\n", p.Id, p.Name)
		if apm.Provider1 == nil && p.GetType() == api.ProviderType_AGENT {
			apm.Provider1 = p
			logger.Warn("Set Provider1!\n")
		}
	}
	mu.Unlock()
}

func (apm *AgentProviderManager) SetProvider2(ps []*api.Provider) {
	mu.Lock()
	for _, p := range ps {
		logger.Debug("Provider2: PID: %v, Name: %v\n", p.Id, p.Name)
		if apm.Provider2 == nil && p.GetType() == api.ProviderType_AGENT {
			apm.Provider2 = p
			logger.Warn("Set Provider2!\n")
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
		fmt.Printf("Get Sp from Worker1\n")

		//time.Sleep(10 * time.Millisecond)
		worker1api.SendSpToWait(sp)
		/*msgId2 := sp.GetSimSupply().GetMsgId()
		// send to worker1
		msgId1 := apm.MsgIdMap[msgId2]
		senderId := myProvider.Id
		agents := sp.GetSimSupply().GetGetAgentResponse()
		worker1api.GetAgentResponse(senderId, targets, msgId1, agents)*/

	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		mu.Lock()
		worker1 = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		mu.Unlock()
		fmt.Printf("regist provider to Worler1 Provider!\n")
	}
}

// Demandのコールバック関数
func demandCallback1(clt *api.SMServiceClient, dm *api.Demand) {

	switch dm.GetSimDemand().GetType() {

	case api.DemandType_GET_AGENT_REQUEST:
		logger.Debug("get agent request\n")
		// 隣接エリアがない場合はそのまま返す
		t1 := time.Now()

		agents := []*api.Agent{}
		senderId := myProvider.Id
		// worker2のagent-providerから取得
		targets2 := []uint64{agentProvider2.Id}
		sps, _ := worker2api.GetAgentRequest(senderId, targets2)

		for _, sp := range sps {
			ags := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
			agents = append(agents, ags...)
		}

		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		worker1api.GetAgentResponse(senderId, targets, msgId, agents)
		logger.Debug("Finish: Get Agent Response to Worker1 %v %v %v\n", targets, msgId)

		t2 := time.Now()
		duration := t2.Sub(t1).Milliseconds()
		logger.Info("Duration: %v, PID: %v", duration, myProvider.Id)

	case api.DemandType_UPDATE_PROVIDERS_REQUEST:
		logger.Debug("update providers request\n")
		ps1 := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()

		//test
		for _, p := range ps1 {
			if agentProvider1 == nil && p.GetType() == api.ProviderType_AGENT {
				mu.Lock()
				agentProvider1 = p
				mu.Unlock()
				logger.Debug("Set Provider1!\n")
			}
		}

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
		fmt.Printf("Get Sp from Worker2\n")
		worker2api.SendSpToWait(sp)

	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		mu.Lock()
		worker2 = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		mu.Unlock()
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

	case api.DemandType_GET_AGENT_REQUEST:
		logger.Debug("get agent request\n", dm.)

		t1 := time.Now()
		// 隣接エリアがない場合はそのまま返す
		agents := []*api.Agent{}
		senderId := myProvider.Id
		// worker2のagent-providerから取得
		targets1 := []uint64{agentProvider1.Id}
		sps, _ := worker1api.GetAgentRequest(senderId, targets1)

		for _, sp := range sps {
			ags := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
			agents = append(agents, ags...)
		}

		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		worker2api.GetAgentResponse(senderId, targets, msgId, agents)
		logger.Debug("Finish: Get Agent Response to Worker2 %v %v\n", targets, msgId)

		t2 := time.Now()
		duration := t2.Sub(t1).Milliseconds()
		logger.Info("Duration: %v, PID: %v", duration, myProvider.Id)
		// ない場合はそのまま返す

	case api.DemandType_UPDATE_PROVIDERS_REQUEST:
		logger.Debug("update providers request\n")
		ps2 := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()

		//test
		for _, p := range ps2 {
			if agentProvider2 == nil && p.GetType() == api.ProviderType_AGENT {
				mu.Lock()
				agentProvider2 = p
				mu.Unlock()
				logger.Debug("Set Provider2!\n")
			}
		}

		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		worker2api.UpdateProvidersResponse(senderId, targets, msgId)
		logger.Info("Finish: Update Providers2 num: %v\n", len(ps2))
	}

}

func registToWorker() {
	// workerへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	worker1api.RegistProviderRequest(senderId, targets, myProvider)
	worker2api.RegistProviderRequest(senderId, targets, myProvider)

	go func() {
		for {
			if worker2 != nil {
				logger.Debug("Regist Success to Worker2!")
				return
			} else {
				logger.Debug("Couldn't Regist Worker2...Retry...\n")
				time.Sleep(2 * time.Second)
				// workerへ登録
				worker2api.RegistProviderRequest(senderId, targets, myProvider)
			}
		}
	}()

	go func() {
		for {
			if worker1 != nil {
				logger.Debug("Regist Success to Worker1!")
				return
			} else {
				logger.Debug("Couldn't Regist Worker1...Retry...\n")
				time.Sleep(2 * time.Second)
				// workerへ登録
				worker1api.RegistProviderRequest(senderId, targets, myProvider)
			}
		}
	}()
}

func main() {
	logger.Info("StartUp Provider")
	fmt.Printf("NumCPU=%d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: providerName,
		Type: api.ProviderType_GATEWAY,
	}

	//////////////////////////////////////////////////
	//////////           worker1             ////////
	////////////////////////////////////////////////

	// Connect to Worker1 Node Server
	nodeapi1 := napi.NewNodeAPI()
	for {
		err := nodeapi1.RegisterNodeName(workerNodeIdAddr1, providerName, false)
		if err == nil {
			logger.Info("connected NodeID server!")
			go nodeapi1.HandleSigInt()
			nodeapi1.RegisterDeferFunction(nodeapi1.UnRegisterNode)
			break
		} else {
			logger.Warn("NodeID Error... reconnecting...")
			time.Sleep(2 * time.Second)
		}
	}

	// Connect to Worker2 Node Server
	nodeapi2 := napi.NewNodeAPI()
	for {
		err := nodeapi2.RegisterNodeName(workerNodeIdAddr2, providerName, false)
		if err == nil {
			logger.Info("connected NodeID server!")
			go nodeapi2.HandleSigInt()
			nodeapi2.RegisterDeferFunction(nodeapi2.UnRegisterNode)
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
	nodeapi1.RegisterDeferFunction(func() { conn.Close() })
	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Gateway}")

	//////////////////////////////////////////////////
	//////////           worker2             ////////
	////////////////////////////////////////////////

	// Connect to Worker2 Synerex Server
	var opts2 []grpc.DialOption
	opts2 = append(opts, grpc.WithInsecure())
	conn2, err2 := grpc.Dial(workerSynerexAddr2, opts2...)
	if err2 != nil {
		log.Fatalf("fail to dial: %v", err2)
	}
	nodeapi2.RegisterDeferFunction(func() { conn2.Close() })
	client2 := api.NewSynerexClient(conn2)
	argJson2 := fmt.Sprintf("{Client:Gateway}")

	// Communicator
	worker1api = api.NewSimAPI()
	worker1api.RegistClients(client, myProvider.Id, argJson)  // channelごとのClientを作成
	worker1api.SubscribeAll(demandCallback1, supplyCallback1) // ChannelにSubscribe

	// Communicator
	worker2api = api.NewSimAPI()
	worker2api.RegistClients(client2, myProvider.Id, argJson2) // channelごとのClientを作成
	worker2api.SubscribeAll(demandCallback2, supplyCallback2)  // ChannelにSubscribe

	time.Sleep(5 * time.Second)

	registToWorker()

	wg := sync.WaitGroup{}
	wg.Add(1)

	wg.Wait()
	nodeapi1.CallDeferFunctions() // cleanup!
	nodeapi2.CallDeferFunctions() // cleanup!

}
