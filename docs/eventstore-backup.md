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

## Restore

The helm chart of the EventStore allows to create a restore job specifying the following:

```yaml
backup:
  restore:
    # -- Enabling this will deploy a job which restores the backup specified in backup.restore.backupIdentifier from the backup.destination.
    enabled: true
    # -- Identifier of the backup to restore.
    backupIdentifier: "some/backup.tar"
    # -- Timeout for restore job
    timeout: 1h
```
