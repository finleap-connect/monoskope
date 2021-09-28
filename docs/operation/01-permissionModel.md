# Monoskope Permission Model

## Commands & Aggregates

|        COMMAND        |    AGGREGATE    |
|-----------------------|-----------------|
| CreateCluster         | Cluster         |
| CreateTenant          | Tenant          |
| CreateUser            | User            |
| CreateUserRoleBinding | UserRoleBinding |
| DeleteCluster         | Cluster         |
| DeleteTenant          | Tenant          |
| DeleteUser            | User            |
| DeleteUserRoleBinding | UserRoleBinding |
| RequestCertificate    | Certificate     |
| UpdateCluster         | Cluster         |
| UpdateTenant          | Tenant          |

## Permissions

The permissions shown here are auto generated from code.

|        COMMAND        |    ROLE     | SCOPE  |
|-----------------------|-------------|--------|
| CreateCluster         | admin       | system |
| CreateTenant          | admin       | system |
| CreateUser            | admin       | system |
| CreateUserRoleBinding | admin       | system |
|                       | admin       | tenant |
| DeleteCluster         | admin       | system |
| DeleteTenant          | admin       | system |
| DeleteUser            | admin       | system |
| DeleteUserRoleBinding | admin       | system |
|                       | admin       | tenant |
| RequestCertificate    | admin       | system |
|                       | k8soperator | system |
| UpdateCluster         | admin       | system |
| UpdateTenant          | admin       | system |
