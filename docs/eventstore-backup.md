**[[Back To Overview]](README.md)**

---

# Monoskope EventStore Backup

Monoskope is based on ES/CQRS.
Thus the events are the only really crucial part of the system which must not be lost.
The ability to create backups automated and regularly along with restoring events in case of desaster easily is what Monoskope provides.
Additionally the backups are AES encrypted and do not depend on the underlying database used to store events.

## Backup

The helm chart of the EventStore allows to schedule automated backups.
At the moment of writing only S3 is available as backup destination.
In the following yaml snippet you can see the available options:

```yaml
# See build/package/helm/eventstore/values.yaml for the full values file.
backup:
  # -- Enables automated backups for the eventstore
  enabled: true
  alerting:
    # -- Enables alerting for failed backups
    enabled: false
    secondsSinceLastSuccessfulBackup: 86400 # 60x60x24 := 24h
    alertAfter: 1h
  # -- CRON expression defining the backup schedule
  schedule: "0 22 * * *"
  # -- Number of most recent backups to keep
  retentionCount: 7
  # -- Secret containing destination specific secrets, e.g. credentials to s3. The secret will be mounted as environment.
  existingSecretName: "my-s3-credentails"
  # -- Prometheus PushGateway to push metrics to
  prometheusPushgatewayUrl: ""
  # -- Timeout for backup job
  timeout: 1h
  # -- Backup destination, e.g. s3
  destination: {}
    s3:
      endpoint: s3.monoskope.io
      bucket: my-backup-bucket
      region: us-east-1
      disableSSL: false
```

With this deployed Monoskope will automatically create backups every day at 10pm with a retention of 7 days to S3.

## Restore

### Drop existing database

Prior to restoring you need to drop the existing eventstore database.
This is done via a job which can be created via the helm chart of Monoskope:

```yaml
# See build/package/helm/monoskope/values.yaml for the full values file.
cockroachdb:
  dropExistingDatabase: false # ATTENTION: If true the existing database will be dropped on crdb init job, only when restoring backup
```

Now you can either:

1. deploy the whole chart
1. template the chart and apply only one yaml:

```bash
$: HELM_VALUES_FILE=examples/01-monoskope-cluster-values.yaml make helm-template-monoskope
$: kubectl apply -f tmp/monoskope/templates/job-crdb-setup.yaml
```

After this job has run through, the database is new and shiny without any events.

### Restore backup

The helm chart of the EventStore allows to create a restore job.
This job restores a backup from the destination set up under `backup.destination`.

```yaml
# See build/package/helm/eventstore/values.yaml for the full values file.
backup:
  restore:
    # -- Enabling this will deploy a job which restores the backup specified in backup.restore.backupIdentifier from the backup.destination.
    enabled: true
    # -- Identifier of the backup to restore.
    backupIdentifier: "some/backup.tar"
    # -- Timeout for restore job
    timeout: 1h
```

Now you can either:

1. deploy the whole chart
1. template the chart and apply only one yaml:

```bash
$: HELM_VALUES_FILE=examples/01-monoskope-cluster-values.yaml make helm-template-monoskope
$: kubectl apply -f tmp/monoskope/charts/eventstore/templates/job-restore-backup.yaml
```

Be aware that you have to roll through all things with a cached state after the restore.
This applies to the QueryHandler deployment for example.
They will be in the state of before the restore.
A simple roll through will let them update their caches.
This applies to other components like Reactors.
To be safe roll through all deployments of Monoskope:

```bash
$: for resource in `kubectl get deploy -l app.kubernetes.io/part-of=monoskope --no-headers -oname`; do kubectl rollout restart $resource; done
```
