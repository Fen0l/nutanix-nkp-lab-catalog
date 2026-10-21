#!/usr/bin/env bash
# Scaffold applications/<app>/<version>/ for a chart that is already in an OCI registry.
# Usage: scripts/add-app.sh <app> <chart-version> <oci-chart-url> [--namespace <ns>] [--multi]
set -euo pipefail

usage() { sed -n '2,3p' "$0"; exit 1; }
[[ $# -ge 3 ]] || usage

app="$1"; version="$2"; url="$3"; shift 3
ns=""; multi=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    --namespace) ns="$2"; shift 2 ;;
    --multi) multi=true; shift ;;
    *) usage ;;
  esac
done

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$repo_root"

nkp generate catalog-repository --repo-dir . --apps "${app}=${version}"

hr="applications/${app}/${version}/helmrelease/helmrelease.yaml"
yq -i "(select(.kind == \"OCIRepository\") | .spec.url) = \"${url}\"" "$hr"
yq -i "(select(.kind == \"OCIRepository\") | .spec.ref.tag) = \"${version}\"" "$hr"

if $multi; then
  yq -i '(select(.kind == "HelmRelease") | .spec.releaseName) = "${releaseName}"' "$hr"
else
  yq -i "(select(.kind == \"HelmRelease\") | .spec.releaseName) = \"${app}\"" "$hr"
fi
[[ -n "$ns" ]] && yq -i "(select(.kind == \"HelmRelease\") | .spec.targetNamespace) = \"${ns}\"" "$hr"

meta="applications/${app}/${version}/metadata.yaml"
yq -i '.type = "custom" | .nkpVersionSupport = ">=2.19.0"' "$meta"
$multi || yq -i '.allowMultipleInstances = false | .scope = ["workspace"]' "$meta"

echo "Scaffolded applications/${app}/${version}. Now edit:"
echo "  - ${meta}  (displayName, description, category, overview, supportLink, icon)"
echo "  - applications/${app}/${version}/helmrelease/cm.yaml  (default values, keep upstream defaults unless required)"
echo "Then: scripts/validate.sh"
