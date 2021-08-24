**[[Back To Overview]](README.md)**

---

# Implementing Repositories

`Repositories` are all about storing and retrieving [`Projections`](04-projections.md). They can be in-memory, a database behind them or whatever.

## Prerequisites

`Repositories` store and load [`Projections`](04-projections.md) from in-memory or database or whatever you like.

## Steps to add a new `Repository`

We're implementing a `User` repository.
Every Repository must implement the interface [`Repository`](../../pkg/eventsourcing/repository.go).

1. Add a new repository in [`Repositories`](../../pkg/domain/repositories) called `UserRepository`.
    * There are basic implementations like an in-memory repository.
    * This can be used as a quick start.
    * For more complex things like a database backed repository you have to implement it yourself.

1. You now implement methods specific for your type like getting a user by it's email address:

    ```go
    // ByEmail searches for the a user projection by it's email address.
    func (r *userRepository) ByEmail(ctx context.Context, email string) (*projections.User, error) {
    ps, err := r.GetAll(ctx, true)
    if err != nil {
        return nil, err
    }

    for _, u := range ps {
        if email == u.Email {
            return u, nil
        }
    }

    return nil, errors.ErrUserNotFound
    }
    ```

1. Register the repository on QueryHandler [startup](../../pkg/domain/queryhandler.go).
