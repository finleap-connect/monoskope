#!/usr/bin/env bash
set -euo pipefail

pwd=$(pwd)
cd docs/flow-charts/

for f in *.mmd; do
    docker run -it -v $pwd/docs/flow-charts:/data minlag/mermaid-cli -p puppeteer-config.json -i /data/$f -o images/$f.png -w 1920 -H 1080 -t neutral
    # mmdc -p puppeteer-config.json -i $f -o images/$f.png -w 1920 -H 1080 -t neutral
done
