package main

import (
	//"flag"
	"log"
	//"math/rand"
	"time"

	//"strings"
	"sync"

	//"github.com/golang/protobuf/jsonpb"
	"runtime"

	api "github.com/synerex/synerex_alpha/api"
	napi "github.com/synerex/synerex_alpha/nodeapi"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"github.com/synerex/synerex_alpha/util"

	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	gosocketio "github.com/mtfelian/golang-socketio"
	"google.golang.org/grpc"
)

var (
	synerexAddr    string
	nodeIdAddr     string
	visAddr        string
	providerName   string
	myProvider     *api.Provider
	workerProvider *api.Provider
	pm             *simutil.ProviderManager
	mu             sync.Mutex
	assetsDir      http.FileSystem
	ioserv         *gosocketio.Server
	logger         *util.Logger
	simapi         *api.SimAPI
	waiter         *api.Waiter
	db             *Database
)

func init() {
	//flag.Parse()
	logger = util.NewLogger()
	synerexAddr = os.Getenv("SYNEREX_SERVER")
	if synerexAddr == "" {
		synerexAddr = "127.0.0.1:10000"
	}
	nodeIdAddr = os.Getenv("NODEID_SERVER")
	if nodeIdAddr == "" {
		nodeIdAddr = "127.0.0.1:9000"
	}

	providerName = os.Getenv("PROVIDER_NAME")
	if providerName == "" {
		providerName = "DatabaseProvider"
	}

	waiter = api.NewWaiter()
	capacity := uint64(3) // 3サイクル分のみ格納
	db = NewDatabase(capacity)
}

////////////////////////////////////////////////////////////
////////////              Database              ///////////
///////////////////////////////////////////////////////////

type TimeData struct {
	Clock: *api.Clock
	Agents: []*api.Agent
}

type Database struct {
	Data     []*TimeData
	Capacity uint64
}

func NewDatabase(capacity uint64) *Database {
	db := &Database{
		Data:     []*TimeData{},
		Capacity: capacity,
	}

	return db
}

func (db *Database) Push(clock *api.Clock, agents []*api.Agent) {
	db.Data = append(db.Data, &TimeData{
		Clock: clock,
		Agents: agents,
	})
	if len(db.Data) > int(db.Capacity) {
		pos := len(db.Data) - int(db.Capacity)
		db.Data = db.Data[pos : len(db.Data)-1]
	}
}

func (db *Database) Get() [][]*api.Agent {
	return db.Data
}

// callbackForwardClockRequest: クロックを進める関数
func forwardClock(dm *api.Demand) {
	t1 := time.Now()
	// エージェントからの可視化リクエスト待ち
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_AGENT,
	})
	//uid, _ := uuid.NewRandom()
	senderId := myProvider.Id
	sps, _ := simapi.GetAgentRequest(senderId, targets)
	//sps, _ := waiter.WaitSp(msgId, targets, 1000)

	allAgents := []*api.Agent{}
	for _, sp := range sps {
		agents := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
		allAgents = append(allAgents, agents...)
	}

	db.Push(allAgents)
	t2 := time.Now()
	duration := t2.Sub(t1).Milliseconds()
	logger.Info("Duration: %v, PID: %v", duration, myProvider.Id)
}

// callback for each Supply
func demandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	switch dm.GetSimDemand().GetType() {

	case api.DemandType_GET_AGENT_REQUEST:
		data := db.Get()
		agents := data[len(data)-1]

		// response
		pId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.GetAgentResponse(pId, targets, msgId, agents)

	case api.DemandType_SET_AGENT_REQUEST:
		
		agents := dm.GetSimSupply().GetSetAgentResponse().GetAgents()
		db.Push(agents)
		// response
		pId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.SetAgentResponse(pId, targets, msgId)

}

// callback for each Supply
func supplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_GET_AGENT_RESPONSE:
		//time.Sleep(10 * time.Millisecond)
		fmt.Printf("get agents response")
		simapi.SendSpToWait(sp)
	case api.SupplyType_REGIST_PROVIDER_RESPONSE:

		mu.Lock()
		workerProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		mu.Unlock()
		fmt.Printf("resist provider response")
	}
}

func main() {
	logger.Info("StartUp Provider %v, %v", synerexAddr, myProvider)
	fmt.Printf("NumCPU=%d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Provider
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: providerName,
		Type: api.ProviderType_DATABASE,
	}
	pm = simutil.NewProviderManager(myProvider)

	// Connect to Node Server
	nodeapi := napi.NewNodeAPI()
	for {
		err := nodeapi.RegisterNodeName(nodeIdAddr, providerName, false)
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
	argJson := fmt.Sprintf("{Client:Visualization}")

	// WorkerAPI作成
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, myProvider.Id, argJson) // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback)  // ChannelにSubscribe

	time.Sleep(5 * time.Second)

	registToWorker()

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	nodeapi.CallDeferFunctions() // cleanup!
}
