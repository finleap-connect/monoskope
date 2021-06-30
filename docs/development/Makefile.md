**[[Back To Overview]](README.md)**

---

# Makefile

When developing, the `Makefile` comes in handy to help you with various tasks.
There are specific `*.mk` files for things like helm, go, etc. which provides targets for developing with those tools.

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
wrote tmp/monoskope/charts/monoskope-gateway/templates/service.yaml
wrote tmp/monoskope/templates/service-crdb-metrics.yaml
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
| `tools` | Install necessary tools to `TOOLS_DIR`, like `ginkgo`, `golangci-lint`, ... |
| `tools-clean` | Removes the tools |
| `echo-<VARIABLENAME>` | Echos the content of `<VARIABLENAME>` |
| *helm* | |
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
| `go-report` | Generates report of commands and permssions |
