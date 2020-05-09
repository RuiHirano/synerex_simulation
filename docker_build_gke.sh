#!/bin/sh

function ConfirmExecution() {

  echo "Please input version"
  echo "ex. 1.2"
  read version

  echo "Please input projectID"
  echo "ex. xxxxx-xxxx"
  read projectID

  echo "----------------------------"
  echo "1. nodeid-server"
  echo "2. synerex-server"
  echo "3. master-provider"
  echo "4. worker-provider"
  echo "5. agent-provider"
  echo "6. visualization-provider"
  echo "7. gateway-provider"
  echo "8. simulator"
  echo "9. all"
  echo "----------------------------"
  echo "Please select build targets"
  echo "ex. 1 2 5 is nodeid, synerex, agent"
  declare -a inputs=()
  read inputs

  echo $inputs

  for input in ${inputs[@]}; do
    if [ $input = '1' ] ; then
        echo "building nodeid server..."
        docker image build -t gcr.io/${projectID}/nodeid-server:${version} -f nodeserv/Dockerfile .

    elif [ $input = '2' ] ; then
        echo "building synerex server..."
        docker image build -t gcr.io/${projectID}/synerex-server:${version} -f server/Dockerfile .

    elif [ $input = '3' ] ; then
        echo "building master provider..."
        docker image build -t gcr.io/${projectID}/master-provider:${version} -f provider/master/Dockerfile .

    elif [ $input = '4' ] ; then
        echo "building worker provider..."
        docker image build -t gcr.io/${projectID}/worker-provider:${version} -f provider/worker/Dockerfile .

    elif [ $input = '5' ] ; then
        echo "building agent provider..."
        docker image build -t gcr.io/${projectID}/agent-provider:${version} -f provider/agent/Dockerfile .

        elif [ $input = '6' ] ; then
        echo "building visualization provider..."
        docker image build -t gcr.io/${projectID}/visualization-provider:${version} -f provider/visualization/Dockerfile .

    elif [ $input = '7' ] ; then
        echo "building gateway provider..."
        docker image build -t gcr.io/${projectID}/gateway-provider:${version} -f provider/gateway/Dockerfile .

    elif [ $input = '8' ] ; then
        echo "building simulator provider..."
        docker image build -t gcr.io/${projectID}/simulator:${version} -f cli/Dockerfile .
    
    elif [ $input = '9' ] ; then
        echo "building all"
        docker image build -t gcr.io/${projectID}/nodeid-server:${version} -f nodeserv/Dockerfile .
        docker image build -t gcr.io/${projectID}/synerex-server:${version} -f server/Dockerfile .
        docker image build -t gcr.io/${projectID}/master-provider:${version} -f provider/master/Dockerfile .
        docker image build -t gcr.io/${projectID}/worker-provider:${version} -f provider/worker/Dockerfile .
        docker image build -t gcr.io/${projectID}/agent-provider:${version} -f provider/agent/Dockerfile .
        docker image build -t gcr.io/${projectID}/visualization-provider:${version} -f provider/visualization/Dockerfile .
        docker image build -t gcr.io/${projectID}/gateway-provider:${version} -f provider/gateway/Dockerfile .
        docker image build -t gcr.io/${projectID}/simulator:${version} -f cli/Dockerfile .

    else
        echo "unknown number ${input}"

    fi
  done

}

ConfirmExecution

echo "----------------------------"
echo "finished!"
