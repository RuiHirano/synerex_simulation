#!/bin/sh

VERSION=$1
echo "version is ${VERSION}"

docker image pull ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/nodeid-server:${VERSION} 
docker image pull ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/synerex-server:${VERSION} 
docker image pull ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/master-provider:${VERSION} 
docker image pull ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/worker-provider:${VERSION} 
docker image pull ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/agent-provider:${VERSION} 
docker image pull ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/visualization-provider:${VERSION} 
docker image pull ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/gateway-provider:${VERSION} 
docker image pull ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/simulator:${VERSION} 

#docker image pull ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/pod-manager:${VERSION} 

echo "pull finished"