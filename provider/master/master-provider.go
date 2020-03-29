package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	gosocketio "github.com/mtfelian/golang-socketio"
	api "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/provider/simutil"

	//"github.com/synerex/synerex_alpha/sxutil"
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
	startFlag     bool
	masterClock   int
	workerHosts   []string
	mu            sync.Mutex
	simapi        *api.SimAPI
	//sim             *Simulator
	//providerManager *simutil.ProviderManager
	//areaManager *simutil.AreaManager
	//pSources map[provider.ProviderType]*provider.Source
	logger *simutil.Logger
)

func init() {
	startFlag = false
	masterClock = 0
	workerHosts = make([]string, 0)
	logger = simutil.NewLogger()
	logger.SetPrefix("Master")
	flag.Parse()
}

var (
	//fcs *geojson.FeatureCollection
	//geofile string
	port      = 9995
	assetsDir http.FileSystem
	server    *gosocketio.Server = nil
)

////////////////////////////////////////////////////////////
//////////////////        Util          ///////////////////
///////////////////////////////////////////////////////////

/*func createRandomCoord(areaCoords []*api.Coord) *api.Coord {

	maxLat, maxLon, minLat, minLon := simutil.GetCoordRange(areaCoords)
	longitude := minLon + (maxLon-minLon)*rand.Float64()
	latitude := minLat + (maxLat-minLat)*rand.Float64()
	coord := &api.Coord{
		Longitude: longitude,
		Latitude:  latitude,
	}

	return coord
}*/

// Agentオブジェクトの変換
func calcRoute() *api.Route {

	//departure := createRandomCoord(mockAreaInfos[*areaId].ControlArea)
	departure := &api.Coord{Latitude: 35.00, Longitude: 136.234}
	//destAreaId := rand.Intn(3)
	destination := &api.Coord{Latitude: 35.00, Longitude: 136.234}
	//destination := createRandomCoord(mockAreaInfos[uint64(destAreaId)].ControlArea)

	/*departure = &api.Coord{
		Latitude:  35.1542,
		Longitude: 136.975231,
	}
	destination = &api.Coord{
		Latitude:  35.1542,
		Longitude: 136.975231,
	}*/

	transitPoints := make([]*api.Coord, 0)
	transitPoints = append(transitPoints, destination)

	route := &api.Route{
		Position:      departure,
		Direction:     0.0001 * rand.Float64(),
		Speed:         10 + 10*rand.Float64(),
		Departure:     departure,
		Destination:   destination,
		TransitPoints: transitPoints,
		NextTransit:   destination,
	}

	return route
}

////////////////////////////////////////////////////////////
////////////     Demand Supply Callback     ////////////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func supplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	// 自分宛かどうか
	// check if supply is match with my demand.
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_SET_CLOCK_RESPONSE:
		//api.SendToWaitCh(sp, sp.GetSimSupply().GetType())
	case api.SupplyType_SET_AGENT_RESPONSE:
		//api.SendToWaitCh(sp, sp.GetSimSupply().GetType())
	case api.SupplyType_START_CLOCK_RESPONSE:
		//api.SendToWaitCh(sp, sp.GetSimSupply().GetType())
	case api.SupplyType_STOP_CLOCK_RESPONSE:
		//api.SendToWaitCh(sp, sp.GetSimSupply().GetType())

	}
}

// Demandのコールバック関数
func demandCallback(clt *api.SMServiceClient, dm *api.Demand) {
	//tid := dm.GetSimDemand().GetSenderId()
	//pid := myProvider.Id
	// check if supply is match with my demand.
	switch dm.GetSimDemand().GetType() {
	case api.DemandType_REGIST_PROVIDER_REQUEST:
		// providerを追加する
		p := dm.GetSimDemand().GetRegistProviderRequest().GetProvider()
		fmt.Printf("regist provider! %v\n", p)
		// 登録完了通知
		//targets := []uint64{tid}
		senderInfo := myProvider.Id
		simapi.RegistProviderResponse(senderInfo)

		logger.Info("Success Update Providers")

	}
}

// setAgents: agentをセットするDemandを出す関数
func setAgents(agentNum uint64) (bool, error) {

	agents := make([]*api.Agent, 0)

	for i := 0; i < int(agentNum); i++ {
		uuid, err := uuid.NewRandom()
		if err == nil {
			agent := &api.Agent{
				Id:    uint64(uuid.ID()),
				Type:  api.AgentType_PEDESTRIAN,
				Route: calcRoute(),
			}
			agents = append(agents, agent)
		}
	}

	// エージェントを設置するリクエスト
	senderId := myProvider.Id
	simapi.SetAgentRequest(senderId, agents)

	logger.Info("Finish Setting Agents \n Add: %v", len(agents))
	return true, nil
}

// startClock:
func startClock() {

	t1 := time.Now()

	senderId := myProvider.Id
	simapi.ForwardClockRequest(senderId)

	// calc next time
	masterClock++
	log.Printf("\x1b[30m\x1b[47m \n Finish: Clock forwarded \n Time:  %v \x1b[0m\n", masterClock)

	t2 := time.Now()
	duration := t2.Sub(t1).Milliseconds()
	logger.Info("Duration: %v", duration)
	if duration > 1000 {
		logger.Error("time cycle delayed...")
	} else {
		// 待機
		logger.Info("wait %v ms", 1000-duration)
		time.Sleep(time.Duration(1000-duration) * time.Millisecond)
	}

	// 次のサイクルを行う
	if startFlag {
		startClock()
	} else {
		log.Printf("\x1b[30m\x1b[47m \n Finish: Clock stopped \n GlobalTime:  %v \x1b[0m\n", masterClock)
		startFlag = false
		return
	}

}

type ClockOptions struct {
	Time int `validate:"required,min=0" json:"time"`
}

func orderSetClock() echo.HandlerFunc {
	return func(c echo.Context) error {
		co := new(ClockOptions)
		if err := c.Bind(co); err != nil {
			return err
		}
		fmt.Printf("time %d\n", co.Time)
		masterClock = co.Time
		return c.String(http.StatusOK, "Set Clock")
	}
}

type AgentOptions struct {
	Num int `validate:"required,min=0,max=10", json:"num"`
}

func orderSetAgent() echo.HandlerFunc {
	return func(c echo.Context) error {
		ao := new(AgentOptions)
		if err := c.Bind(ao); err != nil {
			return err
		}
		fmt.Printf("agent num %d\n", ao.Num)
		setAgents(uint64(ao.Num))
		return c.String(http.StatusOK, "Set Agent")
	}
}

func orderStart() echo.HandlerFunc {
	return func(c echo.Context) error {
		if startFlag == false {
			startFlag = true
			go startClock()
			return c.String(http.StatusOK, "Start")
		} else {
			logger.Warn("Clock is already started.")
			return c.String(http.StatusBadRequest, "Start")
		}
	}
}

func orderStop() echo.HandlerFunc {
	return func(c echo.Context) error {
		startFlag = false
		return c.String(http.StatusOK, "Stop")
	}
}

func startSimulatorServer() {
	fmt.Printf("Starting Simulator Server...")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/order/set/clock", orderSetClock())
	e.POST("/order/set/agent", orderSetAgent())
	e.POST("/order/start", orderStart())
	e.POST("/order/stop", orderStop())

	e.Start(":8000")
}

func main() {

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: "MasterServer",
		Type: api.ProviderType_MASTER,
	}
	//providerManager = simutil.NewProviderManager(myProvider)
	//providerManager.CreateIDMap()

	//AreaManager
	//areaManager = simutil.NewAreaManager(mockAreaInfos[*areaId])

	// CLI, GUIの受信サーバ
	go startSimulatorServer()
	//simulatorServer := NewSimulatorServer()
	//simulatorServer.Run()

	// 初期プロバイダ起動
	//runInitServer()

	// Connect to Node Server
	api.RegisterNodeName(*nodeIdAddr, "MasterProvider", false)
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
	argJson := fmt.Sprintf("{Client:Master}")

	// api
	fmt.Printf("client: %v\n", client)
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, argJson)               // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback) // ChannelにSubscribe*/

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	api.CallDeferFunctions() // cleanup!

}
