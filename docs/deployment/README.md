# Monoskope Deployment

## Prerequisites

* [cert-manager](https://cert-manager.io/docs/) used to generate certificates for
  * Ambassador
  * RabbitMQ
  * CockroachDB
* [vault-operator](https://gitlab.figo.systems/platform/vault-operator) [optional] to get secrets for
  * RabbitMQ
  * Gateway
  from your HashiCorp Vault

## Step-by-step setup

1. Make sure you have the following available in your target cluster:
    * jetstack/cert-manager
    * finleap-connect/vaultoperator [optional]
1. Create m8 PKI.
See [certificate management](01-certificate-management.md).
1. Configure an identity provider.
See [identity provider setup](02-identity-provider-setup.md).
1. Configure m8 Ambassador.
See [DNS and IP address setup](03-dns-and-ip-address-setup.md).
1. Deploy [Helm Chart](../../build/package/helm/monoskope/README.md).
