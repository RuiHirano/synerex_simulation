#!/bin/sh

function ConfirmExecution() {

  echo "Please input projectID"
  echo "ex. latest"
  read projectID

  echo "Please input version"
  echo "ex. latest"
  read version

  echo "----------------------------"
  echo "1. nodeid-server"
  echo "2. synerex-server"
  echo "3. master-provider"
  echo "4. worker-provider"
  echo "5. agent-provider"
  echo "8. simulator"
  echo "10. all"
  echo "----------------------------"
  echo "Please select build targets"
  echo "ex. 1 2 5 is nodeid, synerex, agent"
  declare -a inputs=()
  read inputs

  echo $inputs

  for input in ${inputs[@]}; do
    if [ $input = '1' ] ; then
        echo "building nodeid server..."
        docker image build -t gcr.io/ruirui_synerex_simulation/nodeid-server:${version} -f nodeserv/Dockerfile .

    elif [ $input = '2' ] ; then
        echo "building synerex server..."
        docker image build -t gcr.io/ruirui_synerex_simulation/synerex-server:${version} -f server/Dockerfile .

    elif [ $input = '3' ] ; then
        echo "building master provider..."
        docker image build -t gcr.io/ruirui_synerex_simulation/master-provider:${version} -f provider/master/Dockerfile .

    elif [ $input = '4' ] ; then
        echo "building worker provider..."
        docker image build -t gcr.io/ruirui_synerex_simulation/worker-provider:${version} -f provider/worker/Dockerfile .

    elif [ $input = '5' ] ; then
        echo "building agent provider..."
        docker image build -t gcr.io/ruirui_synerex_simulation/agent-provider:${version} -f provider/agent/Dockerfile .

    elif [ $input = '8' ] ; then
        echo "building simulator provider..."
        docker image build -t gcr.io/ruirui_synerex_simulation/simulator:${version} -f cli/Dockerfile .

    elif [ $input = '10' ] ; then
        echo "building all"
        docker image build -t gcr.io/ruirui_synerex_simulation/nodeid-server:${version} -f nodeserv/Dockerfile .
        docker image build -t gcr.io/ruirui_synerex_simulation/synerex-server:${version} -f server/Dockerfile .
        docker image build -t gcr.io/ruirui_synerex_simulation/master-provider:${version} -f provider/master/Dockerfile .
        docker image build -t gcr.io/ruirui_synerex_simulation/worker-provider:${version} -f provider/worker/Dockerfile .
        docker image build -t gcr.io/ruirui_synerex_simulation/agent-provider:${version} -f provider/agent/Dockerfile .
        docker image build -t gcr.io/ruirui_synerex_simulation/simulator:${version} -f cli/Dockerfile .
    else
        echo "unknown number ${input}"

    fi
  done

}

ConfirmExecution

echo "----------------------------"
echo "finished!"