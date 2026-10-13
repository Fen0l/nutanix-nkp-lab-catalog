#!/usr/bin/env bash
# Add applications/<app>/<new-version>/ as a copy of the latest version with the new chart tag.
# Usage: scripts/bump-app.sh <app> <new-chart-version>
set -euo pipefail
[[ $# -eq 2 ]] || { sed -n '2,3p' "$0"; exit 1; }
app="$1"; new="$2"
cd "$(dirname "$0")/.."

appdir="applications/${app}"
current="$(ls -d "$appdir"/*/ | xargs -n1 basename | grep -v '^\.' | sort -V | tail -1)"
[[ "$current" != "$new" ]] || { echo "${app} is already at ${new}"; exit 0; }
[[ ! -d "$appdir/$new" ]] || { echo "${appdir}/${new} already exists"; exit 1; }

cp -a "$appdir/$current" "$appdir/$new"
hr="$appdir/$new/helmrelease/helmrelease.yaml"
current_tag="$(yq -r 'select(.kind == "OCIRepository") | .spec.ref.tag' "$hr" | head -1)"
new_tag="$new"
[[ "$current_tag" == v* && "$new" != v* ]] && new_tag="v${new}"
sed -i -E "0,/^(\s*tag: )${current_tag//./\\.}$/s//\\1${new_tag}/" "$hr"

if [[ -f "$appdir/.catalog-source.yaml" ]]; then
  repo="$(yq -r '.helmrepoUrl' "$appdir/.catalog-source.yaml")"; chart="$(yq -r '.chart' "$appdir/.catalog-source.yaml")"
  echo "NOTE: mirrored chart. Run: scripts/mirror-chart.sh ${repo} ${chart} ${new}"
fi
sed -i "/^upgradesFrom:/d; /^overview:/i upgradesFrom: \"${current}\"" "$appdir/$new/metadata.yaml"
esc="${current//./\\.}"
sed -i -E "s/^(  - Chart: \`[^\`]+\`) ${esc}\b/\\1 ${new}/" "$appdir/$new/metadata.yaml"
for f in "$appdir/$new/helmrelease/dashboard.yaml" "$appdir/$new/user-inputs.yaml"; do
  [[ -f "$f" ]] && sed -i -E "s/^(\s*version: )\"${esc}\"/\\1\"${new}\"/" "$f"
done
echo "Created ${appdir}/${new} from ${current}. Review helmrelease/cm.yaml and metadata.yaml, then scripts/validate.sh ${app}=${new}"
