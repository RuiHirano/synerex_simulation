package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"io/ioutil"
	"os/exec"

	"github.com/go-yaml/yaml"
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
	config *Config
	podgen *PodGenerator
	proc   *Processor
)

type Config struct {
	Area Config_Area `yaml:"area"`
}

type Config_Area struct {
	SideRange      float64        `yaml:"sideRange"`
	DuplicateRange float64        `yaml:"duplicateRange"`
	DefaultAreaNum Config_AreaNum `yaml:"defaultAreaNum"`
}
type Config_AreaNum struct {
	Row    uint64 `yaml:"row"`
	Column uint64 `yaml:"column"`
}

func readConfig() (*Config, error) {
	var config *Config
	buf, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err)
		return config, err
	}
	// []map[string]string のときと使う関数は同じです。
	// いい感じにマッピングしてくれます。
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		fmt.Println(err)
		return config, err
	}
	fmt.Printf("yaml is %v\n", config)
	return config, nil
}

func init() {
	podgen = NewPodGenerator()
	proc = NewProcessor()
	waiter = api.NewWaiter()
	startFlag = false
	masterClock = 0
	workerHosts = make([]string, 0)
	logger = util.NewLogger()
	logger.SetPrefix("Master")
	flag.Parse()
	//providerManager = NewManager()
	// configを読み取る
	config, _ = readConfig()

	// kubetest
	/*id := "test"
	area := &Area{
		Id:        3,
		Control:   []Coord{{Latitude: 0, Longitude: 0}, {Latitude: 10, Longitude: 0}, {Latitude: 10, Longitude: 10}, {Latitude: 0, Longitude: 10}},
		Duplicate: []Coord{{Latitude: 0, Longitude: 0}, {Latitude: 10, Longitude: 0}, {Latitude: 10, Longitude: 10}, {Latitude: 0, Longitude: 10}},
	}
	go podgen.applyWorker(id, area)
	time.Sleep(4 * time.Second)
	go podgen.deleteWorker(id)*/
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
////////////     Demand Supply Callback     ////////////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func supplyCallback(clt *api.SMServiceClient, sp *api.Supply) {
	// 自分宛かどうか
	// check if supply is match with my demand.
	switch sp.GetSimSupply().GetType() {
	case api.SupplyType_SET_CLOCK_RESPONSE:
		simapi.SendSpToWait(sp)
	case api.SupplyType_SET_AGENT_RESPONSE:
		simapi.SendSpToWait(sp)
	case api.SupplyType_FORWARD_CLOCK_RESPONSE:
		simapi.SendSpToWait(sp)
	case api.SupplyType_FORWARD_CLOCK_INIT_RESPONSE:
		simapi.SendSpToWait(sp)
	case api.SupplyType_UPDATE_PROVIDERS_RESPONSE:
		simapi.SendSpToWait(sp)
	case api.SupplyType_SEND_AREA_INFO_RESPONSE:
		simapi.SendSpToWait(sp)
	}
}

// Demandのコールバック関数
func demandCallback(clt *api.SMServiceClient, dm *api.Demand) {

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

		logger.Info("Success Update Providers! Worker Num: ", len(targets))

	}
}

///////////////////////////////////////////////
////////////  Processor  //////////////////////
///////////////////////////////////////////////

type Processor struct {
	Area        *api.Area            // 全体のエリア
	AreaMap     map[string]*api.Area // [areaid] []areaCoord     エリア情報を表したmap
	NeighborMap map[string][]string  // [areaid] []neighborAreaid   隣接関係を表したmap
}

func NewProcessor() *Processor {
	proc := &Processor{
		Area:        nil,
		AreaMap:     make(map[string]*api.Area),
		NeighborMap: make(map[string][]string),
	}
	return proc
}

// setAgents: agentをセットするDemandを出す関数
func (proc *Processor) setAgents(agentNum uint64) (bool, error) {

	if proc.Area == nil {
		return false, fmt.Errorf("area is nil")
	}

	agents := make([]*api.Agent, 0)
	//minLon, maxLon, minLat, maxLat := 136.971626, 136.989379, 35.152210, 35.161499
	maxLat, maxLon, minLat, minLon := GetCoordRange(proc.Area.ControlArea)
	fmt.Printf("minLon %v, maxLon %v, minLat %v, maxLat %v\n", minLon, maxLon, minLat, maxLat)
	for i := 0; i < int(agentNum); i++ {
		uid, _ := uuid.NewRandom()
		position := &api.Coord{
			Longitude: minLon + (maxLon-minLon)*rand.Float64(),
			Latitude:  minLat + (maxLat-minLat)*rand.Float64(),
		}
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
		fmt.Printf("position %v\n", position)
	}

	// エージェントを設置するリクエスト
	senderId := myProvider.Id
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_WORKER,
	})
	simapi.SetAgentRequest(senderId, targets, agents)

	logger.Info("Finish Setting Agents \n Add: %v", len(agents))
	return true, nil
}

// setAreas: areaをセットするDemandを出す関数
func (proc *Processor) setAreas(areaCoords []*api.Coord) (bool, error) {

	proc.Area = &api.Area{
		Id:            0,
		ControlArea:   areaCoords,
		DuplicateArea: areaCoords,
	}
	//id := "test"

	areas, neighborsMap := proc.divideArea(areaCoords, config.Area)
	for _, area := range areas {
		neighbors := neighborsMap[int(area.Id)]
		go podgen.applyWorker(area, neighbors)
		//defer podgen.deleteWorker(id) // not working...
	}

	// send area info to visualization
	senderId := myProvider.Id
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_VISUALIZATION,
	})
	logger.Debug("Send Area Info to Vis! \n%v\n", targets)
	//areas := []*api.Area{proc.Area}
	simapi.SendAreaInfoRequest(senderId, targets, areas)

	return true, nil
}

// startClock:
func (proc *Processor) startClock() {
	t1 := time.Now()

	senderId := myProvider.Id
	targets := pm.GetProviderIds([]simutil.IDType{
		simutil.IDType_WORKER,
		simutil.IDType_VISUALIZATION,
	})
	logger.Debug("Next Cycle! \n%v\n", targets)
	simapi.ForwardClockInitRequest(senderId, targets)
	simapi.ForwardClockRequest(senderId, targets)

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
		proc.startClock()
	} else {
		log.Printf("\x1b[30m\x1b[47m \n Finish: Clock stopped \n GlobalTime:  %v \x1b[0m\n", masterClock)
		startFlag = false
		return
	}

}

// areaをrow、columnに分割する関数
func (proc *Processor) divideArea(areaCoords []*api.Coord, areaConfig Config_Area) ([]*api.Area, map[int][]string) {
	row := areaConfig.DefaultAreaNum.Row
	column := areaConfig.DefaultAreaNum.Column
	dupRange := areaConfig.DuplicateRange
	areas := []*api.Area{}
	neighborMap := make(map[int][]string)

	maxLat, maxLon, minLat, minLon := GetCoordRange(proc.Area.ControlArea)
	//areaId := 0
	for c := 0; c < int(column); c++ {
		// calc slon, elon
		slon := minLon + (maxLon-minLon)*float64(c)/float64(column)
		elon := minLon + (maxLon-minLon)*float64((c+1))/float64(column)
		for r := 0; r < int(row); r++ {
			areaId := strconv.Itoa(c+1) + strconv.Itoa(r+1)
			areaIdint, _ := strconv.Atoi(strconv.Itoa(c+1) + strconv.Itoa(r+1))
			// calc slat, elat
			slat := minLat + (maxLat-minLat)*float64(r)/float64(row)
			elat := minLat + (maxLat-minLat)*float64((r+1))/float64(row)
			fmt.Printf("test id %v\n", areaId)
			areas = append(areas, &api.Area{
				Id: uint64(areaIdint),
				ControlArea: []*api.Coord{
					{Latitude: slat, Longitude: slon},
					{Latitude: slat, Longitude: elon},
					{Latitude: elat, Longitude: elon},
					{Latitude: elat, Longitude: slon},
				},
				DuplicateArea: []*api.Coord{
					{Latitude: slat - dupRange, Longitude: slon - dupRange},
					{Latitude: slat - dupRange, Longitude: elon + dupRange},
					{Latitude: elat + dupRange, Longitude: elon + dupRange},
					{Latitude: elat + dupRange, Longitude: slon - dupRange},
				},
			})

			// add neighbors 各エリアの右と上を作成すれば全体を満たす
			if c+2 <= int(column) {
				id := strconv.Itoa(c+2) + strconv.Itoa(r+1)
				neighborMap[areaIdint] = append(neighborMap[areaIdint], id)
			}
			if r+2 <= int(row) {
				id := strconv.Itoa(c+1) + strconv.Itoa(r+2)
				neighborMap[areaIdint] = append(neighborMap[areaIdint], id)
			}
		}
	}

	return areas, neighborMap
}

///////////////////////////////////////////////
////////////  Order  //////////////////////
///////////////////////////////////////////////

type Order struct {
}

func NewOrder() *Order {
	order := &Order{}
	return order
}

type ClockOptions struct {
	Time int `validate:"required,min=0" json:"time"`
}

func (or *Order) SetClock() echo.HandlerFunc {
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

func (or *Order) SetAgent() echo.HandlerFunc {
	return func(c echo.Context) error {
		ao := new(AgentOptions)
		if err := c.Bind(ao); err != nil {
			return err
		}
		fmt.Printf("agent num %d\n", ao.Num)
		ok, err := proc.setAgents(uint64(ao.Num))
		fmt.Printf("ok %v, err %v", ok, err)
		return c.String(http.StatusOK, "Set Agent")
	}
}

type AreaOptions struct {
	SLat string `min=0,max=100", json:"slat"`
	SLon string `min=0,max=200", json:"slon"`
	ELat string `min=0,max=100", json:"elat"`
	ELon string `min=0,max=200", json:"elon"`
}

func (or *Order) SetArea() echo.HandlerFunc {
	return func(c echo.Context) error {
		ao := new(AreaOptions)
		if err := c.Bind(ao); err != nil {
			return err
		}
		fmt.Printf("area %d\n", ao)
		slat, _ := strconv.ParseFloat(ao.SLat, 64)
		slon, _ := strconv.ParseFloat(ao.SLon, 64)
		elat, _ := strconv.ParseFloat(ao.ELat, 64)
		elon, _ := strconv.ParseFloat(ao.ELon, 64)
		area := []*api.Coord{
			{Latitude: slat, Longitude: slon},
			{Latitude: slat, Longitude: elon},
			{Latitude: elat, Longitude: elon},
			{Latitude: elat, Longitude: slon},
		}
		proc.setAreas(area)
		return c.String(http.StatusOK, "Set Area")
	}
}

func (or *Order) Start() echo.HandlerFunc {
	return func(c echo.Context) error {
		if startFlag == false {
			startFlag = true
			go proc.startClock()
			return c.String(http.StatusOK, "Start")
		} else {
			logger.Warn("Clock is already started.")
			return c.String(http.StatusBadRequest, "Start")
		}
	}
}

func (or *Order) Stop() echo.HandlerFunc {
	return func(c echo.Context) error {
		startFlag = false
		return c.String(http.StatusOK, "Stop")
	}
}

func startSimulatorServer() {
	fmt.Printf("Starting Simulator Server...")
	order := NewOrder()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/order/set/clock", order.SetClock())
	e.POST("/order/set/agent", order.SetAgent())
	e.POST("/order/set/area", order.SetArea())
	e.POST("/order/start", order.Start())
	e.POST("/order/stop", order.Stop())

	e.Start(":" + port)
}

func main() {
	fmt.Printf("NumCPU=%d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	/*quit := make(chan os.Signal)
	// 受け取るシグナルを設定
	signal.Notify(quit, os.Interrupt)
	<-quit*/

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

	// api
	fmt.Printf("client: %v\n", client)
	simapi = api.NewSimAPI()
	simapi.RegistClients(client, myProvider.Id, argJson) // channelごとのClientを作成
	simapi.SubscribeAll(demandCallback, supplyCallback)  // ChannelにSubscribe*/

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	nodeapi.CallDeferFunctions() // cleanup!

}

//////////////////////////////////
////////// Pod Generator ////////
//////////////////////////////////

type PodGenerator struct {
	RsrcMap map[string][]Resource
}

func NewPodGenerator() *PodGenerator {
	pg := &PodGenerator{
		RsrcMap: make(map[string][]Resource),
	}
	return pg
}

func (pg *PodGenerator) applyWorker(area *api.Area, neighbors []string) error {
	fmt.Printf("applying WorkerPod... %v\n", area.Id)
	areaid := strconv.FormatUint(area.Id, 10)
	rsrcs := []Resource{
		pg.NewWorkerService(areaid),
		pg.NewWorker(areaid),
		pg.NewAgent(areaid, area),
	}
	for _, neiId := range neighbors {
		rsrcs = append(rsrcs, pg.NewGateway(areaid, neiId))
	}
	fmt.Printf("applying WorkerPod2... %v\n", areaid)
	// write yaml
	fileName := "scripts/worker" + areaid + ".yaml"
	for _, rsrc := range rsrcs {
		err := WriteOnFile(fileName, rsrc)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	fmt.Printf("test: %v %v\n", fileName, areaid)
	// apply yaml
	cmd := exec.Command("kubectl", "apply", "-f", fileName)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Command Start Error. %v\n", err)
		return err
	}

	// delete yaml
	/*if err := os.Remove(fileName); err != nil {
		fmt.Println(err)
		return err
	}*/
	fmt.Printf("out: %v\n", string(out))

	// regist resource
	pg.RsrcMap[areaid] = rsrcs

	return nil
}

func (pg *PodGenerator) deleteWorker(areaid string) error {
	fmt.Printf("deleting WorkerPod...")
	rsrcs := pg.RsrcMap[areaid]

	// write yaml
	fileName := "worker" + areaid + ".yaml"
	for _, rsrc := range rsrcs {
		err := WriteOnFile(fileName, rsrc)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	// apply yaml
	cmd := exec.Command("kubectl", "delete", "-f", fileName)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Command Start Error.")
		return err
	}

	// delete yaml
	if err := os.Remove(fileName); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("out: %v\n", string(out))

	// regist resource
	pg.RsrcMap[areaid] = nil

	return nil
}

// gateway
func (pg *PodGenerator) NewGateway(areaId string, neiId string) Resource {
	worker1Name := "worker" + areaId
	worker2Name := "worker" + neiId
	gatewayName := "gateway" + areaId + neiId
	gateway := Resource{
		ApiVersion: "v1",
		Kind:       "Pod",
		Metadata: Metadata{
			Name:   gatewayName,
			Labels: Label{App: gatewayName},
		},
		Spec: Spec{
			Containers: []Container{
				{
					Name:            "gateway-provider",
					Image:           "synerex-simulation/gateway-provider:latest",
					ImagePullPolicy: "Never",
					Env: []Env{
						{
							Name:  "WORKER_SYNEREX_SERVER1",
							Value: worker1Name + ":700",
						},
						{
							Name:  "WORKER_NODEID_SERVER1",
							Value: worker1Name + ":600",
						},
						{
							Name:  "WORKER_SYNEREX_SERVER2",
							Value: worker2Name + ":700",
						},
						{
							Name:  "WORKER_NODEID_SERVER2",
							Value: worker2Name + ":600",
						},
						{
							Name:  "PROVIDER_NAME",
							Value: "GatewayProvider" + areaId + neiId,
						},
					},
					Ports: []Port{{ContainerPort: 9980}},
				},
			},
		},
	}
	return gateway
}

func (pg *PodGenerator) NewAgent(areaid string, area *api.Area) Resource {
	workerName := "worker" + areaid
	agentName := "agent" + areaid
	agent := Resource{
		ApiVersion: "v1",
		Kind:       "Pod",
		Metadata: Metadata{
			Name:   agentName,
			Labels: Label{App: agentName},
		},
		Spec: Spec{
			Containers: []Container{
				{
					Name:            "agent-provider",
					Image:           "synerex-simulation/agent-provider:latest",
					ImagePullPolicy: "Never",
					Env: []Env{
						{
							Name:  "NODEID_SERVER",
							Value: workerName + ":600",
						},
						{
							Name:  "SYNEREX_SERVER",
							Value: workerName + ":700",
						},
						{
							Name:  "VIS_SYNEREX_SERVER",
							Value: "visualization:700",
						},
						{
							Name:  "VIS_NODEID_SERVER",
							Value: "visualization:600",
						},
						{
							Name:  "AREA",
							Value: convertAreaToJson(area),
						},
						{
							Name:  "PROVIDER_NAME",
							Value: "AgentProvider" + areaid,
						},
					},
				},
			},
		},
	}
	return agent
}

// worker
func (pg *PodGenerator) NewWorkerService(areaid string) Resource {
	name := "worker" + areaid
	service := Resource{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata:   Metadata{Name: name},
		Spec: Spec{
			Selector: Selector{App: name},
			Ports: []Port{
				{
					Name:       "synerex",
					Port:       700,
					TargetPort: 10000,
				},
				{
					Name:       "nodeid",
					Port:       600,
					TargetPort: 9000,
				},
			},
		},
	}
	return service
}

func (pg *PodGenerator) NewWorker(areaid string) Resource {
	name := "worker" + areaid
	worker := Resource{
		ApiVersion: "v1",
		Kind:       "Pod",
		Metadata: Metadata{
			Name:   name,
			Labels: Label{App: name},
		},
		Spec: Spec{
			Containers: []Container{
				{
					Name:            "nodeid-server",
					Image:           "synerex-simulation/nodeid-server:latest",
					ImagePullPolicy: "Never",
					Env: []Env{
						{
							Name:  "NODEID_SERVER",
							Value: ":9000",
						},
					},
					Ports: []Port{{ContainerPort: 9000}},
				},
				{
					Name:            "synerex-server",
					Image:           "synerex-simulation/synerex-server:latest",
					ImagePullPolicy: "Never",
					Env: []Env{
						{
							Name:  "NODEID_SERVER",
							Value: ":9000",
						},
						{
							Name:  "SYNEREX_SERVER",
							Value: ":10000",
						},
						{
							Name:  "SERVER_NAME",
							Value: "SynerexServer" + areaid,
						},
					},
					Ports: []Port{{ContainerPort: 10000}},
				},
				{
					Name:            "worker-provider",
					Image:           "synerex-simulation/worker-provider:latest",
					ImagePullPolicy: "Never",
					Env: []Env{
						{
							Name:  "NODEID_SERVER",
							Value: ":9000",
						},
						{
							Name:  "SYNEREX_SERVER",
							Value: ":10000",
						},
						{
							Name:  "MASTER_SYNEREX_SERVER",
							Value: "master:700",
						},
						{
							Name:  "MASTER_NODEID_SERVER",
							Value: "master:600",
						},
						{
							Name:  "PORT",
							Value: "9980",
						},
						{
							Name:  "PROVIDER_NAME",
							Value: "WorkerProvider" + areaid,
						},
					},
					Ports: []Port{{ContainerPort: 9980}},
				},
			},
		},
	}
	return worker
}

// ファイル名とデータをを渡すとyamlファイルに保存してくれる関数です。
func WriteOnFile(fileName string, data interface{}) error {
	// ここでデータを []byte に変換しています。
	buf, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		//エラー処理
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintln(file, string(buf))   //書き込み
	fmt.Fprintln(file, string("---")) //書き込み
	return nil
}

func convertAreaToJson(area *api.Area) string {
	id := area.Id
	duplicateText := `[`
	controlText := `[`
	for i, ctl := range area.ControlArea {
		ctlText := fmt.Sprintf(`{"latitude":%v, "longitude":%v}`, ctl.Latitude, ctl.Longitude)
		//fmt.Printf("ctl %v\n", ctlText)
		if i == len(area.ControlArea)-1 { // 最後は,をつけない
			controlText += ctlText
		} else {
			controlText += ctlText + ","
		}
	}
	for i, dpl := range area.DuplicateArea {
		dplText := fmt.Sprintf(`{"latitude":%v, "longitude":%v}`, dpl.Latitude, dpl.Longitude)
		//fmt.Printf("dpl %v\n", dplText)
		if i == len(area.DuplicateArea)-1 { // 最後は,をつけない
			duplicateText += dplText
		} else {
			duplicateText += dplText + ","
		}
	}

	duplicateText += `]`
	controlText += `]`
	result := fmt.Sprintf(`{"id":%d, "name":"Unknown", "duplicate_area": %s, "control_area": %s}`, id, duplicateText, controlText)
	//result = fmt.Sprintf("%s", result)
	//fmt.Printf("areaJson: %s\n", result)
	return result
}

type Resource struct {
	ApiVersion string   `yaml:"apiVersion,omitempty"`
	Kind       string   `yaml:"kind,omitempty"`
	Metadata   Metadata `yaml:"metadata,omitempty"`
	Spec       Spec     `yaml:"spec,omitempty"`
}

type Spec struct {
	Containers []Container `yaml:"containers,omitempty"`
	Selector   Selector    `yaml:"selector,omitempty"`
	Ports      []Port      `yaml:"ports,omitempty"`
	Type       string      `yaml:"type,omitempty"`
}

type Container struct {
	Name            string `yaml:"name,omitempty"`
	Image           string `yaml:"image,omitempty"`
	ImagePullPolicy string `yaml:"imagePullPolicy,omitempty"`
	Stdin           bool   `yaml:"stdin,omitempty"`
	Tty             bool   `yaml:"tty,omitempty"`
	Env             []Env  `yaml:"env,omitempty"`
	Ports           []Port `yaml:"ports,omitempty"`
}

type Env struct {
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type Selector struct {
	App         string `yaml:"app,omitempty"`
	MatchLabels Label  `yaml:"matchLabels,omitempty"`
}

type Port struct {
	Name          string `yaml:"name,omitempty"`
	Port          int    `yaml:"port,omitempty"`
	TargetPort    int    `yaml:"targetPort,omitempty"`
	ContainerPort int    `yaml:"containerPort,omitempty"`
}

type Metadata struct {
	Name   string `yaml:"name,omitempty"`
	Labels Label  `yaml:"labels,omitempty"`
}

type Label struct {
	App string `yaml:"app,omitempty"`
}

type Area struct {
	Id        int
	Control   []*api.Coord
	Duplicate []*api.Coord
}

type Coord struct {
	Latitude  float64
	Longitude float64
}

func GetCoordRange(coords []*api.Coord) (float64, float64, float64, float64) {
	maxLon, maxLat := math.Inf(-1), math.Inf(-1)
	minLon, minLat := math.Inf(0), math.Inf(0)
	for _, coord := range coords {
		if coord.Latitude > maxLat {
			maxLat = coord.Latitude
		}
		if coord.Longitude > maxLon {
			maxLon = coord.Longitude
		}
		if coord.Latitude < minLat {
			minLat = coord.Latitude
		}
		if coord.Longitude < minLon {
			minLon = coord.Longitude
		}
	}
	return maxLat, maxLon, minLat, minLon
}
