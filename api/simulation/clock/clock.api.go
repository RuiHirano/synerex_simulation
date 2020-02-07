
func NewClock(globalTime float64, timeStep float64, stepNum float64) *Clock{
	c := &Clock{
		GlobalTime: globalTime,
		TimeStep: timeStep,
		StepNum: stepNum,
	}
	return c
}