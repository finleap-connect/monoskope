# DNS & IP Address Setup

For m8 to be available on the internet you have to configure
the domain and IP address it will be hosted at:

```yaml
hosting:
  #-- Configure which issuer should be used for TLS via cert-manager
  issuer: letsencrypt-issuer
  #-- Configure which base domain the control plane will be hosted at
  domain: monoskope.somedomain.io
ambassador:
  service:
    #-- Configure the ip address the hosting.domain will be pointed at
    loadBalancerIP: 1.2.3.4
```

Now you have to setup the DNS entries accordingly:

1. `api.monoskope.somedomain.io` for `monoctl` via TLS
1. `mapi.monoskope.somedomain.io` for m8 operators via mTLS
