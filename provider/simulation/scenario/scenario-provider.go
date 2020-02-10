package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/uuid"
	gosocketio "github.com/mtfelian/golang-socketio"
	pb "github.com/synerex/synerex_alpha/api"
	simapi "github.com/synerex/synerex_alpha/api/simulation"
	agent "github.com/synerex/synerex_alpha/api/simulation/agent"
	area "github.com/synerex/synerex_alpha/api/simulation/area"
	"github.com/synerex/synerex_alpha/api/simulation/clock"
	common "github.com/synerex/synerex_alpha/api/simulation/common"
	provider "github.com/synerex/synerex_alpha/api/simulation/provider"
	"github.com/synerex/synerex_alpha/provider/simulation/simutil"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr      = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodeIdAddr      = flag.String("nodeid_addr", "127.0.0.1:9990", "Node ID Server")
	isStart         bool
	mu              sync.Mutex
	com             *simutil.Communicator
	sim             *Simulator
	providerManager *simutil.ProviderManager
	pSources        map[provider.ProviderType]*provider.Source
	logger          *simutil.Logger
)

const MAX_AGENTS_NUM = 1000

func init() {
	isStart = false
	logger = simutil.NewLogger()
	flag.Parse()
}

var (
	//fcs *geojson.FeatureCollection
	//geofile string
	port            = 9995
	assetsDir       http.FileSystem
	server          *gosocketio.Server = nil
	providerMutex   sync.RWMutex
	providerSources []*Source
	serverSources   []*Source
	orderInfos      []OrderInfo
)

/*func loadGeoJson(fname string) *geojson.FeatureCollection{

	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Print("Can't read file:", err)
		panic("load json")
	}
	fc, _ := geojson.UnmarshalFeatureCollection(bytes)

	return fc
}*/

// エリア分割の閾値(100人超えたら分割)
var areaAgentNum = 100

/*var mockProviderStats = []*ProviderStats{
	{
		AgentNum: 50,
		Time: 0.0,
	},
}*/

func Sub(coord1 *common.Coord, coord2 *common.Coord) *common.Coord {
	return &common.Coord{Latitude: coord1.Latitude - coord2.Latitude, Longitude: coord1.Longitude - coord2.Longitude}
}

func Mul(coord1 *common.Coord, coord2 *common.Coord) float64 {
	return coord1.Latitude*coord2.Latitude + coord1.Longitude*coord2.Longitude
}

func Abs(coord1 *common.Coord) float64 {
	return math.Sqrt(Mul(coord1, coord1))
}

func Add(coord1 *common.Coord, coord2 *common.Coord) *common.Coord {
	return &common.Coord{Latitude: coord1.Latitude + coord2.Latitude, Longitude: coord1.Longitude + coord2.Longitude}
}

func Div(coord *common.Coord, s float64) *common.Coord {
	return &common.Coord{Latitude: coord.Latitude / s, Longitude: coord.Longitude / s}
}

type ByAbs struct {
	Coords []*common.Coord
}

func (b ByAbs) Less(i, j int) bool {
	return Abs(b.Coords[i]) < Abs(b.Coords[j])
}
func (b ByAbs) Len() int {
	return len(b.Coords)
}

func (b ByAbs) Swap(i, j int) {
	b.Coords[i], b.Coords[j] = b.Coords[j], b.Coords[i]
}

/*func updateProvider(newProviderStats []*ProviderStats) []*ProviderStats{
	// write area devide algorithm
	updatedProviderStats := make([]*ProviderStats, 0)
	for _, state := range newProviderStats{
		if state.AgentNum > areaAgentNum{
			// devide area
			//vectors := make([]*common.Coord, 0)

			point1, point2, point3, point4 := state.Area[0], state.Area[1], state.Area[2], state.Area[3]
			point1vecs := []*common.Coord{Sub(point1, point1), Sub(point2, point1), Sub(point3, point1), Sub(point4, point1)}
			// 昇順にする
			sort.Sort(ByAbs{point1vecs})
			devPoint1 := Div(point1vecs[2], 2)	//分割点1
			divPoint2 := Add(Div(point1vecs[2], 2), point1vecs[1]) //分割点2
			// 二つに分割
			coords1 := []*common.Coord{
				Add(point1vecs[0], point1), Add(point1vecs[1], point1), Add(devPoint1, point1), Add(divPoint2, point1),
			}
			coords2 := []*common.Coord{
				Add(point1vecs[2], point1), Add(point1vecs[3], point1), Add(devPoint1, point1), Add(divPoint2, point1),
			}
			// 追加
			state1 := &ProviderStats{
				Area: coords1,
				AgentNum: state.AgentNum/2,
				Time: state.Time,
			}
			state2 := &ProviderStats{
				Area: coords2,
				AgentNum: state.AgentNum/2,
				Time: state.Time,
			}

			updatedProviderStats = append(updatedProviderStats, state1)
			updatedProviderStats = append(updatedProviderStats, state2)
		}else{
			updatedProviderStats = append(updatedProviderStats, state)
		}
	}
	return updatedProviderStats
}

func areaTest(){
	go func(){
		for{
			log.Printf("time: ---\n")
			time.Sleep(1 * time.Second)
			// change provider stats
			newProviderStats := make([]*ProviderStatus, 0)
			for i, state := range mockProviderStats{
				if i == 0{
					log.Printf("agentNum: %v, providerNum: %v\n", state.AgentNum, len(mockProviderStats))
				}
				log.Printf("area: %v\n", state.Area)
				state.AgentNum += 30
				newProviderStats = append(newProviderStats, state)
			}
			// update provider
			mockProviderStats = updateProvider(newProviderStats)

			// startup provider
			//log.Printf("providerStats: %v\n", providerStats)
		}
	}()
}*/

type Option struct {
	Key   string
	Value string
}

type Source struct {
	CmdName     string
	Type        provider.ProviderType
	Description string
	SrcDir      string
	BinName     string
	GoFiles     []string
	Options     []Option
}

type OrderInfo struct {
	CmdName string
	Type    OrderType
	Options []Option
}

func init() {
	//areaTest()
	//geofile = "transit_points.geojson"
	//fcs = loadGeoJson(geofile)
	//runProviders = make(map[uint64]*Provider)
	//providerMap = make(map[string]*exec.Cmd)
	providerMutex = sync.RWMutex{}
	pSources = make(map[provider.ProviderType]*provider.Source)
	pSources[provider.ProviderType_CLOCK] = &provider.Source{
		CmdName: "Clock",
		Type:    provider.ProviderType_CLOCK,
		SrcDir:  "provider/simulation/clock",
		BinName: "clock-provider",
		GoFiles: []string{"clock-provider.go"},
		SubFunc: SendLog,
	}
	pSources[provider.ProviderType_SCENARIO] = &provider.Source{
		CmdName: "Scenario",
		Type:    provider.ProviderType_SCENARIO,
		SrcDir:  "provider/simulation/scenario",
		BinName: "scenario-provider",
		GoFiles: []string{"scenario-provider.go"},
		SubFunc: SendLog,
	}
	pSources[provider.ProviderType_VISUALIZATION] = &provider.Source{
		CmdName: "Visualization",
		Type:    provider.ProviderType_VISUALIZATION,
		SrcDir:  "provider/simulation/visualization",
		BinName: "visualization-provider",
		GoFiles: []string{"visualization-provider.go"},
		SubFunc: SendLog,
	}
	pSources[provider.ProviderType_AGENT] = &provider.Source{
		CmdName: "Agent",
		Type:    provider.ProviderType_AGENT,
		SrcDir:  "provider/simulation/agent",
		BinName: "agent-provider",
		GoFiles: []string{"agent-provider.go"},
		SubFunc: SendLog,
	}
	pSources[provider.ProviderType_NODE_ID] = &provider.Source{
		CmdName: "NodeIDServer",
		Type:    provider.ProviderType_NODE_ID,
		SrcDir:  "nodeserv",
		BinName: "nodeid-server",
		GoFiles: []string{"nodeid-server.go"},
		SubFunc: SendLog,
	}
	pSources[provider.ProviderType_SYNEREX] = &provider.Source{
		CmdName: "SynerexServer",
		Type:    provider.ProviderType_SYNEREX,
		SrcDir:  "server",
		BinName: "synerex-server",
		GoFiles: []string{"synerex-server.go", "message-store.go"},
		SubFunc: SendLog,
	}
	pSources[provider.ProviderType_MONITOR] = &provider.Source{
		CmdName: "MonitorServer",
		Type:    provider.ProviderType_MONITOR,
		SrcDir:  "monitor",
		BinName: "monitor-server",
		GoFiles: []string{"monitor-server.go"},
		SubFunc: SendLog,
	}

	orderInfos = []OrderInfo{
		{
			CmdName: "SetAgents",
			Type:    OrderType_SET_AGENTS,
			Options: []Option{Option{
				Key:   "test",
				Value: "0",
			}},
		},
		{
			CmdName: "SetArea",
			Type:    OrderType_SET_AREA,
			Options: []Option{Option{
				Key:   "test",
				Value: "0",
			}},
		},
		{
			CmdName: "SetClock",
			Type:    OrderType_SET_CLOCK,
			Options: []Option{Option{
				Key:   "test",
				Value: "0",
			}},
		},
		{
			CmdName: "StartClock",
			Type:    OrderType_START_CLOCK,
			Options: []Option{Option{
				Key:   "test",
				Value: "0",
			}},
		},
		{
			CmdName: "StopClock",
			Type:    OrderType_STOP_CLOCK,
			Options: []Option{Option{
				Key:   "test",
				Value: "0",
			}},
		},
	}

}

////////////////////////////////////////////////////////////
//////////////////        Util          ///////////////////
///////////////////////////////////////////////////////////

// providerが変化した場合に更新通知をする
func updateProviderOrder() {
	go func() {
		providerNum := providerManager.GetProviderNum()
		for {
			mu.Lock()
			newProviderNum := providerManager.GetProviderNum()
			if providerNum != newProviderNum {
				providerNum = newProviderNum
				// IDMapを再作成
				providerManager.CreateIDMap()
				// 同期するIDリスト
				idList := providerManager.GetIDList([]simutil.IDType{
					simutil.IDType_CLOCK,
					simutil.IDType_VISUALIZATION,
					simutil.IDType_AGENT,
				})
				pid := providerManager.MyProvider.Id
				com.UpdateProvidersRequest(pid, idList, providerManager.Providers)
				logger.Info("Success Update Providers")
			}
			mu.Unlock()
		}
	}()
}

// providerに変化があった場合にGUIに情報を送る
/*func sendRunnningProviders() {
	providerMutex.RLock()

	//fmt.Printf("providers---------- %v\n", len(runProviders))
	rpJsons := make([]string, 0)
	for _, rp := range providerManager.Providers {
		bytes, _ := json.Marshal(rp)
		rpJson := string(bytes)
		fmt.Printf("provider----------\n")
		//fmt.Printf("Json: %v \n", rpJson)
		rpJsons = append(rpJsons, rpJson)
	}
	//c.Emit("providers", rpJsons)
	server.BroadcastToAll("providers", rpJsons)
	providerMutex.RUnlock()
}*/

// assetsFileHandler for static Data
func assetsFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}

	file := r.URL.Path

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

// Agentオブジェクトの変換
func calcRoute() *agent.Route {

	var departure, destination *common.Coord

	departure = &common.Coord{
		Latitude:  35.12532,
		Longitude: 135.235231,
	}
	destination = &common.Coord{
		Latitude:  35.12532,
		Longitude: 135.235231,
	}

	transitPoints := make([]*common.Coord, 0)
	transitPoints = append(transitPoints, destination)

	route := &agent.Route{
		Position:      departure,
		Direction:     100 * rand.Float64(),
		Speed:         100 * rand.Float64(),
		Departure:     departure,
		Destination:   destination,
		TransitPoints: transitPoints,
		NextTransit:   destination,
	}

	return route
}

func divideArea(areaInfo *area.Area, num uint64) []*area.Area {
	DUPLICATE_RANGE := 5.0
	// エリアを分割する
	// 最初は単純にエリアを半分にする
	//providerStats := mockProviderStats
	//duplicateRate := 0.1	// areaCoordの10%の範囲
	// 二等分にするアルゴリズム
	areaCoord := areaInfo.ControlArea
	point1, point2, point3, point4 := areaCoord[0], areaCoord[1], areaCoord[2], areaCoord[3]
	point1vecs := []*common.Coord{Sub(point1, point1), Sub(point2, point1), Sub(point3, point1), Sub(point4, point1)}
	// 昇順にする
	sort.Sort(ByAbs{point1vecs})
	divPoint1 := Div(point1vecs[2], 2)                     //分割点1
	divPoint2 := Add(Div(point1vecs[2], 2), point1vecs[1]) //分割点2
	// 二つに分割
	control1 := []*common.Coord{
		Add(point1vecs[0], point1), Add(point1vecs[1], point1), Add(divPoint1, point1), Add(divPoint2, point1),
	}
	control2 := []*common.Coord{
		Add(point1vecs[2], point1), Add(point1vecs[3], point1), Add(divPoint1, point1), Add(divPoint2, point1),
	}
	controls := [][]*common.Coord{control1, control2}

	// calc duplicate area
	var duplicates [][]*common.Coord
	for _, control := range controls {
		maxLat, maxLon := math.Inf(-1), math.Inf(-1)
		minLat, minLon := math.Inf(0), math.Inf(0)
		for _, coord := range control {
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
		duplicate := []*common.Coord{
			&common.Coord{Latitude: minLat - DUPLICATE_RANGE, Longitude: minLon - DUPLICATE_RANGE},
			&common.Coord{Latitude: minLat - DUPLICATE_RANGE, Longitude: maxLon + DUPLICATE_RANGE},
			&common.Coord{Latitude: maxLat + DUPLICATE_RANGE, Longitude: maxLon + DUPLICATE_RANGE},
			&common.Coord{Latitude: maxLat + DUPLICATE_RANGE, Longitude: minLon - DUPLICATE_RANGE},
		}
		duplicates = append(duplicates, duplicate)
	}

	// calc areaInfo
	areaInfos := make([]*area.Area, 0)
	for i, control := range controls {
		uid, _ := uuid.NewRandom()
		areaInfos = append(areaInfos, &area.Area{
			Id:            uint64(uid.ID()),
			Name:          "test",
			DuplicateArea: duplicates[i],
			ControlArea:   control,
			NeighborAreaIds: []uint64{}
		})
	}

	return areaInfos
}

////////////////////////////////////////////////////////////
//////////////////     ps Command     ////////////////////
//////////////////////////////////////////////////////////

/*func checkRunning(opt string) []string {
	isLong := false
	if opt == "long" {
		isLong = true
	}
	var procs []string
	i := 0
	providerMutex.RLock()
	if isLong {
		procs = make([]string, len(providerManager.Providers)+2)
		str := fmt.Sprintf("  pid: %-20s : \n", "process name")
		procs[i] = str
		procs[i+1] = "-----------------------------------------------------------------\n"
		i += 2
	} else {
		procs = make([]string, len(providerManager.Providers))
	}
	for _, provider := range providerManager.Providers {
		pid := pSources[provider.Type].Cmd.Process.Pid
		name := provider.Name
		if isLong {
			str := fmt.Sprintf("%5d: %-20s : \n", pid, name)
			procs[i] = str
		} else {
			if i != 0 {
				procs[i] = ", " + name
			} else {
				procs[i] = name
			}
		}
		i++
	}
	providerMutex.RUnlock()
	return procs

}*/

////////////////////////////////////////////////////////////
//////////////         Order Class         /////////////////
///////////////////////////////////////////////////////////

// Order
type OrderType int

const (
	OrderType_SET_AGENTS  OrderType = 0
	OrderType_SET_AREA    OrderType = 1
	OrderType_SET_CLOCK   OrderType = 2
	OrderType_START_CLOCK OrderType = 3
	OrderType_STOP_CLOCK  OrderType = 4
)

type OrderOption struct {
	AgentNum  string
	ClockTime string
}

type Order struct {
	Type   OrderType
	Name   string
	Option *OrderOption
}

func NewOrder(name string, option *OrderOption) (*Order, error) {
	for _, sc := range orderInfos {
		if sc.CmdName == name {
			o := &Order{
				Type:   sc.Type,
				Name:   name,
				Option: option,
			}
			return o, nil
		}
	}
	msg := "invalid OrderName..."
	return nil, fmt.Errorf("Error: %s\n", msg)
}

func (o *Order) Send() string {
	target := o.Name
	fmt.Printf("Target is : %v\n", target)
	switch target {
	case "SetClock":
		fmt.Printf("SetClock\n")
		globalTime := float64(0)
		timeStep := float64(1)
		o.SetClock(globalTime, timeStep)
		return "ok"

	case "SetAgents":
		fmt.Printf("SetAgents\n")
		//agentNum, _ := strconv.Atoi(order.Option)
		agentNum := uint64(1)
		o.SetAgents(agentNum)
		return "ok"

	case "StartClock":
		fmt.Printf("StartClock\n")
		o.StartClock()
		return "ok"

	case "StopClock":
		fmt.Printf("StopClock\n")
		o.StopClock()
		return "ok"

	case "SetArea":
		fmt.Printf("SetArea\n")
		//o.StopClock()
		return "ok"

	default:
		err := "true"
		log.Printf("Can't find command %s", target)
		return err
	}
}

// startClock:
func (o *Order) StartClock() (bool, error) {

	// 同期するIDリスト
	idList := providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_CLOCK,
	})
	// エージェントを設置するリクエスト
	pid := providerManager.MyProvider.Id
	com.StartClockRequest(pid, idList)
	return true, nil
}

// stopClock: Clockを停止する
func (o *Order) StopClock() (bool, error) {
	// 同期するIDリスト
	idList := providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_CLOCK,
	})
	// エージェントを設置するリクエスト
	pid := providerManager.MyProvider.Id
	com.StopClockRequest(pid, idList)

	return true, nil
}

// setAgents: agentをセットするDemandを出す関数
func (o *Order) SetAgents(agentNum uint64) (bool, error) {

	agents := make([]*agent.Agent, 0)

	for i := 0; i < int(agentNum); i++ {
		uuid, err := uuid.NewRandom()
		if err == nil {
			agent := &agent.Agent{
				Id:    uint64(uuid.ID()),
				Type:  agent.AgentType_PEDESTRIAN,
				Route: calcRoute(),
				Data: &agent.Agent_Pedestrian{
					Pedestrian: &agent.Pedestrian{
						Status: &agent.PedStatus{
							Age:  "20",
							Name: "rui",
						},
					},
				},
			}
			agents = append(agents, agent)
		}
	}

	// エージェントを設置するリクエスト
	// 同期するIDリスト
	idList := providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_AGENT,
	})
	pid := providerManager.MyProvider.Id
	com.SetAgentsRequest(pid, idList, agents)

	logger.Info("Finish Setting Agents \n Add: %v", len(agents))
	return true, nil
}

// setClock : クロック情報をDaemonから受け取りセットする
func (o *Order) SetClock(globalTime float64, timeStep float64) (bool, error) {
	// クロックをセット
	clockInfo := clock.NewClock(globalTime, timeStep, 1)
	sim.Clock = clockInfo

	// クロック情報をプロバイダに送信
	idList := providerManager.GetIDList([]simutil.IDType{
		simutil.IDType_CLOCK,
	})
	pid := providerManager.MyProvider.Id
	com.SetClockRequest(pid, idList, sim.Clock)
	logger.Info("Finish Setting Clock. \n GlobalTime:  %v \n TimeStep: %v", sim.Clock.GlobalTime, sim.Clock.TimeStep)
	return true, nil
}

////////////////////////////////////////////////////////////
////////////     Simulator CLI GUI Server    //////////////
//////////////////////////////////////////////////////////

type Log struct {
	ID          uint64
	Description string
}

type SimulatorServer struct{}

func NewSimulatorServer() *SimulatorServer {
	ss := &SimulatorServer{}
	return ss
}

func (ss *SimulatorServer) Run() error {
	go func() {
		log.Printf("Starting.. Synergic Engine:")
		currentRoot, err := os.Getwd()
		if err != nil {
			log.Printf("se-daemon: Can' get registered directory: %s", err.Error())
		}
		d := filepath.Join(currentRoot, "monitor", "build")

		assetsDir = http.Dir(d)
		server = gosocketio.NewServer()

		server.On(gosocketio.OnConnection, func(c *gosocketio.Channel, param interface{}) {
			log.Printf("Connected from %s as %s", c.IP(), c.Id())
			// we need to send providers array
			time.Sleep(1000 * time.Millisecond)
			//sendRunnningProviders()

		})
		server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
			log.Printf("Disconnected from %s as %s", c.IP(), c.Id())
		})

		server.On("ps", func(c *gosocketio.Channel, param interface{}) []string {
			// need to check param short or long
			//opt := param.(string)

			//return checkRunning(opt)
			return []string{"ok"}
		})

		server.On("run", func(c *gosocketio.Channel, param interface{}) string {
			targetName := param.(string)
			log.Printf("Get run command %s", targetName)

			p := provider.NewProvider(targetName, provider.ProviderType_AGENT)
			p.Run(pSources[provider.ProviderType_AGENT])

			//sendRunnningProviders()
			return "ok"
		})

		server.On("order", func(c *gosocketio.Channel, param *Order) string {
			name := param.Name
			log.Printf("Get order command %s\n", name)
			log.Printf("Get order command %v\n", param)
			log.Printf("Get order command %v\n", param.Option)
			order, _ := NewOrder(name, nil)
			order.Send()
			return "ok"
		})

		serveMux := http.NewServeMux()
		serveMux.Handle("/socket.io/", server)
		serveMux.HandleFunc("/", assetsFileHandler)
		log.Println("Serving at localhost:9995...")
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), serveMux); err != nil {
			log.Panic(err)
		}

		return

	}()
	return nil
}

func SendLog(pipe io.ReadCloser, name string) {
	// logを読み取る

	reader := bufio.NewReader(pipe)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			log.Printf("Command [%s] EOF\n")
			break
		} else if err != nil {
			log.Printf("Err %v\n", err)
		}

		logInfo := &Log{
			ID:          uint64(0),
			Description: string(line),
		}

		bytes, err := json.Marshal(logInfo)
		logjson := string(bytes)

		if server != nil {
			server.BroadcastToAll("log", logjson)
		}
		log.Printf("[%s]:\n%s", name, string(line))
	}
}

////////////////////////////////////////////////////////////
////////////     Demand Supply Callback     ////////////////
///////////////////////////////////////////////////////////

// Supplyのコールバック関数
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// 自分宛かどうか
	if sp.GetTargetId() == providerManager.MyProvider.Id {
		// check if supply is match with my demand.
		switch sp.GetSimSupply().GetType() {
		case simapi.SupplyType_UPDATE_PROVIDERS_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_SET_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_SET_AGENTS_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_START_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		case simapi.SupplyType_STOP_CLOCK_RESPONSE:
			com.SendToWaitCh(sp, sp.GetSimSupply().GetType())
		}
	}
}

// Demandのコールバック関数
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	tid := dm.GetSimDemand().GetPid()
	pid := providerManager.MyProvider.Id
	// check if supply is match with my demand.
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_REGIST_PROVIDER_REQUEST:
		// providerを追加する
		p := dm.GetSimDemand().GetRegistProviderRequest().GetProvider()
		providerManager.AddProvider(p)
		// 登録完了通知
		//logger.Debug("RegistProviderRequest: Send from %v to %v\n", pid, tid)
		com.RegistProviderResponse(pid, tid)

		// UpdateRequest
		idList := providerManager.GetIDList([]simutil.IDType{
			simutil.IDType_CLOCK,
			simutil.IDType_VISUALIZATION,
			simutil.IDType_AGENT,
		})
		pid := providerManager.MyProvider.Id
		com.UpdateProvidersRequest(pid, idList, providerManager.Providers)
		logger.Info("Success Update Providers")

	case simapi.DemandType_DIVIDE_PROVIDER_REQUEST:
	case simapi.DemandType_KILL_PROVIDER_REQUEST:
	case simapi.DemandType_SEND_PROVIDER_STATUS_REQUEST:

	}
}

////////////////////////////////////////////////////////////
////////////     Run Initial Provider     ////////////////
///////////////////////////////////////////////////////////

// 担当するAreaの範囲
var mockAreaInfo = &area.Area{
	Id:   uint64(0),
	Name: "Area1",
	DuplicateArea: []*common.Coord{
		{Latitude: 35.156431, Longitude: 136.97285},
		{Latitude: 35.156431, Longitude: 136.981308},
		{Latitude: 35.153578, Longitude: 136.981308},
		{Latitude: 35.153578, Longitude: 136.97285},
	},
	ControlArea: []*common.Coord{
		{Latitude: 35.156431, Longitude: 136.97285},
		{Latitude: 35.156431, Longitude: 136.981308},
		{Latitude: 35.153578, Longitude: 136.981308},
		{Latitude: 35.153578, Longitude: 136.97285},
	},
}

// 担当するAreaの範囲
var mockAreaInfo2 = &area.Area{
	Id:   uint64(0),
	Name: "Area1",
	DuplicateArea: []*common.Coord{
		{Latitude: 1, Longitude: 1},
		{Latitude: 1, Longitude: 100},
		{Latitude: 100, Longitude: 100},
		{Latitude: 100, Longitude: 1},
	},
	ControlArea: []*common.Coord{
		{Latitude: 1, Longitude: 1},
		{Latitude: 1, Longitude: 100},
		{Latitude: 100, Longitude: 100},
		{Latitude: 100, Longitude: 1},
	},
}

func runInitServer() {
	// Run Server and Provider
	nodeServer := provider.NewProvider("NodeIDServer", provider.ProviderType_NODE_ID)
	nodeServer.Run(pSources[nodeServer.Type])

	time.Sleep(100 * time.Millisecond)
	monitorServer := provider.NewProvider("MonitorServer", provider.ProviderType_MONITOR)
	monitorServer.Run(pSources[monitorServer.Type])

	time.Sleep(100 * time.Millisecond)
	synerexServer := provider.NewProvider("SynerexServer", provider.ProviderType_SYNEREX)
	synerexServer.Run(pSources[synerexServer.Type])
	time.Sleep(100 * time.Millisecond)

}

func runInitProvider() {
	m := jsonpb.Marshaler{}
	scenarioJson, _ := m.MarshalToString(providerManager.MyProvider)

	clockProvider := provider.NewProvider("Clock", provider.ProviderType_CLOCK)
	js, _ := m.MarshalToString(clockProvider)
	options := provider.NewProviderOptions(*serverAddr, *nodeIdAddr, js, scenarioJson)
	pSources[clockProvider.Type].Options = options
	clockProvider.Run(pSources[clockProvider.Type])

	time.Sleep(100 * time.Millisecond)
	visProvider := provider.NewProvider("Visualization", provider.ProviderType_VISUALIZATION)
	js, _ = m.MarshalToString(visProvider)
	options = provider.NewProviderOptions(*serverAddr, *nodeIdAddr, js, scenarioJson)
	pSources[visProvider.Type].Options = options
	visProvider.Run(pSources[visProvider.Type])

	var INIT_PROVIDER_NUM = uint64(2)
	var INIT_AGENT_TYPES = map[string]agent.AgentType{
		"Pedestrian": agent.AgentType_PEDESTRIAN,
		"Car":        agent.AgentType_CAR,
	}

	for name, agentType := range INIT_AGENT_TYPES {
		areaInfos := divideArea(mockAreaInfo2, INIT_PROVIDER_NUM)
		//logger.Error("mockAreaInfo: %v\n", mockAreaInfo2.ControlArea)
		for _, areaInfo := range areaInfos {
			//logger.Error("areaInfo: %v\n", areaInfo.ControlArea)
			//logger.Error("duplicateInfo: %v\n", areaInfo.DuplicateArea)
			agentStatus := &provider.AgentStatus{
				Area:      areaInfo,
				AgentType: agentType,
				AgentNum:  0,
			}
			p := provider.NewAgentProvider(name, agentType, agentStatus)
			js, _ = m.MarshalToString(p)
			options = provider.NewProviderOptions(*serverAddr, *nodeIdAddr, js, scenarioJson)

			time.Sleep(100 * time.Millisecond)
			pSources[p.Type].Options = options
			p.Run(pSources[p.Type])

		}
		//logger.Fatal("error")
	}

}

func main() {

	// ProviderManager
	myProvider := provider.NewProvider("ScenarioProvider", provider.ProviderType_SCENARIO)
	providerManager = simutil.NewProviderManager(myProvider)
	providerManager.CreateIDMap()

	// CLI, GUIの受信サーバ
	simulatorServer := NewSimulatorServer()
	simulatorServer.Run()

	// 初期プロバイダ起動
	runInitServer()

	// Connect to Node Server
	sxutil.RegisterNodeName(*nodeIdAddr, "ScenarioProvider", false)
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
	argJson := fmt.Sprintf("{Client:Scenario}")

	// Simulator
	clockInfo := clock.NewClock(0, 1, 1)
	sim = NewSimulator(clockInfo)

	// Communicator
	com = simutil.NewCommunicator()
	com.RegistClients(client, argJson)               // channelごとのClientを作成
	com.SubscribeAll(demandCallback, supplyCallback) // ChannelにSubscribe

	wg := sync.WaitGroup{}
	wg.Add(1)
	updateProviderOrder()
	// 初期プロバイダ起動
	runInitProvider()
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
