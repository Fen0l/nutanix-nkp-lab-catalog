#!/usr/bin/env bash
# Make the mirrored chart packages public on GHCR. GitHub has no API to change package
# visibility; packages pushed by a workflow with the repository token are public, so private
# ones are deleted and re-created by the mirror-charts workflow (pushed to a bootstrap branch
# when the default branch does not carry it yet).
# Usage: scripts/publish-packages.sh [publish|status|manual] [--owner <user>] [--repo <name>] [--ref <branch>]
# Needs gh (repo, workflow, read:packages, delete:packages), git, curl, yq.
set -euo pipefail

mode="publish"
case "${1:-}" in status|manual|publish) mode="$1"; shift ;; esac
owner="Fen0l"; repo="nutanix-nkp-lab-catalog"; ref="ghcr-publish"
while [[ $# -gt 0 ]]; do case "$1" in --owner) owner="$2"; shift 2;; --repo) repo="$2"; shift 2;; --ref) ref="$2"; shift 2;; *) echo "unknown arg $1" >&2; exit 1;; esac; done
owner_lc="$(echo "$owner" | tr '[:upper:]' '[:lower:]')"
root="$(cd "$(dirname "$0")/.." && pwd)"; cd "$root"

enc() { printf '%s' "$1" | sed 's#/#%2F#g'; }
all_packages() { gh api --paginate "/user/packages?package_type=container&per_page=100" --jq '.[].name' 2>/dev/null | grep "^${repo}" || true; }
visibility() { gh api "/user/packages/container/$(enc "$1")" --jq '.visibility' 2>/dev/null || echo "missing"; }
anon_ok() { [[ "$(curl -s -o /dev/null -w '%{http_code}' "https://ghcr.io/token?scope=repository:${owner_lc}/$1:pull&service=ghcr.io")" == "200" ]]; }
settings_url() { echo "https://github.com/users/${owner}/packages/container/$(enc "$1")/settings"; }

status() {
  printf '%-55s %-9s %s\n' PACKAGE VISIBILITY ANONYMOUS-PULL
  for p in $(all_packages); do printf '%-55s %-9s %s\n' "$p" "$(visibility "$p")" "$(anon_ok "$p" && echo ok || echo denied)"; done
}

workflow_in_repo() { gh api "repos/${owner}/${repo}/contents/.github/workflows/mirror-charts.yaml?ref=$1" >/dev/null 2>&1; }

bootstrap_branch() {
  echo "==> pushing bootstrap branch '${ref}' (workflow + chart list only; main untouched)"
  local tmp; tmp="$(mktemp -d)"; trap 'rm -rf "$tmp"' RETURN
  git -C "$tmp" init -q -b "$ref"
  mkdir -p "$tmp/.github/workflows" "$tmp/scripts"
  cp .github/workflows/mirror-charts.yaml "$tmp/.github/workflows/"
  cp scripts/mirror-chart.sh "$tmp/scripts/"
  for src in applications/*/.catalog-source.yaml; do
    app="$(dirname "$src")"
    mkdir -p "$tmp/$app"; cp "$src" "$tmp/$app/"
    for v in "$app"/*/; do mkdir -p "$tmp/$v/helmrelease"; cp "$v/helmrelease/helmrelease.yaml" "$tmp/$v/helmrelease/"; done
  done
  cat > "$tmp/README.md" <<EOF
Bootstrap branch used by scripts/publish-packages.sh to run the mirror-charts workflow before the
catalog's main branch is published. Safe to delete once main carries .github/workflows/mirror-charts.yaml.
EOF
  git -C "$tmp" add -A
  git -C "$tmp" -c user.name="catalog-bot" -c user.email="catalog-bot@users.noreply.github.com" commit -q -m "ci: bootstrap mirror-charts workflow for GHCR publishing"
  git -C "$tmp" push -q --force "https://github.com/${owner}/${repo}.git" "HEAD:refs/heads/${ref}"
  echo "    pushed ${ref}"
}

case "$mode" in
  status) status ;;

  manual)
    for p in $(all_packages); do
      [[ "$(visibility "$p")" == "private" ]] || continue
      url="$(settings_url "$p")"; echo; echo "==> $p"; echo "    $url  (Danger zone -> Change visibility -> Public -> type the package name)"
      if [[ -n "${BROWSER:-}" ]]; then "$BROWSER" "$url" >/dev/null 2>&1 & else gh browse "$url" >/dev/null 2>&1 || true; fi
      until [[ "$(visibility "$p")" == "public" ]]; do sleep 5; done; echo "    public."
    done; echo; status ;;

  publish)
    if ! gh auth status 2>&1 | grep -q "delete:packages"; then
      echo "The gh token lacks the delete:packages scope needed to switch private packages to public." >&2
      echo "Run:  gh auth refresh -h github.com -s delete:packages   (then re-run this script)" >&2
      exit 1
    fi
    default="$(gh api "repos/${owner}/${repo}" --jq '.default_branch // ""')"
    if [[ -n "$default" ]] && workflow_in_repo "$default"; then
      run_ref="$default"; echo "==> workflow found on ${default}"
    else
      bootstrap_branch; run_ref="$ref"
    fi
    echo "==> switching private chart packages to public (delete; the workflow re-creates them public)"
    for p in $(all_packages); do
      [[ "$p" == "${repo}/charts/"* ]] || continue
      if [[ "$(visibility "$p")" == "private" ]]; then echo "    deleting private $p"; gh api -X DELETE "/user/packages/container/$(enc "$p")"; fi
    done
    echo "==> running mirror-charts.yaml on ${run_ref}"
    gh workflow run mirror-charts.yaml --repo "${owner}/${repo}" --ref "${run_ref}"
    sleep 10
    run_id="$(gh run list --repo "${owner}/${repo}" --workflow mirror-charts.yaml --branch "${run_ref}" --limit 1 --json databaseId --jq '.[0].databaseId')"
    gh run watch "$run_id" --repo "${owner}/${repo}" --exit-status
    echo; status
    private_left="$(for p in $(all_packages); do [[ "$(visibility "$p")" == private ]] && echo "$p"; done | wc -l)"
    [[ "$private_left" -eq 0 ]] && echo "All packages public." || echo "${private_left} package(s) still private: run 'scripts/publish-packages.sh manual'." ;;
esac
