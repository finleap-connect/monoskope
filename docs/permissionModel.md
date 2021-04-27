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

|          COMMAND           |     ROLE     | SCOPE  | RESOURCEMATCH |
|----------------------------|--------------|--------|---------------|
| ApproveClusterRegistration | admin        | system | false         |
| CreateCluster              | admin        | system | false         |
| CreateTenant               | admin        | system | false         |
| CreateUser                 | admin        | system | false         |
| CreateUserRoleBinding      | admin        | system | false         |
|                            | admin        | tenant | true          |
| DeleteCluster              | admin        | system | false         |
| DeleteTenant               | admin        | system | false         |
| DeleteUserRoleBinding      | admin        | system | false         |
|                            | admin        | tenant | true          |
| DenyClusterRegistration    | admin        | system | false         |
| RequestClusterRegistration | k8s-operator | *      | false         |
| UpdateTenant               | admin        | system | false         |

* RESOURCEMATCH - means that the permission to execute the command is only granted if the resource/scope combination matches, e.g. a tenant admin can only create a role binding for the same tenant resource.
