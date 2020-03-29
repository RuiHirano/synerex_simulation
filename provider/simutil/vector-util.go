package simutil

import (
	"math"

	api "github.com/synerex/synerex_alpha/api"
)

////////////////////////////////////////////////////////////
////////////        Vector Util            ////////////////
///////////////////////////////////////////////////////////

func Sub(coord1 *api.Coord, coord2 *api.Coord) *api.Coord {
	return &api.Coord{Latitude: coord1.Latitude - coord2.Latitude, Longitude: coord1.Longitude - coord2.Longitude}
}

func Mul(coord1 *api.Coord, coord2 *api.Coord) float64 {
	return coord1.Latitude*coord2.Latitude + coord1.Longitude*coord2.Longitude
}

func Abs(coord1 *api.Coord) float64 {
	return math.Sqrt(Mul(coord1, coord1))
}

func Add(coord1 *api.Coord, coord2 *api.Coord) *api.Coord {
	return &api.Coord{Latitude: coord1.Latitude + coord2.Latitude, Longitude: coord1.Longitude + coord2.Longitude}
}

func Div(coord *api.Coord, s float64) *api.Coord {
	return &api.Coord{Latitude: coord.Latitude / s, Longitude: coord.Longitude / s}
}

type ByAbs struct {
	Coords []*api.Coord
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
