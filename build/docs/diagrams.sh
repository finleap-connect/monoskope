#!/usr/bin/env bash
set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd docs/diagrams

for f in *.mmd; do
    mmdc -p $DIR/puppeteer-config.json -i $f -o $f.png -w 1920 -H 1080 -t neutral
done
