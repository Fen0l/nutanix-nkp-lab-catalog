#!/usr/bin/env bash
# Mirror an unmodified Helm chart from an HTTP Helm repository to this catalog's OCI
# registry. NKP catalog applications require OCI charts; some upstream projects only
# publish an HTTP index (KEDA, Kubescape, Headlamp, ingress-nginx, ...).
#
# Usage:
#   scripts/mirror-chart.sh <helm-repo-url> <chart> <version> [oci-registry]
#
# Example:
#   scripts/mirror-chart.sh https://kedacore.github.io/charts keda 2.21.0
#
# Idempotent: a tag that already exists in the registry is left untouched.
# Requires: helm registry login ghcr.io. From GitHub Actions use the workflow token
# (packages inherit the repository's visibility); from a laptop a PAT with
# write:packages works but the new package starts private.
set -euo pipefail

[[ $# -ge 3 ]] || { sed -n '2,14p' "$0"; exit 1; }
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
