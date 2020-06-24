package main

import (
	//"context"

	"flag"
	"fmt"
	"log"

	//"math/rand"
	"os"
	"sync"

	"time"

	//"runtime"
	"encoding/json"
	"runtime"

	"github.com/google/uuid"
	api "github.com/synerex/synerex_alpha/api"
	napi "github.com/synerex/synerex_alpha/nodeapi"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"github.com/synerex/synerex_alpha/util"
	"google.golang.org/grpc"
)

var (
	synerexAddr    string
	nodeIdAddr     string
	visSynerexAddr string
	visNodeIdAddr  string
	visAddr        string
	providerName   string
	myProvider     *api.Provider
	workerProvider *api.Provider
	visProvider    *api.Provider
	pm             *simutil.ProviderManager
	simapi         *api.SimAPI
	vissimapi      *api.SimAPI
	sim            *Simulator
	logger         *util.Logger
	mu             sync.Mutex
	agentsMessage  *Message
	myArea         *api.Area
	agentType      api.AgentType
)

func init() {
	flag.Parse()
	logger = util.NewLogger()
	agentsMessage = NewMessage()

	synerexAddr = os.Getenv("SYNEREX_SERVER")
	if synerexAddr == "" {
		synerexAddr = "127.0.0.1:10000"
	}
	nodeIdAddr = os.Getenv("NODEID_SERVER")
	if nodeIdAddr == "" {
		nodeIdAddr = "127.0.0.1:9000"
	}

	visSynerexAddr = os.Getenv("VIS_SYNEREX_SERVER")
	if visSynerexAddr == "" {
		visSynerexAddr = "127.0.0.1:10000"
	}
	visNodeIdAddr = os.Getenv("VIS_NODEID_SERVER")
	if visNodeIdAddr == "" {
		visNodeIdAddr = "127.0.0.1:9000"
	}

	providerName = os.Getenv("PROVIDER_NAME")
	if providerName == "" {
		providerName = "AgentProvider"
	}

	visAddr = os.Getenv("VIS_SERVER")
	if visAddr == "" {
		visAddr = "127.0.0.1:8080"
	}

	areaJson := os.Getenv("AREA")
	bytes := []byte(areaJson)
	json.Unmarshal(bytes, &myArea)
	fmt.Printf("myArea: %v\n", myArea)

	agentType = api.AgentType_PEDESTRIAN
}

////////////////////////////////////////////////////////////
////////////            Message Class           ///////////
///////////////////////////////////////////////////////////

type Message struct {
	ready  chan struct{}
	agents []*api.Agent
}

func NewMessage() *Message {
	return &Message{ready: make(chan struct{}), agents: make([]*api.Agent, 0)}
}

func (m *Message) Set(a []*api.Agent) {
	m.agents = a
	close(m.ready)
}

func (m *Message) Get() []*api.Agent {
	select {
	case <-m.ready:
		//case <-time.After(100 * time.Millisecond):
		//	logger.Warn("Timeout Get")
	}

	return m.agents
}

////////////////////////////////////////////////////////////
////////////            Message Class2           ///////////
///////////////////////////////////////////////////////////

type Message2 struct {
	isFinish bool
	agents   []*api.Agent
}

func NewMessage2() *Message2 {
	return &Message2{isFinish: false, agents: make([]*api.Agent, 0)}
}

func (m *Message2) Set(a []*api.Agent) {
	m.agents = a
	m.isFinish = true
}

func (m *Message2) Get() []*api.Agent {
	for {
		if m.isFinish == true {
			time.Sleep(1 * time.Millisecond)
			break
		}
	}

	return m.agents
}

func forwardClock() {
	//senderId := myProvider.Id
	//t1 := time.Now()
	//logger.Debug("1: 同エリアエージェント取得")
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_SAME,
	})
	sameAgents := []*api.Agent{}
	if len(targets) != 0 {
		senderId := myProvider.Id
		sps, _ := simapi.GetAgentRequest(senderId, targets)
		////logger.Debug("1: targets %v\n", targets)
		for _, sp := range sps {
			agents := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
			sameAgents = append(sameAgents, agents...)
		}
	}

	// [2. Calculation]次の時間のエージェントを計算する
	//logger.Debug("2: エージェント計算を行う")
	nextControlAgents := sim.ForwardStep(sameAgents) // agents in control area
	//logger.Debug("2: Set")
	agentsMessage.Set(nextControlAgents)

	// databaseに保存
	/*targets = pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_DATABASE,
	})
	simapi.SetAgentRequest(myProvider.Id, targets, nextControlAgents)*/

	// visに保存
	targets = pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_VISUALIZATION,
	})
	vissimapi.SetAgentRequest(myProvider.Id, targets, nextControlAgents)

	//logger.Debug("3: 隣接エージェントを取得")
	targets = pm.GetProviderIds([]simutil.IDType{
		//simutil.IDType_NEIGHBOR,
		simutil.IDType_GATEWAY,
	})

	neighborAgents := []*api.Agent{}
	if len(targets) != 0 {
		senderId := myProvider.Id
		sps, _ := simapi.GetAgentRequest(senderId, targets)
		////logger.Debug("3: targets %v\n", targets)
		for _, sp := range sps {
			agents := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
			neighborAgents = append(neighborAgents, agents...)
		}
	}

	//logger.Debug("4: エージェントを更新")
	// [4. Update Agents]重複エリアのエージェントを更新する
	nextAgents := sim.UpdateDuplicateAgents(nextControlAgents, neighborAgents)
	// Agentsをセットする
	sim.SetAgents(nextAgents)

	//logger.Info("Finish: Clock Forwarded. AgentNum:  %v", len(nextControlAgents))
	logger.Info("\x1b[32m\x1b[40m [ Agent : %v ] \x1b[0m", len(nextControlAgents))
	//t2 := time.Now()
	//duration := t2.Sub(t1).Milliseconds()
	//logger.Info("Duration: %v, PID: %v", duration, myProvider.Id)
}

// callback for each Supply
func demandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	switch dm.GetSimDemand().GetType() {

	case api.DemandType_UPDATE_PROVIDERS_REQUEST:
		providers := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
		pm.SetProviders(providers)

		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.UpdateProvidersResponse(senderId, targets, msgId)
		//logger.Info("Finish: Update Providers num: %v\n", len(providers))
		//for _, p := range providers {
		//logger.Debug("PID: %v,  Name: %v\n", p.Id, p.Name)
		//}

	case api.DemandType_SET_AGENT_REQUEST:
		// Agentをセットする
		agents := dm.GetSimDemand().GetSetAgentRequest().GetAgents()

		// Agent情報を追加する
		sim.AddAgents(agents)

		// セット完了通知を送る
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.SetAgentResponse(senderId, targets, msgId)
		//logger.Info("\x1b[32m\x1b[40m [ Agent : %v ] \x1b[0m", num)
		////logger.Info("Finish: Agents %v\n", num)

	case api.DemandType_FORWARD_CLOCK_REQUEST:
		// クロックを進める要求
		forwardClock()

		// response
		senderId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.ForwardClockResponse(senderId, targets, msgId)
		//logger.Info("Finish: Forward Clock")

	case api.DemandType_FORWARD_CLOCK_INIT_REQUEST:
		agentsMessage = NewMessage()

		// response
		senderId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.ForwardClockInitResponse(senderId, targets, msgId)
		//logger.Info("Finish: Forward Clock Init")

	case api.DemandType_GET_AGENT_REQUEST:
		////logger.Debug("get agent request %v\n", dm)
		senderId := dm.GetSimDemand().GetSenderId()
		sameAreaIds := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_SAME,
		})
		neighborAreaIds := pm.GetProviderIds([]simutil.IDType{
			//simutil.IDType_NEIGHBOR,
			simutil.IDType_GATEWAY,
		})
		visIds := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_VISUALIZATION,
		})

		agents := []*api.Agent{}
		if util.Contains(sameAreaIds, senderId) {
			// 同じエリアのエージェントプロバイダの場合
			agents = sim.Agents
		} else if util.Contains(neighborAreaIds, senderId) {
			// 隣接エリアのエージェントプロバイダの場合
			////logger.Debug("Get Agent Request from \n%v\n", dm)
			agents = agentsMessage.Get()
		} else if util.Contains(visIds, senderId) {
			// Visプロバイダの場合
			agents = agentsMessage.Get()
		}
		////logger.Debug("get agent request2 %v\n")

		// response
		pId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.GetAgentResponse(pId, targets, msgId, agents)

	}
}

// callback for each Supply
func supplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		//logger.Debug("resist provider response")
		mu.Lock()
		workerProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		mu.Unlock()
	case api.SupplyType_GET_AGENT_RESPONSE:
		//time.Sleep(10 * time.Millisecond)
		////logger.Debug("get agent response \n", sp)
		simapi.SendSpToWait(sp)
	case api.SupplyType_SET_AGENT_RESPONSE:
		////logger.Debug("response set agent")
		//time.Sleep(10 * time.Millisecond)
		////logger.Debug("get agent response \n", sp)
		simapi.SendSpToWait(sp)
	}
}

//////////////////// for VIS ////////////////////////////
// callback for each Supply
func visDemandCallback(clt *api.SMServiceClient, dm *api.Demand) {

}

// callback for each Supply
func visSupplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		//logger.Debug("resist provider response")
		mu.Lock()
		visProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		pm.AddProvider(visProvider)
		mu.Unlock()
	case api.SupplyType_SET_AGENT_RESPONSE:
		////logger.Debug("response set agent from vis")
		//time.Sleep(10 * time.Millisecond)
		////logger.Debug("get agent response \n", sp)
		vissimapi.SendSpToWait(sp)
	}
}

func registToWorker() {
	// workerへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	simapi.RegistProviderRequest(senderId, targets, myProvider)

	go func() {
		for {
			if workerProvider != nil {
				//logger.Debug("Regist Success to Worker!")
				return
			} else {
				//logger.Debug("Couldn't Regist Worker...Retry...\n")
				time.Sleep(2 * time.Second)
				// workerへ登録
				simapi.RegistProviderRequest(senderId, targets, myProvider)
			}
		}
	}()
}

func registToVis() {
	// workerへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	vissimapi.RegistProviderRequest(senderId, targets, myProvider)

	go func() {
		for {
			if visProvider != nil {
				//logger.Debug("Regist Success to Vis!")
				return
			} else {
				//logger.Debug("Couldn't Regist Vis...Retry...\n")
				time.Sleep(2 * time.Second)
				// visへ登録
				vissimapi.RegistProviderRequest(senderId, targets, myProvider)
			}
		}
	}()
}

func main() {
	//logger.Info("StartUp Provider")
	fmt.Printf("NumCPU=%d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: providerName,
		Type: api.ProviderType_AGENT,
		Data: &api.Provider_AgentStatus{
			AgentStatus: &api.AgentStatus{
				Area:      myArea,
				AgentType: api.AgentType_PEDESTRIAN,
			},
		},
	}
	pm = simutil.NewProviderManager(myProvider)

	// Connect to Node Server
	nodeapi := napi.NewNodeAPI()
	for {
		err := nodeapi.RegisterNodeName(nodeIdAddr, myProvider.GetName(), false)
		if err == nil {
			//logger.Info("connected NodeID server!")
			go nodeapi.HandleSigInt()
			nodeapi.RegisterDeferFunction(nodeapi.UnRegisterNode)
			break
		} else {
			logger.Warn("NodeID Error... reconnecting...")
			time.Sleep(2 * time.Second)
		}
	}

	// Connect to Synerex Server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(synerexAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	nodeapi.RegisterDeferFunction(func() { conn.Close() })
	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Agent}")

	// WorkerAPI作成
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, myProvider.Id, argJson) // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback)  // ChannelにSubscribe

	////////////////  for Vis //////////////////////
	// Connect to VisNode Server
	visNodeApi := napi.NewNodeAPI()
	for {
		err := visNodeApi.RegisterNodeName(visNodeIdAddr, myProvider.GetName(), false)
		if err == nil {
			//logger.Info("connected VIS NodeID server!")
			go visNodeApi.HandleSigInt()
			visNodeApi.RegisterDeferFunction(visNodeApi.UnRegisterNode)
			break
		} else {
			logger.Warn("VisNodeID Error... reconnecting..., %v, %v\n", visNodeIdAddr, visSynerexAddr)
			time.Sleep(2 * time.Second)
		}
	}

	// Connect to Synerex Server
	//var opts []grpc.DialOption
	//opts = append(opts, grpc.WithInsecure())
	visconn, err := grpc.Dial(visSynerexAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	visNodeApi.RegisterDeferFunction(func() { visconn.Close() })
	visclient := api.NewSynerexClient(visconn)
	//argJson := fmt.Sprintf("{Client:Agent}")

	// VisAPI作成
	vissimapi = api.NewSimAPI()
	vissimapi.RegistClients(visclient, myProvider.Id, argJson)   // channelごとのClientを作成
	vissimapi.SubscribeAll(visDemandCallback, visSupplyCallback) // ChannelにSubscribe

	//////////////////////////////////////////////////

	// Simulator
	sim = NewSimulator(myArea, api.AgentType_PEDESTRIAN)

	time.Sleep(5 * time.Second)

	registToWorker()
	registToVis()

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	nodeapi.CallDeferFunctions() // cleanup!

}
