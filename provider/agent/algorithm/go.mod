module github.com/synerex/synerex_alpha/provider/agent/algorithm

require (
	github.com/RuiHirano/rvo2-go/src/rvosimulator v0.0.0-20200118052731-21c801eb6c10 // indirect
	github.com/paulmach/orb v0.1.5 // indirect
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/nodeapi v0.0.0-20190904050037-afd3805ee7cc // indirect
	google.golang.org/grpc v1.28.1 // indirect

)

replace github.com/synerex/synerex_alpha/api => ../../../api

go 1.13
