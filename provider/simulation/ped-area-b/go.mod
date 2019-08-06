module github.com/synerex/synerex_alpha/provider/simulation/ped-area-b

require (
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/sxutil v0.0.0
	google.golang.org/grpc v1.17.0
)

replace (
	github.com/synerex/synerex_alpha/api => ./../../../api
	github.com/synerex/synerex_alpha/api/adservice => ../../../api/adservice
	github.com/synerex/synerex_alpha/api/fleet => ../../../api/fleet
	github.com/synerex/synerex_alpha/api/library => ../../../api/library
	github.com/synerex/synerex_alpha/api/ptransit => ../../../api/ptransit
	github.com/synerex/synerex_alpha/api/rideshare => ../../../api/rideshare
	github.com/synerex/synerex_alpha/api/routing => ../../../api/routing
	github.com/synerex/synerex_alpha/api/simulation/agent => ./../../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ./../../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ./../../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/common => ./../../../api/common
	github.com/synerex/synerex_alpha/nodeapi => ./../../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ./../../../sxutil
)