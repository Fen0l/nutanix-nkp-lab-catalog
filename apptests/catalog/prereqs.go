package catalog

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mesosphere/kommander-applications/apptests/environment"
)

// Prereq is something an application needs on the cluster before it can
// install. On NKP these are platform applications (cert-manager,
// gateway-api-crds) or other catalog entries (kyverno); on the throw-away kind
// cluster the tests install them first.
type Prereq struct {
	// Name is the HelmRelease name (chart prereq) or a label (manifest prereq).
	Name string
	// Chart is an OCI chart URL. Together with Version and Namespace it becomes a
	// HelmRelease in the test namespace.
	Chart     string
	Version   string
	Namespace string
	Values    string
	// ManifestURL is a plain Kubernetes manifest to apply (CRDs, for example).
	ManifestURL string
}

// IsHelm reports whether the prerequisite is installed as a HelmRelease.
func (p Prereq) IsHelm() bool { return p.Chart != "" }

var certManager = Prereq{
	Name:      "cert-manager",
	Chart:     "oci://quay.io/jetstack/charts/cert-manager",
	Version:   "v1.21.2",
	Namespace: "cert-manager",
	Values:    "crds:\n  enabled: true\n",
}

// Gateway API CRDs from the upstream release manifest. (A Helm chart is not an
// option: the CRDs exceed the 1 MiB Helm release secret limit.) The manifest is
// sanitised before applying: documents without a kind (comment-only) are dropped.
var gatewayAPICRDs = Prereq{
	Name:        "gateway-api-crds",
	ManifestURL: "https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.6.1/standard-install.yaml",
}

// Traefik CRDs (Middleware etc.). NKP ships Traefik on every cluster; the kind
// test cluster needs the CRDs so apps that add a Traefik Middleware can apply.
var traefikCRDs = Prereq{
	Name:        "traefik-crds",
	ManifestURL: "https://raw.githubusercontent.com/traefik/traefik/v3.7/docs/content/reference/dynamic-configuration/kubernetes-crd-definition-v1.yml",
}

// prereqs maps a catalog application to what must exist before it installs.
// Keep in sync with `dependencies` in the application's metadata.yaml.
var prereqs = map[string][]Prereq{
	"kyverno-policies":             {{Name: "kyverno", Chart: "oci://ghcr.io/kyverno/charts/kyverno", Version: "3.9.1", Namespace: "kyverno"}},
	"trust-manager":                {certManager},
	"capsule":                      {certManager},
	"cert-manager-approver-policy": {certManager},
	"keda-add-ons-http":            {{Name: "keda", Chart: "oci://ghcr.io/fen0l/nutanix-nkp-lab-catalog/charts/keda", Version: "2.21.0", Namespace: "keda"}},
	"nginx-gateway-fabric":         {gatewayAPICRDs},
	// Policy Reporter needs the PolicyReport API (installed by Kyverno) and, with the Trivy
	// plugin on, the Trivy Operator CRDs.
	"policy-reporter": {
		{Name: "kyverno", Chart: "oci://ghcr.io/kyverno/charts/kyverno", Version: "3.9.1", Namespace: "kyverno"},
		{Name: "trivy-operator", Chart: "oci://ghcr.io/aquasecurity/helm-charts/trivy-operator", Version: "0.32.1", Namespace: "trivy-system"},
	},
	"argo-workflows": {traefikCRDs},
}

// Prereqs returns the prerequisites of an application (nil when none).
func Prereqs(app string) []Prereq { return prereqs[app] }

// Install applies the prerequisite to the cluster. Helm prerequisites are
// written as OCIRepository + HelmRelease in the test namespace; the caller
// waits for the HelmRelease named p.Name to become Ready.
func (p Prereq) Install(ctx context.Context, env *environment.Env) error {
	dir, err := os.MkdirTemp("", "catalog-prereq-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, p.Name+".yaml")

	if p.IsHelm() {
		if err := os.WriteFile(path, []byte(p.helmReleaseYAML()), 0o600); err != nil {
			return err
		}
	} else {
		if err := download(ctx, p.ManifestURL, path); err != nil {
			return err
		}
		if err := dropKindlessDocs(path); err != nil {
			return err
		}
	}
	return env.ApplyYAMLWithoutSubstitutions(ctx, path)
}

func (p Prereq) helmReleaseYAML() string {
	values := p.Values
	if values == "" {
		values = "{}\n"
	}
	return fmt.Sprintf(`apiVersion: source.toolkit.fluxcd.io/v1
kind: OCIRepository
metadata:
  name: %[1]s-prereq-chart
  namespace: %[2]s
spec:
  interval: 1h
  ref:
    tag: %[3]s
  url: %[4]s
  layerSelector:
    mediaType: application/vnd.cncf.helm.chart.content.v1.tar+gzip
    operation: copy
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: %[1]s-prereq-values
  namespace: %[2]s
data:
  values.yaml: |
%[6]s
---
apiVersion: helm.toolkit.fluxcd.io/v2
kind: HelmRelease
metadata:
  name: %[1]s
  namespace: %[2]s
spec:
  chartRef:
    kind: OCIRepository
    name: %[1]s-prereq-chart
    namespace: %[2]s
  interval: 15s
  releaseName: %[1]s
  targetNamespace: %[5]s
  install:
    crds: CreateReplace
    createNamespace: true
    remediation:
      retries: 10
  upgrade:
    crds: CreateReplace
  valuesFrom:
    - kind: ConfigMap
      name: %[1]s-prereq-values
`, p.Name, DefaultNamespace, p.Version, p.Chart, p.Namespace, indent(values, "    "))
}

func indent(s, prefix string) string {
	out := ""
	for _, line := range splitLines(s) {
		out += prefix + line + "\n"
	}
	return out
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// dropKindlessDocs rewrites a multi-document YAML file without the documents
// that carry no "kind:" (license headers, trailing separators), which the
// harness would otherwise reject as objects without a kind.
func dropKindlessDocs(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var kept []string
	for _, doc := range strings.Split(string(raw), "\n---") {
		if strings.Contains(doc, "\nkind:") || strings.HasPrefix(strings.TrimSpace(doc), "kind:") {
			kept = append(kept, strings.TrimLeft(doc, "\n"))
		}
	}
	return os.WriteFile(path, []byte(strings.Join(kept, "\n---\n")), 0o600)
}

func download(ctx context.Context, url, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
