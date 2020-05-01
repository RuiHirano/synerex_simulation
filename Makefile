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
		kubectl logs -f worker${arg} -c worker-provider

apply:
		cd kube && bash kube_apply.sh

delete:
		cd kube && bash ./kube_delete.sh

gen-rsc:
		bash ./kube/util/gen-rsc.sh

build:
		bash ./docker_build.sh

push-lab:
		bash ./docker_push_uclab.sh

build-lab:
		bash ./docker_build_uclab.sh

