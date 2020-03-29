package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	gosocketio "github.com/mtfelian/golang-socketio"
	pb "github.com/synerex/synerex_alpha/api"
	simapi "github.com/synerex/synerex_alpha/api/simulation"
	provider "github.com/synerex/synerex_alpha/api/simulation/provider"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr      = flag.String("synerex", "127.0.0.1:10000", "The server address in the format of host:port")
	gatewayAddr     = flag.String("gateway", "", "The server address in the format of host:port")
	nodeIdAddr      = flag.String("nodeid", "127.0.0.1:9990", "Node ID Server")
	simulatorAddr   = flag.String("simulator", "127.0.0.1:9995", "Node ID Server")
	visAddr         = flag.String("vis", "127.0.0.1:9995", "Node ID Server")
	monitorAddr     = flag.String("monitor", "127.0.0.1:9993", "Monitor Server")
	areaId          = flag.Int("areaId", 0, "Area ID")
	mu              sync.Mutex
	com             *simutil.Communicator
	workerClock     int
	providerHosts   []string
	providerManager *simutil.ProviderManager
	areaManager     *simutil.AreaManager
	pSources        map[provider.ProviderType]*provider.Source
	logger          *simutil.Logger
)

const MAX_AGENTS_NUM = 1000

func init() {
	workerClock = 0
	providerHosts = make([]string, 0)
	logger = simutil.NewLogger()
	logger.SetPrefix("Scenario")
	flag.Parse()
}

var (
	//fcs *geojson.FeatureCollection
	//geofile string
	port          = 9995
	assetsDir     http.FileSystem
	server        *gosocketio.Server = nil
	providerMutex sync.RWMutex
)

func init() {
	providerMutex = sync.RWMutex{}
}

////////////////////////////////////////////////////////////
//////////////////        Util          ///////////////////
///////////////////////////////////////////////////////////

// providerに変化があった場合にGUIに情報を送る
/*func sendRunnningProviders() {
	providerMutex.RLock()

	//fmt.Printf("providers---------- %v\n", len(runProviders))
	rpJsons := make([]string, 0)
	for _, rp := range providerManager.Providers {
		bytes, _ := json.Marshal(rp)
		rpJson := string(bytes)
		fmt.Printf("provider----------\n")
		//fmt.Printf("Json: %v \n", rpJson)
		rpJsons = append(rpJsons, rpJson)
	}
	//c.Emit("providers", rpJsons)
	server.BroadcastToAll("providers", rpJsons)
	providerMutex.RUnlock()
}*/

////////////////////////////////////////////////////////////
//////////////////     ps Command     ////////////////////
//////////////////////////////////////////////////////////

/*func checkRunning(opt string) []string {
	isLong := false
	if opt == "long" {
		isLong = true
	}
	var procs []string
	i := 0
	providerMutex.RLock()
	if isLong {
		procs = make([]string, len(providerManager.Providers)+2)
		str := fmt.Sprintf("  pid: %-20s : \n", "process name")
		procs[i] = str
		procs[i+1] = "-----------------------------------------------------------------\n"
		i += 2
	} else {
		procs = make([]string, len(providerManager.Providers))
	}
	for _, provider := range providerManager.Providers {
		pid := pSources[provider.Type].Cmd.Process.Pid
		name := provider.Name
		if isLong {
			str := fmt.Sprintf("%5d: %-20s : \n", pid, name)
			procs[i] = str
		} else {
			if i != 0 {
				procs[i] = ", " + name
			} else {
				procs[i] = name
			}
		}
		i++
	}
	providerMutex.RUnlock()
	return procs

}*/

////////////////////////////////////////////////////////////
////////////     Demand Supply Callback     ////////////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// 自分宛かどうか
	if sp.GetTargetId() == providerManager.MyProvider.Id {
		// check if supply is match with my demand.
		switch sp.GetSimSupply().GetType() {
		case simapi.SupplyType_UPDATE_PROVIDERS_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_SET_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_SET_AGENTS_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_START_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_STOP_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		}
	}
}

// Demandのコールバック関数
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	//tid := dm.GetSimDemand().GetPid()
	//pid := providerManager.MyProvider.Id
	// check if supply is match with my demand.
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_FORWARD_CLOCK_REQUEST:
		fmt.Printf("get forwardClockRequest")
	case simapi.DemandType_SET_AGENTS_REQUEST:
		fmt.Printf("set agent")
		/*case simapi.DemandType_REGIST_PROVIDER_REQUEST:
			// providerを追加する
			p := dm.GetSimDemand().GetRegistProviderRequest().GetProvider()
			providerManager.AddProvider(p)
			providerManager.AddMyProvider(p)
			providerManager.CreateIDMap()
			// 登録完了通知
			targets := []uint64{tid}
			senderInfo := providerManager.MyProvider
			com.RegistProviderResponse(senderInfo, targets, pid, tid)

			// UpdateRequest
			idList := providerManager.GetIDList([]simutil.IDType{
				simutil.IDType_CLOCK,
				simutil.IDType_VISUALIZATION,
				simutil.IDType_AGENT,
			})
			pid := providerManager.MyProvider.Id
			targets = idList
			com.UpdateProvidersRequest(senderInfo, targets, pid, idList, providerManager.Providers)

			logger.Info("Success Update Providers")

		case simapi.DemandType_DIVIDE_PROVIDER_REQUEST:
		case simapi.DemandType_KILL_PROVIDER_REQUEST:
		case simapi.DemandType_SEND_PROVIDER_STATUS_REQUEST:
		case simapi.DemandType_SET_PROVIDERS_REQUEST:
			providers := dm.GetSimDemand().GetSetProvidersRequest().GetProviders()

			logger.Info("Get Providers from Gateway")
			for _, p := range providers {
				if p.Type == provider.ProviderType_AGENT {
					logger.Debug("Provider: %v", p.Id)
					providerManager.AddProvider(p)
				}
			}
			//providerManager.CreateIDMap()

			// UpdateRequest
			idList := providerManager.GetIDList([]simutil.IDType{
				simutil.IDType_CLOCK,
				simutil.IDType_VISUALIZATION,
				simutil.IDType_AGENT,
			})
			pid := providerManager.MyProvider.Id
			targets := idList
			senderInfo := providerManager.MyProvider
			com.UpdateProvidersRequest(senderInfo, targets, pid, idList, providerManager.Providers)

			logger.Info("Success Update Providers")

		case simapi.DemandType_GET_PROVIDERS_REQUEST:
			logger.Info("Get Providers Request")
			targets := []uint64{tid}
			senderInfo := providerManager.MyProvider
			com.GetProvidersResponse(senderInfo, targets, pid, tid, providerManager.MyProviders)*/

	}
}

func main() {

	// ProviderManager
	myProvider := provider.NewProvider("ScenarioProvider", provider.ProviderType_SCENARIO, *serverAddr)
	providerManager = simutil.NewProviderManager(myProvider)
	providerManager.CreateIDMap()

	//AreaManager
	//areaManager = simutil.NewAreaManager(mockAreaInfos[*areaId])

	// Connect to Node Server
	sxutil.RegisterNodeName(*nodeIdAddr, "ScenarioProvider", false)
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
	argJson := fmt.Sprintf("{Client:Scenario}")

	// Communicator
	com = simutil.NewCommunicator()
	com.RegistClients(client, argJson)               // channelごとのClientを作成
	com.SubscribeAll(demandCallback, supplyCallback) // ChannelにSubscribe

	// masterへ登録
	targets = []uint64{*serverAddr}
	pid = myProvider.Id
	com.RegistProviderRequest(myProvider, targets)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
