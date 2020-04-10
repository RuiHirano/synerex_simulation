#!/bin/sh

docker image build -t synerex-simulation/nodeid-server:latest -f nodeserv/Dockerfile .
docker image build -t synerex-simulation/synerex-server:latest -f server/Dockerfile .
docker image build -t synerex-simulation/master-provider:latest -f provider/master/Dockerfile .
docker image build -t synerex-simulation/worker-provider:latest -f provider/worker/Dockerfile .
docker image build -t synerex-simulation/agent-provider:latest -f provider/agent/Dockerfile .
docker image build -t synerex-simulation/visualization-provider:latest -f provider/visualization/Dockerfile .
docker image build -t synerex-simulation/gateway-provider:latest -f provider/gateway/Dockerfile .
docker image build -t synerex-simulation/simulator:latest -f cli/Dockerfile .