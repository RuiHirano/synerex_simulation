package main

import (
	//"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sync"

	//"time"
	//"runtime"
	//"encoding/json"

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
	serverAddr    = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv       = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	areaIdFlag    = flag.Int("areaId", 1, "Area Id")
	pidFlag       = flag.Int("pid", 1, "Provider Id")
	agentTypeFlag = flag.Int("agentType", 1, "Provider Id")
	areaJson      = flag.String("area_json", "", "Area Json")
	areaInfo      *area.Area
	agentType     = agent.AgentType_PEDESTRIAN // PEDESTRIAN
	com           *simutil.Communicator
	sim           *Simulator
	areaId        uint64
	pid           uint64
)

func flagToAreaInfo(areaJson string) *area.Area {
	bytes := []byte(areaJson)
	json.Unmarshal(bytes, &areaInfo)
	return areaInfo
}

func init() {
	flag.Parse()
	areaInfo = flagToAreaInfo(*areaJson)
	log.Printf("\x1b[31m\x1b[47m \nAreaInfo: %v \x1b[0m\n", areaInfo)
	areaId = uint64(*areaIdFlag)
	pid = uint64(*pidFlag)
}

// callbackSetAgents: Agent情報をセットする要求
func setAgents(dm *pb.Demand) {
	agents := dm.GetSimDemand().GetSetAgentsRequest().GetAgents()
	targetId := dm.GetId()

	// Agent情報を追加する
	sim.AddAgents(agents)

	// セット完了通知を送る
	com.SetAgentsResponse(targetId)
	log.Printf("\x1b[30m\x1b[47m \n Finish: Agents information set. \n Total:  %v \n Add: %v \x1b[0m\n", len(sim.GetAgents()), len(agents))
}

// callbackForwardClock: Agentを計算し、クロックを進める要求
func forwardClock(dm *pb.Demand) {

	targetId := dm.GetId()

	// 同じエリアのAgent情報を取得する
	_, sameAreaAgents := com.GetAgentsRequest(nil)

	// 次の時間のエージェントを計算する
	nextControlAgents := sim.ForwardStep(sameAreaAgents)

	// 隣接エリアにエージェントの情報を送信
	_, neighborAreaAgents := com.GetAgentsRequest(nil)

	// 重複エリアのエージェントを更新する
	nextDuplicateAgents := sim.UpdateDuplicateAgents(nextControlAgents, neighborAreaAgents)

	// Agentsをセットする
	sim.SetAgents(nextDuplicateAgents)

	// クロックを進める
	sim.ForwardClock()

	// 可視化プロバイダへ送信
	com.SetAgentsRequest(nil, nextControlAgents)

	// セット完了通知を送る
	com.ForwardClockResponse(targetId)

	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock forwarded. \n Time:  %v \n Agents Num: %v \x1b[0m\n", sim.Clock.GlobalTime, len(nextControlAgents))
}

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	tid := dm.GetId()
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_UPDATE_PROVIDERS_REQUEST:
		// 参加者リストをセットする要求
		//callbackSetParticipantsRequest(dm)
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
		com.GetAgentsResponse(tid, sim.Agents, sim.AgentType, sim.Area.Id)
	}
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {

	switch sp.GetSimSupply().GetType() {
	case simapi.SupplyType_GET_CLOCK_RESPONSE:
		com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
	case simapi.SupplyType_GET_AGENTS_RESPONSE:
		com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
	case simapi.SupplyType_REGIST_PROVIDER_RESPONSE:
		com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
	}

}

func main() {
	log.Printf("\x1b[31m\x1b[47m \n SyneServ: %v, NodeServ: %v, AreaJson: %v Pid: %v   \x1b[0m\n", *serverAddr, *nodesrv, *areaJson, pid)
	//log.Printf("area id is: %v, agent type is %v", areaId, agentType)

	// Connect to Node Server
	sxutil.RegisterNodeName(*nodesrv, "PedProvider", false)
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
	p := provider.NewProvider("PedestrianProvider", provider.ProviderType_AGENT)
	com = simutil.NewCommunicator(p)
	com.RegistClients(client, argJson)               // channelごとのClientを作成
	com.SubscribeAll(demandCallback, supplyCallback) // ChannelにSubscribe

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 新規参加登録
	com.RegistProviderRequest(nil, p)
	log.Printf("\x1b[30m\x1b[47m \n Finish: This provider registered in scenario-provider \x1b[0m\n")

	// Clock情報を取得
	_, clockInfo = com.GetClockRequest(nil)
	sim.Clock = clockInfo
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.Clock.GlobalTime, sim.Clock.TimeStep)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
