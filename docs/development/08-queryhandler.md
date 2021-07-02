# Adding new aggregates to query handler

1. create a new `service` in `api/domain/queryhandler_service.proto`. The appropriate messages should be placed into the relevant `api/domain/eventdata/` and `api/domain/projections/` files.

1. add mapping to the ambassador configuration in `build/packages/helm/monoskope/templates/ambassador-mapping.yaml`. The anem of the new service must be placed after the `/domain.` prefix both for the field `spec.prefix` and `spec.rewrite`

This document must be rewritten once the query handler has been moved to its own package name. See ticket
[FCLOUD-4174](https://finleap-connect.atlassian.net/browse/FCLOUD-4174) for details
