# Upstream version report

Generated 2026-10-21 11:23 UTC by `scripts/check-versions.sh --markdown`.
A new chart version becomes a new `applications/<app>/<version>/` directory (`scripts/bump-app.sh`);
the daily workflow opens one pull request per outdated application.

| Application | Catalog | Upstream chart | Upstream app | Status | Source |
|---|---|---|---|---|---|
| argo-cd | 10.10.5 | 10.10.5 | v3.5.5 | up-to-date | `oci://ghcr.io/argoproj/argo-helm/argo-cd` |
| argo-events | 2.4.27 | 2.4.27 | v1.9.11 | up-to-date | `oci://ghcr.io/argoproj/argo-helm/argo-events` |
| argo-rollouts | 2.43.8 | 2.43.8 | v1.10.1 | up-to-date | `oci://ghcr.io/argoproj/argo-helm/argo-rollouts` |
| argo-workflows | 2.0.13 | 2.0.13 | v4.1.5 | up-to-date | `oci://ghcr.io/argoproj/argo-helm/argo-workflows` |
| backstage | 2.10.2 | 2.10.2 | - | up-to-date | `oci://ghcr.io/backstage/charts/backstage` |
| capsule | 0.14.6 | 0.14.6 | 0.14.6 | up-to-date | `oci://ghcr.io/projectcapsule/charts/capsule` |
| cert-manager-approver-policy | 0.28.0 | 0.28.0 | v0.28.0 | up-to-date | `oci://quay.io/jetstack/charts/cert-manager-approver-policy` |
| crossplane | 2.4.2 | 2.4.2 | 2.4.2 | up-to-date | `https://charts.crossplane.io/stable (crossplane, mirrored)` |
| dapr | 1.18.7 | 1.18.7 | 1.18.7 | up-to-date | `https://dapr.github.io/helm-charts (dapr, mirrored)` |
| descheduler | 0.37.0 | 0.37.0 | 0.37.0 | up-to-date | `https://kubernetes-sigs.github.io/descheduler (descheduler, mirrored)` |
| envoy-gateway | 1.9.2 | 1.9.2 | v1.9.2 | up-to-date | `oci://docker.io/envoyproxy/gateway-helm` |
| falco | 9.2.0 | 9.2.1 | 0.45.1 | outdated :arrow_up: | `oci://ghcr.io/falcosecurity/charts/falco` |
| flagger | 1.45.0 | 1.45.0 | 1.45.0 | up-to-date | `oci://ghcr.io/fluxcd/charts/flagger` |
| gitea | 12.7.0 | 12.7.0 | 1.27.0 | up-to-date | `oci://docker.gitea.com/charts/gitea` |
| headlamp | 0.45.0 | 0.46.0 | 0.46.0 | outdated :arrow_up: | `https://kubernetes-sigs.github.io/headlamp/ (headlamp, mirrored)` |
| ingress-nginx | 4.15.1 | 4.15.1 | 1.15.1 | up-to-date | `https://kubernetes.github.io/ingress-nginx (ingress-nginx, mirrored)` |
| jenkins | 5.9.69 | 5.9.72 | 2.580.1 | outdated :arrow_up: | `oci://ghcr.io/jenkinsci/helm-charts/jenkins` |
| k6-operator | 4.6.0 | 4.7.0 | 1.7.0 | outdated :arrow_up: | `oci://ghcr.io/grafana/helm-charts/k6-operator` |
| kafka-ui | 1.6.5 | 1.6.5 | v1.5.0 | up-to-date | `https://ui.charts.kafbat.io (kafka-ui, mirrored)` |
| keda-add-ons-http | 0.16.0 | 0.16.0 | 0.16.0 | up-to-date | `https://kedacore.github.io/charts (keda-add-ons-http, mirrored)` |
| keda | 2.21.0 | 2.21.0 | 2.21.0 | up-to-date | `https://kedacore.github.io/charts (keda, mirrored)` |
| keycloak | 7.3.2 | 7.4.0 | 26.8.1 | outdated :arrow_up: | `https://codecentric.github.io/helm-charts (keycloakx, mirrored)` |
| kubescape-operator | 1.40.5 | 1.40.5 | 1.40.5 | up-to-date | `https://kubescape.github.io/helm-charts (kubescape-operator, mirrored)` |
| kured | 6.1.0 | 6.1.0 | 1.23.0 | up-to-date | `oci://ghcr.io/kubereboot/charts/kured` |
| kyverno-policies | 3.9.1 | 3.9.2 | v1.19.2 | outdated :arrow_up: | `oci://ghcr.io/kyverno/charts/kyverno-policies` |
| kyverno | 3.9.1 | 3.9.2 | v1.19.2 | outdated :arrow_up: | `oci://ghcr.io/kyverno/charts/kyverno` |
| litmus | 3.30.0 | 3.30.0 | 3.30.0 | up-to-date | `oci://ghcr.io/litmuschaos/litmus-helm/litmus` |
| mongodb-community-operator | 0.13.0 | 0.13.0 | 0.13.0 | up-to-date | `https://mongodb.github.io/helm-charts (community-operator, mirrored)` |
| nginx-gateway-fabric | 2.8.0 | 2.8.0 | 2.8.0 | up-to-date | `oci://ghcr.io/nginx/charts/nginx-gateway-fabric` |
| openbao | 0.30.4 | 0.30.4 | v2.7.1 | up-to-date | `oci://ghcr.io/openbao/charts/openbao` |
| opentelemetry-collector | 0.176.0 | 0.176.1 | 0.162.0 | outdated :arrow_up: | `oci://ghcr.io/open-telemetry/opentelemetry-helm-charts/opentelemetry-collector` |
| podinfo | 6.15.0 | 6.15.0 | 6.15.0 | up-to-date | `oci://ghcr.io/stefanprodan/charts/podinfo` |
| policy-reporter | 3.11.1 | 3.11.1 | 3.11.1 | up-to-date | `oci://ghcr.io/kyverno/charts/policy-reporter` |
| portainer | 245.2.0 | 245.2.0 | ce-latest-ee-2.45.2 | up-to-date | `https://portainer.github.io/k8s (portainer, mirrored)` |
| redis-operator | 0.27.0 | 0.27.0 | 0.27.0 | up-to-date | `https://ot-container-kit.github.io/helm-charts (redis-operator, mirrored)` |
| reflector | 10.0.65 | 10.0.65 | 10.0.65 | up-to-date | `oci://ghcr.io/emberstack/helm-charts/reflector` |
| sealed-secrets | 2.20.0 | 2.20.0 | 0.40.0 | up-to-date | `https://bitnami.github.io/sealed-secrets (sealed-secrets, mirrored)` |
| strimzi-kafka-operator | 1.3.0 | 1.3.0 | 1.3.0 | up-to-date | `oci://quay.io/strimzi-helm/strimzi-kafka-operator` |
| tempo | 1.24.4 | 1.24.4 | 2.9.0 | up-to-date | `oci://ghcr.io/grafana/helm-charts/tempo` |
| tetragon | 1.7.1 | 1.7.1 | 1.7.1 | up-to-date | `https://helm.cilium.io (tetragon, mirrored)` |
| trivy-operator | 0.32.1 | 0.32.1 | 0.30.1 | up-to-date | `oci://ghcr.io/aquasecurity/helm-charts/trivy-operator` |
| trust-manager | 0.25.0 | 0.25.0 | v0.25.0 | up-to-date | `oci://quay.io/jetstack/charts/trust-manager` |
| vcluster | 0.37.3 | 0.37.4 | 0.37.4 | outdated :arrow_up: | `https://charts.loft.sh (vcluster, mirrored)` |
| victoria-logs | 0.13.10 | 0.13.10 | v1.53.0 | up-to-date | `oci://ghcr.io/victoriametrics/helm-charts/victoria-logs-single` |

**9 application(s) outdated.**
