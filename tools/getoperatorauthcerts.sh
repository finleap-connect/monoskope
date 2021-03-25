#!/usr/bin/env bash
set -euo pipefail

echo "Getting secret from K8s..."
kubectl get secret m8dev-monoskope-mtls-operator-auth -oyaml >operatorauth.yaml

echo "Extracting PEMs..."
yq read operatorauth.yaml 'data."ca.crt"' | base64 -d >ca.crt
yq read operatorauth.yaml 'data."tls.crt"' | base64 -d >tls.crt
yq read operatorauth.yaml 'data."tls.key"' | base64 -d >tls.key

rm operatorauth.yaml
echo "Done."
