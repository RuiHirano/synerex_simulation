module github.com/synerex/synerex_alpha/provider/visualizations-provider

go 1.13

require (
	github.com/RuiHirano/rvo2-go/src/rvosimulator v0.0.0-20200118052731-21c801eb6c10 // indirect
	github.com/google/gops v0.3.10 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/synerex/synerex_alpha/provider/agent/algorithm v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/provider/simutil v0.0.0-00010101000000-000000000000 // indirect
	github.com/synerex/synerex_alpha/util v0.0.0-00010101000000-000000000000 // indirect
	google.golang.org/grpc v1.29.1 // indirect
)

replace (
	github.com/synerex/synerex_alpha/api => ./../../api
	github.com/synerex/synerex_alpha/nodeapi => ./../../nodeapi
	github.com/synerex/synerex_alpha/provider/agent/algorithm => ../../provider/agent/algorithm
	github.com/synerex/synerex_alpha/provider/simutil => ../simutil
	github.com/synerex/synerex_alpha/util => ../../util
)
