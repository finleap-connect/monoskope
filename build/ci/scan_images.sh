#!/bin/bash
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


set -e

echo "Downloading roxctl..."
curl -k -L -H "Authorization: Bearer $ROX_API_TOKEN" https://$ROX_CENTRAL_ENDPOINT/api/cli/download/roxctl-linux --output ./roxctl
chmod a+x ./roxctl

echo "Scanning images..."

for file in $(find $CI_PROJECT_DIR/build/package/helm -type f -name values.yaml); do
    CURRENT_IMAGE=$(grep "repository:" $file | cut -d':' -f2- | tr -d '[:space:]' | cut -d':' -f3)
    if [ "$CURRENT_IMAGE" != "${CURRENT_IMAGE#$CI_REGISTRY}" ]; then # Starts with
        echo "Scanning '$CURRENT_IMAGE' ..."
        ./roxctl image scan -e $ROX_CENTRAL_API_ENDPOINT --force --image $CURRENT_IMAGE:$VERSION
        ./roxctl image check -e $ROX_CENTRAL_API_ENDPOINT --image $CURRENT_IMAGE:$VERSION
    fi
done

echo "Scanning images finished."
