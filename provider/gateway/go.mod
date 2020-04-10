module github.com/synerex/synerex_alpha/provider/gateway-provider

require (
	github.com/mtfelian/golang-socketio v1.5.2
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/provider/simutil v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200207204624-4f3edf09f4f6 // indirect
	google.golang.org/grpc v1.28.0
)

replace (
	github.com/synerex/synerex_alpha/api => ./../../api
	github.com/synerex/synerex_alpha/provider/simutil => ../simutil
)

go 1.13
