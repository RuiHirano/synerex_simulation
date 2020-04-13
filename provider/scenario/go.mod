module github.com/synerex/synerex_alpha/provider/scenario-provider

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/logger v1.0.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/labstack/echo v3.3.10+incompatible // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mtfelian/golang-socketio v1.5.2
	github.com/paulmach/orb v0.1.5 // indirect
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/simulation v0.0.0-00010101000000-000000000000
	github.com/synerex/synerex_alpha/api/simulation/agent v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/area v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/clock v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/common v0.0.0-00010101000000-000000000000
	github.com/synerex/synerex_alpha/api/simulation/provider v0.0.0
	github.com/synerex/synerex_alpha/provider/simutil v0.0.0-00010101000000-000000000000
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
	github.com/synerex/synerex_alpha/provider/simutil => ../simutil
	github.com/synerex/synerex_alpha/sxutil => ./../../sxutil
)

go 1.13