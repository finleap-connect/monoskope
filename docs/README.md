**[[Back To Overview]](../README.md)**

---

# Monoskope Documentation Overview

## Docs for development

* [Makefile](Makefile.md)
* [Flow charts](flow-charts/README.md)
* [Development guide](development/README.md)

## Deploying Monoskope

### Quick Start

To install your instance in a running Kubernetes cluster:

* log into you Kubernetes cluster
* Select a commit and the `CI_PIPELINE_IID` (the project internal one!) from a successful build for that commit and concatenate without separator to generate a version identifier
  `export VERSION="0.0.0-$COMMIT_HASH$PIPELINE_ID" ; export KUBE_NAMESPACE=<your namespace> make helm-install-monoskope`

### Prerequisites

The following things must be set up in your target K8s Cluster:

* [cert-manager](https://cert-manager.io/docs/) used to generate certificates for
  * Ambassador
  * RabbitMQ
  * CockroachDB
* [vault-operator](https://gitlab.figo.systems/platform/vault-operator) to generate/gather secrets for
  * RabbitMQ
  * Gateway
  from your HashiCorp Vault

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
