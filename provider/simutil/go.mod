module github.com/synerex/synerex_alpha/provider/simutil

require (
	github.com/RuiHirano/rvo2-go v0.0.0-20191123125933-81940413d701 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/synerex/synerex_alpha/api v0.0.0
)

replace (
	github.com/synerex/synerex_alpha/api => ../../api
)

go 1.13
