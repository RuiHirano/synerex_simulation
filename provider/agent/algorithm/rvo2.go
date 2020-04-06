package algorithm

import (
	//"fmt"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"

	//monitor "github.com/RuiHirano/rvo2-go/monitor"

	"io/ioutil"
	"log"

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	"github.com/synerex/synerex_alpha/api"
)

var (
	sim *rvo.RVOSimulator
	fcs *geojson.FeatureCollection
)

func loadGeoJson(fname string) *geojson.FeatureCollection {

	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Print("Can't read file:", err)
		panic("load json")
	}
	fc, _ := geojson.UnmarshalFeatureCollection(bytes)

	return fc
}

type RVO2Route struct {
	TimeStep   float64
	GlobalTime float64
	Area       *api.Area
	Agents     []*api.Agent
	AgentType  api.AgentType
}

func NewRVO2Route(timeStep float64, globalTime float64, area *api.Area, agentsInfo []*api.Agent, agentType api.AgentType) *RVO2Route {

	// set obstacle
	fcs = loadGeoJson("higashiyama.geojson")

	r := &RVO2Route{
		TimeStep:   timeStep,
		GlobalTime: globalTime,
		Area:       area,
		Agents:     agentsInfo,
		AgentType:  agentType,
	}
	return r
}

// CalcDirectionAndDistance: 目的地までの距離と角度を計算する関数
func (rvo2route *RVO2Route) CalcDirectionAndDistance(startCoord *api.Coord, goalCoord *api.Coord) (float64, float64) {

	r := 6378137 // equatorial radius
	sLat := startCoord.Latitude * math.Pi / 180
	sLon := startCoord.Longitude * math.Pi / 180
	gLat := goalCoord.Latitude * math.Pi / 180
	gLon := goalCoord.Longitude * math.Pi / 180
	dLon := gLon - sLon
	dLat := gLat - sLat
	cLat := (sLat + gLat) / 2
	dx := float64(r) * float64(dLon) * math.Cos(float64(cLat))
	dy := float64(r) * float64(dLat)

	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	direction := float64(0)
	if dx != 0 && dy != 0 {
		direction = math.Atan2(dy, dx) * 180 / math.Pi
	}

	return direction, distance
}

// DecideNextTransit: 次の経由地を求める関数
func (rvo2route *RVO2Route) DecideNextTransit(nextTransit *api.Coord, transitPoint []*api.Coord, distance float64, destination *api.Coord) *api.Coord {
	// 距離が5m以下の場合
	if distance < 150 {
		if nextTransit != destination {
			for i, tPoint := range transitPoint {
				if tPoint.Longitude == nextTransit.Longitude && tPoint.Latitude == nextTransit.Latitude {
					if i+1 == len(transitPoint) {
						// すべての経由地を通った場合、nextTransitをdestinationにする
						nextTransit = destination
					} else {
						// 次の経由地を設定する
						nextTransit = transitPoint[i+1]
					}
				}
			}
		} else {
			//fmt.Printf("arrived!")
		}
	}
	return nextTransit
}

// SetupScenario: Scenarioを設定する関数
func (rvo2route *RVO2Route) SetupScenario() {

	// Set Agent
	for _, agentInfo := range rvo2route.Agents {

		position := &rvo.Vector2{X: agentInfo.Route.Position.Longitude, Y: agentInfo.Route.Position.Latitude}
		goal := &rvo.Vector2{X: agentInfo.Route.NextTransit.Longitude, Y: agentInfo.Route.NextTransit.Latitude}

		// Agentを追加
		id, _ := sim.AddDefaultAgent(position)

		// 目的地を設定
		sim.SetAgentGoal(id, goal)

		// エージェントの速度方向ベクトルを設定
		goalVector := sim.GetAgentGoalVector(id)
		sim.SetAgentPrefVelocity(id, goalVector)
		//sim.SetAgentMaxSpeed(id, float64(api.MaxSpeed))
	}

	// Set Obstacle
	for _, feature := range fcs.Features {
		multiPosition := feature.Geometry.(orb.MultiLineString)[0]
		//fmt.Printf("geometry: ", multiPosition)
		rvoObstacle := []*rvo.Vector2{}

		for _, positionArray := range multiPosition {
			position := &rvo.Vector2{
				X: positionArray[0],
				Y: positionArray[1],
			}

			rvoObstacle = append(rvoObstacle, position)
		}

		sim.AddObstacle(rvoObstacle)
	}

	sim.ProcessObstacles()

	//fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	//fmt.Printf("Running Simulation..\n\n")
}

func (rvo2route *RVO2Route) CalcNextAgents() []*api.Agent {

	nextControlAgents := make([]*api.Agent, 0)
	currentAgents := rvo2route.Agents

	timeStep := rvo2route.TimeStep
	neighborDist := 0.00008 // どのくらいの距離の相手をNeighborと認識するか?Neighborとの距離をどのくらいに保つか？ぶつかったと認識する距離？
	maxneighbors := 10      // 周り何体を計算対象とするか
	timeHorizon := 1.0
	timeHorizonObst := 1.0
	radius := 0.00001   // エージェントの半径
	maxSpeed := 0.00004 // エージェントの最大スピード
	sim = rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})

	// scenario設定
	rvo2route.SetupScenario()

	// Stepを進める
	sim.DoStep()

	// 管理エリアのエージェントのみを抽出
	for rvoId, agentInfo := range currentAgents {
		//nextRVOAgent := sim.GetAgent(int(agentInfo.Id))
		// 計算する前に自エリアにいる場合、次のルートを計算する
		position := agentInfo.Route.Position
		if IsAgentInArea(position, rvo2route.Area.ControlArea) {
			destination := agentInfo.Route.Destination

			// rvoの位置情報を緯度経度に変換する
			rvoAgentPosition := sim.GetAgentPosition(int(rvoId))

			nextCoord := &api.Coord{
				Latitude:  rvoAgentPosition.Y,
				Longitude: rvoAgentPosition.X,
			}

			// 現在の位置とゴールとの距離と角度を求める (度, m))
			direction, distance := rvo2route.CalcDirectionAndDistance(nextCoord, agentInfo.Route.NextTransit)
			// 次の経由地nextTransitを求める
			nextTransit := rvo2route.DecideNextTransit(agentInfo.Route.NextTransit, agentInfo.Route.TransitPoints, distance, destination)

			nextRoute := &api.Route{
				Position:      nextCoord,
				Direction:     direction,
				Speed:         distance,
				Destination:   destination,
				Departure:     agentInfo.Route.Departure,
				TransitPoints: agentInfo.Route.TransitPoints,
				NextTransit:   nextTransit,
				TotalDistance: agentInfo.Route.TotalDistance,
				RequiredTime:  agentInfo.Route.RequiredTime,
			}

			nextControlAgent := &api.Agent{
				Id:    agentInfo.Id,
				Type:  agentInfo.Type,
				Route: nextRoute,
			}

			nextControlAgents = append(nextControlAgents, nextControlAgent)
		}
	}

	return nextControlAgents
}
