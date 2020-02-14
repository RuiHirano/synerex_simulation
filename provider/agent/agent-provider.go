package main

import (
	//"context"

	"flag"
	"fmt"
	"log"
	"sync"

	//"time"
	//"runtime"
	//"encoding/json"

	"strings"

	"github.com/golang/protobuf/jsonpb"
	pb "github.com/synerex/synerex_alpha/api"
	simapi "github.com/synerex/synerex_alpha/api/simulation"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/provider"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr           = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv              = flag.String("nodeid_addr", "127.0.0.1:9990", "Node ID Server")
	providerJson         = flag.String("provider_json", "", "Provider Json")
	scenarioProviderJson = flag.String("scenario_provider_json", "", "Provider Json")
	myProvider           *provider.Provider
	scenarioProvider     *provider.Provider
	com                  *simutil.Communicator
	sim                  *Simulator
	providerManager      *simutil.ProviderManager
	logger               *simutil.Logger
	mu                   sync.Mutex
	agentsMessage        *Message
)

func flagToProviderInfo(pJson string) *provider.Provider {
	pInfo := &provider.Provider{}
	jsonpb.Unmarshal(strings.NewReader(pJson), pInfo)
	return pInfo
}

func init() {
	flag.Parse()
	logger = simutil.NewLogger()
	myProvider = flagToProviderInfo(*providerJson)
	scenarioProvider = flagToProviderInfo(*scenarioProviderJson)
	agentsMessage = NewMessage()
	//log.Printf("\x1b[31m\x1b[47m \nProviderInfo: %v \x1b[0m\n", providerInfo.GetStatus())

}

////////////////////////////////////////////////////////////
////////////            Message Class           ///////////
///////////////////////////////////////////////////////////

type Message struct {
	ready  chan struct{}
	agents []*agent.Agent
}

func NewMessage() *Message {
	return &Message{ready: make(chan struct{})}
}
func (m *Message) Set(a []*agent.Agent) {
	m.agents = a
	close(m.ready)
}

func (m *Message) Get() []*agent.Agent {
	<-m.ready
	return m.agents
}

// callbackForwardClock: Agentを計算し、クロックを進める要求
func forwardClock(dm *pb.Demand) {

	tid := dm.GetSimDemand().GetPid()
	pid := providerManager.MyProvider.Id

	// [1. Get Same Area Agents]同じエリアの異種Agent情報を取得する
	idList := providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_SAME,
	})
	_, sameAreaAgents := com.GetAgentsRequest(pid, idList)

	// [2. Calculation]次の時間のエージェントを計算する
	nextControlAgents := sim.ForwardStep(sameAreaAgents)
	agentsMessage.Set(nextControlAgents)
	logger.Error("ForwardAgents %v", nextControlAgents)

	// [3. Get Neighbor Area Agents]隣接エリアのエージェントの情報を取得
	// 同期するIDリスト
	idList = providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_NEIGHBOR,
	})
	_, neighborAreaAgents := com.GetAgentsRequest(pid, idList)

	// [4. Update Agents]重複エリアのエージェントを更新する
	nextDuplicateAgents := sim.UpdateDuplicateAgents(nextControlAgents, neighborAreaAgents)
	// Agentsをセットする
	sim.SetAgents(nextDuplicateAgents)

	// [5. Forward Clock]クロックを進める
	sim.ForwardClock()

	// [6. Send Finish Forward Response]セット完了通知を送る
	com.ForwardClockResponse(pid, tid)

	//logger.Info("Finish: Clock Forwarded. \n Time:  %v \n Agents Num: %v", sim.Clock.GlobalTime, len(nextControlAgents))
	logger.Info("Finish: Clock Forwarded.  AgentNum:  %v", len(sim.Agents))
	agentsMessage = NewMessage()
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	tid := dm.GetSimDemand().GetPid()
	pid := providerManager.MyProvider.Id
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_UPDATE_PROVIDERS_REQUEST:
		// 参加者リストをセットする要求
		providers := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
		providerManager.UpdateProviders(providers)
		providerManager.CreateIDMap()
		com.UpdateProvidersResponse(pid, tid)

	case simapi.DemandType_SET_AGENTS_REQUEST:
		// Agentをセットする
		agents := dm.GetSimDemand().GetSetAgentsRequest().GetAgents()

		// Agent情報を追加する
		sim.AddAgents(agents)

		// セット完了通知を送る
		com.SetAgentsResponse(pid, tid)
		logger.Info("Finish: Set Agents Add: %v", len(sim.Agents))

	case simapi.DemandType_FORWARD_CLOCK_REQUEST:
		// クロックを進める要求
		forwardClock(dm)
	case simapi.DemandType_UPDATE_CLOCK_REQUEST:
		// Clockをセットする
		clockInfo := dm.GetSimDemand().GetSetClockRequest().GetClock()
		sim.Clock = clockInfo
		logger.Info("Finish Update Clock %v, %v", pid, tid)
		com.UpdateClockResponse(pid, tid)

	case simapi.DemandType_GET_AGENTS_REQUEST:
		// vis, neighborProviderからの場合
		idList := providerManager.GetIDList([]simutil.IDType{
			simutil.IDType_VISUALIZATION,
			simutil.IDType_NEIGHBOR,
		})
		if simutil.Contains(idList, tid) {
			go func() {
				agentsInfo := agentsMessage.Get()
				com.GetAgentsResponse(pid, tid, agentsInfo, sim.AgentType, sim.Area.Id)
				//logger.Error("Finish: Send Agents2 %v", dm.GetSimDemand().GetPid())
			}()

		}
		// sameAreaProviderからの場合
		idList = providerManager.GetIDList([]simutil.IDType{
			simutil.IDType_SAME,
		})
		if simutil.Contains(idList, tid) {
			com.GetAgentsResponse(pid, tid, sim.Agents, sim.AgentType, sim.Area.Id)
		}

		logger.Info("Finish: Send Agents")
	}
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// 自分宛かどうか
	if sp.GetTargetId() == providerManager.MyProvider.Id {
		com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		/*switch sp.GetSimSupply().GetType() {
		case simapi.SupplyType_GET_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_GET_AGENTS_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_SET_AGENTS_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_REGIST_PROVIDER_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		}*/
	}
}

func main() {
	logger.Info("StartUp Provider")
	//logger.Error("AgentType %v, AreaID: %v, NeighborIDs: %v", myProvider.GetAgentStatus().GetAgentType(), myProvider.GetAgentStatus().GetArea().GetId(), myProvider.GetAgentStatus().GetArea().GetNeighborAreaIds())
	//log.Printf("area id is: %v, agent type is %v", areaId, agentType)

	// ProviderManager
	//myProvider := provider.NewProvider("AgentProvider", provider.ProviderType_AGENT)
	providerManager = simutil.NewProviderManager(myProvider)
	providerManager.AddProvider(scenarioProvider)
	providerManager.CreateIDMap()

	// Connect to Node Server
	sxutil.RegisterNodeName(*nodesrv, myProvider.GetName(), false)
	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	// Connect to Synerex Server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })
	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Agent}")

	// Simulator
	clockInfo := clock.NewClock(0, 1, 1)
	areaInfo := myProvider.GetAgentStatus().GetArea()
	agentType := myProvider.GetAgentStatus().GetAgentType()
	sim = NewSimulator(clockInfo, areaInfo, agentType)

	// Communicator
	//p := provider.NewProvider("AgentProvider", provider.ProviderType_AGENT)
	com = simutil.NewCommunicator()
	com.RegistClients(client, argJson)               // channelごとのClientを作成
	com.SubscribeAll(demandCallback, supplyCallback) // ChannelにSubscribe

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)

	// 新規参加登録
	// 同期するIDリスト
	idList := providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_SCENARIO,
	})
	pid := providerManager.MyProvider.Id
	com.RegistProviderRequest(pid, idList, myProvider)
	logger.Info("Finish Provider Registration.")

	// Clock情報を取得
	idList = providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_CLOCK,
	})
	_, clockInfo = com.GetClockRequest(pid, idList)
	sim.Clock = clockInfo
	logger.Info("Finish Setting Clock. %v\n", clockInfo)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
