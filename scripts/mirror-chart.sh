#!/usr/bin/env bash
# Copy an unmodified chart from an HTTP Helm repository to the catalog's OCI registry (idempotent).
# Usage: scripts/mirror-chart.sh <helm-repo-url> <chart> <version> [oci-registry]   (needs helm registry login)
set -euo pipefail

[[ $# -ge 3 ]] || { sed -n '2,3p' "$0"; exit 1; }
repo_url="$1"; chart="$2"; version="$3"
registry="${4:-${OCI_REGISTRY:-oci://ghcr.io/fen0l/nutanix-nkp-lab-catalog/charts}}"

if helm show chart "${registry}/${chart}" --version "${version}" >/dev/null 2>&1; then
  echo "Already mirrored: ${registry}/${chart}:${version}"
  exit 0
fi

workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

helm pull "$chart" --repo "$repo_url" --version "$version" --destination "$workdir"
tgz="$(ls "$workdir"/*.tgz)"
sha256sum "$tgz"
helm push "$tgz" "$registry"
echo "Mirrored ${chart} ${version} -> ${registry}/${chart}:${version}"
