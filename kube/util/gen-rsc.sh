#!/bin/sh

function ConfirmExecution() {



  echo "----------------------------"
  echo "1. for uclab server"
  echo "2. for local server"
  echo "----------------------------"
  echo "Please select target"
  read target

  echo "Please input filename "
  echo "ex. higashiyama-4.yaml"
  read filename

  echo "Please using version "
  echo "ex. latest, 1.5"
  read version

  echo "Please input DevideSquareNum"
  echo "if input 2, 2*2=4 area devided"
  read squareNum

  echo "Please input DuplicateRate "
  echo "if input 0.1, 10% of each area duplicated"
  read duplicateRate

  if [ $input = '1' ] ; then
        go run resource-generator-uclab.go

  else
        go run resource-generator.go

  fi

}
echo arg is $@

ConfirmExecution

echo "----------------------------"
echo "finished!"

