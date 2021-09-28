# Monoskope (m8)

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)
[![Build status](https://github.com/finleap-connect/monoskope/actions/workflows/golang.yaml/badge.svg)](https://github.com/finleap-connect/monoskope/actions/workflows/golang.yaml)

`Monoskope` (short `m8` spelled "mate") implements the management and operation of tenants, users and their roles in a [Kubernetes](https://kubernetes.io/) multi-cluster environment. It fulfills the needs of operators of the clusters as well as the needs of developers using the cloud infrastructure provided by the operators.

## Acknowledgments

* The implementation of CQRS/ES in Monoskope is not cloned, but inspired by [looplab/eventhorizon](https://github.com/looplab/eventhorizon), a CQRS/ES toolkit for Go.
* The implementation of the RabbitMQ client is forked from [wagslane/go-rabbitmq](https://github.com/wagslane/go-rabbitmq), a wrapper of streadway/amqp that provides reconnection logic.

## Documentation

Detailed documentation can be found at the [/docs](docs) directory.
