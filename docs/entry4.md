# Entry 4: Adding Functionality to Shops

## Introduction

In this entry I expand the functionality of shops so they more closely resemble a marketplace where shops can interact and share resources.  
I also dive deeper into the Kubernetes **Service** object that manages communication between different shops.

### Service

Expanding on my previous entry, Kubernetes [Services](https://kubernetes.io/docs/concepts/services-networking/service/) expose an application (running in a Pod) to the rest of the cluster.  
A key feature of a Service is the **stable endpoint** it provides, even when the underlying Pods change.

I use the service's stability directly in the setup phase of the worker application. I will go into more detail about this in the Helm section.

This stability is important when a Service handles multiple Pods, which can be created or destroyed at any time.  
It’s similar to a flagship store in a town square that represents a shop, or an online storefront that fulfills requests from any arbitrary location.

A Service enables seamless access to whichever shop instance (Pod) is currently able to fulfil a request.  

#### A Shop's Service

Each shop type will have its own service. Since a shop can have multiple locations, or in kubernetes terms, multiple pods, the service will be used to route purchase requests to any of the shops that match the request.

A Service uses selectors to determine which pods it should route traffic to. Each shop has its own labels that are used to identify it among all shops in a town.

Services are also namespaced objects, meaning that they live at the same level as other objects in the same namespace.  
In my analogy, this means that services are located in the same part of the Kingdom that the town's shops are in.

#### Services are Storefronts in a Town Square

Like a town square, a Service provides a central point for communication and interaction between different shops.  
All of the services in a town can be thought of representing the various shops in a town square. The storefronts of the shops are the services, and the shops themselves are the pods that are running behind the scenes.

When a customer wants to buy a product, they can go to the storefront in the town square (the Service) and request a product from any of the shops (pods) that are available.  
The buyer doesn't care about which specific shop their purchase is being fulfilled by, as long as the product is available.  

## Project Changes by Tool

### Helm

To support the expanded shop functionality, I have updated the Helm chart to better represent the changes in shop interactions and resource sharing.

I updated the service definitions to correctly map to the port the worker service is listening to.  
It now correctly sets up the service to match the town and shop type labels, allowing for proper routing of requests to the correct pods.

The specification is fairly simple, so I will show how each shop type is defined in the Helm chart.  
Given that each shop type is defined in the `values.yaml` file, I can easily create a Service for each shop type.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ .type }}
  namespace: {{ $kingdom }}
  labels:
    town: {{ $town }}
    shop: {{ .type }}
spec:
  selector:
    town: {{ $town }}
    shop: {{ .type }}
  type: ClusterIP
  ports:
    - name: http
      protocol: TCP
      port: 80
      targetPort: 8080
```

Notice that the `targetPort` is set to `8080`, which is just the port that I am using in my Go worker server.  
I am able to reference these services explicitly by their names. I show how I do this in the next section.

I also updated the `values.yaml` file to include new configurations for the shops.  
This updated how shops are defined and added to how resources are created in each shop with a `directions` field.

```yaml
directions:
  - product: iron
    productInputList:
    - product: stone
      store: http://stoneworker
      amount: 3
    - product: wood
      store: http://woodworker
      amount: 10
```

Notice the `store` field in the `productInputList`. This field is a reference to the service that I created for the shop.  
This url won't change and is used by pods to request resources from the specified shop.  

### Go Worker Logic

I added the same `Store` string to the directions struct. This is used to specify where the Product comes from.  
This matches the info in the `values.yaml` file, which is used to create a config map that is mounted into the worker pod.  
The worker can then use this information during runtime.  

```golang
type ProductInput struct {
    Product string // Product is the name of the product to buy
    Store   string // Store is the URL of the store to buy from
    Amount  int    // Amount is the quantity of the product to buy
}
```

I added logic to the Go worker application to properly request for resources determined by the workers `directions` field.  
Now when a worker wants to create a product, it will first satisfy the required inputs by requesting them from the specified stores.  

The worker will use the reference to the store's service to make a request to the store's API.  

The store will wait until it successfully buys all pre-requisite products before it attempts to create the product itself.  
Each store will sell the product only if enough inventory exists above its configured minimum buffer.  

## Conclusion

In future entries, I hope to expand on services by enabling shops to scale with demand, allowing them to be more available and responsive to other shops requesting products.
I also hope to add richer ways to interact with the shops, such as a useful interface for watching the activity of a shop, town, or entire kingdom.

[Previous Entry](entry3.md)  
[Back to Home](index.md)
