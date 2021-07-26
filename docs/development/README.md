**[[Back To Overview]](../README.md)**

---

# Monoskope Development

## Project Structure

```yaml
├── api # protospec's defining the m8 API
│   ├─ domain # Monoskope domain specific API
│   ├─ eventsourcing # Basic EventSourcing API
│   └─ gateway # Gateway API
├── build
│   ├─ ci # CI/CD specific scripts and yaml
│   └─ package # Gateway API
│   |  ├── helm # Monoskope Helm Charts
│   |  └── *.Dockerfile # Monoskope Dockerfiles
├── cmd # Entry points for all m8 binaries 
├── docs # contains the Monoskope documentation
├── internal # contains parts of the code which are not to be exposed as modules.
├── pkg # contains parts of the code which are exposed as modules to be used by monoctl for example.
│   ├─ api # API code which is generated based on the protospec's from the ~/api directory.
│   ├─ domain # Domain specific implementation of Aggregates, Commands, Projections, Projectors, Repositories etc.
│   |  ├── aggregates # Code implementing aggregates
│   |  ├── commands # Code implementing commands
│   |  ├── constants # Constants for aggregate types, command types, event types, roles, scopes and such.
│   |  ├── errors # Defines several errors as variables to be easily comparable.
│   |  ├── metadata # Code implementing event metadata 
│   |  ├── projections # Code implementing projections
│   |  ├── repositories # Code implementing repositories
│   |  ├── commandhandler.go # Code to set up the whole domain for a command handler
│   |  └── queryhandler.go # Code to set up the whole domain for a query handler
│   ├─ eventsourcing # EventSourcing 'framework' of Monoskope. It contains all things like interfaces and basic implementation necessary for ES/CQRS.
│   |  ├── (...)
│   |  ├── errors # Defines several errors as variables to be easily comparable.
│   |  ├── messaging # Contains the implementation for supported message busses.
│   |  ├── repositories # Interfaces for repositories and a basic in-memory repository implementation.
│   |  ├── storage # Contains the implementation for supported event store storages.
│   |  └── *.go # Defining all the interfaces and basic implementations for the framework.
│   ├─ grpc # Code to create gRPC connections, a generic gRPC server implementation and gRPC error handling.
│   ├─ logger # The logging solution 
│   ├─ metrics # A simple prometheus metrics server
│   ├─ operation # Simple implementation to find the current operation mode based on an environment variable.
│   ├─ usecase # Interface and basic implementation for a UseCase based coding pattern.
│   ├─ util # Contains several helper utilities for various use cases.
└── README.md # The entry point for the project docs.
```

## Prerequisites

### Go

* Execute `make go-tools` to get linter and testing binaries

### gRPC

* [Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/)
* [Quickstart - gRPC in Go](https://grpc.io/docs/languages/go/quickstart/)

## Deploy to local KinD cluster

* stand up a KinD cluster:
  `kind create cluster --name m8kind --config build/deploy/kind/kubeadm_conf.yaml --kubeconfig $HOME/.kube/kind-m8kind-config`
* pick up configuration for kubectl and check nodes of newly created cluster
  'export KUBECONFIG="$KUBECONFIG:$HOME/.kube/kind-m8kind-config"
* ensure the finleap helm repo has been added to your local configuration

  ```bash
  helm3 repo add finleap https://artifactory.figo.systems/artifactory/virtual_helm
  helm3 repo update
  ```

* deploy helm charts
  `VERSION="v0.0.15-dev3" HELM_VALUES_FILE=examples/01-monoskope-cluster-values.yaml make helm-install-from-repo-monoskope`

## Testing

There are three general areas of testing within the code:

* **Unit Tests**, that are co-located with the functions that they are testing. These are implemented using [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega) to aid in readability. These should be implemented using TDD and BDD principles.
* **Integration Tests**, the reconstruct the complete software stack automatically. These should be used as the primary test environment for developers to verify that new modules fit with the rest of the system. They are also implemented using Ginkgo and Gomega. They can be found in `internal/integration_test.go`
* **Acceptance Tests**, they ensure that the business rules are correctly implemented. They are written in Gherkin and use [godog](https://github.com/cucumber/godog) to validate against the code.

### Testing Caveats

Sometimes, when the integration test is run repeatetly, the automatic startup and teardown can get stuck and cause the `BeforeSuite` function of that test to fail with the following error message:

```
Unexpected error:
      <*amqp091.Error | 0xc0000aa1e0>: {
          Code: 501,
          Reason: "read tcp 127.0.0.1:40804->127.0.0.1:5672: read: connection reset by peer",
          Server: false,
          Recover: false,
      }
      Exception (501) Reason: "read tcp 127.0.0.1:40804->127.0.0.1:5672: read: connection reset by peer"
  occurred
```

This can be remedied by stoping the relevant docker container running:

```
$ docker ps | head -n 1 ; docker ps | grep rabbitmq
CONTAINER ID   IMAGE                                         COMMAND                  CREATED              STATUS              PORTS                                                                                                                                                    NAMES
<CONTAINER-ID>   some.repo.example.com/bitnami/rabbitmq:3.8.19         "/opt/bitnami/script…"   About a minute ago   Up About a minute   0.0.0.0:5672->5672/tcp, 0.0.0.0:49923->4369/tcp, 0.0.0.0:49922->5671/tcp, 0.0.0.0:49921->15671/tcp, 0.0.0.0:49920->15672/tcp, 0.0.0.0:49919->25672/tcp   rabbitmq
$ docker stop  <CONTAINER-ID>
```

## Event Sourcing & CQRS

### Reading list

* [Greg Young - CQRS and Event Sourcing - Code on the Beach 2014](https://www.youtube.com/watch?v=JHGkaShoyNs)
* [GOTO 2014 • Event Sourcing • Greg Young](https://www.youtube.com/watch?v=8JKjvY4etTY)
* [Greg Young — A Decade of DDD, CQRS, Event Sourcing](https://www.youtube.com/watch?v=LDW0QWie21s)
* [Event sourcing in practice](https://ookami86.github.io/event-sourcing-in-practice/index.html#title.md)
* [CQRS Documents by Greg Young](https://cqrs.files.wordpress.com/2010/11/cqrs_documents.pdf)

### Glossary

| Term | Description |
| --------- | ----------- |
| Aggregate | An entity in your business model (e.g. `User`) which has a state built from the `EventStream` belonging to it. `Commands` can be applied to it and an `Aggregate` may emit `Event(s)` in reaction to a `Command`. |
| Command | A `Command` describes an action that should be performed; it's always named in the imperative tense such as `CreateUser`. |
| Event | An `Event` is the notification that something happened in the past. You can view an event as the representation of the reaction to a `Command` after being executed. All `Events` should be represented as verbs in the past tense such as `UserCreated`. |
| Projection | `Projections` are state which is reconstructed by a sequence of `Events` called an `EventStream`. |
| Repository | `Repositories` are all about storing and retrieving `Projections`. They can be in-memory, a database behind them or whatever. |
| Projector | `Projectors` contain the logic to process `EventStreams` and build up the state of `Projections` out of it. They use repositories to get cached `Projections` and store them. |
| EventStream | The sequence of `Events` belonging to a single aggregate. |
| EventStore | The storage where the `Events` of the system are persisted. |
| Reactor | A component that reacts to `Events` and does any arbitrary action. For example, think of sending an welcoming email to a user after an `UserCreated` event has been observed. |

### Command / Write Side

* [Events](01-events.md)
* [Commands](02-commands.md)
* [Aggregates](03-aggregates.md)

### Query / Read Side

* [Projections](04-projections.md)
* [Projectors](05-projectors.md)
* [Repositories](06-repositories.md)
* [Query Handler](08-queryhandler.md)

### Reactors

* [Reactors](07-reactors.md)
