**[[Back To Overview]](../README.md)**

---

# Monoskope Documentation Overview

## Prerequisites

The following things must be set up in your target K8s Cluster:

* [cert-manager](https://cert-manager.io/docs/) used to generate certificates for
  * Ambassador
  * RabbitMQ
  * CockroachDB
* [vault-operator](https://gitlab.figo.systems/platform/vault-operator) to generate/gather secrets for
  * RabbitMQ
  * Gateway
  from your HashiCorp Vault
