package main

import (
	"flag"
	"log"
	"sync"

	"fmt"

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
	serverAddr      = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv         = flag.String("nodeid_addr", "127.0.0.1:9990", "Node ID Server")
	mu              sync.Mutex
	com             *simutil.Communicator
	providerManager *simutil.ProviderManager
	logger          *simutil.Logger
)

func init() {
	flag.Parse()
	logger = simutil.NewLogger()
}

func apiTest() {
	pid := providerManager.MyProvider.Id
	idList := []uint64{pid}

	// clock test
	clockInfo := &clock.Clock{}
	logger.Info("1. Clear UpdateClockRequest")
	com.UpdateClockRequest(pid, idList, clockInfo)
	logger.Info("3. Clear UpdateClockRequest")

	logger.Info("1. Clear SetClockRequest")
	com.SetClockRequest(pid, idList, clockInfo)
	logger.Info("3. Clear SetClockRequest")
}

// Supplyのコールバック関数
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// 自分宛かどうか
	if sp.GetTargetId() == providerManager.MyProvider.Id {
		// check if supply is match with my demand.

		com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		//switch sp.GetSimSupply().GetType() {
		//case simapi.SupplyType_UPDATE_CLOCK_RESPONSE:
		//	com.SendToWaitCh(sp, sp.GetSimSupply().GetType())

		/*case simapi.SupplyType_FORWARD_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_BACK_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_SEND_PROVIDER_STATUS_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_REGIST_PROVIDER_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())*/

	}
}

// Demandのコールバック関数
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	tid := dm.GetSimDemand().GetPid()
	pid := providerManager.MyProvider.Id
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_UPDATE_CLOCK_REQUEST:
		logger.Info("2. Clear UpdateClockRequest")
		com.UpdateClockResponse(pid, tid)
	case simapi.DemandType_SET_CLOCK_REQUEST:
		logger.Info("2. Clear SetClockRequest")
		com.SetClockResponse(pid, tid)
		/*case simapi.DemandType_GET_CLOCK_REQUEST:
			// Clock情報を提供する
			com.GetClockResponse(pid, tid, sim.Clock)

		case simapi.DemandType_SET_CLOCK_REQUEST:
			// Clockをセットする
			clock := dm.GetSimDemand().GetSetClockRequest().GetClock()
			sim.Clock = clock
			log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.Clock.GlobalTime, sim.Clock.TimeStep)

		case simapi.DemandType_START_CLOCK_REQUEST:
			// Clockをスタートする
			if isStart == false {
				isStart = true
				go startClock()
			}

		case simapi.DemandType_STOP_CLOCK_REQUEST:
			//Clockをストップする
			isStart = false
			com.StopClockResponse(pid, tid)

		case simapi.DemandType_UPDATE_PROVIDERS_REQUEST:
			providers := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
			providerManager.UpdateProviders(providers)
			com.UpdateProvidersResponse(pid, tid)
			// プロバイダーを更新する
			//setClock(dm)*/

	}

}

func main() {
	logger.Info("StartUp Provider: SyneServ: %v, NodeServ: %v, AreaJson: %v Pinfo: %v ", *serverAddr, *nodesrv)

	// ProviderManager
	myProvider := provider.NewProvider("TestProvider", provider.ProviderType_CLOCK)
	providerManager = simutil.NewProviderManager(myProvider)
	providerManager.CreateIDMap()

	// connect to node server
	sxutil.RegisterNodeName(*nodesrv, "TestProvider1", false)
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
	argJson := fmt.Sprintf("{Client:Test}")

	// Communicator
	com = simutil.NewCommunicator()
	com.RegistClients(client, argJson)               // channelごとのClientを作成
	com.SubscribeAll(demandCallback, supplyCallback) // ChannelにSubscribe

	// Communicatorのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	apiTest()

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
