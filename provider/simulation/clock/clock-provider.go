package main

import (
	"flag"
	"log"
	"strings"
	"sync"

	"fmt"

	"time"

	"github.com/golang/protobuf/jsonpb"
	pb "github.com/synerex/synerex_alpha/api"
	simapi "github.com/synerex/synerex_alpha/api/simulation"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/provider"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

// rvo2適用
// daemonでprovider起動setup
// daemonをサーバとしてscenarioに命令

var (
	serverAddr           = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv              = flag.String("nodeid_addr", "127.0.0.1:9990", "Node ID Server")
	providerJson         = flag.String("provider_json", "", "Provider Json")
	scenarioProviderJson = flag.String("scenario_provider_json", "", "Provider Json")
	port                 = flag.Int("port", 10080, "HarmoVis Provider Listening Port")
	isStart              bool
	myProvider           *provider.Provider
	scenarioProvider     *provider.Provider
	mu                   sync.Mutex
	com                  *simutil.Communicator
	sim                  *Simulator
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
	myProvider = flagToProviderInfo(*providerJson)
	scenarioProvider = flagToProviderInfo(*scenarioProviderJson)
	isStart = false
}

// startClock:
func startClock() {

	// 同期するIDリスト
	idList := providerManager.GetIDList([]simutil.IDType{
		//simutil.IDType_SCENARIO,
		simutil.IDType_VISUALIZATION,
		simutil.IDType_AGENT,
	})
	logger.Info("Send Forward Clock Request %v", idList)
	pid := providerManager.MyProvider.Id
	com.ForwardClockRequest(pid, idList)

	// calc next time
	sim.ForwardStep()
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock forwarded \n Time:  %v \x1b[0m\n", sim.Clock.GlobalTime)

	// 待機
	time.Sleep(time.Duration(sim.Clock.TimeStep) * time.Second)

	// 次のサイクルを行う
	if isStart {
		startClock()
	} else {
		log.Printf("\x1b[30m\x1b[47m \n Finish: Clock stopped \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.Clock.GlobalTime, sim.Clock.TimeStep)
		isStart = false
		// exit goroutin
		return
	}

}

// Supplyのコールバック関数
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// 自分宛かどうか
	//logger.Info("Get Forward Clock Response2")
	if sp.GetTargetId() == providerManager.MyProvider.Id {
		// check if supply is match with my demand.
		switch sp.GetSimSupply().GetType() {
		case simapi.SupplyType_UPDATE_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_FORWARD_CLOCK_RESPONSE:
			//logger.Info("Get Forward Clock Response")
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_BACK_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_SEND_PROVIDER_STATUS_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_REGIST_PROVIDER_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())

		}
	}
}

// Demandのコールバック関数
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	tid := dm.GetSimDemand().GetPid()
	pid := providerManager.MyProvider.Id
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_GET_CLOCK_REQUEST:
		logger.Debug("GetClock: Clock %v\n", sim.Clock)
		// Clock情報を提供する
		com.GetClockResponse(pid, tid, sim.Clock)

	case simapi.DemandType_SET_CLOCK_REQUEST:
		// Clockをセットする
		clockInfo := dm.GetSimDemand().GetSetClockRequest().GetClock()
		sim.Clock = clockInfo
		//log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.Clock.GlobalTime, sim.Clock.TimeStep)

		// Request
		idList := providerManager.GetIDList([]simutil.IDType{
			simutil.IDType_VISUALIZATION,
			simutil.IDType_AGENT,
		})
		logger.Info("Request Update Clock %v", idList)
		com.UpdateClockRequest(pid, idList, clockInfo)
		logger.Info("Finish Update Clock")

		// Response to Scenario
		com.SetClockResponse(pid, tid)

	case simapi.DemandType_START_CLOCK_REQUEST:
		// Clockをスタートする
		if isStart == false {
			isStart = true
			go startClock()
		} else {
			logger.Warn("Clock is already started.")
		}

	case simapi.DemandType_STOP_CLOCK_REQUEST:
		//Clockをストップする
		isStart = false
		com.StopClockResponse(pid, tid)

	case simapi.DemandType_UPDATE_PROVIDERS_REQUEST:
		providers := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
		providerManager.UpdateProviders(providers)
		providerManager.CreateIDMap()
		com.UpdateProvidersResponse(pid, tid)
		// プロバイダーを更新する
		//setClock(dm)

	}

}

func main() {
	logger.Info("StartUp Provider")

	// ProviderManager
	//myProvider := provider.NewProvider("ClockProvider", provider.ProviderType_CLOCK)
	providerManager = simutil.NewProviderManager(myProvider)
	providerManager.AddProvider(scenarioProvider)
	providerManager.CreateIDMap()
	logger.Debug("ClockPID %v \n", myProvider.Id)

	// connect to node server
	sxutil.RegisterNodeName(*nodesrv, "ClockProvider", false)
	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	// connect to synerex server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })
	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Clock}")

	// Simulator
	clockInfo := clock.NewClock(0, 1, 1)
	sim = NewSimulator(clockInfo)

	// Communicator
	//clockProviderInfo := &provider.ClockStatus{}
	//provider := provider.NewClockProvider("ClockProvider", clockProviderInfo)
	com = simutil.NewCommunicator()
	com.RegistClients(client, argJson)               // channelごとのClientを作成
	com.SubscribeAll(demandCallback, supplyCallback) // ChannelにSubscribe

	// Communicatorのsetup
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

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
