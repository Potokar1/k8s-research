# How to run this entry

## Cluster Setup

`./setup.sh` will create a KinD cluster.

## Deploy

`skaffold run` will deploy the Helm chart to the KinD cluster.

## Go CLI

`cd bin` will change to the `bin` directory.

`source <(./civ completion bash)` will enable bash completion (hitting tab in the terminal) for the `civ` CLI.

`./civ kingdoms` will list the kingdoms in the simulation. (from previous entry)

`./civ towns --kingdom <kingdom>` will list the towns in a kingdom.

`./civ workers --kingdom <kingdom> --town <town>` will list the workers in a town.

`./civ logs --kingdom <kingdom> --town <town> --worker <worker>` will list the logs for a worker.

## Cleanup

`./cleanup.sh` will delete the KinD cluster.
