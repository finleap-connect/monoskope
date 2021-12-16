## Context

Flexible Authorization is a mandatory feature for Monoskope(m8).
Flexible role binding will allow more granular control over resource access.
Authorisation process simply consists of 2 steps - evaluation and enforcement

There are currently 2 ways  to implement it:

1. each command and query provide their own policies and enforce them.
2. integration of Open Policy Agent(OPA) for evaluation process and enforce it in a separate component

### Pros and Cons of first approach

#### Pros

* policy definition is closer to the component that needs it, which makes debugging less complex

#### Cons

* non consistent policies for roles
* extra complexity for engineers in support
* mixed policy evaluation and enforcement

### Pros and Cons of OPA integration

Open Policy Agent(OPA) is graduated opensource CNCF project widely adopted by community. It allows us to separate
evaluation process from policy enforcement and provide own language for policy definition(Rego)
This approach will also enable us to utilise Admission Control webhooks of k8s

#### Pros

* responsible only for decision process
* opensource solution with big community
* simplifies further integration in k8s admission control process

#### Cons

* uses internal policy declaration language

## Criteria

* unify the approach to authorization for commands and queries
* allow for maximum extensibility.

## Decision

* integrate OPA into system, enable authorization for commands and queries
* generate and export data layer into OPA policy from reactor based on events
* define policies with Rego
* expose admission control webhook endpoint for further integration with k8s


## Status

Proposed

## Consequences

* Integrate OPA into Monoskope
