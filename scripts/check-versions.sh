#!/usr/bin/env bash
# Compare each application's latest catalog version with the latest upstream chart.
# Usage: scripts/check-versions.sh [--json|--markdown]
set -euo pipefail
cd "$(dirname "$0")/.."

mode="table"
case "${1:-}" in
  --json) mode="json" ;;
  --markdown) mode="markdown" ;;
  "") ;;
  *) echo "usage: $0 [--json|--markdown]" >&2; exit 1 ;;
esac

show_chart() {
  local out
  for _ in 1 2 3; do
    if out="$(timeout 60 helm show chart "$@" 2>/dev/null)"; then echo "$out"; return 0; fi
    sleep 2
  done
  return 1
}

rows=()
outdated=0
for appdir in applications/*/; do
  app="$(basename "$appdir")"
  current="$(find "$appdir" -mindepth 1 -maxdepth 1 -type d ! -name '.*' -printf '%f\n' | sort -V | tail -1)"
  hr="$appdir$current/helmrelease/helmrelease.yaml"
  src="$appdir.catalog-source.yaml"
  if [[ -f "$src" ]]; then
    repo="$(yq -r '.helmrepoUrl' "$src")"; chart="$(yq -r '.chart' "$src")"
    source="${repo} (${chart}, mirrored)"
    meta="$(show_chart "$chart" --repo "$repo" || true)"
  else
    url="$(yq -r 'select(.kind == "OCIRepository") | .spec.url' "$hr" | head -1)"
    source="$url"
    meta="$(show_chart "$url" || true)"
  fi
  latest="$(echo "$meta" | yq -r '.version // "unknown"' 2>/dev/null || echo unknown)"
  latest_app="$(echo "$meta" | yq -r '.appVersion // "-"' 2>/dev/null || echo -)"
  [[ -z "$meta" ]] && latest="unknown"
  latest="${latest#v}"
  status="up-to-date"
  if [[ "$latest" == "unknown" ]]; then status="unknown"
  elif [[ "$latest" != "$current" ]]; then status="outdated"; outdated=$((outdated + 1)); fi
  rows+=("$app|$current|$latest|$latest_app|$status|$source")
done

case "$mode" in
  table)
    printf '%-26s %-12s %-12s %-12s %s\n' APP CATALOG UPSTREAM APP-VERSION STATUS
    for r in "${rows[@]}"; do IFS='|' read -r a c l la s _ <<<"$r"; printf '%-26s %-12s %-12s %-12s %s\n' "$a" "$c" "$l" "$la" "$s"; done
    echo; echo "${outdated} application(s) outdated"
    ;;
  json)
    items=()
    for r in "${rows[@]}"; do IFS='|' read -r a c l la s src <<<"$r"
      items+=("$(jq -cn --arg app "$a" --arg current "$c" --arg latest "$l" --arg appVersion "$la" --arg status "$s" --arg source "$src" '$ARGS.named')")
    done
    printf '[%s]\n' "$(IFS=,; echo "${items[*]}")"
    ;;
  markdown)
    echo "# Upstream version report"
    echo
    echo "Generated $(date -u +'%Y-%m-%d %H:%M UTC') by \`scripts/check-versions.sh --markdown\`."
    echo "A new chart version becomes a new \`applications/<app>/<version>/\` directory (\`scripts/bump-app.sh\`);"
    echo "the daily workflow opens one pull request per outdated application."
    echo
    echo "| Application | Catalog | Upstream chart | Upstream app | Status | Source |"
    echo "|---|---|---|---|---|---|"
    for r in "${rows[@]}"; do IFS='|' read -r a c l la s src <<<"$r"
      mark=""; [[ "$s" == "outdated" ]] && mark=" :arrow_up:"; [[ "$s" == "unknown" ]] && mark=" :grey_question:"
      echo "| $a | $c | $l | $la | ${s}${mark} | \`$src\` |"
    done
    echo
    echo "**${outdated} application(s) outdated.**"
    ;;
esac
