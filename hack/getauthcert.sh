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
