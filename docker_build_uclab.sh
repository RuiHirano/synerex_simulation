#!/bin/sh

VERSION=$1
echo "version is ${VERSION}"

docker image build -t ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/nodeid-server:${VERSION} -f nodeserv/Dockerfile .
docker image build -t ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/synerex-server:${VERSION} -f server/Dockerfile .
docker image build -t ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/master-provider:${VERSION} -f provider/master/Dockerfile .
docker image build -t ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/worker-provider:${VERSION} -f provider/worker/Dockerfile .
docker image build -t ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/agent-provider:${VERSION} -f provider/agent/Dockerfile .
docker image build -t ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/visualization-provider:${VERSION} -f provider/visualization/Dockerfile .
docker image build -t ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/gateway-provider:${VERSION} -f provider/gateway/Dockerfile .
docker image build -t ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/simulator:${VERSION} -f cli/Dockerfile .

echo "build finished"