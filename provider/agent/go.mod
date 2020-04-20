module github.com/synerex/synerex_alpha/provider/agent-provider

require (
	github.com/RuiHirano/rvo2-go v1.1.1 // indirect
	github.com/paulmach/orb v0.1.5
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/provider/agent/algorithm v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simutil v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/util v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200207204624-4f3edf09f4f6 // indirect
	google.golang.org/grpc v1.28.1
)

replace (
	github.com/synerex/synerex_alpha/api => ./../../api
	github.com/synerex/synerex_alpha/nodeapi => ./../../nodeapi
	github.com/synerex/synerex_alpha/provider/agent/algorithm => ./algorithm
	github.com/synerex/synerex_alpha/provider/simutil => ../simutil
	github.com/synerex/synerex_alpha/util => ../../util
)

go 1.13
