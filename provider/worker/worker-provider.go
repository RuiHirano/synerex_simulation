package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	gosocketio "github.com/mtfelian/golang-socketio"
	api "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"google.golang.org/grpc"
)

var (
	myProvider    *api.Provider
	serverAddr    = flag.String("synerex", "127.0.0.1:10000", "The server address in the format of host:port")
	gatewayAddr   = flag.String("gateway", "", "The server address in the format of host:port")
	nodeIdAddr    = flag.String("nodeid", "127.0.0.1:9990", "Node ID Server")
	simulatorAddr = flag.String("simulator", "127.0.0.1:9995", "Node ID Server")
	visAddr       = flag.String("vis", "127.0.0.1:9995", "Node ID Server")
	monitorAddr   = flag.String("monitor", "127.0.0.1:9993", "Monitor Server")
	areaId        = flag.Int("areaId", 0, "Area ID")
	mu            sync.Mutex
	simapi        *api.SimAPI
	workerClock   int
	providerHosts []string
	logger        *simutil.Logger
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
////////////     Demand Supply Callback     ////////////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func supplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	// 自分宛かどうか
	// check if supply is match with my demand.
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		fmt.Printf("regist provider!\n")
	}
}

// Demandのコールバック関数
func demandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	//tid := dm.GetSimDemand().GetPid()
	//pid := providerManager.MyProvider.Id
	// check if supply is match with my demand.
	switch dm.GetSimDemand().GetType() {
	case api.DemandType_FORWARD_CLOCK_REQUEST:
		fmt.Printf("get forwardClockRequest")
	case api.DemandType_SET_AGENT_REQUEST:
		fmt.Printf("set agent")
	}
}

func main() {

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: "WorkerServer",
		Type: api.ProviderType_WORKER,
	}

	//AreaManager
	//areaManager = simutil.NewAreaManager(mockAreaInfos[*areaId])

	// Connect to Node Server
	api.RegisterNodeName(*nodeIdAddr, "WorkerProvider", false)
	go api.HandleSigInt()
	api.RegisterDeferFunction(api.UnRegisterNode)

	// Connect to Synerex Server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	api.RegisterDeferFunction(func() { conn.Close() })
	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Worker}")

	// Communicator
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, myProvider.Id, argJson) // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback)  // ChannelにSubscribe

	// masterへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	simapi.RegistProviderRequest(senderId, targets, myProvider)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	api.CallDeferFunctions() // cleanup!

}
