# Entry 3: RETCON: Towns as Deployments and Pods as Shops

## Introduction

In this entry, I will be changing the analogy of the civilization simulation to better reflect the Kubernetes concepts that I am exploring.

After some thought, I realized that the analogy of Pods as Towns and Containers as Citizens was not the best fit for the concepts that I am trying show.

I will explore the changes that I made to the analogy in this entry.

## Labels

### What are Labels and Selectors?

[Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) are key-value pairs that are attached to Kubernetes objects.  
Each object (such as pods, services, deployments) can be assigned labels that can be used to select objects based on those labels.  
This pairing fits well with how we can group objects in our analogy.  
The labels are arbitrary and adhere to a set of rules that kubernetes provides.  
As long as the labels are consistent, we can use them to group objects together.

Labels and Selectors have a LIST and WATCH API that we will use to designate objects as part of a group.

Labels are metadata and each object can be assigned multiple labels using the `metadata.labels` field.

```yaml
metadata:
  labels:
    town: simple-town
    shop: carpenter
```

Any object can be filtered or selected by matching at least 1 of the labels above.

## Deployments

### What are Deployments?

[Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) are a way to manage the creation and scaling of Pods in Kubernetes.

The deployment is the "source of truth" for the Pods that it manages.  
It is a declarative way to define the desired state of the Pods and a deployment controller will continuously attempt to move the current state of the Pods to the desired state.

An important concept of Deployments is that they are scalable.  
They can be increased or decreased the amount of pods they manage based on a set "replica" count.  
There is support for more advanced scaling strategies, such as autoscaling that can scale based on resource usage.  
With custom controllers, we might be able to scale based on other metrics that we can define to complement our analogy.

Deployments define a "template" that is used to create the Pods associated with the deployment.  
Each pod is created based on the template and is continuously watched and managed by the deployment.  
We are also able to provide metadata (such as labels) to the pods defined by the deployment's template.  

### Deployments as a Town

Due to the scalability of Deployments, I have decided to change the analogy to have a group of Deployments represent a Town.  
In combination with a Label `town=<name of town>` we can group deployments together that are part of the same town.  
The collection of Deployments with the same `town` label are what makes up a town in our analogy.  
This allows us to operate and manage all kubernetes objects that are part of the town by selecting on the `town` label.

## Pods - Revisited

### Pods as Shops

I have decided to change the analogy of Pods to Shops.  
Previously, I had Pods as Towns, but I believe that the Deployments + Labels = Towns analogy is a better fit.  
Since the deployments manage the pods, the pods can be thought of as shops that are part of a town.  
The deployments are like a group of owners that can manage multiple shops in a town.  
The pods are given labels to group them together as part of a town, which can enable us to manage them as a group.  

### Why is this better?

Similar to the ephemeral nature of shops in the real world, coming and going based on demand, pods can be created and destroyed by altering the deployment.  
If a certain resource is in high demand, we can create another "clone" of that shop by increasing the replica count of the deployment.  

I also prefer viewing a pod's container image as the "set of workers and tools" that the shop needs.  

## Containers - Revisited

### New Analogy

Instead of having containers as citizens or workers, they are now a static idea of the necessary set of workers and tools that a shop needs to run efficiently.  
It does not make sense to double the set of workers and tools in a shop, due to size and resource constraints.  
This directly relates to the limits and constrains a pod is given in the kubernetes cluster.  

In future entries, I can explore more ways we can use containers to represent how shops (pods) can be used.  

### Created Container Image

For more control over the shops, I implemented a simple REST API that is built and deployed as the main container image for each shop.

Right now, the API has readiness and liveness endpoints which I will touch on later.  
There are also endpoints that can be used to view inventory and a endpoint for requesting that the shop sell resources to another shop.  

The container image, alongside serving the http server, also has a simple logic for shops to create resources.  
How each shop does this is defined by a set of `directions` that can be passed to the shop as a json configuration file.  

Any example of a direction is:

```json
{
    "product": "wood",
    "amount": 1,
    "minimum": 2,
    "interval": 5
}
```

Product describes the resource that the shop will create.  
Amount describes how much of the resource will be created each time interval.  
Minimum describes the minimum amount of the resource that the shop will keep in inventory and won't be willing to sell.  
Interval describes the time in seconds it takes to create the product.  

## Config Maps

### What are Config Maps?

[Config Maps](https://kubernetes.io/docs/concepts/configuration/configmap/) are a way to store configuration data in Kubernetes.  
It's important to note that config maps are not meant to store sensitive data.  
Config maps can be used by pods to consume configuration data.  
The strength of config maps is that they are separate from pods and can be updated and reused by multiple pods.  

### How I am using Config Maps

I create a config map for each shop to hold the directions that the shop will use to create resources.  
The config map supports operations such as defaults and json conversions.  

The config maps are able to be referenced in the deployments with the end goal of giving the pods a volume that can house the directions file.  

## Volumes and Volume Mounts

### What are Volumes and Volume Mounts?

[Volumes](https://kubernetes.io/docs/concepts/storage/volumes/) are a way to store data in Kubernetes.  
For storage, all I want to discuss is that volumes support taking data from a config map and storing it in a pod as a file.  
Given a config map name, I can pick a certain key-value pair and store that in a file in the pod.  

An example:

```yaml
volumes:
- name: config
  configMap:
  name: shop-directions
```

This will create a volume that stores the **value** of the config map key `shop-directions` in a file in the pod.  

The Volume Mount is how we attach the volume to the pod.  

```yaml
volumeMounts:
- name: config
  mountPath: /config
  readOnly: true
```

This will attach the volume defined above to the pod at the path `/config`.  
The pod will be able to read the file stored in the volume.  

## Services

### What are Services?

[Services](https://kubernetes.io/docs/concepts/services-networking/service/) are a way to expose an application running in a pod to other parts of the cluster.  
I will go into more detail about services in future entries as I add support for shops to interact with each other.

## Project Changes by Tool

### Helm

To support the new analogy, I have updated the Helm chart to include Deployments and remove explicitly defining pods.  
I also changed the container images that the deployments use to create the pods to reference the build image in the skaffold section.  
I have added a config map that is used to hold and pass `directions` defined in the `values.yaml` file to the respective pod file systems.  
I also included a basic service that is used to expose the shop to the rest of the cluster. This will be explored in later entries.  

The values.yaml file now has the directions configurations for each shop.  
This allows us to store the configuration for each shop in a file.  
The config map can be used to give the pods created by a deployment a volume (file storage) that can house the directions file.  

### Skaffold

To support the custom built container image, I added a build stage to the skaffold file.  

```yaml
build:
  artifacts:
    - image: ghcr.io/potokar1/k8s-research/entry3/worker
      ko:
        main: ./cmd/civ
        dependencies:
          paths:
            - "**/*.go"
            - go.mod
  local:
    useBuildkit: true
```

This build stage will build a container image using the `ko` tool, using BuildKit instead of Docker to create the image.  
This configuration will catch any changes made in a go file and rebuild the container image, so each change can be deployed to the cluster quickly.  
This image will be available to the deployments in the Helm chart.  

### Go CLI

I have updated the Go CLI to reflect the new analogy.  
The commands are now focused on the kingdoms, towns, and shops.  
Shops are the new command that I transitioned workers into.  
All of the functionality is still supported but the name changes reflect the new analogy.  

The Go portion of the project will also include the server and server logic for http endpoints.  
When the `serve` command is ran, the http server will be started and the endpoints will be available.  
The "worker" behavior will also start, which will be more flushed out in further entries.  

Right now, the interesting logic being ran in the worker component is the creation of resources and how I am able to use liveness and readiness probes to determine if the shop is healthy and able to accept trade requests.

```golang
// Work is the loop that will run the worker until the context is canceled
func (w *Worker) Work(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            for _, direction := range w.directions {
                select {
                case <-time.After(time.Duration(direction.Interval) * time.Second):
                    w.produce(direction)
                case <-ctx.Done():
                    return
                }
            }
            // sleep for a second to prevent a busy loop
            time.Sleep(1 * time.Second)
        }
    }
}
```

This block runs continuously until the server is shut down.  
As defined by the directions, the worker will produce resources at a set interval.  

```golang
// restReady implements the REST API for the ready check
// The worker is ready if it has an inventory greater than a set min.
func (s *Server) restReady(w http.ResponseWriter, r *http.Request) {
    if s.worker.AboveMinimum() {
        w.WriteHeader(http.StatusOK)
        return
    }
    http.Error(w, "Not ready", http.StatusServiceUnavailable)
}

func (w *Worker) AboveMinimum() bool {
    w.mu.Lock()
    defer w.mu.Unlock()
    for _, direction := range w.directions {
        if w.Inventory[direction.Product] < direction.Minimum {
            return false
        }
    }
    return true
}
```

In combination with the readiness probe defined in the deployment, the worker will only accept trade requests if it has enough resources to sell.  
The `AboveMinimum` function is used to determine if the shop has enough resources to sell.  

We are able to define the probes in the deployment as such:

```yaml
readinessProbe:
httpGet:
    path: /ready
    port: http
initialDelaySeconds: 5
periodSeconds: 10
```

This will check the `/ready` endpoint every 10 seconds after the initial delay of 5 seconds.  
If the shop is not ready, it will not accept trade requests.  
This is a simple yet effective way to determine if the shop is healthy and "open" for business.  

## Conclusion

This entry contains a large amount of changes to the core project structure to prepare for future entries.  
We now have the analogy of Deployments and labels as Towns and Pods as Shops.  
We have a custom container image that is built and deployed to the cluster using skaffold.  

I am excited to expand on the shops and how they can request resources from other shops.  
This will allow networking and services to be expanded upon in future entries.  

The addition of the custom container image also gives this project much more control over the behavior of pods.  
The increase in complexity will be interesting to explore as there will be more moving parts to manage.

[Previous Entry](entry2.md)  
[Back to Home](index.md)
