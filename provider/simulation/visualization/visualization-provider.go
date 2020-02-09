package main

import (
	"flag"
	"log"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	simapi "github.com/synerex/synerex_alpha/api/simulation"
	"github.com/synerex/synerex_alpha/api/simulation/agent"
	"github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	"github.com/synerex/synerex_alpha/api/simulation/provider"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	"github.com/synerex/synerex_alpha/sxutil"

	"fmt"
	"net/http"
	"os"
	"path/filepath"

	gosocketio "github.com/mtfelian/golang-socketio"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	pidFlag    = flag.Int("pid", 2, "Provider Id")
	port       = flag.Int("port", 10080, "HarmoVis Provider Listening Port")
	version    = "0.01"
	mu         sync.Mutex
	assetsDir  http.FileSystem
	ioserv     *gosocketio.Server
	com        *simutil.Communicator
	sim        *Simulator
	pid        uint64
)

func init() {
	flag.Parse()
	pid = uint64(*pidFlag)
}

func sendAreaToHarmowareVis(areas []*area.Area) {
	jsonAreas := make([]string, 0)
	jsonAreas = append(jsonAreas, "test")
	mu.Lock()
	ioserv.BroadcastToAll("area", jsonAreas)
	mu.Unlock()
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
func sendToHarmowareVis(sumAgents []*agent.Agent) {

	if sumAgents != nil {
		jsonAgents := make([]string, 0)
		for _, agentInfo := range sumAgents {

			// agentInfoTypeによってエージェントを取得
			switch agentInfo.Type {
			case agent.AgentType_PEDESTRIAN:
				ped := agentInfo.GetPedestrian()
				mm := &MapMarker{
					mtype: int32(agentInfo.Type),
					id:    int32(agentInfo.Id),
					lat:   float32(ped.Route.Position.Latitude),
					lon:   float32(ped.Route.Position.Longitude),
					angle: float32(ped.Route.Direction),
					speed: int32(ped.Route.Speed),
				}
				jsonAgents = append(jsonAgents, mm.GetJson())

			case agent.AgentType_CAR:
				car := agentInfo.GetCar()
				mm := &MapMarker{
					mtype: int32(agentInfo.Type),
					id:    int32(agentInfo.Id),
					lat:   float32(car.Route.Position.Latitude),
					lon:   float32(car.Route.Position.Longitude),
					angle: float32(car.Route.Direction),
					speed: int32(car.Route.Speed),
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
func forwardClock(dm *pb.Demand) {
	log.Printf("\x1b[30m\x1b[47m \n Start: Clock forwarded \n Time:  %v \x1b[0m\n", sim.Clock.GlobalTime)
	targetId := dm.GetId()

	_, agents := com.GetAgentsRequest(nil)

	// Harmowareに送る
	sendToHarmowareVis(agents)

	// clockを進める
	sim.ForwardStep()

	// セット完了通知を送る
	com.ForwardClockResponse(targetId)
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock forwarded \n Time:  %v \x1b[0m\n", sim.Clock.GlobalTime)
}

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

		//sendAreaToHarmowareVis(make([]*area.Area2, 0))
		// geojsonを送信
		//sendFile2()
		//sendFile()
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

// callback for each Supply
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	//tid := dm.GetId()
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_UPDATE_PROVIDERS_REQUEST:
		// 参加者リストをセットする要求
		//callbackSetParticipantsRequest(dm)
	case simapi.DemandType_FORWARD_CLOCK_REQUEST:
		// クロックを進める
		forwardClock(dm)

	case simapi.DemandType_UPDATE_CLOCK_REQUEST:
		// Clockをセットする
		clock := dm.GetSimDemand().GetSetClockRequest().GetClock()
		sim.Clock = clock
		log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.Clock.GlobalTime, sim.Clock.TimeStep)
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

	// Connect to Node Server
	sxutil.RegisterNodeName(*nodesrv, "VisualizationProvider", false)
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
	argJson := fmt.Sprintf("{Client:Visualization}")

	// Simulator
	clockInfo := clock.NewClock(0, 1, 1)
	sim = NewSimulator(clockInfo)

	// Communicator
	visInfo := &provider.Visualization{}
	provider := provider.NewVisualizationProvider("VisualizationProvider", provider.ProviderType_VISUALIZATION, visInfo)
	com = simutil.NewCommunicator(provider)
	com.RegistClients(client, argJson)               // channelごとのClientを作成
	com.SubscribeAll(demandCallback, supplyCallback) // ChannelにSubscribe

	// Run HarmowareVis Monitor
	ioserv = runServer()
	log.Printf("Running Sio Server..\n")
	if ioserv == nil {
		os.Exit(1)
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", ioserv)
	serveMux.HandleFunc("/", assetsFileHandler)
	log.Printf("Starting Harmoware VIS  Provider %s  on port %d", version, *port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), serveMux)
	if err != nil {
		log.Fatal(err)
	}

	// プロバイダのsetup
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 新規参加登録
	com.RegistProviderRequest(nil, provider)
	log.Printf("\x1b[30m\x1b[47m \n Finish: This provider registered in scenario-provider \x1b[0m\n")

	// Clock情報を取得
	_, clockInfo = com.GetClockRequest(nil)
	sim.Clock = clockInfo
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock information set. \n GlobalTime:  %v \n TimeStep: %v \x1b[0m\n", sim.Clock.GlobalTime, sim.Clock.TimeStep)
	wg.Wait()

	sxutil.CallDeferFunctions() // cleanup!
}
