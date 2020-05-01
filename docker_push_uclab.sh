
#!/bin/sh

function ConfirmExecution() {

  echo "Please input version"
  echo "ex. latest"
  read version

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
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/nodeid-server:${version} 

    elif [ $input = '2' ] ; then
        echo "pusing synerex server..."
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/synerex-server:${version} 

    elif [ $input = '3' ] ; then
        echo "pusing master provider..."
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/master-provider:${version} 

    elif [ $input = '4' ] ; then
        echo "pusing worker provider..."
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/worker-provider:${version} 

    elif [ $input = '5' ] ; then
        echo "pusing agent provider..."
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/agent-provider:${version} 

        elif [ $input = '6' ] ; then
        echo "pusing visualization provider..."
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/visualization-provider:${version} 

    elif [ $input = '7' ] ; then
        echo "pusing gateway provider..."
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/gateway-provider:${version} 

    elif [ $input = '8' ] ; then
        echo "pusing simulator provider..."
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/simulator:${version} 

    elif [ $input = '9' ] ; then
        echo "pusing all"
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/nodeid-server:${version} 
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/synerex-server:${version} 
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/master-provider:${version} 
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/worker-provider:${version} 
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/agent-provider:${version} 
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/visualization-provider:${version} 
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/gateway-provider:${version} 
        docker image push ucl.nuee.nagoya-u.ac.jp/uclab/synerex-simulation/simulator:${version} 

    else
        echo "unknown number ${input}"

    fi
  done

}

ConfirmExecution

echo "----------------------------"
echo "finished!"

