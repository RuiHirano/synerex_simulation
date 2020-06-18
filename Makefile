arg=11

node:
		kubectl get nodes -o wide

pod:
		kubectl get pods -o wide

svc:
		kubectl get svc -o wide

master:
		kubectl logs -f master -c master-provider

simulator:
		kubectl exec -it simulator -c simulator bash

worker:
		kubectl logs -f workertest -c worker-provider

database:
		kubectl logs -f database${arg} -c database-provider

agent:
		kubectl logs -f agenttest -c agent-provider

vis:
		kubectl logs -f visualization -c visualization-provider

gateway:
		kubectl logs -f gateway${arg} -c gateway-provider

apply:
		cd kube && bash kube_apply.sh

delete:
		cd kube && bash ./kube_delete.sh

app:
		kubectl apply -f ./kube/util/master-simulator.yaml

del:
		kubectl delete -f ./kube/util/master-simulator.yaml

gen-rsc:
		bash ./kube/util/gen-rsc.sh

build:
		bash ./docker_build.sh

push-lab:
		bash ./docker_push_uclab.sh

build-lab:
		bash ./docker_build_uclab.sh

desc:
		cd kube && bash kube_describe.sh

log:
		cd kube && bash kube_log.sh

exec:
		cd kube && bash kube_exec.sh
