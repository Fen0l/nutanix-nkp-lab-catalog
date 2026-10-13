#!/usr/bin/env bash
# nkp validate + UI text lint.
# Usage: scripts/validate.sh [app=version ...]
set -euo pipefail
cd "$(dirname "$0")/.."
args=()
for a in "$@"; do args+=(--apps "$a"); done
nkp validate catalog-repository --repo-dir . "${args[@]}"
if [[ $# -eq 0 ]]; then scripts/check-overviews.sh; else scripts/check-overviews.sh "${@%%=*}"; fi
