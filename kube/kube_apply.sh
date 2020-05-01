#!/bin/sh

kubectl apply -f master.yaml
kubectl apply -f worker.yaml
kubectl apply -f worker2.yaml
kubectl apply -f gateway.yaml