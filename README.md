# Monoskope

![Monoskope Logo](assets/logo/monoskope.png)

[![pipeline status](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/pipeline.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)
[![coverage report](https://gitlab.figo.systems/platform/monoskope/monoskope/badges/main/coverage.svg)](https://gitlab.figo.systems/platform/monoskope/monoskope/-/commits/main)

`Monoskope` implements the management and operation of tenants, users and their roles in a [Kubernetes](https://kubernetes.io/) multi-cluster environment. It fullfills the needs of operators of the clusters as well as the needs of developers using the cloud infrastructure provided by the operators.

## Further documentation

Find the detailed documentation at [/docs](docs/Overview.md).

## Development

When developing, the `Makefile` comes in handy to help you with various tasks.
There are specific `*.mk` files for things like helm, kind, go, etc. which provides targets for developing with those tools.

The following example renders the `monoskope` helm chart:

```sh
$ make helm-template
==> Linting build/package/helm/monoskope
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed

wrote tmp/monoskope/charts/dex/templates/serviceaccount.yaml
wrote tmp/monoskope/charts/dex/templates/secret.yaml
wrote tmp/monoskope/charts/dex/templates/service.yaml
wrote tmp/monoskope/charts/dex/templates/deployment.yaml
```

The following targets are defined. Please not that there are variables (uppercase) which can be overriden:

| target | Description |
| --------- | ----------- |
| *general* | |
| `clean` | Cleans everything, tools, tmp dir used, whatever |
| `tools` | Install necessary tools to `TOOLS_DIR`, like `kind`, `ginkgo`, `golangci-lint`, ... |
| `tools-clean` | Removes the tools |
| *helm* | |
| `helm-add-kubism` | Add the kubism helm repository to the local list of repos |
| `helm-template` | Templates the helm chart to `HELM_OUTPUT_DIR/monoskope` |
| `helm-install` | Installs the helm chart to namespace `KUBE_NAMESPACE` with your current `kubecontext` and `HELM` |
| `helm-uninstall` | Uninstalls the helm chart from namespace `KUBE_NAMESPACE` |
| `helm-clean` | Clears `HELM_OUTPUT_DIR` |
| `helm-dep` | Does a helm dep update for `monoskope` |
| `helm-lint` | Does a helm lint for `monoskope` |
| *kind* | |
| `kind-create` | Ramps up a kind cluster `KIND_CLUSTER` |
| `kind-delete` | Deletes the kind cluster `KIND_CLUSTER` |
| `kind-get-kubeconfig` | Gets the kubeconfig to connect the kind cluster `KIND_CLUSTER` |
| *go* | |
| `go-prepare` | Downloads all require go modules |
| `go-lint` | Lints all go code |
| `go-vet` | Vets all go code |
| `go-fmt` | Formats all go code |
| `go-run-*` | Runs the app in `cmd/*`, e.g. `go-run-monoctl` to run `monoctl` from sources |
