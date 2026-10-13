package catalog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fluxcd/flux2/v2/pkg/manifestgen"
	"github.com/fluxcd/flux2/v2/pkg/manifestgen/install"
	runclient "github.com/fluxcd/pkg/runtime/client"
	"github.com/mesosphere/kommander-applications/apptests/environment"
	"github.com/mesosphere/kommander-applications/apptests/flux"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/util/retry"
)

const (
	// FluxNamespace is where NKP runs Flux (and where Kommander stores Helm release secrets).
	FluxNamespace = "kommander-flux"

	// DefaultFluxVersion is the flux2 release whose helm-controller matches the one shipped by
	// NKP 2.18 (helm-controller v1.5.4, source-controller v1.8.x). The upstream harness installs
	// "latest". Override with FLUX_VERSION (a flux2 tag, or "latest").
	DefaultFluxVersion = "v2.8.6"

	// helmControllerFeatureGates mirrors the NKP 2.18 helm-controller flags that change how charts
	// render. DisableChartDigestTracking stops helm-controller from appending the OCI digest to the
	// chart version (`0.37.3+755df9…`), which breaks charts that reuse .Chart.Version as a label or
	// an image tag (LitmusChaos, vCluster). Without it the kind results differ from NKP.
	helmControllerFeatureGates = "--feature-gates=DisableChartDigestTracking=true"
)

// FluxVersion returns the flux2 release to install on the kind cluster.
func FluxVersion() string {
	if v := os.Getenv("FLUX_VERSION"); v != "" {
		return v
	}
	return DefaultFluxVersion
}

// InstallFlux installs the source, kustomize and helm controllers of the given flux2 release in
// the kommander-flux namespace and waits for them to be ready. It mirrors
// environment.Env.InstallLatestFlux with a pinned version.
func InstallFlux(ctx context.Context, env *environment.Env, version string) error {
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: FluxNamespace}}
	if _, err := env.K8sClient.Clientset().CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{}); err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	options := install.MakeDefaultOptions()
	options.Version = version
	options.Namespace = FluxNamespace
	options.Components = []string{"source-controller", "kustomize-controller", "helm-controller"}

	manifest, err := install.Generate(options, "")
	if err != nil {
		return fmt.Errorf("generate flux %s manifests: %w", version, err)
	}
	tmpDir, err := manifestgen.MkdirTempAbs("", FluxNamespace)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	if _, err := manifest.WriteFile(tmpDir); err != nil {
		return fmt.Errorf("write flux manifests: %w", err)
	}

	kubeconfigPath := env.Cluster.KubeconfigFilePath()
	kubeconfigArgs := genericclioptions.NewConfigFlags(true)
	kubeconfigArgs.KubeConfig = &kubeconfigPath
	if _, err := flux.Apply(ctx, kubeconfigArgs, new(runclient.Options), tmpDir, filepath.Join(tmpDir, manifest.Path)); err != nil {
		return fmt.Errorf("apply flux manifests: %w", err)
	}

	if err := setHelmControllerFeatureGates(ctx, env); err != nil {
		return err
	}

	return wait.PollUntilContextTimeout(ctx, PollInterval, 5*time.Minute, true, func(ctx context.Context) (bool, error) {
		deployments, err := env.K8sClient.Clientset().AppsV1().Deployments(FluxNamespace).List(ctx, metav1.ListOptions{
			LabelSelector: manifestgen.PartOfLabelKey + "=" + manifestgen.PartOfLabelValue,
		})
		if err != nil || len(deployments.Items) < len(options.Components) {
			return false, nil
		}
		for _, d := range deployments.Items {
			if d.Generation > d.Status.ObservedGeneration || d.Status.ReadyReplicas < 1 {
				return false, nil
			}
		}
		return true, nil
	})
}

// setHelmControllerFeatureGates adds the NKP helm-controller feature gates to the deployment.
func setHelmControllerFeatureGates(ctx context.Context, env *environment.Env) error {
	deployments := env.K8sClient.Clientset().AppsV1().Deployments(FluxNamespace)
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		d, err := deployments.Get(ctx, "helm-controller", metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get helm-controller deployment: %w", err)
		}
		for i := range d.Spec.Template.Spec.Containers {
			c := &d.Spec.Template.Spec.Containers[i]
			if c.Name != "manager" {
				continue
			}
			for _, a := range c.Args {
				if a == helmControllerFeatureGates {
					return nil
				}
			}
			c.Args = append(c.Args, helmControllerFeatureGates)
		}
		_, err = deployments.Update(ctx, d, metav1.UpdateOptions{})
		return err
	})
}
