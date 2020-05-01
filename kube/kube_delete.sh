#!/bin/sh

files="./pod-test/*.yaml"
array=($files)

i=0
for filepath in $files; do
  echo $i: $filepath
  let i++
done

echo "Please select target"
read ti

echo file is ${array[ti]}
kubectl delete -f ${array[ti]}
