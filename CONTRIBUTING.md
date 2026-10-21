# Contributing

Thanks for adding to the ITCS Catalog for Nutanix NKP. The rules below keep every entry predictable
for the people who deploy it with one click.

## Ground rules

1. **Upstream chart, upstream defaults.** Reference the project's own OCI chart. Only add values the
   chart needs to render or become Ready, plus the NKP integration block (ingress on the cluster
   Traefik, dashboard card). Comment every added value in `helmrelease/cm.yaml` and explain it in
   the overview.
2. **HTTP-only charts are mirrored unchanged** to `oci://ghcr.io/fen0l/nutanix-nkp-lab-catalog/charts/<chart>`
   through the mirror workflow (packages created by the workflow are public). Declare the source in
   `applications/<app>/.catalog-source.yaml` so version tracking keeps working.
3. **New version = new directory.** Never edit a published `applications/<app>/<version>/`; use
   `scripts/bump-app.sh <app> <version>`.
4. **Literal HelmRelease names.** `metadata.name: <app>`; `spec.releaseName: ${releaseName}` for
   multi-instance apps. Kommander's override patch targets the literal name.
5. **Dependencies are declared.** `requiredDependencies` for anything the chart cannot run without
   (CRDs, cert-manager, an operator), `dependencies` for soft suggestions. Platform application names
   (`cert-manager`, `gateway-api-crds`, `kube-prometheus-stack`) and catalog entry names are both valid.
6. **Every entry is tested.** `nkp validate catalog-repository` must pass and the kind install test
   (`cd apptests && go test ./suites/ -v -app=<app>`) must reach Ready. Add prerequisites for the kind
   cluster in `apptests/catalog/prereqs.go` when the app needs CRDs from another entry. The kind
   cluster runs the Flux release line of NKP 2.19, unchanged since 2.18 (`apptests/catalog/flux.go`, `FLUX_VERSION` to
   test another one).

## Adding an application

```bash
scripts/add-app.sh <app> <chart-version> oci://<registry>/<chart> [--namespace <ns>] [--multi]
```

Then fill in:

- `metadata.yaml`: `displayName`, `description`, `category` (one of `gitops`, `developer-tools`, `testing`,
  `security`, `observability`, `networking`, `api-gateway`, `infrastructure`, `data`; lowercase like
  the NKP platform applications so the UI filter stays consistent), `scope`, `allowMultipleInstances`, `supportLink`, dependencies, a base64 SVG
  `icon`, and an `overview`. The NKP UI shows `description` on the application card and renders
  `overview` (Markdown) on the detail page, so both carry as much as they can:
  - `description`: what it is in one or two sentences, then "Installs into namespace `<ns>`." (or
    "One instance per AppDeployment, each in its own namespace.") and "Web UI at `/<path>/` on the
    cluster ingress." / "Web UI on `<app>.<ingress domain>` (asked when enabling)." / "No web UI."
  - `overview`, in this order: intro paragraph · **At a glance** (chart and source, namespace and
    instance model, web UI, requires / works with) · **Changes from upstream** (the `values.yaml`
    block of `cm.yaml`, comments included, so it can be copied into an override ConfigMap; an empty
    override example when the chart runs unchanged) ·
    **Defaults worth knowing** · **Access** · **Requires** / **Works with** · hands-on sections
    (**Try it**, **Example**, **Lab override example**) · **NKP 2.19 notes** (verified behaviour on
    NKP, not upstream folklore) · **Chart source** (mirrored charts) · **Links** (docs, chart, default
    values, releases).
  `scripts/check-overviews.sh` (run by `scripts/validate.sh`) fails when a value set in `cm.yaml`
  is not named under **Changes from upstream** or when the description lacks the install / UI
  sentence.
- `helmrelease/cm.yaml`: upstream defaults (empty) plus the NKP integration block if the app has a UI.
- `.bloodhound.yaml` (per app, only when needed): `skip_types` for custom resources the chart renders
  whose CRDs come from a dependency, or `additional_crds` URLs.

### Exposing a web UI

- The app serves a sub-path natively (`rootpath`, `--prefix`, `base-url`, `relativePath`, context
  path): put the ingress values in `cm.yaml` with `ingressClassName: kommander-traefik`, a host-less
  rule and `traefik.ingress.kubernetes.io/router.tls: "true"`, and add `helmrelease/dashboard.yaml`
  (ConfigMap `${releaseName}-app-dashboard-info`, label `kommander.d2iq.io/application`,
  `dashboardLink: /<path>/`).
- The app needs the prefix stripped (Argo Workflows): add a Traefik `Middleware` **in the app's
  namespace** plus a `Namespace` manifest, and reference it as `<namespace>-<name>@kubernetescrd`.
- The app cannot run under a sub-path: add `user-inputs.yaml` with an `ingressDomain` field whose
  templates render `<app>-config-overrides` (ingress host `<app>.{{ ingressDomain }}`) and
  `<app>-app-dashboard-info` (`dashboardLink: https://<app>.{{ ingressDomain }}/`). See
  `applications/openbao` for the pattern. The overview shows that override ConfigMap verbatim under
  **Ingress with a hostname**, so CLI users can create it by hand.

## Before opening a pull request

```bash
scripts/validate.sh <app>=<version>
cd apptests && go test ./suites/ -v -app=<app>
```

The PR runs the same validation and the install test for the changed applications. The daily
version-check workflow will pick the app up automatically once merged.
