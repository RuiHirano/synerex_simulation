package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	api "github.com/synerex/synerex_alpha/api"
	napi "github.com/synerex/synerex_alpha/nodeapi"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"github.com/synerex/synerex_alpha/util"
	"google.golang.org/grpc"
)

var (
	myProvider  *api.Provider
	synerexAddr string
	nodeIdAddr  string
	port        string
	startFlag   bool
	masterClock int
	workerHosts []string
	mu          sync.Mutex
	simapi      *api.SimAPI
	//providerManager *Manager
	pm     *simutil.ProviderManager
	logger *util.Logger
	waiter *api.Waiter
)

func init() {
	waiter = api.NewWaiter()
	startFlag = false
	masterClock = 0
	workerHosts = make([]string, 0)
	logger = util.NewLogger()
	flag.Parse()
	//providerManager = NewManager()

	synerexAddr = os.Getenv("SYNEREX_SERVER")
	if synerexAddr == "" {
		synerexAddr = "127.0.0.1:10000"
	}
	nodeIdAddr = os.Getenv("NODEID_SERVER")
	if nodeIdAddr == "" {
		nodeIdAddr = "127.0.0.1:9000"
	}
	port = os.Getenv("PORT")
	if port == "" {
		port = "9990"
	}
}

////////////////////////////////////////////////////////////
////////////     Demand Supply Callback     ////////////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func supplyCallback(clt *api.SMServiceClient, sp *api.Supply) {

}

// Demandのコールバック関数
func demandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	//tid := dm.GetSimDemand().GetSenderId()
	//pid := myProvider.Id
	// check if supply is match with my demand.
	switch dm.GetSimDemand().GetType() {
	case api.DemandType_CREATE_POD_REQUEST:
		// providerを追加する
		cpr := dm.GetSimDemand().GetCreatePodRequest()
		fmt.Printf("get CreatePodRequest %v\n", cpr)
		// 登録完了通知
		senderInfo := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.CreatePodResponse(senderInfo, targets, msgId)

		logger.Info("Finish: Create Pod")
	case api.DemandType_DELETE_POD_REQUEST:
		// providerを追加する
		cpr := dm.GetSimDemand().GetDeletePodRequest()
		fmt.Printf("get DeletePodRequest %v\n", cpr)
		// 登録完了通知
		senderInfo := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.DeletePodResponse(senderInfo, targets, msgId)

		logger.Info("Finish: Delete Pod")
	}
}

func main() {
	fmt.Printf("NumCPU=%d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: "PodManager",
		Type: api.ProviderType_MASTER,
	}
	pm = simutil.NewProviderManager(myProvider)

	// Connect to Node Server
	nodeapi := napi.NewNodeAPI()
	for {
		err := nodeapi.RegisterNodeName(nodeIdAddr, "PodManager", false)
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
	argJson := fmt.Sprintf("{Client:PodManager}")

	// api
	fmt.Printf("client: %v\n", client)
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, myProvider.Id, argJson) // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback)  // ChannelにSubscribe*/
	logger.Info("Connected Synerex Server!\n")

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	nodeapi.CallDeferFunctions() // cleanup!

}
