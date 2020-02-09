package main

import (
	"flag"
	"log"
	"sync"

	"fmt"

	"time"

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
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	pidFlag    = flag.Int("pid", 1, "Provider Id")
	port       = flag.Int("port", 10080, "HarmoVis Provider Listening Port")
	isStart    bool
	pid        uint64
	mu         sync.Mutex
	com        *simutil.Communicator
	sim        *Simulator
)

func init() {
	flag.Parse()
	isStart = false
	pid = uint64(*pidFlag)
}

// startClock:
func startClock() {

	com.ForwardClockRequest(nil)

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
	// check if supply is match with my demand.
	switch sp.GetSimSupply().GetType() {
	case simapi.SupplyType_UPDATE_CLOCK_RESPONSE:
		com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
	case simapi.SupplyType_FORWARD_CLOCK_RESPONSE:
		com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
	case simapi.SupplyType_BACK_CLOCK_RESPONSE:
		com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
	case simapi.SupplyType_SEND_PROVIDER_STATUS_RESPONSE:
		com.SendToWaitCh(sp, sp.GetSimSupply().GetType())

	}
}

// Demandのコールバック関数
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	tid := dm.GetId()
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_GET_CLOCK_REQUEST:
		// Clock情報を提供する
		com.GetClockResponse(tid, sim.Clock)

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
		com.StopClockResponse(tid)

	case simapi.DemandType_UPDATE_PROVIDERS_REQUEST:
		// プロバイダーを更新する
		//setClock(dm)

	}

}

func main() {

	log.Printf("\x1b[31m\x1b[47m \n SyneServ: %v, NodeServ: %v, Pid: %v   \x1b[0m\n", *serverAddr, *nodesrv, pid)

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
	clockProviderInfo := &provider.Clock{}
	provider := provider.NewClockProvider("ClockProvider", provider.ProviderType_CLOCK, clockProviderInfo)
	com = simutil.NewCommunicator(provider)
	com.RegistClients(client, argJson)               // channelごとのClientを作成
	com.SubscribeAll(demandCallback, supplyCallback) // ChannelにSubscribe

	// Communicatorのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 新規参加登録
	com.RegistProviderRequest(nil, provider)
	log.Printf("\x1b[30m\x1b[47m \n Finish: This provider registered in scenario-provider \x1b[0m\n")

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
