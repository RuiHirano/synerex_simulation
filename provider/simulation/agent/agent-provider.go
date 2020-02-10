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
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/provider"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr           = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv              = flag.String("nodeid_addr", "127.0.0.1:9990", "Node ID Server")
	providerJson         = flag.String("provider_json", "", "Provider Json")
	scenarioProviderJson = flag.String("scenario_provider_json", "", "Provider Json")
	areaInfo             *area.Area
	providerInfo         *provider.Provider
	scenarioProviderInfo *provider.Provider
	agentType            = agent.AgentType_PEDESTRIAN // PEDESTRIAN
	com                  *simutil.Communicator
	sim                  *Simulator
	areaId               uint64
	providerManager      *simutil.ProviderManager
	logger               *simutil.Logger
)

func flagToProviderInfo(pJson string) *provider.Provider {
	pInfo := &provider.Provider{}
	jsonpb.Unmarshal(strings.NewReader(pJson), pInfo)
	return pInfo
}

func init() {
	flag.Parse()
	logger = simutil.NewLogger()
	providerInfo = flagToProviderInfo(*providerJson)
	scenarioProviderInfo = flagToProviderInfo(*scenarioProviderJson)
	log.Printf("\x1b[31m\x1b[47m \nProviderInfo: %v \x1b[0m\n", providerInfo.GetStatus())

}

// callbackSetAgents: Agent情報をセットする要求
func setAgents(dm *pb.Demand) {
	agents := dm.GetSimDemand().GetSetAgentsRequest().GetAgents()
	targetId := dm.GetId()

	// Agent情報を追加する
	sim.AddAgents(agents)

	// セット完了通知を送る
	pid := providerManager.MyProvider.Id
	com.SetAgentsResponse(pid, targetId)
	log.Printf("\x1b[30m\x1b[47m \n Finish: Agents information set. \n Total:  %v \n Add: %v \x1b[0m\n", len(sim.GetAgents()), len(agents))
}

// callbackForwardClock: Agentを計算し、クロックを進める要求
func forwardClock(dm *pb.Demand) {

	targetId := dm.GetSimDemand().GetPid()
	pid := providerManager.MyProvider.Id

	// 同じエリアのAgent情報を取得する
	// 同期するIDリスト
	idList := providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_SAME,
	})
	_, sameAreaAgents := com.GetAgentsRequest(pid, idList)

	// 次の時間のエージェントを計算する
	nextControlAgents := sim.ForwardStep(sameAreaAgents)

	// 隣接エリアのエージェントの情報を取得
	// 同期するIDリスト
	idList = providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_NEIGHBOR,
	})
	_, neighborAreaAgents := com.GetAgentsRequest(pid, idList)

	// 重複エリアのエージェントを更新する
	nextDuplicateAgents := sim.UpdateDuplicateAgents(nextControlAgents, neighborAreaAgents)

	// Agentsをセットする
	sim.SetAgents(nextDuplicateAgents)

	// クロックを進める
	sim.ForwardClock()

	// 可視化プロバイダへ送信
	// 同期するIDリスト
	idList = providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_VISUALIZATION,
	})
	com.SetAgentsRequest(pid, idList, nextControlAgents)

	// セット完了通知を送る
	com.ForwardClockResponse(pid, targetId)

	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock forwarded. \n Time:  %v \n Agents Num: %v \x1b[0m\n", sim.Clock.GlobalTime, len(nextControlAgents))
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	tid := dm.GetSimDemand().GetPid()
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_UPDATE_PROVIDERS_REQUEST:
		// 参加者リストをセットする要求
	case simapi.DemandType_SET_AGENTS_REQUEST:
		// Agentをセットする
		setAgents(dm)
	case simapi.DemandType_FORWARD_CLOCK_REQUEST:
		// クロックを進める要求
		forwardClock(dm)
	case simapi.DemandType_UPDATE_CLOCK_REQUEST:
		// Clockをセットする
		clock := dm.GetSimDemand().GetSetClockRequest().GetClock()
		sim.Clock = clock
		log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.Clock.GlobalTime, sim.Clock.TimeStep)

	case simapi.DemandType_GET_AGENTS_REQUEST:
		// エージェント情報を送る
		pid := providerManager.MyProvider.Id
		com.GetAgentsResponse(pid, tid, sim.Agents, sim.AgentType, sim.Area.Id)

	}
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// 自分宛かどうか
	if sp.GetTargetId() == providerManager.MyProvider.Id {
		switch sp.GetSimSupply().GetType() {
		case simapi.SupplyType_GET_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_GET_AGENTS_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_REGIST_PROVIDER_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		}
	}
}

func main() {
	logger.Info("StartUp Provider: SyneServ: %v, NodeServ: %v, AreaJson: %v Pinfo: %v ", *serverAddr, *nodesrv, providerInfo)
	//log.Printf("area id is: %v, agent type is %v", areaId, agentType)

	// ProviderManager
	myProvider := provider.NewProvider("AgentProvider", provider.ProviderType_AGENT)
	providerManager = simutil.NewProviderManager(myProvider)
	providerManager.AddProvider(scenarioProviderInfo)
	providerManager.CreateIDMap()

	// Connect to Node Server
	sxutil.RegisterNodeName(*nodesrv, providerInfo.GetName(), false)
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
	argJson := fmt.Sprintf("{Client:Ped, AreaId: %d}", areaId)

	// Simulator
	clockInfo := clock.NewClock(0, 1, 1)
	sim = NewSimulator(clockInfo, areaInfo, agentType)

	// Communicator
	p := provider.NewProvider("AgentProvider", provider.ProviderType_AGENT)
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
	com.RegistProviderRequest(pid, idList, p)
	logger.Info("Finish Provider Registration.")

	// Clock情報を取得
	idList = providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_CLOCK,
	})
	_, clockInfo = com.GetClockRequest(pid, idList)
	sim.Clock = clockInfo
	logger.Info("Finish Setting Clock. \n GlobalTime:  %v \n TimeStep: %v", sim.Clock.GlobalTime, sim.Clock.TimeStep)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
