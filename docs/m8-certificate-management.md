# Monoskope Certificate Management

## Prerequisites

* [step-cli](https://smallstep.com/cli/) installed locally
* [cert-manager](https://cert-manager.io) running in target cluster
* [OPTIONAL] [Vault](https://www.vaultproject.io/) accessable by `cert-manager`

## Overview

Monoskope needs it's own PKI because components like the `m8 Operator` use [mTLS](https://en.wikipedia.org/wiki/Mutual_authentication) to communicate and authenticate with the Monoskope control plane.

## Create a trust anchor

Create a root CA certificate which we call the trust anchor:

```bash
step certificate create root.monoskope.cluster.local ca.crt ca.key \
  --profile root-ca --no-password --insecure
```

This trust anchor must be made available to [cert-manager](https://cert-manager.io) to let it issue certificates based on that trust anchor.

This can be done in two ways:

* via the [CA issuer](https://cert-manager.io/docs/configuration/ca/)
* via the [Vault issuer](https://cert-manager.io/docs/configuration/vault/)

### Using the CA issuer

If you're using the CA Issuer you have to set the following in the `values.yaml` when deploying Monoskope:

```yaml
pki:
  enabled: true
  issuer:
    ca:
      enabled: true
      existingTrustAnchorSecretName: "m8-trust-anchor" # name of secret in K8s where you have to provide the root ca
```

Create secret containing the generated trust anchor as in the namespace you're about to deploy Monoskope:

```bash
kubectl -n monoskope create secret tls m8-trust-anchor --cert=ca.crt --key=ca.key
```

## Rotating the trust anchor

Rotating the trust anchor without downtime is a multi-step process: you must generate a new trust anchor, bundle it with the old one, rotate the issuer certificate and key pair, and finally remove the old trust anchor from the bundle. If you simply need to rotate the issuer certificate and key pair, you can skip directly to Rotating the identity issuer certificate and ignore the trust anchor rotation steps.

### Create a new trust anchor

```bash
step certificate create root.monoskope.cluster.local ca-new.crt ca-new.key \
  --profile root-ca --no-password --insecure
```

### Create Bundle with old and new CA cert

```bash
step certificate bundle ca-new.crt ca-old.crt bundle.crt
```
