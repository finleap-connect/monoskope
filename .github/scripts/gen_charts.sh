#!/usr/bin/env bash
# Copyright 2022 Monoskope Authors
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

pwd=$(pwd)
cd docs/flow-charts/

for f in *.mmd; do
    docker run -it -v $pwd/docs/flow-charts:/data minlag/mermaid-cli -p puppeteer-config.json -i /data/$f -o images/$f.png -w 1920 -H 1080 -t neutral
    # mmdc -p puppeteer-config.json -i $f -o images/$f.png -w 1920 -H 1080 -t neutral
done
