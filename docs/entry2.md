# Entry 2: Pods as Towns and Containers as Citizens

## Introduction

In this entry, I will add the kubernetes concepts of `Pods` and `Containers` to our civilization simulation analogy.

## Pods

### What are Pods?

[Pods](https://kubernetes.io/docs/concepts/workloads/pods/) are the smallest deployable units that are created and managed in Kubernetes.  
Deploying a pod is how you run an application on a Kubernetes cluster.  
Pods are a collection of containers that share a network and storage (resources).  

Pods are run on a node in the cluster. A node is a worker machine in Kubernetes. Think of this as a single computer that is able to run applications.

### Pods as Towns

In our civilization analogy, we can think of a Pod as a Town.  
A Town is a collection of Citizens (Containers) that share resources and work together to accomplish tasks.  

## Containers

### What are Containers?

A [Container](https://kubernetes.io/docs/concepts/containers/) is the executable [container image](https://kubernetes.io/docs/concepts/containers/images/) that contains all of the dependencies and software needed to run an application.  
The important concept to understand regarding containers is that they are standardized and are expected to have the same behavior regardless of the environment they are run in.  
This promise is what makes kubernetes powerful, as it can run applications in a consistent way across different environments.  
It allows for applications to be scaled up and down easily by creating more instances of the container.  
It also allows for redundancy and fault tolerance by running multiple instances of the container across different nodes in the cluster.  

Pods usually have a single container, and to scale a container image, you would create more instances of the pod.  
Our analogy is going to have the more advanced use case of having multiple containers in a pod.  

### Containers as Citizens

In our civilization analogy, we can think of a Container as a Citizen, or a single worker that performs a single task for the Town(Pod).  
A Town(Pod) can have multiple Citizens(Containers) that work to accomplish tasks.  
In my analogy, each Citizen is a Worker that performs a single task. These tasks are creating resources, such as wood, iron, or general goods.  
The Citizens(Containers) share resources such that there cannot be infinite Workers in a Town(Pod).  

In future entries, I will explore how Workers(Containers) can interact with each other through different networking concepts in Kubernetes.

## Project Changes by Tool

### Helm

I have updated the Helm chart to include a Pod and a Container. I also made some changes to the verbiage in the `values.yaml` file to reflect the civilization analogy.  

```yaml
towns:
  - name: town-of-peas
    kingdom: kingdom-of-foobar
    workers:
      - name: woodworker
        message: "I made 1 wood"
      - name: ironworker
        message: "I made 1 iron"
      - name: craftsman
        message: "I made 1 good"
```

I made this value change, which creates a single town in the list of towns. This town has 3 workers.

When it is rendered and applied to the cluster, it will create a Pod with 3 containers.

```yaml
{{- range .Values.towns }}
apiVersion: v1
kind: Pod
metadata:
  name: {{ .name }}
  namespace: {{ .kingdom }}
spec:
    containers:
    {{- range .workers }}
    - name: {{ .name }}
      image: busybox
      command: ['sh', '-c', 'while true; do echo "$(date): {{ .message }}"; sleep 5; done']
    {{- end }}
{{- end }}
```

Given this template, since we only have one town, it will create a single pod in the namespace `kingdom-of-foobar` with 3 containers.  
I have used the `busybox` image for the containers, which is a minimal image that is useful for debugging and testing.  
Right now the only thing each container does is print a predefined message every 5 seconds.

### Go CLI

The `civ` Go CLI has been updated with new commands that allow us to interact with the new resources.  

- `bin/civ towns --kingdom <kingdom>` will list the towns in a kingdom.
- `bin/civ workers --kingdom <kingdom> --town <town>` will list the workers in a town.
- `bin/civ logs --kingdom <kingdom> --town <town> --worker <worker>` will list the logs for a worker.

By utilizing the `clientset` struct from the `k8s.io/client-go/kubernetes` package, I can easily get all the information I need to get data on each resource.

```golang
// ListContainers returns the names of containers in a Pod
func ListContainers(ctx context.Context, namespace, podName string) ([]string, error) {
    clientset := GetClientSet()

    pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }

    var containerNames []string
    for _, container := range pod.Spec.Containers {
        containerNames = append(containerNames, container.Name)
    }
    return containerNames, nil
}
```

For example, this function uses the `clientset` to get the containers in a pod by getting the pod object given the namespace and pod name.  
Then it iterates over the containers in the pod and appends the names to a slice of strings.

## Conclusion

In this entry, we expanded our civilization simulation analogy by introducing Pods as Towns and Containers as Citizens.  
Pods, like Towns, are collections of Containers that share resources and work together to accomplish tasks.  
Containers, like Citizens, are individual workers performing specific roles within a Pod.  
This analogy illustrates how Kubernetes orchestrates applications by grouping containers and managing them efficiently.

[Previous Entry](entry1.md)
[Back to Home](index.md)
