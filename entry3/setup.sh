#!/usr/bin/env bash

set -e
# set -x

if kind get clusters | grep -q "civ-cluster"; then
    echo "Deleting the existing kind cluster"
    kind delete cluster --name=civ-cluster
fi

echo "Creating a kind cluster"
kind create cluster --name=civ-cluster --wait=30s
