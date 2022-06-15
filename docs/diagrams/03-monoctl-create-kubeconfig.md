# `monoctl create kubeconfig`

```mermaid
sequenceDiagram
    participant U as User
    participant M as monoctl
    participant M8 as Monoskope
    U-->>+M: monoctl create kubeconfig
    M-->>+M8: GetClustersAccessibleToMe
    M8-->>-M: return Clusters
    M-->>M: load users kubeconfig
    M-->>M: set clusters,context,authinfo on kubeconfig
    M-->>M: save kubeconfig
    M-->>-U: return successful
```

## General assumptions

* `monoctl create kubeconfig` adds/updates all clusters which are accessible to the user authenticated to the users kubeconfig.
* No auth token is stored in the kubeconfig by this command.
* Authentication happens when the user uses a kubecontext configured by m8 using [client-go-credential-plugins](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins) authentication.
See the sequence diagram [here](cluster-auth-flow.md).

## Useful links

* K8s docs on [client-go-credential-plugins](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins).
* [zalando/go-keyring](https://github.com/zalando/go-keyring)
