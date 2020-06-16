#!/bin/sh

function ConfirmExecution() {

  echo "Please input version"
  echo "ex. latest"
  read version

  echo "----------------------------"
  echo "1. masters-provider"
  echo "2. workers-provider"
  echo "3. visualization-provider"
  echo "4. simulator"
  echo "5. gateway-provider"
  echo "6. all"
  echo "----------------------------"
  echo "Please select build targets"
  echo "ex. 1 2 5 is nodeid, synerex, agent"
  declare -a inputs=()
  read inputs

  echo $inputs

  for input in ${inputs[@]}; do
    if [ $input = '1' ] ; then
        echo "building masters provider..."
        docker image build -t synerex-simulation/masters-provider:${version} -f provider/masters/Dockerfile .

    elif [ $input = '2' ] ; then
        echo "building workers provider..."
        docker image build -t synerex-simulation/workers-provider:${version} -f provider/workers/Dockerfile .

    elif [ $input = '3' ] ; then
        echo "building visualizations provider..."
        docker image build -t synerex-simulation/visualizations-provider:${version} -f provider/visualizations/Dockerfile .

    elif [ $input = '4' ] ; then
        echo "building simulator provider..."
        docker image build -t synerex-simulation/simulator:${version} -f cli/Dockerfile .
    
    elif [ $input = '5' ] ; then
        echo "building gateway provider..."
        docker image build -t synerex-simulation/gateway-provider:${version} -f provider/gateway/Dockerfile .
    
    elif [ $input = '6' ] ; then
        echo "building all"
        docker image build -t synerex-simulation/masters-provider:${version} -f provider/masters/Dockerfile .
        docker image build -t synerex-simulation/gateway-provider:${version} -f provider/gateway/Dockerfile .
        docker image build -t synerex-simulation/workers-provider:${version} -f provider/workers/Dockerfile .
        docker image build -t synerex-simulation/simulator:${version} -f cli/Dockerfile .
        docker image build -t synerex-simulation/visualizations-provider:${version} -f provider/visualizations/Dockerfile .
    else
        echo "unknown number ${input}"

    fi
  done

}

ConfirmExecution

echo "----------------------------"
echo "finished!"