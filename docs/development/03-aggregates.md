# Implementing `Aggregates`

An `Aggregate` is an entity in your business model (e.g. `User`) which has a state built from the `EventStream` belonging to it.
[`Commands`](02-commands.md) can be applied to it and an `Aggregate` may emit [`Event(s)`](01-events.md) in reaction to a `Aggregate`.

## Prerequisites

[`Aggregates`](03-aggregates.md) emit [`Event(s)`](01-events.md).
So before adding a new `Aggregate` you might have a look at the docs about [`Events`](01-events.md) first.

## Steps to add a new `Aggregate`

1. Add a new constant for the new `AggregateType` in the file [`pkg/domain/constants/aggregates/aggregates.go`](../../pkg/domain/constants/aggregates/aggregates.go):

    ```go
    User es.AggregateType = "User"
    ```

1. Add a new `Aggregate` implementation at [`pkg/domain/aggregates`](../../pkg/domain/aggregates):
    * The new `Aggregate` must implement the interface defined in [`pkg/eventsourcing/aggregate.go`](../../pkg/eventsourcing/aggregate.go)

    ```go
    package aggregates

    import (...)

    // UserAggregate is an aggregate for Users.
    type UserAggregate struct {
        DomainAggregateBase
        Email            string
        Name             string
    }

    // NewUserAggregate creates a new UserAggregate
    func NewUserAggregate(id uuid.UUID) es.Aggregate {
        return &UserAggregate{
            DomainAggregateBase: DomainAggregateBase{
                BaseAggregate: es.NewBaseAggregate(aggregates.User, id),
            },
        }
    }

    // HandleCommand implements the HandleCommand method of the Aggregate interface.
    func (a *UserAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
        return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
    }

    // ApplyEvent implements the ApplyEvent method of the Aggregate interface.
    func (a *UserAggregate) ApplyEvent(event es.Event) error {
        return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
    }
    ```

1. Register your new `Aggregate` with the [default registry](../../pkg/domain/commandhandler.go).

1. Implement [`Commands`](02-commands.md) and [`Events`](01-events.md) to actually have some logic.

<!-- 
## To create a new aggregate

1. Add command messages (command data, to be more specific) to separate file in `api/domain/commanddata/` folder
1. add code to handle new command to separate file in `pkg/domain/commands/` folder. Ideally copy and apapt existing examples.
1. add aggregate to separate file in `pkg/domain/aggregates/` folder.
1. Add service with query functions to `api/domain/queryhandler_service.proto`
1. Add messages for projection in `api/domain/projections` (ideally in separate `.proto` file)
1. Implement query functiosn in new projection repository in `pkg/domain/repositories/`.
1. Implment projector in `pkg/domain/projectors/` folder. There should be at least one projector per aggregate, but there may be multiple projectors. In order to have one projector handle events by multiple Aggregate types, simply register multiple matchers on the same projector. See `pkg/domain/queryhandler.go` for details. (TODO: create more elaborate examples for later use) -->
