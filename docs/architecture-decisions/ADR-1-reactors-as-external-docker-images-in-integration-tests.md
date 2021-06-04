# ADR-1: Reactors as external docker images in integration tests

## Context

Reactors should be loosely coupled as part of the overall CQRS design. But they can and should
become part of the testing worklow for more complex operations. 

Reactors may be implmeneted either in-tree, i.e as part of the Monoskope source repository, or out
of tree, as separate software projects.
As of the time of writing this, there is only one reactor within the source code of this repository
(`ClusterBoostrapReactor`). We expect the number of both to increase as Monoskope gains functionality.

There are two basic options for implementing integration test for these reactors within Monoskope:

1. a stub/mock framework with a separate in-memory messaging solution for testing
2. building the reactor code into Docker images and executing as Docker containers them as part of
   the test environment setup.

### Pros and Cons of stub framework

#### Pros

* it is simply to set breakpoints and step into and through the code of the reactors, as they are
  part of the single test executable build by `go-test`
* turnaround during development is faster, as there is a single command that will build and compose
  all components needed for the test (`go test`). No need to wait for extenal containers to be
  started. 
  
#### Cons

* out-of-tree reactors cannot be tested by the integration test.

### Pros and Cons of separate container images 

It is worthwhile to note that a certain amount of management code to start and stop the reactor
contains during tests needs to be build.

Also it is unclear whether debuging using breakpoints is possible when separate containers are
used. Any solution utilized may or may not be specifig to a given tooling or IDE.

#### Pros

* the integration test represents a complete setup as found in production deployment as well.
* out-of-tree reactors may be used.

#### Cons

* possibly no debug.

## Criteria

* keep integration test as close to production requirements and context as possible
* allow for maximum extensibility.

## Decision

* implement an extension to the test environment, analog to `pkg/eventsourcing/messaging/testenv.go`
  that will start the reactor conatainers
* automatically update Docker images executed by test
** ensure that in-tree reactors are rebuild from latest source before starting integration test.
** pass version information to test environment for starting the correct version of the container
image. Do not use `latest` as development may be conducted by multiple developers concurrently.

Optionally: build internal API for reactors to register themself to the control plane and add
command to `monoctl` to allow admins to query the registered reactors.

## Status

Proposed

## Consequences

