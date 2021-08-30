# Monoskope (m8)

[![pipeline status](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/pipeline.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)
[![coverage report](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/coverage.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)

![Monoskope Logo](assets/logo/monoskope.png)

`Monoskope` (short `m8` spelled "mate") implements the management and operation of tenants, users and their roles in a [Kubernetes](https://kubernetes.io/) multi-cluster environment. It fulfills the needs of operators of the clusters as well as the needs of developers using the cloud infrastructure provided by the operators.

## Acknowledgments

* The implementation of CQRS/ES in Monoskope is not cloned, but inspired by [Event Horizon](https://github.com/looplab/eventhorizon) a CQRS/ES toolkit for Go.
Event Horizon is licensed under Apache License 2.0. A copy of the license is available [here](EVENTHORIZON_LICENSE).
* The implementation of the auto-reconnect RabbitMQ Client is forked from [wagslane/go-rabbitmq](https://github.com/wagslane/go-rabbitmq).

## Documentation

* [Makefile documentation](Makefile.md)
* [Detailed documentation](docs)
