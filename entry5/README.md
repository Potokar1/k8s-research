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

`./civ shops --kingdom <kingdom> --town <town>` will list the shops in a town.

`./civ logs --kingdom <kingdom> --town <town> --shop <shop>` will list the logs for a shop.

## Testing with test.http

`test.http` is a file that can be used with the REST Client extension in VS Code to test the API.

Before running the tests, each pod will need to be port forwarded to the local machine.  
This is easier to do all at once with `k9s` but can be done individually with the following commands.  
Each need to be run in a separate terminal.  

```bash
kubectl port-forward deployment/craftsman 8080:8080 -n kingdom-of-foobar
kubectl port-forward deployment/ironworker 8081:8080 -n kingdom-of-foobar
kubectl port-forward deployment/stoneworker 8082:8080 -n kingdom-of-foobar
kubectl port-forward deployment/woodworker 8083:8080 -n kingdom-of-foobar
```

Navigate to the `test.http` file in VS Code. and run tests by clicking the `Send Request` link above each request.

## Cleanup

`./cleanup.sh` will delete the KinD cluster.
