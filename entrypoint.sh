#!/bin/sh

cd provider/scenario
go build scenario-provider.go simulator.go
cd ../..

cd provider/agent
go build agent-provider.go simulator.go
cd ../..

cd provider/clock
go build clock-provider.go simulator.go
cd ../..

cd provider/gateway
go build gateway-provider.go
cd ../..

cd provider/visualization
go build visualization-provider.go simulator.go
cd ../..

cd server
go build synerex-server.go message-store.go
cd ..

cd nodeserv
go build nodeid-server.go
cd ..

cd monitor
go build monitor-server.go
cd ..

cd cli
go build monitor-server.go
cd ..

cd provider/scenario
./scenario-provider -synerex 127.0.0.1:9000 -nodeid 127.0.0.1:9100 -monitor 127.0.0.1:9400 -simulator 127.0.0.1:3000 -vis 127.0.0.1:9300 -areaId 0