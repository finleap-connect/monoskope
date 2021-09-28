# monoskope

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.1](https://img.shields.io/badge/AppVersion-0.0.1-informational?style=flat-square)

Monoskope implements the management and operation of tenants, users and their roles in a Kubernetes multi-cluster environment.

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../cluster-bootstrap-reactor | cluster-bootstrap-reactor | 0.0.1-local |
| file://../commandhandler | commandhandler | 0.0.1-local |
| file://../eventstore | eventstore | 0.0.1-local |
| file://../gateway | gateway | 0.0.1-local |
| file://../queryhandler | queryhandler | 0.0.1-local |
| https://charts.bitnami.com/bitnami | rabbitmq | 8.17.0 |
| https://charts.cockroachdb.com/ | cockroachdb | 6.1.2 |
| https://getambassador.io | ambassador | 6.7.11 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| ambassador.agent.enabled | bool | `false` |  |
| ambassador.crds.create | bool | `false` |  |
| ambassador.crds.enabled | bool | `false` |  |
| ambassador.deploy | bool | `true` |  |
| ambassador.enableAES | bool | `false` |  |
| ambassador.enabled | bool | `true` |  |
| ambassador.image.repository | string | `"datawire/ambassador"` |  |
| ambassador.image.tag | string | `"1.14.1"` |  |
| ambassador.rbac.create | bool | `false` |  |
| ambassador.replicaCount | int | `1` |  |
| ambassador.resources.limits.cpu | int | `4` |  |
| ambassador.resources.limits.memory | string | `"1000Mi"` |  |
| ambassador.resources.requests.cpu | string | `"100m"` |  |
| ambassador.resources.requests.memory | string | `"512Mi"` |  |
| ambassador.scope.singleNamespace | bool | `true` |  |
| ambassador.serviceAccount.create | bool | `false` |  |
| cluster-bootstrap-reactor.enabled | bool | `true` |  |
| cluster-bootstrap-reactor.keySecret.name | string | `"m8-authentication"` |  |
| cluster-bootstrap-reactor.messageBus.configSecret | string | `"m8-messagebus-client-config"` |  |
| cluster-bootstrap-reactor.messageBus.tlsSecret | string | `"m8-messagebus-client-auth-cert"` |  |
| cluster-bootstrap-reactor.replicaCount | int | `1` |  |
| cockroachdb.conf.cache | string | `"25%"` |  |
| cockroachdb.conf.maxSQLMemory | string | `"25%"` |  |
| cockroachdb.dropExistingDatabase | bool | `false` |  |
| cockroachdb.enabled | bool | `true` |  |
| cockroachdb.image.pullPolicy | string | `"Always"` |  |
| cockroachdb.image.repository | string | `"cockroachdb/cockroach"` |  |
| cockroachdb.image.tag | string | `"v21.1.4"` |  |
| cockroachdb.init.annotations."linkerd.io/inject" | string | `"disabled"` |  |
| cockroachdb.serviceMonitor.annotations | object | `{}` |  |
| cockroachdb.serviceMonitor.enabled | bool | `false` |  |
| cockroachdb.serviceMonitor.interval | string | `"1m"` |  |
| cockroachdb.serviceMonitor.labels.release | string | `"monitoring"` |  |
| cockroachdb.serviceMonitor.scrapeTimeout | string | `"10s"` |  |
| cockroachdb.statefulset.budget.maxUnavailable | int | `1` |  |
| cockroachdb.statefulset.replicas | int | `3` |  |
| cockroachdb.statefulset.resources.limits.cpu | int | `2` |  |
| cockroachdb.statefulset.resources.limits.memory | string | `"2Gi"` |  |
| cockroachdb.statefulset.resources.requests.cpu | string | `"500m"` |  |
| cockroachdb.statefulset.resources.requests.memory | string | `"1Gi"` |  |
| cockroachdb.storage.persistentVolume.size | string | `"1Gi"` |  |
| cockroachdb.tls.certs.certManager | bool | `true` |  |
| cockroachdb.tls.certs.certManagerIssuer.kind | string | `"Issuer"` |  |
| cockroachdb.tls.certs.certManagerIssuer.name | string | `"m8-root-ca-issuer"` |  |
| cockroachdb.tls.certs.provided | bool | `true` |  |
| cockroachdb.tls.certs.useCertManagerV1CRDs | bool | `true` |  |
| cockroachdb.tls.enabled | bool | `true` |  |
| commandhandler.enabled | bool | `true` |  |
| commandhandler.replicaCount | int | `1` |  |
| eventstore.enabled | bool | `true` |  |
| eventstore.messageBus.configSecret | string | `"m8-messagebus-client-config"` |  |
| eventstore.messageBus.tlsSecret | string | `"m8-messagebus-client-auth-cert"` |  |
| eventstore.replicaCount | int | `1` |  |
| eventstore.storeDatabase.configSecret | string | `"m8-db-client-config"` |  |
| eventstore.storeDatabase.tlsSecret | string | `"m8-db-client-auth-cert"` |  |
| fullnameOverride | string | `""` |  |
| gateway.auth.identityProviderName | string | `""` | The identifier of the issuer, e.g. DEX or whatever identifies your identities upstream |
| gateway.auth.identityProviderURL | string | `""` | The URL of the issuer to use for OIDC |
| gateway.auth.selfURL | string | `""` | The URL of the issuer to Gateway itself |
| gateway.enabled | bool | `true` |  |
| gateway.keySecret | object | `{"name":"m8-authentication"}` | The secret containing private key for signing JWTs. |
| gateway.keySecret.name | string | `"m8-authentication"` | Name of the secret to be used by the gateway, required |
| gateway.oidcSecret | object | `{"name":"m8-gateway-oidc"}` | The secret where the gateway finds the OIDC secrets. If vaultOperator.enabled:true the secret must be available at vaultOperator.basePath/gateway/oidc and must contain the fields oidc-clientsecret, oidc-clientid. The oidc-nonce is generated automatically. |
| gateway.replicaCount | int | `1` |  |
| global.imagePullSecrets | list | `[]` |  |
| global.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| hosting.domain | string | `""` |  |
| hosting.issuer | string | `""` |  |
| messageBus.clientAuthCertSecretName | string | `"m8-messagebus-client-auth-cert"` |  |
| messageBus.clientConfigSecretName | string | `"m8-messagebus-client-config"` |  |
| messageBus.routingKeyPrefix | string | `"m8"` |  |
| name | string | `"monoskope"` |  |
| nameOverride | string | `""` |  |
| pki.authentication.keySecretName | string | `"m8-authentication"` |  |
| pki.certificates.certManagerApiVersion | string | `"v1"` | Specify which apiVersion cert-manager resources must have. |
| pki.certificates.duration | string | `"2160h"` |  |
| pki.certificates.renewBefore | string | `"1440h"` |  |
| pki.enabled | bool | `true` |  |
| pki.issuer.ca.enabled | bool | `true` |  |
| pki.issuer.ca.existingTrustAnchorSecretName | string | `"m8-trust-anchor"` |  |
| pki.issuer.ca.secretVersion | int | `1` |  |
| pki.issuer.name | string | `"m8-root-ca-issuer"` |  |
| pki.issuer.vault.enabled | bool | `false` |  |
| queryhandler.enabled | bool | `true` |  |
| queryhandler.messageBus.configSecret | string | `"m8-messagebus-client-config"` |  |
| queryhandler.messageBus.tlsSecret | string | `"m8-messagebus-client-auth-cert"` |  |
| queryhandler.replicaCount | int | `1` |  |
| rabbitmq.auth.existingErlangSecret | string | `"m8-rabbitmq-erlang-cookie"` | Name of the secret containing the erlang secret If vaultOperator.enabled:true the secret will eb auto generated |
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
| rabbitmq.image.tag | string | `"3.8.9"` |  |
| rabbitmq.loadDefinition.enabled | bool | `true` |  |
| rabbitmq.loadDefinition.existingSecret | string | `"m8-rabbitmq-load-definition"` |  |
| rabbitmq.metrics.grafanaDashboard.enabled | bool | `false` |  |
| rabbitmq.persistence.enabled | bool | `false` |  |
| rabbitmq.rbac.create | bool | `false` |  |
| rabbitmq.replicaCount | int | `3` |  |
| rabbitmq.service.tlsPort | int | `5671` |  |
| rabbitmq.serviceAccount.create | bool | `false` |  |
| vaultOperator.basePath | string | `"app/{{ .Release.Namespace }}"` |  |
| vaultOperator.enabled | bool | `false` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.4.0](https://github.com/norwoodj/helm-docs/releases/v1.4.0)
