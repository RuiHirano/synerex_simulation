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
	providerName   string
	myProvider     *api.Provider
	workerProvider *api.Provider
	pm             *simutil.ProviderManager
	waiter         *api.Waiter
	simapi         *api.SimAPI
	//com                  *simutil.Communicator
	sim *Simulator2
	//providerManager      *simutil.ProviderManager
	logger        *util.Logger
	mu            sync.Mutex
	agentsMessage *Message
	myArea        *api.Area
	agentType     api.AgentType
)

func init() {
	flag.Parse()
	logger = util.NewLogger()
	waiter = api.NewWaiter()
	//myProvider = flagToProviderInfo(*providerJson)
	//scenarioProvider = flagToProviderInfo(*scenarioProviderJson)
	agentsMessage = NewMessage()

	synerexAddr = os.Getenv("SYNEREX_SERVER")
	if synerexAddr == "" {
		synerexAddr = "127.0.0.1:10000"
	}
	nodeIdAddr = os.Getenv("NODEID_SERVER")
	if nodeIdAddr == "" {
		nodeIdAddr = "127.0.0.1:9000"
	}

	providerName = os.Getenv("PROVIDER_NAME")
	if providerName == "" {
		providerName = "AgentProvider"
	}

	areaJson := os.Getenv("AREA")
	bytes := []byte(areaJson)
	//var area *api.Area
	json.Unmarshal(bytes, &myArea)
	fmt.Printf("myArea: %v\n", myArea)

	agentType = api.AgentType_PEDESTRIAN
}

////////////////////////////////////////////////////////////
////////////            Message Class           ///////////
///////////////////////////////////////////////////////////

type Message struct {
	ready     chan struct{}
	agents    []*api.Agent
	senderIds []uint64
}

func NewMessage() *Message {
	return &Message{ready: make(chan struct{}), agents: make([]*api.Agent, 0), senderIds: []uint64{}}
}

func (m *Message) AddSenderId(senderId uint64) {
	m.senderIds = append(m.senderIds, senderId)
}

func (m *Message) FinishSend(targets []uint64) bool {
	for _, tgt := range targets {
		isExist := false
		for _, sid := range m.senderIds {
			if tgt == sid {
				isExist = true
			}
		}
		if isExist == false {
			return false
		}
	}
	return true
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

func forwardClock() {
	//senderId := myProvider.Id
	var com1, com2 int64
	t1 := time.Now()
	logger.Debug("1: 同エリアエージェント取得")
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_SAME,
	})
	sameAgents := []*api.Agent{}
	if len(targets) != 0 {
		senderId := myProvider.Id
		comt1 := time.Now()
		sps, _ := simapi.GetAgentRequest(senderId, targets)
		comt2 := time.Now()
		com1 = comt2.Sub(comt1).Milliseconds()
		//logger.Debug("1: targets %v\n", targets)
		//sps, _ := waiter.WaitSp(msgId, targets, 1000)
		for _, sp := range sps {
			agents := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
			sameAgents = append(sameAgents, agents...)
		}
	}

	// [2. Calculation]次の時間のエージェントを計算する
	logger.Debug("2: エージェント計算を行う")
	nextControlAgents := sim.ForwardStep(sameAgents) // agents in control area
	//logger.Debug("2: Set")
	agentsMessage.Set(nextControlAgents)

	logger.Debug("3: 隣接エージェントを取得")
	targets = pm.GetProviderIds([]simutil.IDType{
		//simutil.IDType_NEIGHBOR,
		simutil.IDType_GATEWAY,
	})

	neighborAgents := []*api.Agent{}
	if len(targets) != 0 {
		senderId := myProvider.Id
		comt1 := time.Now()
		sps, _ := simapi.GetAgentRequest(senderId, targets)
		comt2 := time.Now()
		com2 = comt2.Sub(comt1).Milliseconds()
		//logger.Debug("3: targets %v\n", targets)
		//sps, _ := waiter.WaitSp(msgId, targets, 1000)
		for _, sp := range sps {
			agents := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
			neighborAgents = append(neighborAgents, agents...)
		}
	}

	logger.Debug("4: エージェントを更新")
	// [4. Update Agents]重複エリアのエージェントを更新する
	nextAgents := sim.UpdateDuplicateAgents(nextControlAgents, neighborAgents)
	// Agentsをセットする
	sim.SetAgents(nextAgents)

	// [5. Forward Clock]クロックを進める
	//logger.Debug("6: クロックを進める")
	//agentsMessage = NewMessage()
	//sim.ForwardClock()

	logger.Info("Finish: Clock Forwarded. AgentNum:  %v", len(nextControlAgents))
	t2 := time.Now()
	duration := t2.Sub(t1).Milliseconds()
	comDuration := com2 - com1
	lpDuration := duration-comDuration
	logger.Info("Total: %v, ComDuration: %v, LpDuration: %v", duration, comDuration, lpDuration)
}

// callback for each Supply
func demandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	switch dm.GetSimDemand().GetType() {
	case api.DemandType_READY_PROVIDER_REQUEST:
		/*provider := dm.GetSimDemand().GetReadyProviderRequest().GetProvider()
		//pm.SetProviders(providers)

		// workerへ登録
		senderId := myProvider.Id
		targets := []uint64{provider.GetId()}
		simapi.RegistProviderRequest(senderId, targets, myProvider)
		//waiter.WaitSp(msgId, targets, 1000)

		// response
		targets = []uint64{dm.GetSimDemand().GetSenderId()}
		senderId = myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.ReadyProviderResponse(senderId, targets, msgId)
		logger.Info("Finish: Regist Provider from ready ")*/

	case api.DemandType_UPDATE_PROVIDERS_REQUEST:
		providers := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
		pm.SetProviders(providers)

		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.UpdateProvidersResponse(senderId, targets, msgId)
		logger.Info("Finish: Update Providers num: %v\n", len(providers))
		for _, p := range providers {
			logger.Debug("PID: %v,  Name: %v\n", p.Id, p.Name)
		}

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
		logger.Info("Finish: Set Agents Add %v\n", len(agents))

	case api.DemandType_FORWARD_CLOCK_REQUEST:
		// クロックを進める要求
		forwardClock()

		// response
		senderId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.ForwardClockResponse(senderId, targets, msgId)
		logger.Info("Finish: Forward Clock")

	case api.DemandType_FORWARD_CLOCK_INIT_REQUEST:
		agentsMessage = NewMessage()

		// response
		senderId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.ForwardClockInitResponse(senderId, targets, msgId)
		logger.Info("Finish: Forward Clock Init")

	case api.DemandType_GET_AGENT_REQUEST:
		//logger.Debug("get agent request %v\n", dm)
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
			//logger.Debug("Get Agent Request from \n%v\n", dm)
			agents = agentsMessage.Get()
		} else if util.Contains(visIds, senderId) {
			// Visプロバイダの場合
			agents = agentsMessage.Get()
		}

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
		logger.Debug("resist provider response")
		mu.Lock()
		workerProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		mu.Unlock()
	case api.SupplyType_GET_AGENT_RESPONSE:
		//time.Sleep(10 * time.Millisecond)
		//logger.Debug("get agent response \n", sp)
		simapi.SendSpToWait(sp)
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
				Direction:     0,
				Speed:         0,
				Departure:     departure,
				Destination:   destination,
				TransitPoints: transitPoints,
				NextTransit:   destination,
			},
		})
	}
}

func forward() {
	sim.AddAgents(mockAgents)
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("send agents")
		forwardClock()
	}
}*/

func registToWorker() {
	// workerへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	simapi.RegistProviderRequest(senderId, targets, myProvider)

	go func() {
		for {
			if workerProvider != nil {
				logger.Debug("Regist Success to Worker!")
				return
			} else {
				logger.Debug("Couldn't Regist Worker...Retry...\n")
				time.Sleep(2 * time.Second)
				// workerへ登録
				simapi.RegistProviderRequest(senderId, targets, myProvider)
			}
		}
	}()
}

func main() {
	logger.Info("StartUp Provider")

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
			logger.Info("connected NodeID server!")
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

	// Simulator
	//clockInfo := &api.Clock{GlobalTime: 0}
	//areaInfo := myProvider.GetAgentStatus().GetArea()
	//agentType := api.AgentType_PEDESTRIAN
	//sim = NewSimulator(clockInfo, areaInfo, agentType)
	sim = NewSimulator2(myArea, api.AgentType_PEDESTRIAN)

	time.Sleep(5 * time.Second)

	// WorkerAPI作成
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, myProvider.Id, argJson) // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback)  // ChannelにSubscribe

	time.Sleep(5 * time.Second)

	registToWorker()

	// workerへ登録
	/*logger.Debug("regist to worker")
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	simapi.RegistProviderRequest(senderId, targets, myProvider)*/
	//sps := waiter.WaitSp(msgId, targets, 1000)

	// test
	//forward()
	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	nodeapi.CallDeferFunctions() // cleanup!

}
