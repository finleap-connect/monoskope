**[[Back To Overview]](README.md)**

---

# Implementing `Aggregates`

An `Aggregate` is an entity in your business model (e.g. `User`) which has a state built from the `EventStream` belonging to it.
`Aggregates` can be applied to it and an `Aggregate` may emit `Event(s)` in reaction to a `Aggregate`.

## Prerequisites

[`Aggregates`](aggregates.md) emit [`Event(s)`](events.md).
So before adding a new `Aggregate` you might have a look at the docs about [`Events`](events.md) first.

## Steps to add a new `Aggregate`

1. Add a new constant for the new `AggregateType` in the file [`pkg/domain/constants/aggregates/aggregates.go`](../../pkg/domain/constants/aggregates/aggregates.go):

    ```go
    User es.AggregateType = "User"
    ```

1. Add a new `Aggregate` implementation at [`pkg/domain/aggregates`](./../pkg/domain/aggregates):
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

1. Implement [`Commands`](commands.md) and [`Events`](events.md) to actually have some logic.
