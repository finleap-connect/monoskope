**[[Back To Overview]](../README.md)**

---

# Monoskope Permission Model

## Commands & Aggregates

|          COMMAND           |      AGGREGATE      |
|----------------------------|---------------------|
| ApproveClusterRegistration | ClusterRegistration |
| CreateCluster              | Cluster             |
| CreateTenant               | Tenant              |
| CreateUser                 | User                |
| CreateUserRoleBinding      | UserRoleBinding     |
| DeleteCluster              | Cluster             |
| DeleteTenant               | Tenant              |
| DeleteUserRoleBinding      | UserRoleBinding     |
| DenyClusterRegistration    | ClusterRegistration |
| RequestClusterRegistration | ClusterRegistration |
| UpdateTenant               | Tenant              |

## Permissions

The permissions shown here are auto generated from code.

|          COMMAND           |     ROLE     | SCOPE  |
|----------------------------|--------------|--------|
| ApproveClusterRegistration | admin        | system |
| CreateCluster              | admin        | system |
| CreateTenant               | admin        | system |
| CreateUser                 | admin        | system |
| CreateUserRoleBinding      | admin        | system |
|                            | admin        | tenant |
| DeleteCluster              | admin        | system |
| DeleteTenant               | admin        | system |
| DeleteUserRoleBinding      | admin        | system |
|                            | admin        | tenant |
| DenyClusterRegistration    | admin        | system |
| RequestClusterRegistration | k8s-operator | *      |
| UpdateTenant               | admin        | system |
