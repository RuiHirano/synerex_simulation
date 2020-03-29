module github.com/synerex/synerex_alpha/api

require (
	github.com/golang/protobuf v1.3.2
	github.com/stretchr/testify v1.2.2
	github.com/synerex/synerex_alpha/api/simulation/agent v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/area v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/clock v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/provider v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
)

replace (
	github.com/synerex/synerex_alpha/api/common => ./common
	github.com/synerex/synerex_alpha/api/simulation/agent => ./simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ./simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ./simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/provider => ./simulation/provider
	github.com/synerex/synerex_alpha/api/simulation => ./simulation
)

go 1.13
