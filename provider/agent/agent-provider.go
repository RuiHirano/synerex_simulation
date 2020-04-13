package main

import (
	//"context"

	"flag"
	"fmt"
	"log"

	//"math/rand"
	"os"
	"sync"

	//"time"

	//"runtime"
	//"encoding/json"

	"github.com/google/uuid"
	api "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"google.golang.org/grpc"
)

var (
	synerexAddr    string
	nodeIdAddr     string
	myProvider     *api.Provider
	workerProvider *api.Provider
	pm             *simutil.ProviderManager
	waiter         *api.Waiter
	simapi         *api.SimAPI
	//com                  *simutil.Communicator
	sim *Simulator2
	//providerManager      *simutil.ProviderManager
	logger        *simutil.Logger
	mu            sync.Mutex
	agentsMessage *Message
)

/*func flagToProviderInfo(pJson string) *api.Provider {
	pInfo := &api.Provider{}
	jsonapi.Unmarshal(strings.NewReader(pJson), pInfo)
	return pInfo
}*/

func init() {
	flag.Parse()
	logger = simutil.NewLogger()
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

// callbackForwardClock: Agentを計算し、クロックを進める要求
/*func forwardClock(dm *api.Demand) {
	senderId := myProvider.Id

	logger.Debug("1: 同エリアエージェント取得")
	//targets := []uint64{}
	//_, sameAreaAgents := simapi.GetAgentRequest(senderId, targets)

	// [2. Calculation]次の時間のエージェントを計算する
	logger.Debug("2: エージェント計算を行う")
	nextControlAgents := sim.ForwardStep(sameAreaAgents)


	logger.Debug("3: 隣接エージェントを取得")
	// [3. Get Neighbor Area Agents]隣接エリアのエージェントの情報を取得
	_, neighborAreaAgents := simapi.GetAgentRequest(senderId, targets)

	logger.Debug("4: エージェントを更新")
	// [4. Update Agents]重複エリアのエージェントを更新する
	nextDuplicateAgents := sim.UpdateDuplicateAgents(nextControlAgents, neighborAreaAgents)
	// Agentsをセットする
	sim.SetAgents(nextDuplicateAgents)

	// [5. Forward Clock]クロックを進める
	logger.Debug("5: クロックを進める")
	sim.ForwardClock()

	logger.Info("Finish: Clock Forwarded. pid %v,  AgentNum:  %v", senderId, len(nextControlAgents))
}*/

func forwardClock() {
	//senderId := myProvider.Id

	logger.Debug("1: 同エリアエージェント取得")
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_AGENT,
	})
	senderId := myProvider.Id
	msgId := simapi.GetAgentRequest(senderId, targets)
	logger.Debug("1: targets %v\n", targets)
	waiter.WaitSp(msgId, targets)
	//targets := []uint64{}
	//_, sameAreaAgents := simapi.GetAgentRequest(senderId, targets)

	// [2. Calculation]次の時間のエージェントを計算する
	logger.Debug("2: エージェント計算を行う")
	nextAgents := sim.ForwardStep()
	agentsMessage.Set(nextAgents)

	logger.Debug("3: 隣接エージェントを取得")
	targets = pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_AGENT,
		simutil.IDType_GATEWAY,
	})
	senderId = myProvider.Id
	msgId = simapi.GetAgentRequest(senderId, targets)
	waiter.WaitSp(msgId, targets)
	// [3. Get Neighbor Area Agents]隣接エリアのエージェントの情報を取得
	//_, neighborAreaAgents := simapi.GetAgentRequest(senderId, targets)

	logger.Debug("4: エージェントを更新")
	// [4. Update Agents]重複エリアのエージェントを更新する
	//nextDuplicateAgents := sim.UpdateDuplicateAgents(nextControlAgents, neighborAreaAgents)
	// Agentsをセットする
	sim.SetAgents(nextAgents)

	// [5. Forward Clock]クロックを進める
	logger.Debug("6: クロックを進める")
	//agentsMessage = NewMessage()
	//sim.ForwardClock()

	logger.Info("Finish: Clock Forwarded. AgentNum:  %v", len(nextAgents))
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
		logger.Info("Finish: Update Providers ")

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
		logger.Info("forward Clock")
		// クロックを進める要求
		forwardClock()

		// response
		senderId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.ForwardClockResponse(senderId, targets, msgId)
		logger.Info("Finish: Forward Clock")

	case api.DemandType_GET_AGENT_REQUEST:
		logger.Debug("get agent request")
		senderId := dm.GetSimDemand().GetSenderId()
		sameAreaIds := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_AGENT,
		})
		neighborAreaIds := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_AGENT,
			simutil.IDType_GATEWAY,
		})
		visIds := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_VISUALIZATION,
		})

		agents := []*api.Agent{}
		if simutil.Contains(sameAreaIds, senderId) {
			// 同じエリアのエージェントプロバイダの場合
			agents = sim.Agents
		} else if simutil.Contains(neighborAreaIds, senderId) {
			// 隣接エリアのエージェントプロバイダの場合
			agents = agentsMessage.Get()
		} else if simutil.Contains(visIds, senderId) {
			// Visプロバイダの場合
			agents = agentsMessage.Get()
		}

		// 全てのプロバイダにmessageを送信し終えたらMessageを初期化する
		agentsMessage.AddSenderId(senderId)
		if agentsMessage.FinishSend(append(neighborAreaIds, visIds...)) {
			logger.Debug("init Message")
			agentsMessage = NewMessage()
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
	// 自分宛かどうか
	switch sp.GetSimSupply().GetType() {
	//case api.SupplyType_GET_AGENT_RESPONSE:
	//	fmt.Printf("get agents response")
	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		workerProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		fmt.Printf("resist provider request")
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

func main() {
	logger.Info("StartUp Provider")

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: "AgentProvider",
		Type: api.ProviderType_AGENT,
	}
	pm = simutil.NewProviderManager(myProvider)

	// Connect to Node Server
	api.RegisterNodeName(nodeIdAddr, myProvider.GetName(), false)
	go api.HandleSigInt()
	api.RegisterDeferFunction(api.UnRegisterNode)

	// Connect to Synerex Server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(synerexAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	api.RegisterDeferFunction(func() { conn.Close() })
	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Agent}")

	// Simulator
	//clockInfo := &api.Clock{GlobalTime: 0}
	//areaInfo := myProvider.GetAgentStatus().GetArea()
	//agentType := api.AgentType_PEDESTRIAN
	//sim = NewSimulator(clockInfo, areaInfo, agentType)
	sim = NewSimulator2()

	// WorkerAPI作成
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, myProvider.Id, argJson) // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback)  // ChannelにSubscribe

	// workerへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	simapi.RegistProviderRequest(senderId, targets, myProvider)

	// test
	//forward()
	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	api.CallDeferFunctions() // cleanup!

}