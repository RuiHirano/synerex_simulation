package main

import (
	"flag"
	"fmt"
	"log"

	//"math/rand"
	"os"
	"sync"

	"encoding/json"
	"time"

	"github.com/google/uuid"
	api "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"google.golang.org/grpc"
)

var (
	myProvider        *api.Provider
	masterProvider    *api.Provider
	workerSynerexAddr string
	workerNodeIdAddr  string
	masterNodeIdAddr  string
	masterSynerexAddr string
	mu                sync.Mutex
	masterapi         *api.SimAPI
	workerapi         *api.SimAPI
	workerClock       int
	providerHosts     []string
	logger            *simutil.Logger
	//providerManager   *Manager
	pm     *simutil.ProviderManager
	waiter *api.Waiter
)

const MAX_AGENTS_NUM = 1000

func init() {
	workerClock = 0
	providerHosts = make([]string, 0)
	logger = simutil.NewLogger()
	logger.SetPrefix("Scenario")
	flag.Parse()

	workerSynerexAddr = os.Getenv("SYNEREX_SERVER")
	if workerSynerexAddr == "" {
		workerSynerexAddr = "127.0.0.1:10000"
	}
	workerNodeIdAddr = os.Getenv("NODEID_SERVER")
	if workerNodeIdAddr == "" {
		workerNodeIdAddr = "127.0.0.1:9000"
	}
	masterSynerexAddr = os.Getenv("MASTER_SYNEREX_SERVER")
	if masterSynerexAddr == "" {
		masterSynerexAddr = "127.0.0.1:10000"
	}
	masterNodeIdAddr = os.Getenv("MASTER_NODEID_SERVER")
	if masterNodeIdAddr == "" {
		masterNodeIdAddr = "127.0.0.1:9000"
	}
	areaJson := os.Getenv("AREA")
	areaJson = "[{\"latitude\": 3, \"longitude\": 4},{\"latitude\": 3, \"longitude\": 4},{\"latitude\": 3, \"longitude\": 4},{\"latitude\": 3, \"longitude\": 4}]"
	bytes := []byte(areaJson)
	var coords []*api.Coord
	json.Unmarshal(bytes, &coords)
	fmt.Printf("coords: %v\n", coords)
	if areaJson == "" {
		areaJson = "127.0.0.1:9000"
	}

	//providerManager = NewManager()
	waiter = api.NewWaiter()
}

////////////////////////////////////////////////////////////
//////////////////        Manager          ///////////////////
///////////////////////////////////////////////////////////

type Manager struct {
	Providers []*api.Provider
}

func NewManager() *Manager {
	m := &Manager{
		Providers: make([]*api.Provider, 0),
	}
	return m
}

func (m *Manager) AddProvider(provider *api.Provider) {
	m.Providers = append(m.Providers, provider)
}

func (m *Manager) GetProviderIds() []uint64 {
	pids := make([]uint64, 0)
	for _, p := range m.Providers {
		pids = append(pids, p.GetId())
	}
	return pids
}

////////////////////////////////////////////////////////////
////////////     Demand Supply Callback     ////////////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func masterSupplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	// 自分宛かどうか
	// check if supply is match with my demand.
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		masterProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		fmt.Printf("regist provider to Master Provider!\n")

	}
}

// Demandのコールバック関数
func masterDemandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	senderId := myProvider.Id
	switch dm.GetSimDemand().GetType() {
	case api.DemandType_FORWARD_CLOCK_REQUEST:
		fmt.Printf("get forwardClockRequest")

		// request to worker providers
		targets := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_AGENT,
			simutil.IDType_VISUALIZATION,
		})
		msgId := workerapi.ForwardClockRequest(senderId, targets)
		waiter.WaitSp(msgId, targets, 1000)

		// response to master
		targets = []uint64{dm.GetSimDemand().GetSenderId()}
		msgId = dm.GetSimDemand().GetMsgId()
		masterapi.ForwardClockResponse(senderId, targets, msgId)

	case api.DemandType_SET_AGENT_REQUEST:
		fmt.Printf("set agent")
		// request to providers
		agents := dm.GetSimDemand().GetSetAgentRequest().GetAgents()
		targets := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_AGENT,
		})
		msgId := workerapi.SetAgentRequest(senderId, targets, agents)
		waiter.WaitSp(msgId, targets, 1000)

		// response to master
		targets = []uint64{dm.GetSimDemand().GetSenderId()}
		msgId = dm.GetSimDemand().GetMsgId()
		masterapi.SetAgentResponse(senderId, targets, msgId)
	}
}

////////////////////////////////////////////////////////////
////////////     Worker Demand Supply Callback    //////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func workerSupplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	// 自分宛かどうか
	// check if supply is match with my demand.
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_UPDATE_PROVIDERS_RESPONSE:
		logger.Info("get sp: %v\n", sp)
		waiter.SendSpToWait(sp)
	case api.SupplyType_SET_CLOCK_RESPONSE:
		logger.Info("get sp: %v\n", sp)
		waiter.SendSpToWait(sp)
	case api.SupplyType_SET_AGENT_RESPONSE:
		logger.Info("get sp: %v\n", sp)
		waiter.SendSpToWait(sp)
	case api.SupplyType_FORWARD_CLOCK_RESPONSE:
		logger.Info("get sp: %v\n", sp)
		waiter.SendSpToWait(sp)
	}
}

// Demandのコールバック関数
func workerDemandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	switch dm.GetSimDemand().GetType() {
	case api.DemandType_REGIST_PROVIDER_REQUEST:
		// providerを追加する
		p := dm.GetSimDemand().GetRegistProviderRequest().GetProvider()
		pm.AddProvider(p)
		fmt.Printf("regist request from agent of vis provider! %v\n")
		// 登録完了通知
		senderId := myProvider.Id
		targets := []uint64{p.GetId()}
		msgId := dm.GetSimDemand().GetMsgId()
		workerapi.RegistProviderResponse(senderId, targets, msgId, myProvider)

		logger.Info("Success Regist Agent or Vis Providers", targets)

		// 参加プロバイダの更新命令
		// request to worker providers
		targets = pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_GATEWAY,
			simutil.IDType_AGENT,
			simutil.IDType_VISUALIZATION,
		})
		providers := pm.GetProviders()
		msgId = workerapi.UpdateProvidersRequest(senderId, targets, providers)
		logger.Info("Wait response from &v\n", targets)
		waiter.WaitSp(msgId, targets, 1000)

	}
}

///////////////////////////
/////    test      ////////
///////////////////////////
/*var mockAgents []*api.Agent

func init() {
	mockAgents = []*api.Agent{}
	for i := 0; i < 100; i++ {
		uid, _ := uuid.NewRandom()
		departure := &api.Coord{
			Longitude: 136.87285 + rand.Float64()*0.01,
			Latitude:  35.17333 + rand.Float64()*0.01,
		}
		destination := &api.Coord{
			Longitude: 136.92285 + rand.Float64()*0.01,
			Latitude:  35.19333 + rand.Float64()*0.01,
		}
		transitPoints := []*api.Coord{destination}
		mockAgents = append(mockAgents, &api.Agent{
			Type: api.AgentType_PEDESTRIAN,
			Id:   uint64(uid.ID()),
			Route: &api.Route{
				Position: &api.Coord{
					Longitude: 136.97285 + rand.Float64()*0.01,
					Latitude:  35.15333 + rand.Float64()*0.01,
				},
				Direction:     30,
				Speed:         60,
				Departure:     departure,
				Destination:   destination,
				TransitPoints: transitPoints,
				NextTransit:   destination,
			},
		})
	}
}

func forwardCLock() {
	time.Sleep(5 * time.Second) // 5s以内にregist providerすること
	senderId := myProvider.Id
	agents := mockAgents
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_AGENT,
	})
	msgId := workerapi.SetAgentRequest(senderId, targets, agents)
	waiter.WaitSp(msgId, targets, 1000)
	fmt.Printf("finish set agents")
	for {
		time.Sleep(1 * time.Second)
		// request to worker providers
		targets := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_AGENT,
			simutil.IDType_VISUALIZATION,
		})
		msgId := workerapi.ForwardClockRequest(senderId, targets)
		waiter.WaitSp(msgId, targets, 1000)
		fmt.Printf("finish forward clock")
	}
}*/

func main() {

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: "WorkerServer",
		Type: api.ProviderType_WORKER,
	}
	pm = simutil.NewProviderManager(myProvider)

	// For Master
	// Connect to Node Server
	for {
		err := api.RegisterNodeName(masterNodeIdAddr, "WorkerProvider", false)
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

	// Connect to Synerex Server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(masterSynerexAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	api.RegisterDeferFunction(func() { conn.Close() })
	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Worker}")

	// Communicator
	masterapi = api.NewSimAPI()
	masterapi.RegistClients(client, myProvider.Id, argJson)            // channelごとのClientを作成
	masterapi.SubscribeAll(masterDemandCallback, masterSupplyCallback) // ChannelにSubscribe

	// masterへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	masterapi.RegistProviderRequest(senderId, targets, myProvider)

	// For Worker
	// Connect to Node Server
	for {
		err := api.RegisterNodeName(workerNodeIdAddr, "WorkerProvider", false)
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

	// Connect to Synerex Server
	var wopts []grpc.DialOption
	wopts = append(wopts, grpc.WithInsecure())
	wconn, werr := grpc.Dial(workerSynerexAddr, wopts...)
	if werr != nil {
		log.Fatalf("fail to dial: %v", werr)
	}
	api.RegisterDeferFunction(func() { wconn.Close() })
	wclient := api.NewSynerexClient(wconn)
	wargJson := fmt.Sprintf("{Client:Worker}")

	// Communicator
	workerapi = api.NewSimAPI()
	workerapi.RegistClients(wclient, myProvider.Id, wargJson)          // channelごとのClientを作成
	workerapi.SubscribeAll(workerDemandCallback, workerSupplyCallback) // ChannelにSubscribe

	// test
	//go forwardCLock()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	fmt.Printf("clean up!")
	api.CallDeferFunctions() // cleanup!

}
