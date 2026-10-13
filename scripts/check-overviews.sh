#!/usr/bin/env bash
# Lint the UI texts: fixed overview sections, every cm.yaml value in the "Changes from upstream"
# block, install and web-UI sentence in the description.
# Usage: scripts/check-overviews.sh [app ...]
set -euo pipefail
cd "$(dirname "$0")/.."
apps=("$@"); [[ ${#apps[@]} -eq 0 ]] && apps=($(ls applications))
rc=0
for app in "${apps[@]}"; do
  ver="$(ls -d applications/"$app"/*/ | sort -V | tail -1)"
  meta="${ver}metadata.yaml"; cm="${ver}helmrelease/cm.yaml"
  overview="$(yq -r '.overview' "$meta")"; desc="$(yq -r '.description' "$meta")"
  for s in "**At a glance**" "**Changes from upstream**" "**Links**"; do
    grep -qF -- "$s" <<<"$overview" || { echo "$app: overview lacks the section $s"; rc=1; }
  done
  grep -qE "Installs into namespace|One instance per AppDeployment" <<<"$desc" || { echo "$app: description does not say where it installs"; rc=1; }
  grep -qE "Web UI|No web UI" <<<"$desc" || { echo "$app: description does not say whether it has a web UI"; rc=1; }
  changes="$(awk '/^\*\*Changes from upstream\*\*/{f=1;next} /^\*\*/{f=0} f' <<<"$overview")"
  while read -r key; do
    [[ -z "$key" || "$key" == "null" ]] && continue
    grep -qE -- "^${key}:" <<<"$changes" || { echo "$app: value '$key' set in cm.yaml but missing from the YAML block under 'Changes from upstream'"; rc=1; }
  done < <(yq -r 'select(.kind == "ConfigMap" and .data["values.yaml"] != null) | .data["values.yaml"]' "$cm" | yq -r 'keys | .[]' 2>/dev/null || true)
done
[[ $rc -eq 0 ]] && echo "check-overviews: ${#apps[@]} application(s) OK"
exit $rc
