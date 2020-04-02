#!/bin/sh

kubectl delete -f master.yaml
kubectl delete -f worker.yaml
kubectl delete -f simulator.yaml