# API Token Authentication

Monoskope supports generating scoped API tokens for authentication.
API tokens can be generated for any user (existing or not) but only by system administrators.

## Validity

The default validity is 24h.
The validity can be specified by the operator generating the token.

## Scopes

There are several scopes for specific use-cases:

 * NONE              // Dummy to prevent accidents
 * API               // Read-write for the complete API
 * WRITE_SCIM        // Read-write for endpoints with path prefix "/scim"
 * WRITE_K8SOPERATOR // Read-write for K8sOperator endpoints

## Generate

Use `monoctl` generate an API token:

```bash
$ monoctl create api-token --help
Retrieve an API token issued by the m8 control plane.

Usage:
  monoctl create api-token [flags]

Flags:
  -h, --help                help for api-token
  -s, --scopes strings      Specify the scopes for which the token should be valid.
                            Available scopes: NONE, API, WRITE_SCIM, WRITE_K8SOPERATOR
  -u, --user string         Specify the name or UUID of the user for whom the token should be issued. If not a UUID it will be treated as username.
  -v, --validity duration   Specify the validity period of the token. (default 24h0m0s)

Global Flags:
      --command-timeout duration   Timeout for long running commands (default 10s)
      --monoconfig string          Path to explicit monoskope config file to use for CLI requests
```