#!/bin/sh

function ConfirmExecution() {
    docker tag synerex-simulation/nodeid-server:latest gcr.io/ruirui_synerex_simulation/nodeid-server:latest
    docker tag synerex-simulation/synerex-server:latest gcr.io/ruirui_synerex_simulation/synerex-server:latest
    docker tag synerex-simulation/master-provider:latest gcr.io/ruirui_synerex_simulation/master-provider:latest
    docker tag synerex-simulation/worker-provider:latest gcr.io/ruirui_synerex_simulation/worker-provider:latest
    docker tag synerex-simulation/agent-provider:latest gcr.io/ruirui_synerex_simulation/agent-provider:latest
    docker tag synerex-simulation/simulator:latest  gcr.io/ruirui_synerex_simulation/simulator:latest

    docker -- push gcr.io/ruirui_synerex_simulation/nodeid-server:latest
    docker -- push gcr.io/ruirui_synerex_simulation/synerex-server:latest
    docker -- push gcr.io/ruirui_synerex_simulation/master-provider:latest
    docker -- push gcr.io/ruirui_synerex_simulation/worker-provider:latest
    docker -- push gcr.io/ruirui_synerex_simulation/agent-provider:latest
    docker -- push gcr.io/ruirui_synerex_simulation/simulator:latest

}

ConfirmExecution

echo "----------------------------"
echo "finished!"

| ---: | ---: |
| 1 |95000  | 
| 4 |240000 |  
| 9| 540000 |  
| 16| 760000| 
|  25 |1380000  | 