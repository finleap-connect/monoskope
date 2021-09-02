# Adding new aggregates to query handler

1. create a new `service` in [`api/domain/queryhandler_service.proto`](../../api/domain/queryhandler_service.proto). The appropriate messages should be placed into the relevant [`api/domain/eventdata/`](../../api/domain/eventdata) and [`api/domain/projections/`](../../api/domain/projections) files.

1. add mapping to the ambassador configuration in [`build/package/helm/monoskope/templates/ambassador/ambassador-mapping.yaml`](../../build/package/helm/monoskope/templates/ambassador/ambassador-mapping.yaml). The anem of the new service must be placed after the `/domain.` prefix both for the field `spec.prefix` and `spec.rewrite`
