# Monoskope (m8)

[![Build status](https://github.com/finleap-connect/monoskope/actions/workflows/golang.yaml/badge.svg)](https://github.com/finleap-connect/monoskope/actions/workflows/golang.yaml)
[![Coverage Status](https://coveralls.io/repos/github/finleap-connect/monoskope/badge.svg?branch=main)](https://coveralls.io/github/finleap-connect/monoskope?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/finleap-connect/monoskope)](https://goreportcard.com/report/github.com/finleap-connect/monoskope)
[![Go Reference](https://pkg.go.dev/badge/github.com/finleap-connect/monoskope.svg)](https://pkg.go.dev/github.com/finleap-connect/monoskope)
[![GitHub release](https://img.shields.io/github/release/finleap-connect/monoskope.svg)](https://github.com/finleap-connect/monoskope/releases)

`Monoskope` (short `m8` spelled "mate") implements the management and operation of tenants, users and their roles in a [Kubernetes](https://kubernetes.io/) multi-cluster environment. It fulfills the needs of operators of the clusters as well as the needs of developers using the cloud infrastructure provided by the operators.

## Documentation

See the [installation guide](https://github.com/finleap-connect/monoskope/blob/main/INSTALL.md) for getting started.
Detailed documentation can be found at the [/docs](docs) directory.

## Acknowledgments

* The implementation of CQRS/ES in Monoskope is not cloned, but inspired by [looplab/eventhorizon](https://github.com/looplab/eventhorizon), a CQRS/ES toolkit for Go.
* The implementation of the RabbitMQ client is forked from [wagslane/go-rabbitmq](https://github.com/wagslane/go-rabbitmq), a wrapper of streadway/amqp that provides reconnection logic.
