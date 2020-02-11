package area

import (
	"github.com/synerex/synerex_alpha/api/simulation/common"
)

func NewArea(name string, duplicateArea []*common.Coord, controlArea []*common.Coord) *Area {
	a := &Area{
		Id:              0,
		Name:            name,
		DuplicateArea:   duplicateArea,
		ControlArea:     controlArea,
		NeighborAreaIds: []uint64{},
	}
	return a
}

func (a *Area) Dummy() {

}
