#!/bin/sh

kind create cluster --name tobs

clean_up() {
  kind delete cluster --name tobs
  exit
}

trap clean_up SIGHUP SIGINT SIGTERM
go test -v ./tests/external-db-tests
clean_up