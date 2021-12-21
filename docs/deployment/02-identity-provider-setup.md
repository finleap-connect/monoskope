# Identity Provider Setup

Monoskope can connect to any identity provider (IDP) supporting OIDC/OAuth2
like [Dex](https://dexidp.io/) for example.

To connect m8 to your IDP you have to do the following:

1. Register a new client application with the IDP.
For Dex there is documentation on how to do this [here](https://dexidp.io/docs/using-dex/#configuring-your-app).
Find the docs for Gitlab as an IDP [here](https://docs.gitlab.com/ee/integration/oauth_provider.html).
You have to configure the following URLs as valid callback URL:

    * <http://localhost:8000>
    * <http://localhost:18000>

1. Grant the following scopes to the new application:

    * openid
    * profile
    * email

1. Aquire the `clientid` and `clientsecret` and provide it to Monoskope:

    1. Either by a K8s secret you need to create:

    ```bash
    kubectl create secret generic m8-gateway-oidc --from-literal=oidc-clientid=<clientid> --from-literal=oidc-clientsecret=<clientsecret> --from-literal=oidc-nonce=<somerandomstring>
    ```

    1. Or using the finleap/vaultoperator. It will create the secret automatically but you need to put the values in the right place in your Hashicorp Vault.

Set the following configuration in the values file for the Gateway component of m8 like in this example:

```yaml
gateway:
    auth:
        # -- The URL of the issuer to use for OIDC
        identityProviderURL: "https://idp.your-domain.com"
    # -- The secret where the gateway finds the OIDC secrets.
    # If vaultOperator.enabled:true the secret must be available at vaultOperator.basePath/gateway/oidc
    # and must contain the fields oidc-clientsecret, oidc-clientid. The oidc-nonce is generated automatically.
    oidcSecret:
        name: m8-gateway-oidc
```
