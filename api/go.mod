module github.com/synerex/synerex_alpha/api

require (
	github.com/golang/protobuf v1.3.2
	github.com/stretchr/testify v1.2.2
	github.com/synerex/synerex_alpha/api/simulation/agent v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/area v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/clock v0.0.0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	google.golang.org/grpc v1.22.1
)

replace (
	github.com/synerex/synerex_alpha/api/common => ./common
	github.com/synerex/synerex_alpha/api/simulation/agent => ./simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ./simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ./simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/simulation => ./simulation/simulation
	github.com/synerex/synerex_alpha/api/simulation/provider => ./simulation/provider
)
