# Create a new reactor

## in-tree

To make build reactors easier, they can be implemented in-tree

1. implement the behavior in [`pkg/domain/reactors`](../../pkg/domain/reactors)

1. implement an executable command as subfolder of [`cmd`](../../cmd).

1. define a helm chart for deployment and make it a dependency for the Monoskope helm chart

## out-of-tree

TBD

## Adding Reactors to integration test

TBD
