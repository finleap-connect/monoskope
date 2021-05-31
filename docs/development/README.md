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

## Go

* Execute `make go-tools` to get linter and testing binaries

### gRPC

* [Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/)
* [Quickstart - gRPC in Go](https://grpc.io/docs/languages/go/quickstart/)

## To create a new aggregate

1. Add command messages (command data, to be more specific) to separate file in `api/domain/commanddata/` folder
1. add code to handle new command to separate file in `pkg/domain/commands/` folder. Ideally copy and apapt existing examples.
1. add aggregate to separate file in `pkg/domain/aggregates/` folder.
1. Add service with query functions to `api/domain/queryhandler_service.proto`
1. Add messages for projection in `api/domain/projections` (ideally in separate `.proto` file)
1. Implement query functiosn in new projection repository in `pkg/domain/repositories/`.
1. Implment projector in `pkg/domain/projectors/` folder. There should be at least one projector per aggregate, but there may be multiple projectors. In order to have one projector handle events by multiple Aggregate types, simply register multiple matchers on the same projector. See `pkg/domain/queryhandler.go` for details. (TODO: create more elaborate examples for later use)

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

### Command/Write Side

* [Events](events.md)
* [Commands](commands.md)
* [Aggregates](aggregates.md)

### Query/Read Side

* [Projections](projections.md)
* [Projectors](projectors.md)
* [Repositories](repositories.md)
