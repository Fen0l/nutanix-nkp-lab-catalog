#!/usr/bin/env bash
# Write the README badges (docs/badges/*.json): application count, upstream status.
# Usage: scripts/catalog-stats.sh [versions.json]   (output of scripts/check-versions.sh --json)
set -euo pipefail
cd "$(dirname "$0")/.."
apps="$(ls -d applications/*/ | wc -l | tr -d ' ')"
printf '{"schemaVersion":1,"label":"applications","message":"%s","color":"blue"}\n' "$apps" > docs/badges/applications.json
if [[ -n "${1:-}" && -s "$1" ]]; then
  outdated="$(jq '[.[] | select(.status == "outdated")] | length' "$1")"
  if [[ "$outdated" -eq 0 ]]; then msg="all up to date"; color="brightgreen"; else msg="${outdated} outdated"; color="orange"; fi
  printf '{"schemaVersion":1,"label":"upstream","message":"%s","color":"%s"}\n' "$msg" "$color" > docs/badges/upstream.json
fi
echo "badges: ${apps} applications$( [[ -n "${1:-}" ]] && echo ", ${msg}")"
