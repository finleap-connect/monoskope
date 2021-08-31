**[[Back To Overview]](../README.md)**

---

## Register your cluster with Monoskope

To be able to login you have to register the cluster with the m8 control plane:

```shell
$ monoctl create cluster --help
Creates a cluster. The name and display name are derived from the KubeAPIServer address given. They can be overridden by flags.

Usage:
  monoctl create cluster <KUBE_API_SERVER_ADDRESS> <CA_CERT_FILE> [flags]

Flags:
  -d, --display-name string   Display name of the cluster
  -h, --help                  help for cluster
  -n, --name string           Name of the cluster

Global Flags:
      --command-timeout duration   Timeout for long running commands (default 10s)
      --monoconfig string          Path to explicit monoskope config file to use for CLI requests
```

The `<KUBE_API_SERVER_ADDRESS>` is the address of your KubeAPIServer with protocol, e.g.
`https://api.kubernetes.your.domain`.

The CA certificate bundle has to be the CA of your KubeAPIServer so when
`monoctl` updates your `kubeconfig` the CA is known to `kubectl` when talking
to your KubeAPIServer.
