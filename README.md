# ITCS Catalog for Nutanix NKP

[![Validate](https://github.com/Fen0l/nutanix-nkp-lab-catalog/actions/workflows/validate.yaml/badge.svg)](https://github.com/Fen0l/nutanix-nkp-lab-catalog/actions/workflows/validate.yaml)
[![E2E install tests](https://github.com/Fen0l/nutanix-nkp-lab-catalog/actions/workflows/e2e.yaml/badge.svg)](https://github.com/Fen0l/nutanix-nkp-lab-catalog/actions/workflows/e2e.yaml)
[![Publish](https://github.com/Fen0l/nutanix-nkp-lab-catalog/actions/workflows/publish.yaml/badge.svg)](https://github.com/Fen0l/nutanix-nkp-lab-catalog/actions/workflows/publish.yaml)
[![Applications](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2FFen0l%2Fnutanix-nkp-lab-catalog%2Fmain%2Fdocs%2Fbadges%2Fapplications.json)](#applications)
[![Upstream](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2FFen0l%2Fnutanix-nkp-lab-catalog%2Fmain%2Fdocs%2Fbadges%2Fupstream.json)](docs/VERSIONS.md)
[![NKP 2.18](https://img.shields.io/badge/NKP-2.18-7a3ff2)](https://portal.nutanix.com/page/documents/details?targetId=Nutanix-Kubernetes-Platform-v2_18)
[![Kubernetes 1.35](https://img.shields.io/badge/Kubernetes-1.35-326ce5?logo=kubernetes&logoColor=white)](https://kubernetes.io/)
[![OCI catalog](https://img.shields.io/badge/catalog-OCI%20%C2%B7%20ghcr.io-orange)](https://github.com/Fen0l/nutanix-nkp-lab-catalog/pkgs/container/nutanix-nkp-lab-catalog%2Fcollection)
[![License](https://img.shields.io/github/license/Fen0l/nutanix-nkp-lab-catalog)](LICENSE)

44 widely used Kubernetes tools for the **Nutanix Kubernetes Platform (NKP)**, one click away in the
NKP UI. Every entry is the **upstream Helm chart with upstream defaults**; the catalog only adds the
NKP metadata, the ingress and dashboard card for web UIs, and the few values a chart cannot start
without. Everything is tested on NKP 2.18 and on kind before it is published.

> Not an official Nutanix product, not supported by Nutanix or the upstream projects. Labs and demos
> only. For Nutanix products use the [Nutanix Product Catalog](https://github.com/nutanix-cloud-native/nkp-nutanix-product-catalog)
> and the [Partner Catalog](https://github.com/nutanix-cloud-native/nkp-partner-catalog).

## Quick start

NKP **2.18**, **Ultimate** licence, management cluster able to reach `ghcr.io`.

```bash
nkp create catalog-collection \
  --url oci://ghcr.io/fen0l/nutanix-nkp-lab-catalog/collection \
  --tag 2.18 \
  --workspace kommander-workspace
```

Then **Workspace → Applications** in the NKP UI, filter by category (`gitops`, `security`,
`observability`, ...). Pick one, **Enable**, **Deploy**. Apps with a web UI get a card with an open button
on the cluster page. Each application page shows the exact values the catalog sets, how to reach
the tool and what to know on NKP 2.18.

The CLI equivalent:

```bash
kubectl get apps -n kommander
nkp create appdeployment podinfo --app podinfo-6.15.0 --workspace <workspace>
```

Some UIs cannot run under a sub-path (OpenBao, Argo Rollouts, Gitea, Backstage, Policy Reporter,
Falco, Litmus, Portainer): enabling them asks for the cluster's **ingress domain** and serves them as
`<app>.<domain>`. Point a wildcard DNS record `*.<domain>` at the NKP Traefik address, or use
`<traefik-ip>.nip.io`.

## Applications

Grouped by the category shown in the NKP UI. Chart = Helm chart version, App = application version. Instances: *one per cluster* entries are
enabled from a workspace; *per project / workspace* entries can be deployed several times, each in
its own namespace. Web UI is served by the NKP Traefik on the cluster ingress.

### `gitops`

| | Application | Chart | App | Instances | Web UI | Requires |
|---|---|---|---|---|---|---|
| <img src="docs/icons/argo-cd.svg" width="20" height="20" alt=""> | [Argo CD](applications/argo-cd) | 10.10.3 | v3.5.5 | one per cluster | `/argo-cd/` |  |
| <img src="docs/icons/argo-workflows.svg" width="20" height="20" alt=""> | [Argo Workflows](applications/argo-workflows) | 2.0.12 | v4.1.5 | one per cluster | `/argo-workflows/` |  |
| <img src="docs/icons/argo-events.svg" width="20" height="20" alt=""> | [Argo Events](applications/argo-events) | 2.4.27 | v1.9.11 | one per cluster |  |  |
| <img src="docs/icons/argo-rollouts.svg" width="20" height="20" alt=""> | [Argo Rollouts](applications/argo-rollouts) | 2.43.6 | v1.10.0 | one per cluster | `rollouts.<domain>` |  |
| <img src="docs/icons/flagger.svg" width="20" height="20" alt=""> | [Flagger](applications/flagger) | 1.45.0 | 1.45.0 | one per cluster |  |  |
| <img src="docs/icons/sealed-secrets.svg" width="20" height="20" alt=""> | [Sealed Secrets](applications/sealed-secrets) | 2.20.0 | 0.40.0 | one per cluster |  |  |

### `developer-tools`

| | Application | Chart | App | Instances | Web UI | Requires |
|---|---|---|---|---|---|---|
| <img src="docs/icons/gitea.svg" width="20" height="20" alt=""> | [Gitea](applications/gitea) | 12.7.0 | 1.27.0 | per project / workspace | `gitea.<domain>` |  |
| <img src="docs/icons/jenkins.svg" width="20" height="20" alt=""> | [Jenkins](applications/jenkins) | 5.9.68 | 2.580.1 | one per cluster | `/jenkins/` |  |
| <img src="docs/icons/backstage.svg" width="20" height="20" alt=""> | [Backstage](applications/backstage) | 2.10.2 | demo image | per project / workspace | `backstage.<domain>` |  |
| <img src="docs/icons/podinfo.svg" width="20" height="20" alt=""> | [Podinfo](applications/podinfo) | 6.15.0 | 6.15.0 | per project / workspace | `/<instance>/` |  |
| <img src="docs/icons/headlamp.svg" width="20" height="20" alt=""> | [Headlamp](applications/headlamp) | 0.45.0 | 0.45.0 | one per cluster | `/headlamp/` |  |
| <img src="docs/icons/portainer.svg" width="20" height="20" alt=""> | [Portainer CE](applications/portainer) | 245.2.0 | 2.45.2 | one per cluster | `portainer.<domain>` |  |
| <img src="docs/icons/vcluster.svg" width="20" height="20" alt=""> | [vCluster](applications/vcluster) | 0.37.3 | 0.37.3 | per project / workspace |  |  |

### `testing`

| | Application | Chart | App | Instances | Web UI | Requires |
|---|---|---|---|---|---|---|
| <img src="docs/icons/k6-operator.svg" width="20" height="20" alt=""> | [k6 Operator](applications/k6-operator) | 4.6.0 | 1.6.0 | one per cluster |  |  |
| <img src="docs/icons/litmus.svg" width="20" height="20" alt=""> | [LitmusChaos](applications/litmus) | 3.30.0 | 3.30.0 | one per cluster | `litmus.<domain>` |  |

### `security`

| | Application | Chart | App | Instances | Web UI | Requires |
|---|---|---|---|---|---|---|
| <img src="docs/icons/falco.svg" width="20" height="20" alt=""> | [Falco](applications/falco) | 9.2.0 | 0.45.0 | one per cluster | `falco.<domain>` |  |
| <img src="docs/icons/tetragon.svg" width="20" height="20" alt=""> | [Tetragon](applications/tetragon) | 1.7.1 | 1.7.1 | one per cluster |  |  |
| <img src="docs/icons/trivy-operator.svg" width="20" height="20" alt=""> | [Trivy Operator](applications/trivy-operator) | 0.32.1 | 0.30.1 | one per cluster |  |  |
| <img src="docs/icons/kubescape-operator.svg" width="20" height="20" alt=""> | [Kubescape Operator](applications/kubescape-operator) | 1.40.5 | 1.40.5 | one per cluster |  |  |
| <img src="docs/icons/kyverno.svg" width="20" height="20" alt=""> | [Kyverno](applications/kyverno) | 3.9.1 | v1.19.1 | one per cluster |  |  |
| <img src="docs/icons/kyverno-policies.svg" width="20" height="20" alt=""> | [Kyverno Policies (Pod Security)](applications/kyverno-policies) | 3.9.1 | v1.19.1 | one per cluster |  | kyverno |
| <img src="docs/icons/policy-reporter.svg" width="20" height="20" alt=""> | [Policy Reporter](applications/policy-reporter) | 3.11.0 | 3.11.0 | one per cluster | `policy-reporter.<domain>` | kyverno, trivy-operator |
| <img src="docs/icons/openbao.svg" width="20" height="20" alt=""> | [OpenBao](applications/openbao) | 0.30.3 | v2.7.1 | one per cluster | `openbao.<domain>` |  |
| <img src="docs/icons/keycloak.svg" width="20" height="20" alt=""> | [Keycloak](applications/keycloak) | 7.3.2 | 26.7.4 | one per cluster | `/keycloak/` |  |
| <img src="docs/icons/trust-manager.svg" width="20" height="20" alt=""> | [trust-manager](applications/trust-manager) | 0.25.0 | v0.25.0 | one per cluster |  | cert-manager |
| <img src="docs/icons/cert-manager-approver-policy.svg" width="20" height="20" alt=""> | [cert-manager approver-policy](applications/cert-manager-approver-policy) | 0.28.0 | v0.28.0 | one per cluster |  | cert-manager |
| <img src="docs/icons/capsule.svg" width="20" height="20" alt=""> | [Capsule](applications/capsule) | 0.14.6 | 0.14.6 | one per cluster |  | cert-manager |

### `observability`

| | Application | Chart | App | Instances | Web UI | Requires |
|---|---|---|---|---|---|---|
| <img src="docs/icons/opentelemetry-collector.svg" width="20" height="20" alt=""> | [OpenTelemetry Collector](applications/opentelemetry-collector) | 0.175.1 | 0.161.0 | per project / workspace |  |  |
| <img src="docs/icons/tempo.svg" width="20" height="20" alt=""> | [Grafana Tempo](applications/tempo) | 1.24.4 | 2.9.0 | per project / workspace |  |  |
| <img src="docs/icons/victoria-logs.svg" width="20" height="20" alt=""> | [VictoriaLogs](applications/victoria-logs) | 0.13.10 | v1.53.0 | one per cluster | `/victoria-logs/select/vmui/` |  |

### `networking`

| | Application | Chart | App | Instances | Web UI | Requires |
|---|---|---|---|---|---|---|
| <img src="docs/icons/ingress-nginx.svg" width="20" height="20" alt=""> | [NGINX Ingress Controller](applications/ingress-nginx) | 4.15.1 | 1.15.1 | one per cluster |  |  |

### `api-gateway`

| | Application | Chart | App | Instances | Web UI | Requires |
|---|---|---|---|---|---|---|
| <img src="docs/icons/envoy-gateway.svg" width="20" height="20" alt=""> | [Envoy Gateway](applications/envoy-gateway) | 1.9.2 | v1.9.2 | one per cluster |  |  |
| <img src="docs/icons/nginx-gateway-fabric.svg" width="20" height="20" alt=""> | [NGINX Gateway Fabric](applications/nginx-gateway-fabric) | 2.7.2 | 2.7.2 | one per cluster |  | gateway-api-crds |

### `infrastructure`

| | Application | Chart | App | Instances | Web UI | Requires |
|---|---|---|---|---|---|---|
| <img src="docs/icons/keda.svg" width="20" height="20" alt=""> | [KEDA](applications/keda) | 2.21.0 | 2.21.0 | one per cluster |  |  |
| <img src="docs/icons/keda-add-ons-http.svg" width="20" height="20" alt=""> | [KEDA HTTP Add-on](applications/keda-add-ons-http) | 0.16.0 | 0.16.0 | one per cluster |  | keda |
| <img src="docs/icons/descheduler.svg" width="20" height="20" alt=""> | [Descheduler](applications/descheduler) | 0.37.0 | 0.37.0 | one per cluster |  |  |
| <img src="docs/icons/kured.svg" width="20" height="20" alt=""> | [Kured (Reboot Daemon)](applications/kured) | 6.1.0 | 1.23.0 | one per cluster |  |  |
| <img src="docs/icons/reflector.svg" width="20" height="20" alt=""> | [Reflector](applications/reflector) | 10.0.65 | 10.0.65 | one per cluster |  |  |
| <img src="docs/icons/crossplane.svg" width="20" height="20" alt=""> | [Crossplane](applications/crossplane) | 2.4.2 | 2.4.2 | one per cluster |  |  |
| <img src="docs/icons/dapr.svg" width="20" height="20" alt=""> | [Dapr](applications/dapr) | 1.18.7 | 1.18.7 | one per cluster |  |  |

### `data`

| | Application | Chart | App | Instances | Web UI | Requires |
|---|---|---|---|---|---|---|
| <img src="docs/icons/strimzi-kafka-operator.svg" width="20" height="20" alt=""> | [Strimzi (Apache Kafka)](applications/strimzi-kafka-operator) | 1.3.0 | 1.3.0 | one per cluster |  |  |
| <img src="docs/icons/kafka-ui.svg" width="20" height="20" alt=""> | [Kafka UI (Kafbat)](applications/kafka-ui) | 1.6.5 | v1.5.0 | one per cluster | `/kafka-ui/` |  |
| <img src="docs/icons/mongodb-community-operator.svg" width="20" height="20" alt=""> | [MongoDB Community Operator](applications/mongodb-community-operator) | 0.13.0 | 0.13.0 | one per cluster |  |  |
| <img src="docs/icons/redis-operator.svg" width="20" height="20" alt=""> | [Redis Operator (OT)](applications/redis-operator) | 0.27.0 | 0.27.0 | one per cluster |  |  |

Versions versus upstream, refreshed daily: [docs/VERSIONS.md](docs/VERSIONS.md).

**Already in NKP, therefore not duplicated here**: cert-manager, kube-prometheus-stack, Grafana,
Loki, Fluent Bit, Velero, Harbor, Istio, Kiali, Jaeger, Gatekeeper, External Secrets, External DNS,
CloudNativePG, Traefik, Rook Ceph, Knative, NVIDIA GPU Operator, Kubernetes Dashboard, Gateway API
CRDs (platform applications), plus NDK, NAI, KServe, OpenTelemetry Operator, Envoy Gateway NAI,
Kueue and the AI stack (Nutanix product and AI catalogs). Hubble relay already runs on NKP Cilium
clusters; Hubble UI is a value toggle on the existing Cilium release, not a separate chart.

## How it works

- `nkp create catalog-collection` registers the OCI collection; Kommander creates one `App` per
  `applications/<app>/<version>/` and the UI lists them.
- Enabling an app creates an `AppDeployment`; Kommander applies the entry's Flux `Kustomization`
  to the target cluster, which pulls the chart and reconciles the `HelmRelease`.
- Overrides are a ConfigMap `<app>-config-overrides` (key `values.yaml`) referenced by the
  AppDeployment; the catalog's own values are the ConfigMap `<app>-config-defaults`.
- New chart version = new `<version>/` directory, never an in-place edit; a daily workflow opens
  one pull request per outdated application.


## License

Apache-2.0 for the catalog manifests and tooling. Each packaged application keeps its own
upstream licence.
