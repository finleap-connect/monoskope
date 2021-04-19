#!/bin/bash

set -e

echo "Downloading roxctl..."
curl -k -L -H "Authorization: Bearer $ROX_API_TOKEN" https://$ROX_CENTRAL_ENDPOINT/api/cli/download/roxctl-linux --output ./roxctl
chmod a+x ./roxctl

echo "Scanning images..."

for file in $(find ${CI_PROJECT_DIR}/build/package/helm -type f -name values.yaml); do
    CURRENT_IMAGE=$(grep "repository:" $file | cut -d':' -f2- | tr -d '[:space:]' | cut -d':' -f3)
    if [[ $CURRENT_IMAGE == $CI_REGISTRY* ]]; then
        echo "Scanning '$CURRENT_IMAGE' ..."
        ./roxctl image scan -e $ROX_CENTRAL_API_ENDPOINT --force --image $CURRENT_IMAGE:$CI_COMMIT_TAG
        ./roxctl image check -e $ROX_CENTRAL_API_ENDPOINT --image $CURRENT_IMAGE:$CI_COMMIT_TAG
    fi
done

echo "Scanning images finished."
