package main

import (
	//"flag"
	"log"
	//"math/rand"
	"time"

	//"strings"
	"sync"

	//"github.com/golang/protobuf/jsonpb"
	api "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/provider/simutil"

	"fmt"
	"net/http"
	"os"

	"path/filepath"

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
	logger         *simutil.Logger
	simapi         *api.SimAPI
	waiter         *api.Waiter
)

func init() {
	//flag.Parse()
	logger = simutil.NewLogger()
	synerexAddr = os.Getenv("SYNEREX_SERVER")
	if synerexAddr == "" {
		synerexAddr = "127.0.0.1:10000"
	}
	nodeIdAddr = os.Getenv("NODEID_SERVER")
	if nodeIdAddr == "" {
		nodeIdAddr = "127.0.0.1:9000"
	}
	visAddr = os.Getenv("VIS_ADDRESS")
	if visAddr == "" {
		visAddr = "127.0.0.1:9500"
	}

	providerName = os.Getenv("PROVIDER_NAME")
	if providerName == "" {
		providerName = "VisProvider"
	}

	waiter = api.NewWaiter()
}

////////////////////////////////////////////////////////////
////////////           Harmovis server           ///////////
///////////////////////////////////////////////////////////

func runServer() *gosocketio.Server {

	currentRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	d := filepath.Join(currentRoot, "mclient", "build")

	assetsDir = http.Dir(d)
	log.Println("AssetDir:", assetsDir)

	assetsDir = http.Dir(d)
	server := gosocketio.NewServer()

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected from %s as %s", c.IP(), c.Id())
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected from %s as %s", c.IP(), c.Id())
	})

	return server
}

// assetsFileHandler for static Data
func assetsFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}

	file := r.URL.Path
	//	log.Printf("Open File '%s'",file)
	if file == "/" {
		file = "/index.html"
	}
	f, err := assetsDir.Open(file)
	if err != nil {
		log.Printf("can't open file %s: %v\n", file, err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Printf("can't open file %s: %v\n", file, err)
		return
	}
	http.ServeContent(w, r, file, fi.ModTime(), f)
}

type MapMarker struct {
	mtype int32   `json:"mtype"`
	id    int32   `json:"id"`
	lat   float32 `json:"lat"`
	lon   float32 `json:"lon"`
	angle float32 `json:"angle"`
	speed int32   `json:"speed"`
	area  int32   `json:"area"`
}

// GetJson: json化する関数
func (m *MapMarker) GetJson() string {
	s := fmt.Sprintf("{\"mtype\":%d,\"id\":%d,\"lat\":%f,\"lon\":%f,\"angle\":%f,\"speed\":%d,\"area\":%d}",
		m.mtype, m.id, m.lat, m.lon, m.angle, m.speed, m.area)
	return s
}

// sendToHarmowareVis: harmowareVisに情報を送信する関数
func sendToHarmowareVis(agents []*api.Agent) {

	if agents != nil {
		jsonAgents := make([]string, 0)
		for _, agentInfo := range agents {

			// agentInfoTypeによってエージェントを取得
			switch agentInfo.Type {
			case api.AgentType_PEDESTRIAN:
				//ped := agentInfo.GetPedestrian()
				mm := &MapMarker{
					mtype: int32(agentInfo.Type),
					id:    int32(agentInfo.Id),
					lat:   float32(agentInfo.Route.Position.Latitude),
					lon:   float32(agentInfo.Route.Position.Longitude),
					angle: float32(agentInfo.Route.Direction),
					speed: int32(agentInfo.Route.Speed),
				}
				jsonAgents = append(jsonAgents, mm.GetJson())

			case api.AgentType_CAR:
				//car := agentInfo.GetCar()
				mm := &MapMarker{
					mtype: int32(agentInfo.Type),
					id:    int32(agentInfo.Id),
					lat:   float32(agentInfo.Route.Position.Latitude),
					lon:   float32(agentInfo.Route.Position.Longitude),
					angle: float32(agentInfo.Route.Direction),
					speed: int32(agentInfo.Route.Speed),
				}
				jsonAgents = append(jsonAgents, mm.GetJson())
			}
		}
		mu.Lock()
		ioserv.BroadcastToAll("event", jsonAgents)
		mu.Unlock()
	}
}

// callbackForwardClockRequest: クロックを進める関数
func forwardClock(dm *api.Demand) {
	//log.Printf("\x1b[30m\x1b[47m \n Start: Clock forwarded \n Time:  %v \x1b[0m\n", sim.Clock.GlobalTime)
	//senderId := myProvider.Id
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

	//targets := []uint64{}
	//_, sameAreaAgents := simapi.GetAgentRequest(senderId, targets)
	//agents := []*api.Agent{}
	// Harmowareに送る
	sendToHarmowareVis(allAgents)

	t2 := time.Now()
	duration := t2.Sub(t1).Milliseconds()
	logger.Info("Duration: %v, PID: %v", duration, myProvider.Id)
}

// callback for each Supply
func demandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	switch dm.GetSimDemand().GetType() {
	case api.DemandType_READY_PROVIDER_REQUEST:
		/*provider := dm.GetSimDemand().GetReadyProviderRequest().GetProvider()
		//pm.SetProviders(providers)

		// workerへ登録
		senderId := myProvider.Id
		targets := []uint64{provider.GetId()}
		simapi.RegistProviderRequest(senderId, targets, myProvider)
		//waiter.WaitSp(msgId, targets, 1000)

		// response
		targets = []uint64{dm.GetSimDemand().GetSenderId()}
		senderId = myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.ReadyProviderResponse(senderId, targets, msgId)
		logger.Info("Finish: Regist Provider from ready ")*/

	case api.DemandType_UPDATE_PROVIDERS_REQUEST:
		providers := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
		pm.SetProviders(providers)

		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.UpdateProvidersResponse(senderId, targets, msgId)
		logger.Info("Finish: Update Providers num: %v\n", len(providers))

	case api.DemandType_FORWARD_CLOCK_REQUEST:
		// クロックを進める
		forwardClock(dm)

		// response
		senderId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.ForwardClockResponse(senderId, targets, msgId)
		logger.Info("Finish: Forward Clock")
		/*case api.DemandType_SET_AGENT_REQUEST:

		//waiter.SendDmToWait(msgId)

		// セット完了通知を送る
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.SetAgentResponse(senderId, targets, msgId)
		logger.Info("Finish: Set Agents Add ")*/
	case api.DemandType_FORWARD_CLOCK_INIT_REQUEST:

		// response
		senderId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.ForwardClockInitResponse(senderId, targets, msgId)
		logger.Info("Finish: Forward Clock Init")
	}

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

///////////////////////////
/////    test      ////////
///////////////////////////
/*var mockAgents []*api.Agent

func init() {
	mockAgents = []*api.Agent{}
	for i := 0; i < 100; i++ {
		uid, _ := uuid.NewRandom()
		mockAgents = append(mockAgents, &api.Agent{
			Type: api.AgentType_PEDESTRIAN,
			Id:   uint64(uid.ID()),
			Route: &api.Route{
				Position: &api.Coord{
					Longitude: 136.97285 + rand.Float64()*0.01,
					Latitude:  35.15333 + rand.Float64()*0.01,
				},
				Direction: 0,
				Speed:     0,
			},
		})
	}
}

func sendAgents() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("send agents")
		sendToHarmowareVis(mockAgents)
	}
}*/

func registToWorker() {
	// workerへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	simapi.RegistProviderRequest(senderId, targets, myProvider)

	go func() {
		for {
			if workerProvider != nil {
				logger.Debug("Regist Success to Worker!")
				return
			} else {
				logger.Debug("Couldn't Regist Worker...Retry...\n")
				time.Sleep(2 * time.Second)
				// workerへ登録
				simapi.RegistProviderRequest(senderId, targets, myProvider)
			}
		}
	}()
}

func main() {
	logger.Info("StartUp Provider %v, %v", synerexAddr, myProvider)
	// Provider
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: providerName,
		Type: api.ProviderType_VISUALIZATION,
	}
	pm = simutil.NewProviderManager(myProvider)

	// Connect to Node Server
	for {
		err := api.RegisterNodeName(nodeIdAddr, providerName, false)
		if err == nil {
			logger.Info("connected NodeID server!")
			go api.HandleSigInt()
			api.RegisterDeferFunction(api.UnRegisterNode)
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
	api.RegisterDeferFunction(func() { conn.Close() })
	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Visualization}")

	time.Sleep(5 * time.Second)

	// WorkerAPI作成
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, myProvider.Id, argJson) // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback)  // ChannelにSubscribe

	time.Sleep(5 * time.Second)

	registToWorker()

	// workerへ登録
	/*senderId := myProvider.Id
	targets := make([]uint64, 0)
	simapi.RegistProviderRequest(senderId, targets, myProvider)*/

	// test
	//go sendAgents()

	// Run HarmowareVis Monitor
	ioserv = runServer()
	log.Printf("Running Sio Server..\n")
	if ioserv == nil {
		os.Exit(1)
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", ioserv)
	serveMux.HandleFunc("/", assetsFileHandler)
	log.Printf("Starting Harmoware VIS  Provider on %v", visAddr)
	err = http.ListenAndServe(visAddr, serveMux)
	if err != nil {
		log.Fatal(err)
	}

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	api.CallDeferFunctions() // cleanup!
}
