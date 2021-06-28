**[[Back To Overview]](../README.md)**

---

# Monoskope Documentation Overview

## Installation Quickstart

Monoskope requires a Kubernetes cluster with a number of componenents already installed:

* Cert-Manager
* Prometheus-Operator
* Grafana-Operator
* Vault-Operator

It will install a number of packages as a dedicated instance:

* Ambassador, provided the required CRDs are already installed. Set `ambassador.crds.create=true` in the helm `values.yaml` to roll them out with monoskope.
* CockroachDB
* RabbitMQ


## Docs for development

* [Makefile](Makefile.md)
* [Flow charts](flow-charts/README.md)
* [Get started](development/README.md)

## Using Monoskope

* [Permission Model](permissionModel.md)
* [EventStore Backup & Restore](eventstore-backup.md)
* [Certificate Management](certificate-management.md)

## Docs on Monoskope Helm Charts

* `gateway` helm chart [readme](build/package/helm/gateway/README.md)
* `eventstore` helm chart [readme](build/package/helm/eventstore/README.md)
* `commandhandler` helm chart [readme](build/package/helm/commandhandler/README.md)
* `queryhandler` helm chart [readme](build/package/helm/queryhandler/README.md)
* `monoskope` helm chart [readme](build/package/helm/monoskope/README.md)
