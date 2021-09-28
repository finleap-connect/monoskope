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
export PWD=$(pwd)

set -eu

# chartpress is used to publish the charts
pip install chartpress==0.6.*

# Get a private SSH key from the Github secrets. It will
# be used to establish an identity with rights to push to the git repository
# hosting our Helm charts: https://github.com/finleap-connect/charts
echo "$CHART_DEPLOY_KEY" >$PWD/deploy_key
chmod 0400 $PWD/deploy_key

# Activate logging of bash commands now that the sensitive stuff is done
set -x

# As chartpress uses git to push to our Helm chart repository, we configure
# git ahead of time to use the identity we decrypted earlier.
git config --global user.email $GIT_EMAIL
git config --global user.name $GIT_USER
export GIT_SSH_COMMAND="ssh -i $PWD/deploy_key"

echo "Publishing chart via chartpress..."
cp .github/scripts/chartpress.yaml . # chartpress doesn't support custom path for config
chartpress --skip-build --publish-chart --tag "$VERSION"

# Let us log the changes chartpress did, it should include replacements for
# fields in values.yaml, such as what tag for various images we are using.
echo "Changes from chartpress:"
git --no-pager diff
