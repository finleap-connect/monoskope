# TDR-01: Restructure events to improve aggregation

## Context

The current event structure lacks some important information like the event id, affected aggregate, and more. Even if it has the needed information the way it’s structured and aggregated doesn't allow for full utilization.

The current event structure looks as follows:

```json
{
  "type": "EventType",
  "timestamp": "TimeStamp",
  "aggregateId": "UUID",
  "aggregateType": "AggregateType",
  "aggregateVersion": "Integer",
  "data": "Base64 encoded",
  "metadata": {
    "component_commit": "Last 7 digits of the commit hash",
    "component_name": "ComponentName",
    "component_version": "vMajor.Minor.Patch",
    "x-auth-email": "Issuer Email",
    "x-auth-id": "Issuer ID",
    "x-auth-issuer": "Issuing source",
    "x-auth-name": "Issuer name"
  }
}
```

## Proposals

### Extending events structure #PROP01

This should help improve events aggregation by adding the following attributes:

* `EventID` identifies a specific event
* `IssuerID` identifies the issuer of the event
* `AffectedAggregateID` identifies the affected aggregate
* `AffectedAggregateType` identifies the type of the affected aggregate to help looking up the affected aggregate by its ID

#### Pros

* ease identifying events when more details are needed.
* ease identifying the issuer. See [UC01](#getting-audit-log-of-users-actions-uc01)
* ease identifying the affected aggregate. See [UC02](#getting-audit-log-of-a-user-uc02).

#### Cons

* increases the event size.
* issuer id is duplicated in `metadata` under `x-auth-id`. See [#PROP01-ENH01](#enhancement-of-extending-events-structure-prop01-enh01)

### Enhancement of extending events structure #PROP01-ENH01

To avoid duplicating the issuer id a store restructure is needed.

Currently, the metadata are stored as a json blob, which prevents a direct use of its attributes to e.g. filter/query by the issuer id.

By saving the metadata in it's owen table and referencing it in the corresponding event, one can easily query the events by the metadata attributes.

#### Cons

* added schema complexity
* joining on each query

## Use-Cases

### Getting audit-log of user's actions #UC01

When getting the audit-log of user's actions (the events caused by the user commands) it is currently not possible to filter/query the events by the issuer id directly as it is stored under the metadata, which is currently stored as a json blob.

[PROP01](#extending-events-structure-prop01) introduces a solution by utilising the attribute `issuerID`.

For more information see [PR #90](https://github.com/finleap-connect/monoskope/pull/90)

### Getting audit-log of a user #UC02

When getting the audit-log of a user (the events affecting the user aggregate) it is currently not possible to filter/query the events by the affected aggregate directly. The only workaround is to aggregate the different event types separately and based on that decode the data (when needed e.g. UserRoleBinding event) to filter the events affecting the considered user then merge and sort the different events collected in one stream.

[PROP01](#extending-events-structure-prop01) introduces a solution by utilising the attributes `AffectedAggregateID` and `AffectedAggregateType`.

## Decision

Modifying events is prohibited in event-sourcing, however, in this case, no data will be changed, only extended. This should persist the state of the system without violating the event-sourcing guidelines.

[PROP01](#extending-events-structure-prop01) is accepted with the following restrictions:

* `EventID` accepted.
* `IssuerID` accepted, by moving it up from the `metadata` to the event itself.
* `AffectedAggregateID` rejected, as one event can affect multiple aggregates.
* `AffectedAggregateType` rejected, same as above.

## Criteria

Modifying events should be avoided as much as possible and is considered (in this case) as an exception therefore no migration feature will be introduced.

* The migration should happen at the store level.
* Upgrading from older versions should require the migration job to run.
* A migration job should be created and documented.

## Status

Approved

## Consequences

* breaking change that requires migration without downtime