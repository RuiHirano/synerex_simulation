package simutil

import (
	"math"

	common "github.com/synerex/synerex_alpha/api/simulation/common"
)

////////////////////////////////////////////////////////////
////////////        Vector Util            ////////////////
///////////////////////////////////////////////////////////

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
