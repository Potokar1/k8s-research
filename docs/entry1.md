# Entry 1: Introduction and Kingdoms (Namespaces)

## Introduction

Each entry to this project will progress through the development of a civilization simulation.

The simulation will be built using Kubernetes concepts and hosted on a Kubernetes cluster.

## Namespaces

### Namespaces as Kingdoms

The first concept I will introduce is the concept of [`Namespaces`](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) in Kubernetes.

Namespaces are a way to divide cluster resources between multiple users. They are intended for use in environments with many users spread across multiple teams or projects.  
In my civilization analogy, resources in kubernetes can be thought of as an abstraction of the resources that a kingdom has access to.  
Examples of resources include CPU, memory, storage, and network bandwidth.  
For now we will focus on the concept of Namespaces and save the discussion of resources for future entries.

In the context of our civilization simulation, we will use Namespaces to represent different kingdoms.  
The analogy is that each kingdom has its own resources and is isolated from other kingdoms.  
We can think of the resources that are "namespace-scoped" as belonging to that kingdom.  
A resource is "namespace-scoped" if it is created within a specific namespace and is only accessible from within that namespace.  

In future entries, we will explore different namespaced resources and how they can be used to expand our civilization simulation analogy.

### Namespaces Conclusion

Namespaces are a way to divide cluster resources. Resources that are "namespace-scoped" are only accessible from within that namespace.  
The civilization analogy is that kingdoms(Namespaces) are isolated from each other and have their own resources.  

## Project Overview

### Tools Used

#### KinD

I will be using [KinD](https://kind.sigs.k8s.io/) (Kubernetes in Docker) to create a local Kubernetes cluster.  
For this entry, we are not concerned with what KinD is doing under the hood or what is created when we run `kind create cluster`.  
We will be using the default KinD configuration for now and work with the cluster that is created.

We use KinD in our [setup.sh](../entry1/setup.sh) script to create a local Kubernetes cluster.  
The command we will rely on is `kind create cluster --name=civ-cluster` which creates a cluster named `civ-cluster`.  

#### Helm

I will be using [Helm](https://helm.sh/), specifically Helm Charts` to manage the installation of the kubernetes concepts that we will be exploring.  
Helm charts are a way to define kubernetes resources in a templated way.  
Our current helm chart directory structure is as follows:

```shell
charts/
  civ/
    Chart.yaml
    values.yaml
    templates/
      namespace.yaml
```

The `namespace.yaml` file is a kubernetes resource definition for a namespace.  
It references the `values.yaml` file for the namespace name.  
When we install the helm chart, the files in the `templates/` directory are rendered (templated) and applied to the Kubernetes cluster.  

Currently, the only resource we are creating is a namespace which when rendered and applied to the cluster will add one kingdom(namespace) to our civilization simulation:

```yaml
apiVersion: v1
kind: Namespace
metadata:
    name: kingdom-of-foobar
```

### Skaffold

I will be using [Skaffold](https://skaffold.dev/) to automate the deployment of our helm charts.  
Skaffold is a tool that automates the workflow for building, pushing, and deploying applications to a Kubernetes cluster.  
This tool allows us to define a repeatable deployment process that can be run with a single command.  

To use Skaffold, we define a `skaffold.yaml` file in our first entry.  
Right now, the only thing skaffold is doing is deploying our helm chart to the Kubernetes cluster.  

```yaml
apiVersion: skaffold/v4beta11
kind: Config
metadata:
  name: civ
deploy:
  helm:
    releases:
    - name: civ
      chartPath: charts/civ
```

By referencing the helm chart in the `skaffold.yaml` file, we can use `skaffold run` to deploy the helm chart to the Kubernetes cluster with a single command.

### Golang

I will be using Golang(Go) to write any accompanying code for the simulation.  
For now, I will just be writing a simple command-line application that will interact with the Kubernetes cluster in the terms of our civilization analogy.  

The code will be compiled into an executable that can be run from the command line.

The first command I will implement is `kingdoms` which will list the kingdoms(namespaces) in the Kubernetes cluster.  
It will filter out any namespaces that are not part of our civilization simulation.  

I will also be focusing on using Go because the majority of the Kubernetes client libraries are written in Go.  
This allows me to interact with the Kubernetes API in an idiomatic way.  

The `kingdoms` command will use the Kubernetes client library to list the namespaces in the cluster.

```go
func GetClientSet() *kubernetes.Clientset {
    // load kubeconfig

    kubeconfig := ""
    homeDir, err := os.UserHomeDir()
    if err != nil {
        // slog error
        slog.Error("no home directory found, using in-cluster config or default config if no in-cluster config found")
    } else {
        kubeconfig = filepath.Join(homeDir, ".kube", "config")
    }

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

    // create clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    return clientset
}

// ListNamespaces returns a list of namespaces
func ListNamespaces(ctx context.Context) ([]string, error) {
    // get clientset
    clientset := GetClientSet()

    // list namespaces
    namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }

    // return namespaces
    var ns []string
    for _, namespace := range namespaces.Items {
        ns = append(ns, namespace.Name)
    }
    return ns, nil
}
```

The main logic of the `kingdoms` command is to list the namespaces in the Kubernetes cluster.  
This is done by using `clientset` to interact with the Kubernetes API.  
It does this by requesting the list of namespaces from the Kubernetes API.  
I plan on utilizing `clientset` to interact with the Kubernetes API in the future as well.

[Back to Home](index.md)
