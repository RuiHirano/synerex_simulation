

func NewArea(name string, duplicateArea []*common.Coord, controlArea []*common.Coord) *Area{
	a := &Area{
		Id: 0,
		Name: name,
		DuplicateArea: duplicateArea,
		ControlArea: controlArea,
	}
	return a
}

func (a *Area)Dummy(){

}

