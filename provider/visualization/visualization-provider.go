package main

import (
	//"flag"
	"encoding/json"
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

	"path/filepath"

	"github.com/google/uuid"
	gosocketio "github.com/mtfelian/golang-socketio"
	"google.golang.org/grpc"
)

var (
	synerexAddr       string
	nodeIdAddr        string
	masterNodeIdAddr  string
	masterSynerexAddr string
	visAddr           string
	providerName      string
	myProvider        *api.Provider
	masterProvider    *api.Provider
	pm                *simutil.ProviderManager
	mu                sync.Mutex
	masterapi         *api.SimAPI
	visapi            *api.SimAPI
	assetsDir         http.FileSystem
	ioserv            *gosocketio.Server
	logger            *util.Logger
	waiter            *api.Waiter
	agentsMessage     *Message
)

func init() {

	logger = util.NewLogger()
	synerexAddr = os.Getenv("SYNEREX_SERVER")
	if synerexAddr == "" {
		synerexAddr = "127.0.0.1:10000"
	}
	nodeIdAddr = os.Getenv("NODEID_SERVER")
	if nodeIdAddr == "" {
		nodeIdAddr = "127.0.0.1:9000"
	}
	masterSynerexAddr = os.Getenv("MASTER_SYNEREX_SERVER")
	if masterSynerexAddr == "" {
		masterSynerexAddr = "127.0.0.1:10000"
	}
	masterNodeIdAddr = os.Getenv("MASTER_NODEID_SERVER")
	if masterNodeIdAddr == "" {
		masterNodeIdAddr = "127.0.0.1:9000"
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
////////////            Message Class           ///////////
///////////////////////////////////////////////////////////

type Message struct {
	ready     chan struct{}
	agents    []*api.Agent
	senderIds []uint64
	targets   []uint64
}

func NewMessage(targets []uint64) *Message {
	return &Message{ready: make(chan struct{}), agents: make([]*api.Agent, 0), senderIds: []uint64{}, targets: targets}
}

func (m *Message) IsFinish() bool {
	for _, tgt := range m.targets {
		isExist := false
		for _, sid := range m.senderIds {
			if tgt == sid {
				isExist = true
			}
		}
		if isExist == false {
			return false
		}
	}
	return true
}

func (m *Message) Set(a []*api.Agent, senderId uint64) {
	m.agents = append(m.agents, a...)
	m.senderIds = append(m.senderIds, senderId)
	if m.IsFinish() {
		close(m.ready)
	}
}

func (m *Message) Get() []*api.Agent {
	select {
	case <-m.ready:
		//case <-time.After(100 * time.Millisecond):
		//	logger.Warn("Timeout Get")
	}

	return m.agents
}

////////////////////////////////////////////////////////////
////////////           Harmovis server           ///////////
///////////////////////////////////////////////////////////

func runServer() *gosocketio.Server {

	currentRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	d := filepath.Join(currentRoot, "monitor", "build")

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

/////////////////////////////////////////
////////  Send Agent to HamowareVIS /////
////////////////////////////////////////
// GetJson: json化する関数
func (m *MapMarker) GetJson() string {
	s := fmt.Sprintf("{\"mtype\":%d,\"id\":%d,\"lat\":%f,\"lon\":%f,\"angle\":%f,\"speed\":%d,\"area\":%d}",
		m.mtype, m.id, m.lat, m.lon, m.angle, m.speed, m.area)
	return s
}

// sendAgentToHarmowareVis: harmowareVisに情報を送信する関数
func sendAgentToHarmowareVis(agents []*api.Agent) {

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
		ioserv.BroadcastToAll("agents", jsonAgents)
		mu.Unlock()
	}
}

/////////////////////////////////////////
////////  Send Agent to HamowareVIS /////
////////////////////////////////////////

// sendAreaToHarmowareVis: harmowareVisに情報を送信する関数
func sendAreaToHarmowareVis(areas []*api.Area) {

	if areas != nil {
		jsonAreas := make([]string, 0)
		for _, areaInfo := range areas {
			areaJson, _ := json.Marshal(areaInfo)
			jsonAreas = append(jsonAreas, string(areaJson))
		}
		mu.Lock()
		ioserv.BroadcastToAll("areas", jsonAreas)
		mu.Unlock()
	}

}

// callbackForwardClockRequest: クロックを進める関数
func forwardClock(dm *api.Demand) {
	/*// Databaseへ取得する
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_DATABASE,
	})
	//uid, _ := uuid.NewRandom()
	senderId := myProvider.Id
	sps, _ := visapi.GetAgentRequest(senderId, targets)
	//sps, _ := waiter.WaitSp(msgId, targets, 1000)

	allAgents := []*api.Agent{}
	for _, sp := range sps {
		agents := sp.GetSimSupply().GetGetAgentResponse().GetAgents()
		allAgents = append(allAgents, agents...)
	}*/

	agents := agentsMessage.Get()

	// Harmowareに送る
	sendAgentToHarmowareVis(agents)

	logger.Info("Agents %v\n", len(agents))
}

// callback for each Supply
func demandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	switch dm.GetSimDemand().GetType() {

	case api.DemandType_REGIST_PROVIDER_REQUEST:
		// providerを追加する
		p := dm.GetSimDemand().GetRegistProviderRequest().GetProvider()
		pm.AddProvider(p)
		fmt.Printf("regist request from agent of vis provider! %v\n", p)
		// 登録完了通知
		senderId := myProvider.Id
		targets := []uint64{p.GetId()}
		msgId := dm.GetSimDemand().GetMsgId()
		visapi.RegistProviderResponse(senderId, targets, msgId, myProvider)

		logger.Info("Success Regist Agent or Vis Providers", targets)

	case api.DemandType_SET_AGENT_REQUEST:

		agents := dm.GetSimDemand().GetSetAgentRequest().GetAgents()
		logger.Info("get Agents: %v\n", len(agents))
		//sendAgentToHarmowareVis(agents)
		agentsMessage.Set(agents, dm.GetSimDemand().GetSenderId())
		//db.Push(agents)
		// response
		pId := myProvider.Id
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		visapi.SetAgentResponse(pId, targets, msgId)

	}

}

// callback for each Supply
func supplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	switch sp.GetSimSupply().GetType() {
	}
}

////////////////////////////////////////////////////////////
////////////     Demand Supply Callback     ////////////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func masterSupplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_REGIST_PROVIDER_RESPONSE:
		mu.Lock()
		masterProvider = sp.GetSimSupply().GetRegistProviderResponse().GetProvider()
		mu.Unlock()
		fmt.Printf("regist provider to Master Provider!\n")

	}
}

// Demandのコールバック関数
func masterDemandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	senderId := myProvider.Id
	switch dm.GetSimDemand().GetType() {

	case api.DemandType_FORWARD_CLOCK_REQUEST:
		//fmt.Printf("get forwardClockRequest")
		t1 := time.Now()
		forwardClock(dm)
		t2 := time.Now()
		duration := t2.Sub(t1).Milliseconds()
		logger.Info("Duration: %v, PID: %v", duration, myProvider.Id)
		// response to master
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		logger.Debug("Response to master pid %v, msgId%v\n", myProvider.Id, msgId)
		masterapi.ForwardClockResponse(senderId, targets, msgId)

	case api.DemandType_FORWARD_CLOCK_INIT_REQUEST:
		targets := pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_AGENT,
		})
		agentsMessage = NewMessage(targets)

		// response
		targets = []uint64{dm.GetSimDemand().GetSenderId()}
		msgId := dm.GetSimDemand().GetMsgId()
		masterapi.ForwardClockInitResponse(senderId, targets, msgId)
		logger.Info("Finish: Forward Clock Init")

	case api.DemandType_UPDATE_PROVIDERS_REQUEST:
		providers := dm.GetSimDemand().GetUpdateProvidersRequest().GetProviders()
		//pm.SetProviders(providers)

		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		masterapi.UpdateProvidersResponse(senderId, targets, msgId)
		logger.Info("Finish: Update Workers num: %v\n", len(providers))

	case api.DemandType_SEND_AREA_INFO_REQUEST:
		areas := dm.GetSimDemand().GetSendAreaInfoRequest().GetAreas()
		sendAreaToHarmowareVis(areas)
		// response
		targets := []uint64{dm.GetSimDemand().GetSenderId()}
		senderId := myProvider.Id
		msgId := dm.GetSimDemand().GetMsgId()
		masterapi.SendAreaInfoResponse(senderId, targets, msgId)
		logger.Info("Finish: GetAreaInfo: %v\n", areas)
	}
}

func registToMaster() {
	// masterへ登録
	senderId := myProvider.Id
	targets := make([]uint64, 0)
	masterapi.RegistProviderRequest(senderId, targets, myProvider)

	go func() {
		for {
			if masterProvider != nil {
				logger.Debug("Regist Success to Master!")
				return
			} else {
				logger.Debug("Couldn't Regist Master...Retry...\n")
				time.Sleep(2 * time.Second)
				// masterへ登録
				masterapi.RegistProviderRequest(senderId, targets, myProvider)
			}
		}
	}()
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
		Type: api.ProviderType_VISUALIZATION,
	}
	pm = simutil.NewProviderManager(myProvider)

	// For Master
	// Connect to Node Server
	nodeapi1 := napi.NewNodeAPI()
	for {
		err := nodeapi1.RegisterNodeName(masterNodeIdAddr, providerName, false)
		if err == nil {
			logger.Info("connected NodeID server!")
			go nodeapi1.HandleSigInt()
			nodeapi1.RegisterDeferFunction(nodeapi1.UnRegisterNode)
			break
		} else {
			logger.Warn("NodeID Error... reconnecting...")
			time.Sleep(2 * time.Second)
		}
	}

	// for visNodeID
	// Connect to Node Server
	nodeapi2 := napi.NewNodeAPI()
	for {
		err := nodeapi2.RegisterNodeName(nodeIdAddr, providerName, false)
		if err == nil {
			logger.Info("connected NodeID server!")
			go nodeapi2.HandleSigInt()
			nodeapi2.RegisterDeferFunction(nodeapi2.UnRegisterNode)
			break
		} else {
			logger.Warn("NodeID Error... reconnecting...")
			time.Sleep(2 * time.Second)
		}
	}

	// Connect to Synerex Server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(masterSynerexAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	nodeapi1.RegisterDeferFunction(func() { conn.Close() })
	client := api.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Visualization}")

	// Connect to Synerex Server
	var wopts []grpc.DialOption
	wopts = append(wopts, grpc.WithInsecure())
	wconn, werr := grpc.Dial(synerexAddr, wopts...)
	if werr != nil {
		log.Fatalf("fail to dial: %v", werr)
	}
	nodeapi2.RegisterDeferFunction(func() { wconn.Close() })
	wclient := api.NewSynerexClient(wconn)
	wargJson := fmt.Sprintf("{Client:Visualization}")

	// Communicator
	masterapi = api.NewSimAPI()
	masterapi.RegistClients(client, myProvider.Id, argJson)            // channelごとのClientを作成
	masterapi.SubscribeAll(masterDemandCallback, masterSupplyCallback) // ChannelにSubscribe

	// Communicator
	visapi = api.NewSimAPI()
	visapi.RegistClients(wclient, myProvider.Id, wargJson) // channelごとのClientを作成
	visapi.SubscribeAll(demandCallback, supplyCallback)    // ChannelにSubscribe

	time.Sleep(5 * time.Second)

	registToMaster()

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
	nodeapi1.CallDeferFunctions() // cleanup!
	nodeapi2.CallDeferFunctions() // cleanup!
}
