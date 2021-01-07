# monoskope

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.1](https://img.shields.io/badge/AppVersion-0.0.1-informational?style=flat-square)

Monoskope implements the management and operation of tenants, users and their roles in a Kubernetes multi-cluster environment.

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../eventstore | eventstore |  |
| file://../gateway | gateway |  |
| https://artifactory.figo.systems/artifactory/virtual_helm | cockroachdb | 5.0.2 |
| https://charts.bitnami.com/bitnami | rabbitmq | 8.6.1 |
| https://kubism.github.io/charts | dex | 1.0.18 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| cockroachdb.conf.cache | string | `"25%"` |  |
| cockroachdb.conf.maxSQLMemory | string | `"25%"` |  |
| cockroachdb.enabled | bool | `true` |  |
| cockroachdb.image.imagePullPolicy | string | `"Always"` |  |
| cockroachdb.image.repository | string | `"gitlab.figo.systems/platform/dependency_proxy/containers/cockroachdb/cockroach"` |  |
| cockroachdb.image.tag | string | `"v20.2.2"` |  |
| cockroachdb.init.annotations."linkerd.io/inject" | string | `"disabled"` |  |
| cockroachdb.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| cockroachdb.serviceMonitor.annotations | object | `{}` |  |
| cockroachdb.serviceMonitor.enabled | bool | `true` |  |
| cockroachdb.serviceMonitor.interval | string | `"30s"` |  |
| cockroachdb.serviceMonitor.labels | object | `{}` |  |
| cockroachdb.statefulset.annotations."linkerd.io/inject" | string | `"disabled"` |  |
| cockroachdb.statefulset.maxUnavailable | int | `1` |  |
| cockroachdb.statefulset.replicas | int | `3` |  |
| cockroachdb.statefulset.resources.limits.cpu | int | `1` |  |
| cockroachdb.statefulset.resources.limits.memory | string | `"2Gi"` |  |
| cockroachdb.statefulset.resources.requests.cpu | string | `"500m"` |  |
| cockroachdb.statefulset.resources.requests.memory | string | `"1Gi"` |  |
| cockroachdb.storage.persistentVolume.size | string | `"20Gi"` |  |
| cockroachdb.tls.certs.clientRootSecret | string | `"monoskope-crdb-root"` |  |
| cockroachdb.tls.certs.nodeSecret | string | `"monoskope-crdb-node"` |  |
| cockroachdb.tls.certs.provided | bool | `true` |  |
| cockroachdb.tls.certs.tlsSecret | bool | `true` |  |
| cockroachdb.tls.enabled | bool | `true` |  |
| dex.certs.grpc.create | bool | `false` |  |
| dex.certs.web.create | bool | `false` |  |
| dex.config.connectors | list | `[]` |  |
| dex.config.enablePasswordDB | bool | `false` |  |
| dex.config.existingSecret | string | `"monoskope-dex-config"` | Name of the secret containing the dex config |
| dex.config.issuer | string | `"https://monoskope.io/dex"` | Domain of the issuer (dex) |
| dex.config.logger.level | string | `"debug"` |  |
| dex.config.oauth2.alwaysShowLoginScreen | bool | `false` |  |
| dex.config.oauth2.skipApprovalScreen | bool | `true` |  |
| dex.config.staticClients[0].id | string | `"gateway"` |  |
| dex.config.staticClients[0].name | string | `"Monoskope Gateway"` |  |
| dex.config.staticClients[0].redirectURIs[0] | string | `"http://localhost:8000"` |  |
| dex.config.staticClients[0].redirectURIs[1] | string | `"http://localhost:18000"` |  |
| dex.config.staticClients[0].secret | string | `"{{ .gatewayAppSecret }}"` |  |
| dex.config.storage.config.database | string | `"dex_db"` |  |
| dex.config.storage.config.port | int | `26257` |  |
| dex.config.storage.config.secret | string | `"monoskope-crdb-client-dex"` | Secret containing the certificates to communicate with the storage backend |
| dex.config.storage.config.ssl.caFile | string | `"/etc/dex/certs/ca.crt"` |  |
| dex.config.storage.config.ssl.certFile | string | `"/etc/dex/certs/client.crt"` |  |
| dex.config.storage.config.ssl.keyFile | string | `"/etc/dex/certs/client.key"` |  |
| dex.config.storage.config.ssl.mode | string | `"verify-ca"` |  |
| dex.config.storage.config.user | string | `"dex"` |  |
| dex.config.storage.type | string | `"postgres"` |  |
| dex.config.web.address | string | `"0.0.0.0"` |  |
| dex.crd.present | bool | `true` |  |
| dex.enabled | bool | `false` |  |
| dex.grpc | bool | `false` |  |
| dex.https | bool | `false` |  |
| dex.image | string | `"ghcr.io/dexidp/dex"` |  |
| dex.imagePullPolicy | string | `"Always"` |  |
| dex.imageTag | string | `"v2.27.0"` |  |
| dex.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| dex.ports.web.containerPort | int | `5556` |  |
| dex.rbac.create | bool | `false` |  |
| dex.replicas | int | `1` |  |
| dex.resources.limits.cpu | string | `"500m"` |  |
| dex.resources.limits.memory | string | `"100Mi"` |  |
| dex.resources.requests.cpu | string | `"100m"` |  |
| dex.resources.requests.memory | string | `"50Mi"` |  |
| dex.serviceAccount.create | bool | `false` |  |
| dex.telemetry | bool | `true` |  |
| eventstore.config.existingSecret | string | `"monoskope-eventstore-config"` |  |
| eventstore.enabled | bool | `true` |  |
| eventstore.nameOverride | string | `"eventstore"` |  |
| eventstore.replicaCount | int | `1` |  |
| fullnameOverride | string | `""` |  |
| gateway.auth.allowRootToken | bool | `false` |  |
| gateway.auth.issuerURL | string | `"https://monoskope.io/dex"` |  |
| gateway.enabled | bool | `true` |  |
| gateway.nameOverride | string | `"gateway"` |  |
| gateway.replicaCount | int | `1` |  |
| global.imagePullSecrets | list | `[]` |  |
| global.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| ingress.enabled | bool | `false` |  |
| ingress.host | string | `"monoskope.io"` |  |
| monitoring.tenant | string | `"finleap-cloud"` |  |
| name | string | `"monoskope"` |  |
| nameOverride | string | `""` |  |
| rabbitmq.auth.existingErlangSecret | string | `"monoskope-rabbitmq-erlang-cookie"` |  |
| rabbitmq.auth.password | string | `""` |  |
| rabbitmq.auth.tls.enabled | bool | `true` |  |
| rabbitmq.auth.tls.existingSecret | string | `"monoskope-rabbitmq-leaf"` |  |
| rabbitmq.auth.tls.failIfNoPeerCert | bool | `true` |  |
| rabbitmq.auth.tls.sslOptionsVerify | string | `"verify_peer"` |  |
| rabbitmq.auth.username | string | `"admin"` |  |
| rabbitmq.enabled | bool | `true` |  |
| rabbitmq.extraConfiguration | string | `"load_definitions = /app/rabbitmq-definitions.json"` |  |
| rabbitmq.image.pullPolicy | string | `"Always"` |  |
| rabbitmq.image.repository | string | `"gitlab.figo.systems/platform/dependency_proxy/containers/bitnami/rabbitmq"` |  |
| rabbitmq.image.tag | string | `"3.8.9"` |  |
| rabbitmq.loadDefinition.enabled | bool | `true` |  |
| rabbitmq.loadDefinition.existingSecret | string | `"monoskope-rabbitmq-load-definition"` |  |
| rabbitmq.metrics.enabled | bool | `true` |  |
| rabbitmq.metrics.grafanaDashboard.enabled | bool | `true` |  |
| rabbitmq.metrics.grafanaDashboard.extraLabels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| rabbitmq.metrics.grafanaDashboard.extraLabels.tenant | string | `"finleap-cloud"` |  |
| rabbitmq.metrics.serviceMonitor.additionalLabels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| rabbitmq.metrics.serviceMonitor.additionalLabels.tenant | string | `"finleap-cloud"` |  |
| rabbitmq.metrics.serviceMonitor.enabled | bool | `true` |  |
| rabbitmq.persistence.enabled | bool | `false` |  |
| rabbitmq.podAnnotations."linkerd.io/inject" | string | `"disabled"` |  |
| rabbitmq.podLabels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| rabbitmq.rbac.create | bool | `false` |  |
| rabbitmq.replicaCount | int | `3` |  |
| rabbitmq.service.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| rabbitmq.service.tlsPort | int | `5671` |  |
| rabbitmq.serviceAccount.create | bool | `false` |  |
| rabbitmq.statefulsetLabels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.4.0](https://github.com/norwoodj/helm-docs/releases/v1.4.0)
