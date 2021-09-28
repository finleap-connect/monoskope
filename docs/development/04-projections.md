# Implementing Projections

`Projections` are state which is reconstructed by a sequence of [`Events`](01-events.md) called an `EventStream` by [`Projectors`](05-projectors.md).
In some terms a `Projection` is pretty similar to an [`Aggregate`](03-aggregates.md), but lacks the logic that allows to handle [`Commands`](02-commands.md).
A `Projection` might be often pretty similar to an [`Aggregate`](03-aggregates.md), but there can be different `Projections` for the same [`Aggregate`](03-aggregates.md) based on the context it is used in.
The [`Aggregate`](03-aggregates.md) can change but the `Projections` may not and the other way around.
`Projections` are stored and retrieved from [`Repositories`](06-repositories.md).

## Prerequisites

`Projections` are based on [`Events`](01-events.md)  and the [`Aggregate`](03-aggregates.md) which emit them as a result of [`Commands`](02-commands.md).
So have a look at the other topics to get the whole picture.

## Steps to add a new `Projection`

Let's stick to the User example. You want to create a `Projection` for users.

1. Your new projection must implement the interface [`Projection`](../../pkg/eventsourcing/projection.go).

1. Create a new Projection message for the API at [`api/domain/projections`](../../api/domain/projections).

    ```protobuf
    // User within Monoskope
    message User {
        // Unique identifier of the user (UUID 128-bit number)
        string id = 1;
        // Name of the user
        string name = 2;
        // Email address of the user
        string email = 3;
    }
    ```

1. Create a new `user.go` at [`pkg/domain/projections`](../../pkg/domain/projections)

    ```go
    package projections

    import (...)

    type User struct {
        *DomainProjection // Basic implementation of a projection
        *projections.User // API type for User projections
    }

    func NewUserProjection(id uuid.UUID) eventsourcing.Projection {
        dp := NewDomainProjection()
        return &User{
            DomainProjection: dp,
            User: &projections.User{
                Id: id.String(),
            },
        }
    }

    // ID implements the ID method of the Projection interface.
    func (p *User) ID() uuid.UUID {
        return uuid.MustParse(p.Id)
    }
    ```

1. Now that you have a projection for users you need to add a proper [Projector](05-projectors.md).
