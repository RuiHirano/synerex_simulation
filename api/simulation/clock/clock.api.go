package clock

//sxutil "github.com/synerex/synerex_alpha/sxutil"

func NewClock(globalTime float64, timeStep float64, stepNum uint64) *Clock {
	c := &Clock{
		GlobalTime: globalTime,
		TimeStep:   timeStep,
		StepNum:    stepNum,
	}
	return c
}

func (c *Clock) Forward() {
	c.GlobalTime += c.TimeStep * float64(c.StepNum)
}

func (c *Clock) Backward() {
	c.GlobalTime -= c.TimeStep * float64(c.StepNum)
}
