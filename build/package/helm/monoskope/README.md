# monoskope

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.1](https://img.shields.io/badge/AppVersion-0.0.1-informational?style=flat-square)

Monoskope implements the management and operation of tenants, users and their roles in a Kubernetes multi-cluster environment.

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../cluster-bootstrap-reactor | cluster-bootstrap-reactor |  |
| file://../commandhandler | commandhandler |  |
| file://../eventstore | eventstore |  |
| file://../gateway | gateway |  |
| file://../queryhandler | queryhandler |  |
| https://artifactory.figo.systems/artifactory/virtual_helm | cockroachdb | 5.0.2 |
| https://charts.bitnami.com/bitnami | rabbitmq | 8.6.1 |
| https://getambassador.io | ambassador | 6.7.11 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| ambassador.crds.create | bool | `false` |  |
| ambassador.enableAES | bool | `false` |  |
| ambassador.enabled | bool | `true` |  |
| ambassador.image.repository | string | `"gitlab.figo.systems/platform/dependency_proxy/containers/datawire/ambassador"` |  |
| ambassador.image.tag | string | `"1.12.4"` |  |
| ambassador.metrics.serviceMonitor.enabled | bool | `true` |  |
| ambassador.metrics.serviceMonitor.selector.release | string | `"monitoring"` |  |
| ambassador.metrics.serviceMonitor.selector.tenant | string | `"finleap-cloud"` |  |
| ambassador.rbac.create | bool | `false` |  |
| ambassador.replicaCount | int | `3` |  |
| ambassador.resources.limits.cpu | int | `4` |  |
| ambassador.resources.limits.memory | string | `"1000Mi"` |  |
| ambassador.resources.requests.cpu | string | `"100m"` |  |
| ambassador.resources.requests.memory | string | `"512Mi"` |  |
| ambassador.scope.singleNamespace | bool | `true` |  |
| ambassador.serviceAccount.create | bool | `false` |  |
| cluster-bootstrap-reactor.enabled | bool | `true` |  |
| cluster-bootstrap-reactor.keySecret.name | string | `"m8-authentication"` |  |
| cluster-bootstrap-reactor.messageBus.existingSecret | string | `"m8-messagebus-client-config"` |  |
| cockroachdb.conf.cache | string | `"25%"` |  |
| cockroachdb.conf.maxSQLMemory | string | `"25%"` |  |
| cockroachdb.dropExistingDatabase | bool | `false` |  |
| cockroachdb.enabled | bool | `true` |  |
| cockroachdb.image.imagePullPolicy | string | `"Always"` |  |
| cockroachdb.image.repository | string | `"gitlab.figo.systems/platform/dependency_proxy/containers/cockroachdb/cockroach"` |  |
| cockroachdb.image.tag | string | `"v20.2.2"` |  |
| cockroachdb.init.annotations."linkerd.io/inject" | string | `"disabled"` |  |
| cockroachdb.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| cockroachdb.serviceMonitor.annotations | object | `{}` |  |
| cockroachdb.serviceMonitor.enabled | bool | `true` |  |
| cockroachdb.serviceMonitor.interval | string | `"1m"` |  |
| cockroachdb.serviceMonitor.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| cockroachdb.serviceMonitor.labels.release | string | `"monitoring"` |  |
| cockroachdb.serviceMonitor.labels.tenant | string | `"finleap-cloud"` |  |
| cockroachdb.serviceMonitor.scrapeTimeout | string | `"10s"` |  |
| cockroachdb.statefulset.annotations."linkerd.io/inject" | string | `"disabled"` |  |
| cockroachdb.statefulset.maxUnavailable | int | `1` |  |
| cockroachdb.statefulset.replicas | int | `3` |  |
| cockroachdb.statefulset.resources.limits.cpu | int | `2` |  |
| cockroachdb.statefulset.resources.limits.memory | string | `"2Gi"` |  |
| cockroachdb.statefulset.resources.requests.cpu | string | `"500m"` |  |
| cockroachdb.statefulset.resources.requests.memory | string | `"1Gi"` |  |
| cockroachdb.storage.persistentVolume.size | string | `"1Gi"` |  |
| cockroachdb.tls.certs.clientRootSecret | string | `"m8-crdb-root"` |  |
| cockroachdb.tls.certs.nodeSecret | string | `"m8-crdb-node"` |  |
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
| gateway.auth.issuerURL | string | `"https://your-idp.com"` |  |
| gateway.enabled | bool | `true` |  |
| gateway.keySecret.name | string | `"m8-authentication"` |  |
| gateway.replicaCount | int | `1` |  |
| global.imagePullSecrets | list | `[]` |  |
| global.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| hosting.domain | string | `"monoskope.io"` |  |
| hosting.issuer | string | `""` |  |
| messageBus.clientConfigSecretName | string | `"m8-messagebus-client-config"` |  |
| messageBus.routingKeyPrefix | string | `"m8"` |  |
| monitoring.tenant | string | `"finleap-cloud"` |  |
| name | string | `"monoskope"` |  |
| nameOverride | string | `""` |  |
| pki.authentication.keySecretName | string | `"m8-authentication"` |  |
| pki.certificates.duration | string | `"2160h"` |  |
| pki.certificates.renewBefore | string | `"1440h"` |  |
| pki.enabled | bool | `true` |  |
| pki.issuer.ca.enabled | bool | `true` |  |
| pki.issuer.ca.existingTrustAnchorSecretName | string | `"m8-trust-anchor"` |  |
| pki.issuer.ca.secretVersion | int | `1` |  |
| pki.issuer.vault.enabled | bool | `false` |  |
| pki.issuer.vault.path | string | `"pki_int/sign/example-dot-com"` |  |
| pki.issuer.vault.server | string | `"https://vault.local"` |  |
| queryhandler.enabled | bool | `true` |  |
| queryhandler.messageBus.existingSecret | string | `"m8-messagebus-client-config"` |  |
| queryhandler.replicaCount | int | `1` |  |
| rabbitmq.auth.existingErlangSecret | string | `"m8-rabbitmq-erlang-cookie"` |  |
| rabbitmq.auth.password | string | `"w1!!b3r3pl4c3d"` |  |
| rabbitmq.auth.tls.enabled | bool | `true` |  |
| rabbitmq.auth.tls.existingSecret | string | `"m8-rabbitmq-leaf"` |  |
| rabbitmq.auth.tls.failIfNoPeerCert | bool | `true` |  |
| rabbitmq.auth.tls.sslOptionsVerify | string | `"verify_peer"` |  |
| rabbitmq.auth.username | string | `"eventstore"` |  |
| rabbitmq.enabled | bool | `true` |  |
| rabbitmq.extraConfiguration | string | `"load_definitions = /app/rabbitmq-definitions.json\nauth_mechanisms.1 = EXTERNAL\nssl_cert_login_from = common_name\nssl_options.depth = 2"` |  |
| rabbitmq.extraPlugins | string | `"rabbitmq_auth_mechanism_ssl"` |  |
| rabbitmq.image.pullPolicy | string | `"Always"` |  |
| rabbitmq.image.registry | string | `"gitlab.figo.systems/platform/dependency_proxy/containers"` |  |
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
| rabbitmq.replicaCount | int | `1` |  |
| rabbitmq.service.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| rabbitmq.service.tlsPort | int | `5671` |  |
| rabbitmq.serviceAccount.create | bool | `false` |  |
| rabbitmq.statefulsetLabels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.4.0](https://github.com/norwoodj/helm-docs/releases/v1.4.0)
