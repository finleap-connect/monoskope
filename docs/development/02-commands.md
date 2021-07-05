**[[Back To Overview]](README.md)**

---

# Implementing `Commands`

A `Command` describes an action that should be performed; it's always named in the imperative tense such as `CreateUser`.

The interface for `Commands` is defined in [`pkg/eventsourcing/command.go`](../../pkg/eventsourcing/command.go) along with a basic implementation of a `Command`.
The `Command` types are defined as constants of type [`CommandType`](../../pkg/eventsourcing/command.go) in the file [`pkg/domain/constants/commands/commands.go`](../../pkg/domain/constants/commands/commands.go).

## Prerequisites

`Commands` are executed on [`Aggregates`](aggregates.md) which emit [`Event(s)`](events.md).
So before adding a new `Command` you might have a look at the docs about them first.

## Steps to add a new `Command`

1. Add a new constant for the new `CommandType` in the file [`pkg/domain/constants/commands/commands.go`](../../pkg/domain/constants/commands/commands.go):
    * The `Command` name should reflect what should happened to what.
    * So your `Command` is about updating the name of an `Aggregate` called `User` the `Command` should be called something like `UpdateUserName`:

    ```go
    UpdateUserName es.CommandType = "UpdateUserName"
    ```

1. Add a new `CommandData` definition for your new `Command`:
    * For the new `Commmand` called `UpdateUserName` you need `CommandData` containing the new name of the user.
    * `CommandData` are defined as [proto files](https://developers.google.com/protocol-buffers/docs/proto3), see [`api/domain/commanddata`](../../api/domain/commanddata).
    * Use the [Makefile](../Makefile.md) to generate code from the [proto files](https://developers.google.com/protocol-buffers/docs/proto3) with `make go-protobuf` after you've changed something here.

    ```protobuf
    message UpdateUserName {
        // New name of the user
        string name = 1;
    }
    ```

1. Add a new `Command` implementation at [`pkg/domain/commands`](./../pkg/domain/commands):
    * The new `Command` must implement the interface defined in [`pkg/eventsourcing/command.go`](../../pkg/eventsourcing/command.go)

    ```go
    package commands

    import (...)

    // Register the command with default command registry
    func init() {
        es.DefaultCommandRegistry.RegisterCommand(NewUpdateUserNameCommand)
    }

    // UpdateUserNameCommand is a command for updating a users name.
    type UpdateUserNameCommand struct {
        *es.BaseCommand // Basic command implementation from framework to need less code
        cmdData.UpdateUserName // The command data defined before
    }

    // NewUpdateUserNameCommand creates an UpdateUserNameCommand.
    func NewUpdateUserNameCommand(id uuid.UUID) es.Command {
        return &UpdateUserNameCommand{
            BaseCommand: es.NewBaseCommand(id, aggregates.User, commands.UpdateUserName), // Define that this command acts on Users and define it's type.
        }
    }

    // Generic data unmarshalled to type specific command data
    func (c *UpdateUserNameCommand) SetData(a *anypb.Any) error {
        return a.UnmarshalTo(&c.UpdateUserName)
    }

    // Policies returns the Role/Scope combination allowed to execute.
    func (c *UpdateUserNameCommand) Policies(ctx context.Context) []es.Policy {
        return []es.Policy{
            es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins to update a user name
        }
    }
    ```

1. Add the handling for the new `Command` on your [Aggregate](aggregates.md):
    * When an `UpdateUserName` `Command` is executed, the [`User`](../../pkg/domain/aggregates/user.go) `Aggregate` needs to check authorization and validate the new name.

    ```go
    // HandleCommand implements the HandleCommand method of the Aggregate interface.
    func (a *UserAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
        // Authorization and validation is removed for clarity.
        switch cmd := cmd.(type) {
        // Other cases are removed for readability
        case *commands.UpdateUserNameCommand:
            eventData := es.ToEventDataFromProto(&eventdata.UserNameChanged{
                Name:  cmd.GetName(),
            })
            _ = a.AppendEvent(ctx, events.UserNameChanged, eventData)
            return nil
        }
        return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
    }
    ```
