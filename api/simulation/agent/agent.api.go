package agent

import (
	fmt "fmt"
	math "math"

	common "github.com/synerex/synerex_alpha/api/simulation/common"
)

//sxutil "github.com/synerex/synerex_alpha/sxutil"

func NewAgent(agentType AgentType) *Agent {
	a := &Agent{
		Id:   0,
		Type: agentType,
	}
	return a
}

func NewPedestrian(agentType AgentType, ped *Pedestrian) *Agent {
	a := &Agent{
		Id:   0,
		Type: agentType,
	}
	a.WithPedestrian(ped)
	return a
}

func (a *Agent) WithPedestrian(p *Pedestrian) *Agent {
	a.Data = &Agent_Pedestrian{p}
	return a
}

func NewCar(agentType AgentType, car *Car) *Agent {
	a := &Agent{
		Id:   0,
		Type: agentType,
	}
	a.WithCar(car)
	return a
}

func (a *Agent) WithCar(c *Car) *Agent {
	a.Data = &Agent_Car{c}
	return a
}

// CLEAR
// エージェントがエリアの中にいるかどうか
// 各頂点との角度の和が360度となれば内包
func (a *Agent) IsInArea(areaCoords []*common.Coord) bool {
	lat := a.Route.Position.Latitude
	lon := a.Route.Position.Longitude
	deg := 0.0
	for i, coord := range areaCoords {
		p2lat := coord.Latitude
		p2lon := coord.Longitude
		p3lat := areaCoords[i+1].Latitude
		p3lon := areaCoords[i+1].Longitude
		if i == len(areaCoords)-1 {
			p3lat = areaCoords[0].Latitude
			p3lon = areaCoords[0].Longitude
		}
		alat := p2lat - lat
		alon := p2lon - lon
		blat := p3lat - lat
		blon := p3lon - lon
		cos := (alat*blat + alon*blon) / (math.Sqrt(alat*alat+alon+alon) * math.Sqrt(blat*blat+blon+blon))
		deg += math.Acos(cos) * float64(180) / math.Pi
	}
	if math.Round(deg) == 360 {
		return true
	} else {
		return false
	}
}

// CLEAR
// ある座標への距離と角度を返す関数
func (a *Agent) GetDirectionAndDistance(goal *common.Coord) (float64, float64) {
	var direction, distance float64
	r := 6378137 // equatorial radius (m)
	sLat := a.Route.Position.Latitude * math.Pi / 180
	sLon := a.Route.Position.Longitude * math.Pi / 180
	gLat := goal.Latitude * math.Pi / 180
	gLon := goal.Longitude * math.Pi / 180
	dLon := gLon - sLon
	dLat := gLat - sLat
	cLat := (sLat + gLat) / 2
	dx := float64(r) * float64(dLon) * math.Cos(float64(cLat))
	dy := float64(r) * float64(dLat)

	distance = math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	direction = 0
	if dx != 0 && dy != 0 {
		direction = math.Atan2(dy, dx) * 180 / math.Pi
	}

	return direction, distance
}

// CLEAR
// ある座標に到着したかどうか
func (a *Agent) IsReachedGoal(goal *common.Coord, radius float64) bool {

	_, distance := a.GetDirectionAndDistance(goal)
	// 距離がradius m以下の場合
	if distance < math.Abs(radius) {
		return true
	} else {
		return false
	}
}

// CLEAR
// 次の目的地を決める関数
func (a *Agent) DecideNextTransit() {
	// 距離が5m以下の場合
	radius := 5.0
	if a.IsReachedGoal(a.Route.NextTransit, radius) {
		if a.Route.NextTransit != a.Route.Destination {
			for i, point := range a.Route.TransitPoints {
				if point.Longitude == a.Route.NextTransit.Longitude && point.Latitude == a.Route.NextTransit.Latitude {
					if i+1 == len(a.Route.TransitPoints) {
						// すべての経由地を通った場合、nilにする
						a.Route.NextTransit = a.Route.Destination
					} else {
						// 次の経由地を設定する
						a.Route.NextTransit = a.Route.TransitPoints[i+1]
					}
				}
			}
		} else {
			fmt.Printf("\x1b[30m\x1b[47m Arrived Destination! \x1b[0m\n")
		}
	}
}
