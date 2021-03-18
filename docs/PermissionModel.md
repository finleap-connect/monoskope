# Monoskope Permission Model

## Commands & Aggregates

|        COMMAND        |    AGGREGATE    |
|-----------------------|-----------------|
| CreateUserRoleBinding | UserRoleBinding |
| CreateUser            | User            |

## Permissions

|        COMMAND        | ROLE  | SCOPE  | RESOURCE | SUBJECT |
|-----------------------|-------|--------|----------|---------|
| CreateUser            | *     | *      | *        | self    |
|                       | admin | system | *        | *       |
| CreateUserRoleBinding | admin | system | *        | *       |
|                       | admin | tenant | same     | *       |

* same - means that the permission to execute the command is only granted if the resource/scope combination matches, e.g. a tenant admin can only create a role binding for the same tenant resource.
* self - means that the subject in the command must match the subject executing the command, e.g. a user can create itself but not other users.
