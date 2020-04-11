# Synerex Simulation

Simulation Services for Person and Traffic Trip using synerex-alpha

# Requirements

go 1.11 or later (we use go.mod files for module dependencies)
nodejs(10.13.0) / npm(6.4.1) / yarn(1.12.1) for web client development.

# How to start

## for Kubenetes

1. build images

at root directory, run command below.
```
$ bash docker_build.sh
```

2. run master-pod

```
cd kube
kubectl apply -f master.yaml
```

3. run worker-pod

```
cd kube
kubectl apply -f worker.yaml
```

4. run simulator-pod

```
cd kube
$ kubectl apply -f simulator.yaml
```

5. check if running pod normally

```
kubectl get pod
```

6. view visualization-map 
```
kubectl get svc
NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                   AGE
vis-monitor           NodePort    10.110.7.17      <none>        80:31788/TCP              43s
```

Look at vis-monitor PORT. You can view visualized map in http://127.0.0.1:31788

7. send some command

```
$ kubectl get pod
NAME              READY   STATUS    RESTARTS   AGE
simulator-swf6b   1/1     Running   0          156m

// into simulator pod
$ kubectl exec -it simulator-swf6b -c simulator bash
```

- send agents

```
root@simulator-swf6b:/synerex_simulation/cli# /simulator order set agent -n 5
```

- start

```
root@simulator-swf6b:/synerex_simulation/cli# /simulator order start
```

- stop

```
root@simulator-swf6b:/synerex_simulation/cli# /simulator order stop
```

### other tips

- check pod logs

```
$ kubectl get pod
NAME              READY   STATUS    RESTARTS   AGE
master-tpq6k      3/3     Running   0          19m

$ kubectl logs -f master-tpq6k -c master-provider
$ kubectl logs -f master-tpq6k -c synerex-server
$ kubectl logs -f master-tpq6k -c nodeid-server
```