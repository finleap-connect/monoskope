#!/bin/bash

set -e

if [ ! -z "$CI_COMMIT_TAG" ]; then
    VERSION=$CI_COMMIT_TAG
elif [ ! -z "$CI_PIPELINE_IID" ]; then
    VERSION="0.0.0-$CI_COMMIT_SHORT_SHACI_PIPELINE_IID"
else
    VERSION="local-$(git rev-parse HEAD)"
fi

echo $VERSION
