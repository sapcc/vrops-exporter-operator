# vrops-exporter-operator
The Vrops Exporter operator provides Kubernetes-native deployment and management of Vrops Exporter and related monitoring components. The purpose of this project is to simplify and automate the configuration of a Vrops Exporter based monitoring stack for Kubernetes clusters.

## Description
The Vrops Exporter operator includes the following features, among others:

* Kubernetes Custom Resources: Use Kubernetes custom resources to deploy and manage vrops exporters, vrops inventory, and related components.

* Simplified Deployment Configuration: Configure the fundamentals of vrops exporter like versions, exporter types, scrape intervals, and replicas from a native Kubernetes resource.

* Vrops Exporter Target Configuration: Automatically generate monitoring target configurations based on familiar Kubernetes label queries; no need to learn a vrops exporter specific configuration.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/vrops-exporter-operator:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/vrops-exporter-operator:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing
Feel free to contribute. Open issues and pull requests - we're happy to talk about them. 

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022 SAP SE.

Licensed under the Apache License, Version 2.0 
