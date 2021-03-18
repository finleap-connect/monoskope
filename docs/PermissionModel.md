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
|                       | admin | tenant | self     | *       |
