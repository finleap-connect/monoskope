# Monoskope (m8)

![Monoskope Logo](assets/logo/monoskope.png)

`Monoskope` (short `m8` spelled "mate") implements the management and operation of tenants, users and their roles in a [Kubernetes](https://kubernetes.io/) multi-cluster environment. It fulfills the needs of operators of the clusters as well as the needs of developers using the cloud infrastructure provided by the operators.

## Build status

| `main` | `develop` |
| -- | -- |
|[![pipeline status](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/pipeline.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)|[![pipeline status](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/develop/pipeline.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/develop)
|[![coverage report](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/coverage.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)|[![coverage report](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/develop/coverage.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/develop)|

## Acknowledgments

The implementation of CQRS/ES in Monoskope is not cloned, but inspired by [Event Horizon](https://github.com/looplab/eventhorizon) a CQRS/ES toolkit for Go.
Event Horizon is licensed under Apache License 2.0. A copy of the license is available [here](EVENTHORIZON_LICENSE).

## Documentation

* [Detailed documentation](docs/README.md)
* Architecture and more in [GDrive](https://drive.google.com/drive/folders/1QEewDHF0LwSLr6aUVoHvMWrFgaJfJLty)
* [Makefile documentation](Makefile.md)

### Helm Charts

* `gateway` helm chart [readme](build/package/helm/gateway/README.md)
* `eventstore` helm chart [readme](build/package/helm/eventstore/README.md)
* `commandhandler` helm chart [readme](build/package/helm/commandhandler/README.md)
* `queryhandler` helm chart [readme](build/package/helm/queryhandler/README.md)
* `cluster-bootstrap-reactor` helm chart [readme](build/package/helm/cluster-bootstrap-reactor/README.md)
* `monoskope` helm chart [readme](build/package/helm/monoskope/README.md)
