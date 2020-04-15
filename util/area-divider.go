package main

import (
	"fmt"
	"math"
)

type Area struct {
	Control   []*Coord
	Duplicate []*Coord
}

type Coord struct {
	Latitude  float64
	Longitude float64
}

func main() {
	areaCoords := []*Coord{
		{Longitude: 136.972139, Latitude: 35.161764},
		{Longitude: 136.972139, Latitude: 35.151326},
		{Longitude: 136.990549, Latitude: 35.151326},
		{Longitude: 136.990549, Latitude: 35.161764},
	}
	devideSquereNum := uint64(2) //2 * 2 = 4 area
	duplicateRate := 0.1         // 10%
	AreaDivider(areaCoords, devideSquereNum, duplicateRate)

}

func AreaDivider(areaCoords []*Coord, divideSquareNum uint64, duplicateRate float64) {
	var areas []*Area
	maxLat, maxLon, minLat, minLon := GetCoordRange(areaCoords)
	for i := 0; i < int(divideSquareNum); i++ { // 横方向
		leftlon := (maxLon-minLon)*float64(i)/float64(divideSquareNum) + minLon
		rightlon := (maxLon-minLon)*(float64(i)+1)/float64(divideSquareNum) + minLon

		for k := 0; k < int(divideSquareNum); k++ { // 縦方向
			bottomlat := (maxLat-minLat)*float64(i)/float64(divideSquareNum) + minLat
			toplat := (maxLat-minLat)*(float64(i)+1)/float64(divideSquareNum) + minLat
			area := &Area{
				Control: []*Coord{
					{Longitude: leftlon, Latitude: toplat},
					{Longitude: leftlon, Latitude: bottomlat},
					{Longitude: rightlon, Latitude: bottomlat},
					{Longitude: rightlon, Latitude: toplat},
				},
				Duplicate: []*Coord{
					{Longitude: leftlon - (rightlon-leftlon)*duplicateRate, Latitude: toplat + (toplat-bottomlat)*duplicateRate},
					{Longitude: leftlon - (rightlon-leftlon)*duplicateRate, Latitude: bottomlat - (toplat-bottomlat)*duplicateRate},
					{Longitude: rightlon + (rightlon-leftlon)*duplicateRate, Latitude: bottomlat - (toplat-bottomlat)*duplicateRate},
					{Longitude: rightlon + (rightlon-leftlon)*duplicateRate, Latitude: toplat + (toplat-bottomlat)*duplicateRate},
				},
			}
			areas = append(areas, area)
		}
	}

	for i, area := range areas {
		fmt.Printf("--------- area %d ---------\n", i)
		duplicateText := `[`
		controlText := `[`
		for _, ctl := range area.Control {
			ctlText := fmt.Sprintf(`{\"latitude\":%v, \"longitude\":%v}`, ctl.Latitude, ctl.Longitude)
			//fmt.Printf("ctl %v\n", ctlText)
			controlText += ctlText
		}
		for _, dpl := range area.Duplicate {
			dplText := fmt.Sprintf(`{\"latitude\":%v, \"longitude\":%v}`, dpl.Latitude, dpl.Longitude)
			//fmt.Printf("dpl %v\n", dplText)
			duplicateText += dplText
		}

		duplicateText += `]`
		controlText += `]`
		fmt.Printf(`"{\"id\":%d, \"name\":\"Nagoya\", \"duplicate_area\": %s, \"control_area\": %s}"`, i, duplicateText, controlText)

	}

}

func GetCoordRange(coords []*Coord) (float64, float64, float64, float64) {
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
