#!/usr/bin/env bash
#
# This script uses kubectl to get the secert containing the private/public key
# with which JWTs are signed and validated by the Gateway.
set -euo pipefail

echo "Getting secrets from K8s..."
kubectl get secret -n platform-monoskope-monoskope m8-authentication -oyaml >m8auth.yaml

echo "Extracting PEMs..."
yq read m8auth.yaml 'data."tls.crt"' | base64 -d >tmp/tls.crt
yq read m8auth.yaml 'data."tls.key"' | base64 -d >tmp/tls.key

rm m8auth.yaml

echo "Done."
