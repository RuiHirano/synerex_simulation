#!/bin/sh

#kubectl delete -f volume.yaml
#kubectl delete -f simulator.yaml
kubectl delete -f master.yaml
kubectl delete -f worker.yaml
kubectl delete -f worker2.yaml
kubectl delete -f gateway.yaml

#kubectl apply -f volume.yaml
#kubectl apply -f simulator.yaml
kubectl apply -f master.yaml
sleep 1
kubectl apply -f worker.yaml
kubectl apply -f worker2.yaml
sleep 1
kubectl apply -f gateway.yaml