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

After storing the trust anchor in a K8s secret you can delete your local copy or store it in a save location.

## Rotating the trust anchor

Rotating the trust anchor without downtime is a multi-step process:
you must generate a new trust anchor, bundle it with the old one, rotate all certificates derived from the old one, and finally remove the old trust anchor from the bundle.

Create a new trust anchor:

```bash
step certificate create root.monoskope.cluster.local ca-new.crt ca-new.key \
  --profile root-ca --no-password --insecure
```

Create Bundle with old and new CA cert:

```bash
step certificate bundle ca-new.crt ca-old.crt bundle.crt
```

## Issueing mTLS certificates

When issueing certificates for components like the m8 Operator there are some expectations which must be met:

1. The `commonName` must be a subdomain of `monoskope.cluster.local`, e.g. `operator.monoskope.cluster.local`. It should be unique throughout the system.
1. Set `X509v3 Subject Alternative Name` DNS to the same as for `commonName` and add a unique email address as user information, e.g.:

    ```bash
    DNS:operator.monoskope.cluster.local, email:operator@monoskope.io
    ```

1. Set the organization to `Monoskope`.

See the default operator auth `cert-manager` certificat resource definition for this.
This will be deployed along with the m8 control plane:

```yaml
apiVersion: cert-manager.io/v1alpha3
kind: Certificate
metadata:
  name: m8dev-monoskope-mtls-operator-auth
  namespace: platform-monoskope-monoskope
spec:
  commonName: operator.monoskope.cluster.local
  dnsNames:
  - operator.monoskope.cluster.local
  emailSANs:
  - operator@monoskope.io
  issuerRef:
    kind: Issuer
    name: m8dev-monoskope-identity-issuer
  secretName: m8dev-monoskope-mtls-operator-auth
  subject:
    organizations:
    - Monoskope
  usages:
  - client auth
```
