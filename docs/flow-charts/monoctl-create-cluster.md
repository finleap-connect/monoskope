**[[Back To Overview]](..)**

# Create Cluster

```mermaid
sequenceDiagram
 participant C as Client
 participant G as Gateway
 participant CH as CommandHandler
 participant CBR as ClusterBootstrapReactor
 participant P as Projector
 participant RM as ReadModel
 participant O as M8 Operator

 C-->>+G: send CreateCluster command
 G-->>+CH: forward CreateCluster command
 CH-->>+CBR: emit ClusterCreated
 CH-->>+P: emit ClusterCreated
 CBR-->>+P: emit UserCreated
 CBR-->>+P: emit UserRoleBindingCreated
 O->>RM: query JWT
 O-->>+G: send RegisterCluster command
 G-->>+CH: forward RegisterCluster command
```

The process not only creates an aggregate for the cluster, but also a (machine) user with the DNS name of the cluster as the user name as well as an appropriate role binding.
