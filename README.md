# Monoskope (m8)

![Monoskope Logo](assets/logo/monoskope.png)

`Monoskope` implements the management and operation of tenants, users and their roles in a [Kubernetes](https://kubernetes.io/) multi-cluster environment. It fullfills the needs of operators of the clusters as well as the needs of developers using the cloud infrastructure provided by the operators.

## Build status

| `main` | `develop` |
| -- | -- |
|[![pipeline status](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/pipeline.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)|[![pipeline status](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/develop/pipeline.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/develop)
|[![coverage report](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/coverage.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)|[![coverage report](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/develop/coverage.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/develop)|

## Prerequisites

The following things must be set up in your target K8s Cluster:

* [cert-manager](https://cert-manager.io/docs/) used to generate certificates for
  * Ambassador
  * RabbitMQ
  * CockroachDB
* [vault-operator](https://gitlab.figo.systems/platform/vault-operator) to generate/gather secrets for
  * RabbitMQ
  * Gateway

## Documentation

### General

* Architecture and more in [GDrive](https://drive.google.com/drive/folders/1QEewDHF0LwSLr6aUVoHvMWrFgaJfJLty)
* The [m8 docs](docs/overview.md)

### Helm Charts

* `gateway` helm chart [readme](build/package/helm/gateway/README.md)
* `eventstore` helm chart [readme](build/package/helm/eventstore/README.md)
* `commandhandler` helm chart [readme](build/package/helm/commandhandler/README.md)
* `queryhandler` helm chart [readme](build/package/helm/queryhandler/README.md)
* `monoskope` helm chart [readme](build/package/helm/monoskope/README.md)
