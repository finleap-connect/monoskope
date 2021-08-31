#!/usr/bin/env bash
# Copyright 2021 Monoskope Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

echo "Getting secrets from K8s..."
kubectl get secret -n platform-monoskope-monoskope m8dev-monoskope-mtls-operator-auth -oyaml >operatorauth.yaml

echo "Extracting PEMs..."
yq read operatorauth.yaml 'data."ca.crt"' | base64 -d >tmp/ca.crt
yq read operatorauth.yaml 'data."tls.crt"' | base64 -d >tmp/tls.crt
yq read operatorauth.yaml 'data."tls.key"' | base64 -d >tmp/tls.key

rm operatorauth.yaml

echo "Done."
