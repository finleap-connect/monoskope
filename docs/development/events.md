**[[Back To Overview]](README.md)**

---

# Implementing `Events`

An `Event` is the notification that something happened in the past. You can view an `Event` as the representation of the reaction to a `Command` after being executed. All `Events` should be represented as verbs in the past tense such as `UserCreated`.

The interface for `Events` is defined in [`pkg/eventsourcing/event.go`](../../pkg/eventsourcing/event.go) along with the internal implementation of an event.
The `Event` types are defined as constants of type [`EventType`](../../pkg/eventsourcing/event.go) in the file [`pkg/domain/constants/events/events.go`](../../pkg/domain/constants/events/events.go).

## Prerequisites

`Events` alone won't make the system work. They only make sense in combination with [`Commands`](commands.md) executed on [`Aggregates`](aggregates.md) which emit [`Event(s)`](events.md).
So before adding a new `Event` you might have a look at the docs about them first.

## Steps to add a new `Event`

In the following step-by-step example we expect you're working on an `Aggregate` called `User`.
You want to add an `Event` when the name of the user has been updated.

When adding a new `Event` you need to do the following:

1. Add new constant for the new `EventType` in the file [`pkg/domain/constants/events/events.go`](../../pkg/domain/constants/events/events.go):
    * The `Event` name should reflect what has happened to what.
    * So your `Event` is about updating the name of an `Aggregate` called `User` the event should be called `UserNameChanged`:

    ```go
    UserNameChanged es.EventType = "UserNameChanged"
    ```

1. Add a new `EventData` definition for your new `Event`:
    * For the `UserNameChanged` event you need `EventData` containing the new name of the user.
    * `EventData` are defined as [proto files](https://developers.google.com/protocol-buffers/docs/proto3), see [`api/domain/eventdata`](../../api/domain/eventdata).
    * Use the [Makefile](../Makefile.md) to generate code from the [proto files](https://developers.google.com/protocol-buffers/docs/proto3) with `make go-protobuf` after you've changed something here.
    * Attention: Not every new `Event` you introduce needs `EventData`. For example the deletion of a user can be free of additional `EventData`.

    ```protobuf
    message UserNameChanged {
        // New name of the user
        string name = 1;
    }
    ```

1. Add the handling for the new `Event` on your [Aggregate](aggregates.md):
    * When an `UserNameChanged` event occurs, the [`User`](../../pkg/domain/aggregates/user.go) `Aggregate` needs to update the name field.
    * Make yourself clear that applying an event to an `Aggregate` needs no validation whatsoever since this already happened before the `Command` on the `Aggregate` has been executed.

    ```go
    // ApplyEvent implements the ApplyEvent method of the Aggregate interface.
    func (a *UserAggregate) ApplyEvent(event es.Event) error {
        switch event.EventType() {
        // other cases are removed for readability
        case events.UserNameChanged:
            data := &eventdata.UserNameChanged{} // The newly introduced EventData
            if err := event.Data().ToProto(data); err != nil {
                return err
            }
            a.Name = data.Name // The new name of the user
        default:
            return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
        }
        return nil
    }
    ```

1. Add the handling for the new `Event` on your `Projector`:
    * When the `UserProjector` handles the new `UserNameChanged` event, it needs to update the name of the `UserProjection`.
    * See the docs on [Projections](projections.md) and [Projectors](projectors.md) for this.

1. Now when will this `Event` be emitted?
Right!
Until now nowhere.
You need to create a new `Command` which changes the user's name.
Head over to the docs about creating [Commands](commands.md) for this.
