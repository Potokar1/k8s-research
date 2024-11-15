#!/usr/bin/env bash

set -e
# set -x

if kind get clusters | grep -q "civ-cluster"; then
    echo "Deleting the existing kind cluster"
    kind delete cluster --name=civ-cluster
fi