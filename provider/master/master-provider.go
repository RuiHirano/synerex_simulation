package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	logger.SetPrefix("Master")
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
//////////////////        Manager          ///////////////////
///////////////////////////////////////////////////////////

type Manager struct {
	Providers []*api.Provider
}

func NewManager() *Manager {
	m := &Manager{
		Providers: make([]*api.Provider, 0),
	}
	return m
}

func (m *Manager) AddProvider(provider *api.Provider) {
	m.Providers = append(m.Providers, provider)
}

func (m *Manager) GetProviderIds() []uint64 {
	pids := make([]uint64, 0)
	for _, p := range m.Providers {
		pids = append(pids, p.GetId())
	}
	return pids
}

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
		//logger.Info("get sp: %v\n", sp)
		//time.Sleep(10 * time.Millisecond)
		simapi.SendSpToWait(sp)
	case api.SupplyType_SET_AGENT_RESPONSE:
		//logger.Info("get sp: %v\n", sp)
		//time.Sleep(10 * time.Millisecond)
		simapi.SendSpToWait(sp)
	case api.SupplyType_FORWARD_CLOCK_RESPONSE:
		//logger.Info("get sp: %v\n", sp)
		//time.Sleep(10 * time.Millisecond)
		simapi.SendSpToWait(sp)
	case api.SupplyType_UPDATE_PROVIDERS_RESPONSE:
		//logger.Info("get sp: %v\n", sp)
		//time.Sleep(10 * time.Millisecond)
		simapi.SendSpToWait(sp)
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
		//providerManager.AddProvider(p)
		pm.AddProvider(p)
		fmt.Printf("regist provider! %v\n", p.GetId())
		// 登録完了通知
		//targets := []uint64{tid}
		senderInfo := myProvider.Id
		targets := []uint64{p.GetId()}
		msgId := dm.GetSimDemand().GetMsgId()
		simapi.RegistProviderResponse(senderInfo, targets, msgId, pm.MyProvider)

		// update provider to worker
		targets = pm.GetProviderIds([]simutil.IDType{
			simutil.IDType_WORKER,
		})
		simapi.UpdateProvidersRequest(senderInfo, targets, pm.GetProviders())

		logger.Info("Success Update Providers", targets)

	}
}

// setAgents: agentをセットするDemandを出す関数
func setAgents(agentNum uint64) (bool, error) {

	agents := make([]*api.Agent, 0)

	minLon, maxLon, minLat, maxLat := 136.971626, 136.989379, 35.152210, 35.161499
	for i := 0; i < int(agentNum); i++ {
		uid, _ := uuid.NewRandom()
		position := &api.Coord{
			Longitude: minLon + (maxLon-minLon)*rand.Float64(),
			Latitude:  minLat + (maxLat-minLat)*rand.Float64(),
		}
		/*departure := &api.Coord{
			Longitude: 136.975685 + rand.Float64()*0.001,
			Latitude:  35.154533 + rand.Float64()*0.001,
		}*/
		destination := &api.Coord{
			Longitude: minLon + (maxLon-minLon)*rand.Float64(),
			Latitude:  minLat + (maxLat-minLat)*rand.Float64(),
		}
		transitPoints := []*api.Coord{destination}
		agents = append(agents, &api.Agent{
			Type: api.AgentType_PEDESTRIAN,
			Id:   uint64(uid.ID()),
			Route: &api.Route{
				Position:      position,
				Direction:     30,
				Speed:         60,
				Departure:     position,
				Destination:   destination,
				TransitPoints: transitPoints,
				NextTransit:   destination,
			},
		})
	}

	// エージェントを設置するリクエスト
	senderId := myProvider.Id
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_WORKER,
	})
	simapi.SetAgentRequest(senderId, targets, agents)
	//waiter.WaitSp(msgId, targets, 1000)

	logger.Info("Finish Setting Agents \n Add: %v", len(agents))
	return true, nil
}

// startClock:
func startClock() {
	t1 := time.Now()

	senderId := myProvider.Id
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_WORKER,
	})
	logger.Debug("Next Cycle! \n%v\n", targets)
	simapi.ForwardClockRequest(senderId, targets)
	//waiter.WaitSp(msgId, targets, 1000)

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

	e.Start(":" + port)
}

func main() {

	// ProviderManager
	uid, _ := uuid.NewRandom()
	myProvider = &api.Provider{
		Id:   uint64(uid.ID()),
		Name: "MasterServer",
		Type: api.ProviderType_MASTER,
	}
	pm = simutil.NewProviderManager(myProvider)

	// CLI, GUIの受信サーバ
	go startSimulatorServer()

	// Connect to Node Server
	nodeapi := napi.NewNodeAPI()
	for {
		err := nodeapi.RegisterNodeName(nodeIdAddr, "MasterProvider", false)
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
	argJson := fmt.Sprintf("{Client:Master}")

	time.Sleep(3 * time.Second)

	// api
	fmt.Printf("client: %v\n", client)
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, myProvider.Id, argJson) // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback)  // ChannelにSubscribe*/
	//logger.Info("Connected Synerex Server!\n")

	// ready provider request
	//senderId := myProvider.Id
	//targets := make([]uint64, 0)
	//simapi.ReadyProviderRequest(senderId, targets, myProvider)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	nodeapi.CallDeferFunctions() // cleanup!

}
