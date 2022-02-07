# Monoskope SCIM support

The System for Cross-domain Identity Management ([SCIM](http://www.simplecloud.info/)) specification is designed to make managing user identities in cloud-based applications and services easier.

Monoskope implements SCIM and by that allows provisioning of users and rolebindings (scope `system` only) from a 3rd party identity provider of your choice (which must implement SCIM too).

## Setup SCIM

To activate SCIM support of Monoskope, adjust the helm chart values:

```yaml
scimserver:
  enabled: true
```

This will deploy an additional service called `SCIMServer`.

### Using OneLogin as Identity Provider

To configure OneLogin to provision users to Monoskope you can follow the [guide](https://developers.onelogin.com/scim/create-app) provided by OneLogin with small adjustments.

#### Create App

1. Access OneLogin and go to `Applications > Add App`.
2. Search for and select `SCIM Provisioner with SAML (SCIM v2 Enterprise)`
3. Give your SCIM app a `Display Name` value that will help you recognize it.
4. Select `Save`.

#### Configure App

1. Select the `Configuration` tab
1. Provide your `SCIM Base URL` value. This is the address that points OneLogin to Monoskope's SCIM API server. Example: https://api.monoskope.example.com/scim
1. Provide Monoskope's `SCIM JSON Template` value:
```json
{
    "schemas": [
        "urn:scim:schemas:core:2.0"
    ],
    "userName": "{$user.email}",
    "displayName": "{$user.display_name}"
}
```
1. Provide your SCIM Bearer Token value 

    * Create a token with `monoctl` and adjust the values according to your needs:
  
      `monoctl create api-token -u yourscimclient -s WRITE_SCIM -v 8760h`
    * Use the resulting token and put into the `SCIM Bearer Token` form field.

2. Select `Enable`. The app will attempt to make an initial connection to the SCIM base URL defined for your SCIM test app.
3. Select `Save`

#### Provisioning

1. Select the `Provisioning` tab
1. Select `Enable provisioning`

#### Groups

1. Select the `Parameters` tab
1. Select `Groups` from the table
1. Select `Include in User Provisioning` in section `Flags`
