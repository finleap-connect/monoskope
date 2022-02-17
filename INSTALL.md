# Install Monoskope

## Prerequisites

### Required 

* [jetstack/cert-manager](https://cert-manager.io/docs/) used to generate certificates
  * Make sure to configure `pki.certificates.certManagerApiVersion` in the helm chart according to the needed version of the deployed cert-manager version
  * Minimum version required is `v1.0.2`
* [step-cli](https://smallstep.com/cli/) installed (you can use something similar too) for the [PKI](https://en.wikipedia.org/wiki/Public_key_infrastructure) setup

### Optional 

* [finleap-connect/vaultoperator](https://github.com/finleap-connect/vaultoperator) to collect secrets from your [HashiCorp Vault](https://www.vaultproject.io/) for
  * RabbitMQ
  * Gateway
  
## Step-by-step setup

1. Make sure you have the following available in your target cluster:
    * [jetstack/cert-manager](https://cert-manager.io/docs/) [required]
    * [finleap-connect/vaultoperator](https://github.com/finleap-connect/vaultoperator) [optional]
1. Create m8 PKI. See [certificate management](docs/deployment/01-certificate-management.md) for details.
1. Configure an identity provider. See [identity provider setup](docs/deployment/02-identity-provider-setup.md) for details.
1. Configure m8 Ambassador. See [DNS and IP address setup](docs/deployment/03-dns-and-ip-address-setup.md) for details.
1. Deploy [Helm Chart](build/package/helm/monoskope/README.md) via finleap-connect [chart repo](https://finleap-connect.github.io/charts/) and adjust the values to your needs.

## Optional steps

* Provision users and roles to your Monoskope instance via [SCIM](http://www.simplecloud.info/) if your identity provider supports it. See [Monoskope SCIM support](docs/deployment/04-configure-scim.md) for details.