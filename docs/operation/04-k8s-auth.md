# Kubernetes Authentication

Monoskope can be utilized as an OAuth2 provider to authenticate with Kubernetes clusters via [OpenID Connect](https://openid.net/connect/).
See the official docs of [Kubernetes authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) to find additional information on the topic.
The following diagram shows how the authentication flow looks like when Monoskope comes into play:

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

* `monoctl` eventually does it's normal authentication flow when `kubectl` is used to get nodes.
This depends on the current authentication state of `monoctl`.
If a token is available which hasn't expired yet, no authentication flow is necessary here.
* `monoctl cluster auth` call may return immediately without talking to the control plane if there is a cached token available.

## Prerequisites

You need to [register your cluster](03-register-k8s-cluster.md) with the m8 control plane first.

## You have control over the KubeAPIServer

If you have control over your KubeAPIServer refer to the [official docs](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#configuring-the-api-server) of Kubernetes for detailed explanation.
To enable Monoskope as OAuth2 provider, configure the following flags on the API server:

* Set Monoskope as the issuer:

    `--oidc-issuer-url "https://api.monoskope.your-domain.io"`

* Tokens issued by Monoskope will have `k8sauth` as the audience:

    `--oidc-client-id "k8sauth"`

* Claim containing the user name:

    `--oidc-username-claim cluster_username`

* Claim containing the role/group:

    `--oidc-groups-claim cluster_role`

* Claim which must be present in the token. Prohibits valid tokens for different cluster used to auth:

    `--oidc-required-claim cluster_name=your_cluster`

## You are using SAP Gardener

If you are using SAP Gardener you might not be able to directly configure this
on the KubeAPIServer, but there are different ways to do this:

* [ClusterOpenIDConnectPreset](https://github.com/gardener/gardener/blob/master/docs/usage/openidconnect-presets.md#clusteropenidconnectpreset) and [OpenIDConnectPreset](https://github.com/gardener/gardener/blob/master/docs/usage/openidconnect-presets.md#openidconnectpreset) resources can be used.
* Shoots can be configured directly via the `core.gardener.cloud/v1beta1/Shoot` [resource](https://github.com/gardener/gardener/blob/master/example/90-shoot.yaml#L137).

The configuration for both is pretty similar for both cases.
If you're configuring it directly via the Gardener Shoot resource the following has to be at `spec.kubernetes.kubeAPIServer.oidcConfig`.
When you're using the Presets this has to be put under `spec.server`.

```yaml
issuerURL: "<https://api.monoskope.your-domain.io>"
clientID: k8sauth
usernameClaim: cluster_username
groupsClaim: cluster_role
requiredClaims:
    cluster_name: <your_cluster> 
caBundle: |-
    #   -----BEGIN CERTIFICATE-----
    #   CA which issues the cert for https://api.monoskope.your-domain.io
    #   -----END CERTIFICATE-----
```

## Connect the cluster to monoskope

See [Register your Cluster](03-register-k8s-cluster.md)

## Do a login

```shell
$ monoctl create kubeconfig
Your kubeconfig has been generated/updated.
Use `kubectl config get-contexts` to see available contexts.
Use `kubectl config use-context <CONTEXTNAME>` to switch between clusters.
$ kubectl config use-context test-cluster-default
Switched to context "test-cluster-default".
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"16", GitVersion:"v1.16.15", GitCommit:"2adc8d7091e89b6e3ca8d048140618ec89b39369", GitTreeState:"clean", BuildDate:"2020-09-02T11:40:00Z", GoVersion:"go1.13.15", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"20", GitVersion:"v1.20.7", GitCommit:"132a687512d7fb058d0f5890f07d4121b3f0a2e2", GitTreeState:"clean", BuildDate:"2021-05-12T12:32:49Z", GoVersion:"go1.15.12", Compiler:"gc", Platform:"linux/amd64"}
```

You're good to go!

## Certificate rotation

The certificate used by Monoskope to sign and verify k8s tokens has a long expire date by design. 

Rotating it can be done easily using the [cert-manager CLI](https://cert-manager.io/docs/reference/cmctl/#renew)

```shell
cmctl renew m8-monoskope-tls-cert
```

For more information see [here](https://cert-manager.io/docs/usage/certificate#actions-triggering-private-key-rotation)
