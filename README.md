# Monoskope

![Monoskope Logo](assets/logo/monoskope.png)

## Build status

| `main` | `develop` |
| -- | -- |
|[![pipeline status](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/pipeline.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)|[![pipeline status](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/develop/pipeline.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/develop)
|[![coverage report](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/coverage.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)|[![coverage report](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/develop/coverage.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/develop)|

## Helm chart
![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square)
![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)
![AppVersion: 0.0.1](https://img.shields.io/badge/AppVersion-0.0.1-informational?style=flat-square)

`Monoskope` implements the management and operation of tenants, users and their roles in a [Kubernetes](https://kubernetes.io/) multi-cluster environment. It fullfills the needs of operators of the clusters as well as the needs of developers using the cloud infrastructure provided by the operators.

## Further documentation

Find the detailed documentation at [/docs](docs/Overview.md).

## Development

When developing, the `Makefile` comes in handy to help you with various tasks.
There are specific `*.mk` files for things like helm, kind, go, etc. which provides targets for developing with those tools.

The following example renders the `monoskope` helm chart:

```sh
$ make helm-template-monoskope
install.go:172: [debug] Original chart version: ""
install.go:189: [debug] CHART PATH: /home/jsteffen/dev/src/platform/monoskope/monoskope/build/package/helm/monoskope

coalesce.go:196: warning: cannot overwrite table with non table for connectors (map[])
coalesce.go:196: warning: cannot overwrite table with non table for connectors (map[])
coalesce.go:196: warning: cannot overwrite table with non table for connectors (map[])
wrote tmp/monoskope/charts/cockroachdb/templates/poddisruptionbudget.yaml
wrote tmp/monoskope/charts/cockroachdb/templates/serviceaccount.yaml
wrote tmp/monoskope/charts/cockroachdb/templates/role.yaml
wrote tmp/monoskope/charts/cockroachdb/templates/rolebinding.yaml
wrote tmp/monoskope/charts/cockroachdb/templates/service.discovery.yaml
wrote tmp/monoskope/charts/cockroachdb/templates/service.public.yaml
wrote tmp/monoskope/charts/dex/templates/service.yaml
wrote tmp/monoskope/charts/monoskope-gateway/templates/service.yaml
wrote tmp/monoskope/templates/service-crdb-metrics.yaml
wrote tmp/monoskope/charts/dex/templates/deployment.yaml
wrote tmp/monoskope/charts/monoskope-gateway/templates/deployment.yaml
wrote tmp/monoskope/charts/cockroachdb/templates/statefulset.yaml
wrote tmp/monoskope/charts/cockroachdb/templates/job.init.yaml
wrote tmp/monoskope/templates/ingress.yaml
wrote tmp/monoskope/templates/cert-crdb.yaml
wrote tmp/monoskope/templates/cert-crdb.yaml
wrote tmp/monoskope/templates/cert-crdb.yaml
wrote tmp/monoskope/templates/cert-crdb.yaml
wrote tmp/monoskope/templates/cert-crdb.yaml
wrote tmp/monoskope/templates/servicemonitor-crdb-metrics.yaml
wrote tmp/monoskope/charts/monoskope-gateway/templates/vaultsecret.yaml
wrote tmp/monoskope/templates/vaultsecret-dex.yaml
wrote tmp/monoskope/charts/cockroachdb/templates/tests/client.yaml

ATTENTION:
If you want to have the latest dependencies (e.g. gateway chart changes)
execute the following command prior to the current command:
$ make helm-dep-monoskope

```

The following targets are defined. Please not that there are variables (uppercase) which can be overriden:

| target | Description |
| --------- | ----------- |
| *general* | |
| `clean` | Cleans everything, tools, tmp dir used, whatever |
| `diagrams` | Generates mermaidjs diagrams below `docs/flow-charts` |
| `tools` | Install necessary tools to `TOOLS_DIR`, like `kind`, `ginkgo`, `golangci-lint`, ... |
| `tools-clean` | Removes the tools |
| `echo-<VARIABLENAME>` | Echos the content of `<VARIABLENAME>` |
| *helm* | |
| `helm-add-kubism` | Add the kubism helm repository to the local list of repos |
| `helm-template-<CHARTNAME>` | Templates the helm chart `<CHARTNAME>` to `HELM_OUTPUT_DIR/<CHARTNAME>` |
| `helm-install-<CHARTNAME>` | Installs the helm chart `<CHARTNAME>` to namespace `KUBE_NAMESPACE` with your current `kubecontext` and `HELM` |
| `helm-install-from-repo-<CHARTNAME>` | Installs the helm chart `<CHARTNAME>` to namespace `KUBE_NAMESPACE` from `HELM_REGISTRY_ALIAS` in version `VERSION` |
| `helm-uninstall-<CHARTNAME>` | Uninstalls the helm chart `<CHARTNAME>` from namespace `KUBE_NAMESPACE` |
| `helm-clean` | Clears `HELM_OUTPUT_DIR` |
| `helm-dep-<CHARTNAME>` | Does a helm dep update for `<CHARTNAME>` |
| `helm-lint-<CHARTNAME>` | Does a helm lint for `<CHARTNAME>` |
| *go* | |
| `go-mod` | Downloads all require go modules |
| `go-fmt` | Formats all `*.go` files |
| `go-vet` | Vets all go code |
| `go-lint` | Lints all go code |
| `go-run-*` | Runs the app in `cmd/*`, e.g. `go-run-monoctl` to run `monoctl` from sources |
| `go-test` | Runs all go tests |
| `go-protobuf` | Generates code for all proto specs in `api` folder and it's children |

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../gateway | monoskope-gateway |  |
| https://artifactory.figo.systems/artifactory/virtual_helm | cockroachdb | 5.0.2 |
| https://charts.bitnami.com/bitnami | rabbitmq | 8.5.2 |
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
| dex.config.existingSecret | string | `"monoskope-dex-config"` |  |
| dex.config.issuer | string | `"https://monoskope.io/dex"` |  |
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
| dex.config.storage.config.secret | string | `"m8dev-monoskope-crdb-client-dex"` |  |
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
| dex.replicas | int | `3` |  |
| dex.resources.limits.cpu | string | `"500m"` |  |
| dex.resources.limits.memory | string | `"100Mi"` |  |
| dex.resources.requests.cpu | string | `"100m"` |  |
| dex.resources.requests.memory | string | `"50Mi"` |  |
| dex.serviceAccount.create | bool | `false` |  |
| dex.telemetry | bool | `true` |  |
| fullnameOverride | string | `""` |  |
| ingress.enabled | bool | `false` |  |
| ingress.host | string | `"monoskope.io"` |  |
| monitoring.tenant | string | `"finleap-cloud"` |  |
| monoskope-gateway.auth.allowRootToken | bool | `false` |  |
| monoskope-gateway.auth.issuerURL | string | `"https://monoskope.io/dex"` |  |
| monoskope-gateway.enabled | bool | `true` |  |
| monoskope-gateway.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| monoskope-gateway.nameOverride | string | `"gateway"` |  |
| monoskope-gateway.replicaCount | int | `3` |  |
| name | string | `"monoskope"` |  |
| nameOverride | string | `""` |  |
| rabbitmq.auth.existingErlangSecret | string | `"rabbitmq-erlang-cookie"` |  |
| rabbitmq.auth.password | string | `"foo"` |  |
| rabbitmq.auth.tls.enabled | bool | `true` |  |
| rabbitmq.auth.tls.existingSecret | string | `"rabbitmq-leaf"` |  |
| rabbitmq.auth.tls.failIfNoPeerCert | bool | `true` |  |
| rabbitmq.auth.tls.sslOptionsVerify | string | `"verify_peer"` |  |
| rabbitmq.auth.username | string | `"admin"` |  |
| rabbitmq.enabled | bool | `true` |  |
| rabbitmq.extraConfiguration | string | `"load_definitions = /app/rabbitmq-definitions.json"` |  |
| rabbitmq.image.pullPolicy | string | `"Always"` |  |
| rabbitmq.image.repository | string | `"gitlab.figo.systems/platform/dependency_proxy/containers/bitnami/rabbitmq"` |  |
| rabbitmq.image.tag | string | `"3.8.9"` |  |
| rabbitmq.loadDefinition.enabled | bool | `true` |  |
| rabbitmq.loadDefinition.existingSecret | string | `"rabbitmq-load-definition"` |  |
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
| rabbitmq.replicaCount | int | `3` |  |
| rabbitmq.service.labels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |
| rabbitmq.service.port | int | `5672` |  |
| rabbitmq.service.portName | string | `"amqp"` |  |
| rabbitmq.statefulsetLabels."app.kubernetes.io/part-of" | string | `"monoskope"` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.4.0](https://github.com/norwoodj/helm-docs/releases/v1.4.0)