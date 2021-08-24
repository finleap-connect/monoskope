# Create a new reactor

## in-tree

To make build reactors easier, they can be implemented in-tree

1. implement the behaviour in [`pkg/domain/reactors`](../../pkg/domain/reactors)

1. implement an executable command as subfolder of [`cmd`](../../cmd). As an example check [`cmd/clusterbootstrapreactor/serve.go`](../../cmd/clusterbootstrapreactor/serve.go). 
   The actual code beyond the commandline interface is implemented [`internal/clusterbootstrapreactor/setup.go`](../../internal/clusterbootstrapreactor/setup.go)
   
1. define a helm chart for deployment and make it a dependency for the Monoskope helm chart

Example: the *ClusterBootstrapReactor* is implemented in
* [`cmd/clusterbootstrapreactor/serve.go`](../../cmd/clusterbootstrapreactor/serve.go)
* [`internal/clusterbootstrapreactor/setup.go`](../../internal/clusterbootstrapreactor/setup.go)
* [`pkg/domain/reactors/cluster_bootstrap_reactor.go`](../../pkg/domain/reactors/cluster_bootstrap_reactor.go)
* [`build/package/helm/cluster-bootstrap-reactor/`](../../build/package/helm/cluster-bootstrap-reactor)
* [`build/package/helm/monoskope/Chart.yaml`](../../build/package/helm/monoskope/Chart.yaml) Line 38-40
* [`build/package/helm/monoskope/values.yaml`](../../build/package/helm/monoskope/values.yaml) Line 208ff

## out-of-tree

TBD

## Adding Reactors to integration test

TBD
