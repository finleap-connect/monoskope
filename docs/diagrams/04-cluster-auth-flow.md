# Cluster Authentication Flow

## General assumptions

* `monoctl` eventually does it's normal authentication flow when `kubectl` is used to get nodes.
This depends on the current authentication state of `monoctl`.
If a token is available which hasn't expired yet, no authentication flow is necessary here.
* `monoctl cluster auth` call may return immediately without talking to the control plane if there is a cached token available.

## Diagram

```mermaid
sequenceDiagram
    participant U as User
    participant K as kubectl
    participant m8ctl as monoctl
    participant A as KubeApiServer
    participant G as Monoskope Gateway
    U-->>m8ctl: monoctl create kubeconfig
    m8ctl-->>m8ctl: update users kubeconfig
    U-->>+K: kubectl config set-context <br>test-cluster
    U-->>+K: kubectl get nodes
    Note right of K: kubectl is configured by monoctl<br>to call monoctl to get auth token.
    K-->>+m8ctl: monoctl get cluster-credentials <br>test-cluster default
    m8ctl-->>+G: get token for k8s auth
    G-->>-m8ctl: returns token for k8s auth
    m8ctl-->>-K: return token
    K-->>+A: calls get nodes
    A-->>A: do webhook authentication
    A-->>G: query JWKs
    A-->>A: validate JWT
    A-->>-K: return nodes of cluster
    K-->>-U: shows nodes of cluster
```

## Useful links

* K8s docs on [client-go-credential-plugins](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins).
* Information on [JWKs](https://auth0.com/docs/tokens/json-web-tokens/json-web-key-sets).
* [gardener/oidc-webhook-authenticator](https://github.com/gardener/oidc-webhook-authenticator)
* [zalando/go-keyring](https://github.com/zalando/go-keyring)
