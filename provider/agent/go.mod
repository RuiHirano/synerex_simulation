module github.com/synerex/synerex_alpha/provider/agent-provider

require (
	github.com/RuiHirano/rvo2-go v1.1.1 // indirect
	github.com/RuiHirano/rvo2-go/src/rvosimulator v0.0.0-20200118052731-21c801eb6c10 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/paulmach/orb v0.1.5
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/simulation v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/api/simulation/common v0.0.0-00010101000000-000000000000
	github.com/synerex/synerex_alpha/api/simulation/provider v0.0.0
	github.com/synerex/synerex_alpha/provider/agent/algorithm v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simutil v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200207204624-4f3edf09f4f6 // indirect
	google.golang.org/grpc v1.27.1
)

replace (
	github.com/synerex/synerex_alpha/api => ./../../api
	github.com/synerex/synerex_alpha/api/common => ./../../api/common
	github.com/synerex/synerex_alpha/api/simulation => ./../../api/simulation
	github.com/synerex/synerex_alpha/api/simulation/agent => ./../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ./../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ./../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/common => ./../../api/simulation/common
	github.com/synerex/synerex_alpha/api/simulation/provider => ./../../api/simulation/provider
	github.com/synerex/synerex_alpha/nodeapi => ./../../nodeapi
	github.com/synerex/synerex_alpha/provider/agent/algorithm => ./algorithm
	github.com/synerex/synerex_alpha/provider/simutil => ../simutil
	github.com/synerex/synerex_alpha/sxutil => ./../../sxutil
)

go 1.13
