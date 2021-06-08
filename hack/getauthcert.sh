#!/usr/bin/env bash
set -euo pipefail

echo "Getting secrets from K8s..."
kubectl get secret -n platform-monoskope-monoskope m8-authentication -oyaml >m8auth.yaml

echo "Extracting PEMs..."
yq read m8auth.yaml 'data."tls.crt"' | base64 -d >tmp/tls.crt
yq read m8auth.yaml 'data."tls.key"' | base64 -d >tmp/tls.key

rm m8auth.yaml

echo "Done."
