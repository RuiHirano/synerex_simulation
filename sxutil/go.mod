module github.com/synerex/synerex_alpha/sxutil

require (
	cloud.google.com/go v0.34.0 // indirect
	github.com/bwmarrin/snowflake v0.0.0-20180412010544-68117e6bbede
	github.com/golang/lint v0.0.0-20181217174547-8f45f776aaf1 // indirect
	github.com/golang/mock v1.2.0 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/synerex/synerex_alpha/api v0.0.0
	github.com/synerex/synerex_alpha/api/common v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/agent v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/area v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/clock v0.0.0
	github.com/synerex/synerex_alpha/api/simulation/provider v0.0.0
	github.com/synerex/synerex_alpha/nodeapi v0.0.0
	golang.org/x/oauth2 v0.0.0-20181203162652-d668ce993890 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/genproto v0.0.0-20181221175505-bd9b4fb69e2f // indirect
	google.golang.org/grpc v1.22.1
)

replace (
	github.com/synerex/synerex_alpha/api => ../api
	github.com/synerex/synerex_alpha/api/common => ../api/common
	github.com/synerex/synerex_alpha/api/simulation/agent => ../api/simulation/agent
	github.com/synerex/synerex_alpha/api/simulation/area => ../api/simulation/area
	github.com/synerex/synerex_alpha/api/simulation/clock => ../api/simulation/clock
	github.com/synerex/synerex_alpha/api/simulation/synerex => ../api/simulation/synerex
	github.com/synerex/synerex_alpha/api/simulation/provider => ../api/simulation/provider
	github.com/synerex/synerex_alpha/monitor/monitorapi => ../monitor/monitorapi
	github.com/synerex/synerex_alpha/nodeapi => ../nodeapi
	github.com/synerex/synerex_alpha/sxutil => ../sxutil
)
