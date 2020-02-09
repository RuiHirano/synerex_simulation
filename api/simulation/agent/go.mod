module github.com/synerex/synerex_alpha/api/simulation/agent

require (
	github.com/golang/protobuf v1.3.2
	github.com/synerex/synerex_alpha/api/simulation/common v0.0.0-00010101000000-000000000000
)

replace (
	github.com/synerex/synerex_alpha/api => ../../../api
	github.com/synerex/synerex_alpha/api/common => ../../../api/common
	github.com/synerex/synerex_alpha/api/simulation/agent => ../../../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ../../../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ../../../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/common => ../../../api/simulation/common
	github.com/synerex/synerex_alpha/api/simulation/provider => ../../../api/simulation/provider
	github.com/synerex/synerex_alpha/nodeapi => ../../../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../../../sxutil
)

go 1.13
