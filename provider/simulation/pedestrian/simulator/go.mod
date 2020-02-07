module github.com/synerex/synerex_alpha/provider/simulation/pedestrian/simulator

require (
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/agent v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/area v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/clock v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/provider v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0

)

replace (
	github.com/synerex/synerex_alpha/api => ../../../../api
	github.com/synerex/synerex_alpha/api/simulation/agent => ../../../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ../../../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ../../../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/common => ../../../../api/simulation/common
	github.com/synerex/synerex_alpha/api/simulation/provider => ../../../../api/simulation/provider
	github.com/synerex/synerex_alpha/provider/simulation/simutil/agent => ../agent
	github.com/synerex/synerex_alpha/nodeapi => ../../../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../../../../sxutil
)
