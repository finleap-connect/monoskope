# K8s AuthZ via RBAC Reconciling

Monoskope supports creating Kubernetes ClusterRoleBindings based on the role bindings a user has within Monoskope.
To apply these ClusterRoleBindings to your clusters you can use your GitOps tool of choice (e.g. ArgoCD, Flux, ...).

The configuration of this feature is done on the QueryHandler which projects changes to roles/users/cluster/etc to those git repo:

```yaml
queryhandler:
  # -- K8sAuthZ Configuration
  k8sAuthZ:
    # -- Enable external git repo reconciliation
    enabled: true
    # -- Configure secret provided as env vars
    existingSecret: monoskope-k8sauthz
    config:
      # -- Configure repos
      repository:
        url: git@your.domain/repo.git
        # -- Supported ssh, basic
        authType: ssh
        # -- Prefix for secrets consumed from env
        envPrefix: git
        author:
          name: someauthor
          email: someauthor@your.domain
      # -- Configure ClusterRole mapping
      mappings:
        - scope: CLUSTER
          role: admin
          clusterRole: someclusterolename
      # -- Prefix used for users in clusters
      usernamePrefix: "oidc:"
      # -- Reconcile loop interval
      interval: 5m
      # -- Put RBAC for all clusters into that repo
      allClusters: true
      # -- Put RBAC for only the following cluster into that repo
      # clusters:
      #   - test
```

What the secret must contain depends on the `envPrefix` specified by you and the `authType`.
The existing secret must contains the following fields matching the configuration above:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: monoskope-k8sauthz
  namespace: monoskope
type: Opaque
data:
  git.ssh.known_hosts: <base64-known-hosts-containing-your-registry>
  git.ssh.password: "<optional-password>"
  git.ssh.privateKey: <base64-private-key>
```
