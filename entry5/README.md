# How to run this entry

## Cluster Setup

`./setup.sh` will create a KinD cluster.

## Deploy

`skaffold run` will deploy the Helm chart to the KinD cluster.

## Go CLI

`bin/civ watch --kingdom kingdom-of-foobar` will start the CLI in watch mode.

## Cleanup

`./cleanup.sh` will delete the KinD cluster.
