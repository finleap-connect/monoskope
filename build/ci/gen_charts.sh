#!/usr/bin/env bash
set -euo pipefail

cd docs/flow-charts/

for f in *.mmd; do
    mmdc -p puppeteer-config.json -i $f -o images/$f.png
done
