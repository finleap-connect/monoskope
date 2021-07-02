**[[Back To Overview]](../README.md)**

---

# Overview

## Installation Quickstart

Monoskope requires a Kubernetes cluster with a number of componenents already installed:

* Cert-Manager
* Prometheus-Operator
* Grafana-Operator
* Vault-Operator

It will install a number of packages as a dedicated instance:

* Ambassador, provided the required CRDs are already installed. Set `ambassador.crds.create=true` 
  in the helm `values.yaml` to roll them out with monoskope.
* CockroachDB
* RabbitMQ


## Docs for development
* [Architecture decisions](architecture-decisions/)
* [Diagrams](diagrams/)
* [Development](development/)
* [Deployment](deployment/)
* [Operation](operation/)

## Helm Charts documentation

* `gateway` helm chart [readme](build/package/helm/gateway/README.md)
* `eventstore` helm chart [readme](build/package/helm/eventstore/README.md)
* `commandhandler` helm chart [readme](build/package/helm/commandhandler/README.md)
* `queryhandler` helm chart [readme](build/package/helm/queryhandler/README.md)
* `cluster-bootstrap-reactor` helm chart [readme](build/package/helm/cluster-bootstrap-reactor/README.md)
* `monoskope` helm chart [readme](build/package/helm/monoskope/README.md)
