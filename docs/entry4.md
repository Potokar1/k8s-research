# Entry 4: Adding Functionality to Shops

## Introduction

In this entry, I expand upon the functionality of shops to more closely align with the concept of a marketplace where various shops can interact and share resources.  
I also dive deeper into the kubernetes `Service` for managing communication between different shops.

## Project Changes by Tool

### Helm

To support the expanded shop functionality, I have updated the Helm chart to better represent the changes in shop interactions and resource sharing.

I updated the `values.yaml` file to include new configurations for the shops.  
I updated the service definitions to correctly map to the port the worker service is listening to.

### Go Worker Logic

I added a `Store` string to the directions struct. This is used to specify where the Product comes from.

```golang
type ProductInput struct {
    Product string // Product is the name of the product to buy
    Store   string // Store is the URL of the store to buy from
    Amount  int    // Amount is the quantity of the product to buy
}
```

## Conclusion

[Previous Entry](entry3.md)  
[Back to Home](index.md)
