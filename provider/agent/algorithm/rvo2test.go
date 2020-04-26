package algorithm

import (
	//"fmt"
	"math"

	//"github.com/paulmach/orb"
	//"github.com/paulmach/orb/geojson"

	//monitor "github.com/RuiHirano/rvo2-go/monitor"

	//"io/ioutil"
	//"log"

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	"github.com/synerex/synerex_alpha/api"
)

var (
	sim *rvo.RVOSimulator

//fcs *geojson.FeatureCollection
)

type RVO2Route2 struct {
	Agents []*api.Agent
	Area   *api.Area
}

func NewRVO2Route2(agents []*api.Agent, area *api.Area) *RVO2Route2 {

	r := &RVO2Route2{
		Agents: agents,
		Area:   area,
	}
	return r
}

// CalcDirectionAndDistance: 目的地までの距離と角度を計算する関数
func (rvo2route *RVO2Route2) CalcDirectionAndDistance(startCoord *api.Coord, goalCoord *api.Coord) (float64, float64) {

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
func (rvo2route *RVO2Route2) DecideNextTransit(nextTransit *api.Coord, transitPoint []*api.Coord, distance float64, destination *api.Coord) *api.Coord {
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
func (rvo2route *RVO2Route2) SetupScenario() {
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
}

func (rvo2route *RVO2Route2) CalcNextAgents() []*api.Agent {

	currentAgents := rvo2route.Agents

	timeStep := 0.1
	neighborDist := 0.00008 // どのくらいの距離の相手をNeighborと認識するか?Neighborとの距離をどのくらいに保つか？ぶつかったと認識する距離？
	maxneighbors := 10      // 周り何体を計算対象とするか
	timeHorizon := 1.0
	timeHorizonObst := 1.0
	radius := 0.00001  // エージェントの半径
	maxSpeed := 0.0004 // エージェントの最大スピード
	sim = rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})

	// scenario設定
	rvo2route.SetupScenario()

	// Stepを進める
	sim.DoStep()

	// 管理エリアのエージェントのみを抽出
	nextControlAgents := make([]*api.Agent, 0)
	for rvoId, agentInfo := range currentAgents {
		// 管理エリア内のエージェントのみ抽出
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
			//direction, distance := rvo2route.CalcDirectionAndDistance(nextCoord, agentInfo.Route.NextTransit)
			// 次の経由地nextTransitを求める
			//nextTransit := rvo2route.DecideNextTransit(agentInfo.Route.NextTransit, agentInfo.Route.TransitPoints, distance, destination)
			nextTransit := agentInfo.Route.NextTransit
			goalVector := sim.GetAgentGoalVector(int(rvoId))
			direction := math.Atan2(goalVector.Y, goalVector.X)
			speed := agentInfo.Route.Speed

			nextRoute := &api.Route{
				Position:      nextCoord,
				Direction:     direction,
				Speed:         speed,
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

	return currentAgents
}

// エージェントがエリアの中にいるかどうか
func IsAgentInArea(position *api.Coord, areaCoords []*api.Coord) bool {
	lat := position.Latitude
	lon := position.Longitude
	maxLat, maxLon, minLat, minLon := GetCoordRange(areaCoords)
	if minLat < lat && lat < maxLat && minLon < lon && lon < maxLon {
		return true
	} else {
		return false
	}
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
