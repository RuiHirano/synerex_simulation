#!/bin/sh

while :
do
    sleep 1
    kubectl top node
    echo '------------------'
done