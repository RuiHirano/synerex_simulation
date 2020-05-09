#!/bin/sh

function ConfirmExecution() {

  echo "Please input version"
  echo "ex. latest"
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
  echo "Please select push targets"
  echo "ex. 1 2 5 is nodeid, synerex, agent"
  declare -a inputs=()
  read inputs

  echo $inputs

  for input in ${inputs[@]}; do
    if [ $input = '1' ] ; then
        echo "pusing nodeid server..."
        gcloud docker -- push gcr.io/${projectID}/nodeid-server:${version} 

    elif [ $input = '2' ] ; then
        echo "pusing synerex server..."
        gcloud docker -- push gcr.io/${projectID}/synerex-server:${version} 

    elif [ $input = '3' ] ; then
        echo "pusing master provider..."
        gcloud docker -- push gcr.io/${projectID}/master-provider:${version} 

    elif [ $input = '4' ] ; then
        echo "pusing worker provider..."
        gcloud docker -- push gcr.io/${projectID}/worker-provider:${version} 

    elif [ $input = '5' ] ; then
        echo "pusing agent provider..."
        gcloud docker -- push gcr.io/${projectID}/agent-provider:${version} 

        elif [ $input = '6' ] ; then
        echo "pusing visualization provider..."
        gcloud docker -- push gcr.io/${projectID}/visualization-provider:${version} 

    elif [ $input = '7' ] ; then
        echo "pusing gateway provider..."
        gcloud docker -- push gcr.io/${projectID}/gateway-provider:${version} 

    elif [ $input = '8' ] ; then
        echo "pusing simulator provider..."
        gcloud docker -- push gcr.io/${projectID}/simulator:${version} 

    elif [ $input = '9' ] ; then
        echo "pusing all"
        gcloud docker -- push gcr.io/${projectID}/nodeid-server:${version} 
        gcloud docker -- push gcr.io/${projectID}/synerex-server:${version} 
        gcloud docker -- push gcr.io/${projectID}/master-provider:${version} 
        gcloud docker -- push gcr.io/${projectID}/worker-provider:${version} 
        gcloud docker -- push gcr.io/${projectID}/agent-provider:${version} 
        gcloud docker -- push gcr.io/${projectID}/visualization-provider:${version} 
        gcloud docker -- push gcr.io/${projectID}/gateway-provider:${version} 
        gcloud docker -- push gcr.io/${projectID}/simulator:${version} 

    else
        echo "unknown number ${input}"

    fi
  done

}

ConfirmExecution

echo "----------------------------"
echo "finished!"
