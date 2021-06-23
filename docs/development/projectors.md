**[[Back To Overview]](README.md)**

---

# Implementing Projectors

`Projectors` contain the logic to process `EventStreams` and build up the state of [`Projections`](projections.md) out of it. They use [`Repositories`](repositories.md) to get cached [`Projections`](projections.md) and store them.

## Prerequisites

`Projectors` handle [`Projections`](projections.md).
So have a look at them first.

## Steps to add a new `Projector`

In the guide we create a `Projector` for user [`Projections`](projections.md).

1. Add implementation for your new projector at [`pkg/domain/projectors`](../../pkg/domain/projectors):

    ```go
    package projectors

    import (...)

    type userProjector struct {
        *domainProjector
    }

    func NewUserProjector() es.Projector {
        return &userProjector{
            domainProjector: NewDomainProjector(), // Base projector implementation
        }
    }

    func (u *userProjector) NewProjection(id uuid.UUID) es.Projection {
        return projections.NewUserProjection(id)
    }

    // Project updates the state of the projection according to the given event.
    func (u *userProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
        // Get the actual projection type
        p, ok := projection.(*projections.User)
        if !ok {
            return nil, errors.ErrInvalidProjectionType
        }

        // Apply the changes for the event.
        switch event.EventType() {
        case events.UserCreated:
            data := &eventdata.UserCreated{}
            if err := event.Data().ToProto(data); err != nil {
                return projection, err
            }

            p.Id = event.AggregateID().String()
            p.Email = data.GetEmail()
            p.Name = data.GetName()

            if err := u.projectCreated(event, p.DomainProjection); err != nil {
                return nil, err
            }
        default:
            return nil, errors.ErrInvalidEventType
        }

        p.IncrementVersion()

        return p, nil
    }
    ```

1. Register the projector on QueryHandler [startup](../../pkg/domain/queryhandler.go).
