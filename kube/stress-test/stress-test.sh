#!/bin/bash

kubectl delete -f ../master.yaml
sleep 1
kubectl delete -f worker.yaml
sleep 1
kubectl delete -f ../simulator.yaml

kubectl apply -f ../master.yaml
sleep 1
kubectl apply -f worker.yaml
sleep 1
kubectl apply -f ../simulator.yaml