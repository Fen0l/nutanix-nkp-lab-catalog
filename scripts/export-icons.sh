#!/usr/bin/env bash
# Export each application's icon from metadata.yaml to docs/icons/<app>.svg (used by the README).
set -euo pipefail
cd "$(dirname "$0")/.."
mkdir -p docs/icons
for app in applications/*/; do
  a="$(basename "$app")"; ver="$(ls -d "$app"*/ | sort -V | tail -1)"
  yq -r '.icon' "${ver}metadata.yaml" | base64 -d > "docs/icons/${a}.svg"
done
echo "exported $(ls docs/icons/*.svg | wc -l | tr -d ' ') icons"
