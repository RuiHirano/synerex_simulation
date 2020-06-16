package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
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
)

type Config struct {
	Area Config_Area `yaml:"area"`
}

type Config_Area struct {
	SideRange float64 `yaml:"sideRange"`
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
	id := "test"
	area := &Area{
		Id:        3,
		Control:   []Coord{{Latitude: 0, Longitude: 0}, {Latitude: 10, Longitude: 0}, {Latitude: 10, Longitude: 10}, {Latitude: 0, Longitude: 10}},
		Duplicate: []Coord{{Latitude: 0, Longitude: 0}, {Latitude: 10, Longitude: 0}, {Latitude: 10, Longitude: 10}, {Latitude: 0, Longitude: 10}},
	}
	go podgen.applyWorker(id, area)
	time.Sleep(4 * time.Second)
	go podgen.deleteWorker(id)
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

	logger.Info("Finish Setting Agents \n Add: %v", len(agents))
	return true, nil
}

// startClock:
func startClock() {
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

func (pg *PodGenerator) applyWorker(areaid string, area *Area) error {
	fmt.Printf("applying WorkerPod...")
	rsrcs := []Resource{
		pg.NewWorkerService(areaid),
		pg.NewWorker(areaid),
		pg.NewAgent(areaid, area),
	}

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
	cmd := exec.Command("kubectl", "apply", "-f", fileName)
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

func (pg *PodGenerator) NewAgent(areaid string, area *Area) Resource {
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

func convertAreaToJson(area *Area) string {
	id := area.Id
	duplicateText := `[`
	controlText := `[`
	for i, ctl := range area.Control {
		ctlText := fmt.Sprintf(`{"latitude":%v, "longitude":%v}`, ctl.Latitude, ctl.Longitude)
		//fmt.Printf("ctl %v\n", ctlText)
		if i == len(area.Control)-1 { // 最後は,をつけない
			controlText += ctlText
		} else {
			controlText += ctlText + ","
		}
	}
	for i, dpl := range area.Duplicate {
		dplText := fmt.Sprintf(`{"latitude":%v, "longitude":%v}`, dpl.Latitude, dpl.Longitude)
		//fmt.Printf("dpl %v\n", dplText)
		if i == len(area.Duplicate)-1 { // 最後は,をつけない
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
	Control   []Coord
	Duplicate []Coord
}

type Coord struct {
	Latitude  float64
	Longitude float64
}
