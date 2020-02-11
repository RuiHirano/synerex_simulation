package simutil

import (
	"math"
	"sort"

	"github.com/google/uuid"
	area "github.com/synerex/synerex_alpha/api/simulation/area"
	common "github.com/synerex/synerex_alpha/api/simulation/common"
)

////////////////////////////////////////////////////////////
/////////////        Area Manager Class        ////////////
///////////////////////////////////////////////////////////

type AreaManager struct {
	Areas []*area.Area
}

func NewAreaManager(initArea *area.Area) *AreaManager {
	am := &AreaManager{
		Areas: []*area.Area{initArea},
	}
	//am.CreateNeighborIds()
	return am
}

// エリアを取得する
func (am *AreaManager) GetArea(areaId uint64) *area.Area {
	for _, areaInfo := range am.Areas {
		if areaInfo.Id == areaId {
			return areaInfo
		}
	}
	return nil
}

// エリアをセットする
func (am *AreaManager) SetArea(areaInfo *area.Area) {
	for i, ai := range am.Areas {
		if ai.Id == areaInfo.Id {
			am.Areas[i] = areaInfo
		}
	}
}

// エリアの追加
func (am *AreaManager) AddArea(areaInfo *area.Area) {
	am.Areas = append(am.Areas, areaInfo)
}

// エリアの削除
func (am *AreaManager) DeleteArea(areaId uint64) {
	newAreas := make([]*area.Area, 0)
	for _, ai := range am.Areas {
		if ai.Id != areaId {
			newAreas = append(newAreas, ai)
		}
	}
	am.Areas = newAreas
}

func (am *AreaManager) DivideArea(areaInfo *area.Area) []*area.Area {
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
	dividedAreaInfos := make([]*area.Area, 0)
	for i, control := range controls {
		uid, _ := uuid.NewRandom()
		dividedAreaInfos = append(dividedAreaInfos, &area.Area{
			Id:              uint64(uid.ID()),
			Name:            "test",
			DuplicateArea:   duplicates[i],
			ControlArea:     control,
			NeighborAreaIds: []uint64{},
		})
	}
	dividedAreaInfos[0].NeighborAreaIds = []uint64{dividedAreaInfos[1].Id}
	dividedAreaInfos[1].NeighborAreaIds = []uint64{dividedAreaInfos[0].Id}

	// 分割されたエリアを追加
	am.AddDevidedArea(areaInfo, dividedAreaInfos)
	// 分割前のエリアを削除
	am.DeleteArea(areaInfo.Id)

	// neighbor更新後の分割されたエリア情報を取得
	areaInfos := make([]*area.Area, 0)
	for _, aInfo := range dividedAreaInfos {
		areaInfos = append(areaInfos, am.GetArea(aInfo.Id))
		logger.Debug("AreaInfo: %v, %v\n", am.GetArea(aInfo.Id).Id, am.GetArea(aInfo.Id).NeighborAreaIds)
	}

	return areaInfos
}

// 分割されたエリアを加える
func (am *AreaManager) AddDevidedArea(sourceArea *area.Area, dividedArea []*area.Area) {
	// 分割されたエリアの隣接エリアを対象に隣接判定を行う
	for _, area1 := range dividedArea {
		for _, areaInfo := range am.Areas {
			if Contains(areaInfo.NeighborAreaIds, sourceArea.Id) {
				isNeighbor := isNeighbor(area1, areaInfo)
				if isNeighbor {
					// 隣接していた場合、加える
					area1.NeighborAreaIds = append(area1.NeighborAreaIds, areaInfo.Id)
					areaInfo.NeighborAreaIds = append(areaInfo.NeighborAreaIds, area1.Id)
				}
			}
			am.SetArea(areaInfo)
		}
		am.AddArea(area1)
	}

}

func (am *AreaManager) CreateNeighborIds() {
	if len(am.Areas) != 0 {
		for i, area1 := range am.Areas {
			var area2 *area.Area
			if i == 0 {
				area2 = am.Areas[len(am.Areas)-1]
			} else {
				area2 = am.Areas[i-1]
			}
			isNeighbor := isNeighbor(area1, area2)
			if isNeighbor {
				// 隣接していた場合、加える
				area1.NeighborAreaIds = append(area1.NeighborAreaIds, area2.Id)
				area2.NeighborAreaIds = append(area2.NeighborAreaIds, area1.Id)
			}
		}
	}
}



// 隣接しているかどうか
func isNeighbor(area1 *area.Area, area2 *area.Area) bool {
	// latかlonが等しい時に、逆(latならlon,lonならlat)が重なっていれば隣接している
	maxLat1, maxLon1, minLat1, minLon1 := GetCoordRange(area1.ControlArea)
	maxLat2, maxLon2, minLat2, minLon2 := GetCoordRange(area2.ControlArea)

	for _, coord1 := range area1.ControlArea {
		for _, coord2 := range area2.ControlArea {
			if coord1.Latitude == coord2.Latitude {
				if (minLon1 < maxLon2 && minLon1 > minLon2) || (maxLon1 < maxLon2 && maxLon1 > minLon2) {
					return true
				}
			}
			if coord1.Longitude == coord2.Longitude {
				if (minLat1 < maxLat2 && minLat1 > minLat2) || (maxLat1 < maxLat2 && maxLat1 > minLat2) {
					return true
				}
			}
		}
	}
	return false
}

func GetCoordRange(coords []*common.Coord) (float64, float64, float64, float64) {
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
