# monoskope

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.1](https://img.shields.io/badge/AppVersion-0.0.1-informational?style=flat-square)

Monoskope implements the management and operation of tenants, users and their roles in a Kubernetes multi-cluster environment.

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../commandhandler | commandhandler |  |
| file://../eventstore | eventstore |  |
| file://../gateway | gateway |  |
| file://../queryhandler | queryhandler |  |
| https://artifactory.figo.systems/artifactory/virtual_helm | cockroachdb | 5.0.2 |
| https://charts.bitnami.com/bitnami | rabbitmq | 8.6.1 |
| https://k8s.ory.sh/helm/charts | oathkeeper | 0.5.3 |
| https://www.getambassador.io | ambassador | 6.5.18 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| ambassador.crds.create | bool | `false` |  |
| ambassador.enableAES | bool | `false` |  |
| ambassador.enabled | bool | `true` |  |
| ambassador.image.repository | string | `"gitlab.figo.systems/platform/dependency_proxy/containers/datawire/ambassador"` |  |
| ambassador.rbac.create | bool | `false` |  |
| ambassador.replicaCount | int | `3` |  |
| ambassador.resources.limits.cpu | int | `2` |  |
| ambassador.resources.limits.memory | string | `"512Mi"` |  |
| ambassador.resources.requests.cpu | string | `"100m"` |  |
| ambassador.resources.requests.memory | string | `"256Mi"` |  |
| ambassador.scope.singleNamespace | bool | `true` |  |
| ambassador.serviceAccount.create | bool | `false` |  |
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
| commandhandler.enabled | bool | `true` |  |
| commandhandler.replicaCount | int | `1` |  |
| eventstore.enabled | bool | `true` |  |
| eventstore.messageBus.existingSecret | string | `"m8-messagebus-client-config"` |  |
| eventstore.replicaCount | int | `1` |  |
| eventstore.storeDatabase.existingSecret | string | `"m8-eventstore-db-config"` |  |
| fullnameOverride | string | `""` |  |
| gateway.auth.issuerURL | string | `"https://monoskope.io/auth"` |  |
| gateway.enabled | bool | `true` |  |
| gateway.replicaCount | int | `1` |  |
| global.imagePullSecrets | list | `[]` |  |
| global.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| ingress.enabled | bool | `false` |  |
| ingress.host | string | `"monoskope.io"` |  |
| messageBus.clientConfigSecretName | string | `"m8-messagebus-client-config"` |  |
| messageBus.routingKeyPrefix | string | `"m8"` |  |
| monitoring.tenant | string | `"finleap-cloud"` |  |
| name | string | `"monoskope"` |  |
| nameOverride | string | `""` |  |
| oathkeeper.enabled | bool | `true` |  |
| oathkeeper.maester.enabled | bool | `true` |  |
| oathkeeper.oathkeeper.config.authenticators.jwt.config.allowed_algorithms[0] | string | `"RS256"` |  |
| oathkeeper.oathkeeper.config.authenticators.jwt.config.jwks_urls[0] | string | `"https://monoskope.dev.finleap.cloud/dex/keys"` |  |
| oathkeeper.oathkeeper.config.authenticators.jwt.config.scope_strategy | string | `"none"` |  |
| oathkeeper.oathkeeper.config.authenticators.jwt.config.trusted_issuers[0] | string | `"https://monoskope.dev.finleap.cloud/dex"` |  |
| oathkeeper.oathkeeper.config.authenticators.jwt.enabled | bool | `true` |  |
| oathkeeper.oathkeeper.config.authorizers.allow.enabled | bool | `true` |  |
| oathkeeper.oathkeeper.config.mutators.noop.enabled | bool | `true` |  |
| oathkeeper.oathkeeper.config.serve.api.port | int | `4456` |  |
| oathkeeper.oathkeeper.image.repository | string | `"gitlab.figo.systems/platform/dependency_proxy/containers/oryd/oathkeeper"` |  |
| oathkeeper.oathkeeper.image.tag | string | `"v0.38.6-beta.1"` |  |
| oathkeeper.service.proxy.enabled | bool | `false` |  |
| queryhandler.enabled | bool | `true` |  |
| queryhandler.messageBus.existingSecret | string | `"m8-messagebus-client-config"` |  |
| queryhandler.replicaCount | int | `1` |  |
| rabbitmq.auth.existingErlangSecret | string | `"monoskope-rabbitmq-erlang-cookie"` |  |
| rabbitmq.auth.password | string | `"w1!!b3r3pl4c3d"` |  |
| rabbitmq.auth.tls.enabled | bool | `true` |  |
| rabbitmq.auth.tls.existingSecret | string | `"monoskope-rabbitmq-leaf"` |  |
| rabbitmq.auth.tls.failIfNoPeerCert | bool | `true` |  |
| rabbitmq.auth.tls.sslOptionsVerify | string | `"verify_peer"` |  |
| rabbitmq.auth.username | string | `"admin"` |  |
| rabbitmq.enabled | bool | `true` |  |
| rabbitmq.extraConfiguration | string | `"load_definitions = /app/rabbitmq-definitions.json\nauth_mechanisms.1 = EXTERNAL\nssl_cert_login_from = common_name"` |  |
| rabbitmq.extraPlugins | string | `"rabbitmq_auth_mechanism_ssl"` |  |
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
