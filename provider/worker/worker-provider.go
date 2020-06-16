package main

import (
	"flag"
	"fmt"
	"log"

	//"math/rand"
	"os"
	"sync"

	"encoding/json"
	"runtime"
	"time"

	"github.com/google/uuid"
	api "github.com/synerex/synerex_alpha/api"
	napi "github.com/synerex/synerex_alpha/nodeapi"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"github.com/synerex/synerex_alpha/util"
	"google.golang.org/grpc"
)

var (
	myProvider        *api.Provider
	masterProvider    *api.Provider
	workerSynerexAddr string
	workerNodeIdAddr  string
	masterNodeIdAddr  string
	masterSynerexAddr string
	providerName      string
	mu                sync.Mutex
	masterapi         *api.SimAPI
	workerapi         *api.SimAPI
	workerClock       int
	logger            *util.Logger
	pm                *simutil.ProviderManager
	waiter            *api.Waiter
)

const MAX_AGENTS_NUM = 1000

func init() {
	workerClock = 0
	logger = util.NewLogger()
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
		masterSynerexAddr = "master:700"
	}
	masterNodeIdAddr = os.Getenv("MASTER_NODEID_SERVER")
	if masterNodeIdAddr == "" {
		masterNodeIdAddr = "master:600"
	}
	providerName = os.Getenv("PROVIDER_NAME")
	if providerName == "" {
		providerName = "WorkerProvider"
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
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		mu.Lock()
		masterProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		mu.Unlock()
		fmt.Printf("regist provider to Master Provider!\n")

	}
}

// Demandのコールバック関数
func masterDemandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	senderId := myProvider.Id
	switch dm.GetSimDemand().GetType() {

	case api.DemandType_FORWARD_CLOCK_REQUEST:
		fmt.Printf("get forwardClockRequest")
		t1 := time.Now()

		// request to worker providers
		targets := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_AGENT,
		})

		// init
		workerapi.ForwardClockInitRequest(senderId, targets)

		// forward
		workerapi.ForwardClockRequest(senderId, targets)

		t2 := time.Now()
		duration := t2.Sub(t1).Milliseconds()
		logger.Info("Duration: %v, PID: %v", duration, myProvider.Id)
		// response to master
		targets = []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		logger.Debug("Response to master pid %v, msgId%v\n", myProvider.Id, msgId)
		masterapi.ForwardClockResponse(senderId, targets, msgId)

	case api.DemandType_FORWARD_CLOCK_INIT_REQUEST:
		//agentsMessage = NewMessage()

		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		masterapi.ForwardClockInitResponse(senderId, targets, msgId)
		logger.Info("Finish: Forward Clock Init")

	case api.DemandType_SET_AGENT_REQUEST:
		fmt.Printf("set agent")
		// request to providers
		agents := dm.GetSimDemand().GetSetAgentRequest().GetAgents()
		targets := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_AGENT,
		})
		workerapi.SetAgentRequest(senderId, targets, agents)

		// response to master
		targets = []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		masterapi.SetAgentResponse(senderId, targets, msgId)

	case api.DemandType_UPDATE_PROVIDERS_REQUEST:
		providers := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
		//pm.SetProviders(providers)

		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		masterapi.UpdateProvidersResponse(senderId, targets, msgId)
		logger.Info("Finish: Update Workers num: %v\n", len(providers))
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
		workerapi.SendSpToWait(sp)
	case api.SupplyType_SET_CLOCK_RESPONSE:
		workerapi.SendSpToWait(sp)
	case api.SupplyType_SET_AGENT_RESPONSE:
		workerapi.SendSpToWait(sp)
	case api.SupplyType_FORWARD_CLOCK_RESPONSE:
		workerapi.SendSpToWait(sp)
	case api.SupplyType_FORWARD_CLOCK_INIT_RESPONSE:
		workerapi.SendSpToWait(sp)
	}
}

// Demandのコールバック関数
func workerDemandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	switch dm.GetSimDemand().GetType() {
	case api.DemandType_REGIST_PROVIDER_REQUEST:
		// providerを追加する
		p := dm.GetSimDemand().GetRegistProviderRequest().GetProvider()
		pm.AddProvider(p)
		fmt.Printf("regist request from agent of vis provider! %v\n", p)
		// 登録完了通知
		senderId := myProvider.Id
		targets := []uint64{p.GetId()}
		msgId := dm.GetSimDemand().GetMsgId()
		workerapi.RegistProviderResponse(senderId, targets, msgId, myProvider)

		logger.Info("Success Regist Agent or Vis Providers", targets)

		// 参加プロバイダの更新命令
		targets = pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_GATEWAY,
			simutil.IDType_AGENT,
		})
		providers := pm.GetProviders()
		workerapi.UpdateProvidersRequest(senderId, targets, providers)
		logger.Info("Update Providers! Provider Num %v \n", len(targets))

	}
}

func registToMaster() {
	// masterへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	masterapi.RegistProviderRequest(senderId, targets, myProvider)

	go func() {
		for {
			if masterProvider != nil {
				logger.Debug("Regist Success to Master!")
				return
			} else {
				logger.Debug("Couldn't Regist Master...Retry...\n")
				time.Sleep(2 * time.Second)
				// masterへ登録
				masterapi.RegistProviderRequest(senderId, targets, myProvider)
			}
		}
	}()
}

func main() {
	fmt.Printf("NumCPU=%d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: providerName,
		Type: api.ProviderType_WORKER,
	}
	pm = simutil.NewProviderManager(myProvider)

	// For Master
	// Connect to Node Server
	nodeapi1 := napi.NewNodeAPI()
	for {
		err := nodeapi1.RegisterNodeName(masterNodeIdAddr, providerName, false)
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

	// Connect to Node Server
	nodeapi2 := napi.NewNodeAPI()
	for {
		err := nodeapi2.RegisterNodeName(workerNodeIdAddr, providerName, false)
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

	// Connect to Synerex Server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(masterSynerexAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	nodeapi1.RegisterDeferFunction(func() { conn.Close() })
	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Worker}")

	// Connect to Synerex Server
	var wopts []grpc.DialOption
	wopts = append(wopts, grpc.WithInsecure())
	wconn, werr := grpc.Dial(workerSynerexAddr, wopts...)
	if werr != nil {
		log.Fatalf("fail to dial: %v", werr)
	}
	nodeapi2.RegisterDeferFunction(func() { wconn.Close() })
	wclient := api.NewSynerexClient(wconn)
	wargJson := fmt.Sprintf("{Client:Worker}")

	// Communicator
	masterapi = api.NewSimAPI()
	masterapi.RegistClients(client, myProvider.Id, argJson)            // channelごとのClientを作成
	masterapi.SubscribeAll(masterDemandCallback, masterSupplyCallback) // ChannelにSubscribe

	// Communicator
	workerapi = api.NewSimAPI()
	workerapi.RegistClients(wclient, myProvider.Id, wargJson)          // channelごとのClientを作成
	workerapi.SubscribeAll(workerDemandCallback, workerSupplyCallback) // ChannelにSubscribe

	time.Sleep(3 * time.Second)

	registToMaster()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	fmt.Printf("clean up!")
	nodeapi1.CallDeferFunctions() // cleanup!
	nodeapi2.CallDeferFunctions() // cleanup!

}
