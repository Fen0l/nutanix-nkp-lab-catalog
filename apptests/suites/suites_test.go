package suites

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/Fen0l/nutanix-nkp-lab-catalog/apptests/catalog"
	fluxhelmv2 "github.com/fluxcd/helm-controller/api/v2"
	"github.com/mesosphere/kommander-applications/apptests/client"
	"github.com/mesosphere/kommander-applications/apptests/docker"
	"github.com/mesosphere/kommander-applications/apptests/environment"
	"github.com/mesosphere/kommander-applications/apptests/flux"
	"github.com/mesosphere/kommander-applications/apptests/kind"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	genericClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	env       *environment.Env
	ctx       context.Context
	network   *docker.NetworkResource
	k8sClient genericClient.Client

	useExistingCluster bool
)

var (
	appName    = flag.String("app", "", "Catalog application to test (directory name under applications/)")
	appVersion = flag.String("app-version", "", "Application version to test (defaults to the latest version directory)")
)

var _ = BeforeSuite(func() {
	Expect(*appName).ToNot(BeEmpty(), "-app flag is required")

	log.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	ctx = context.Background()

	if kubeconfig := os.Getenv("E2E_KUBECONFIG"); kubeconfig != "" {
		// Reuse an existing cluster (must already run Flux source/kustomize/helm controllers).
		useExistingCluster = true
		env = &environment.Env{}

		typedClient, err := client.NewClient(kubeconfig)
		Expect(err).ShouldNot(HaveOccurred())
		env.K8sClient = typedClient

		scheme := flux.NewScheme()
		_ = fluxhelmv2.AddToScheme(scheme)
		k8sClient, err = genericClient.New(typedClient.Config(), genericClient.Options{Scheme: scheme})
		Expect(err).ShouldNot(HaveOccurred())
		env.Client = k8sClient
		return
	}

	var err error
	network, err = kind.EnsureDockerNetworkExist(ctx, "", false)
	Expect(err).ShouldNot(HaveOccurred())
	env = &environment.Env{Network: network}
})

func TestCatalogApplications(t *testing.T) {
	RegisterFailHandler(Fail)
	suiteConfig, reporterConfig := GinkgoConfiguration()
	RunSpecs(t, "NKP Lab Catalog Application Suite", suiteConfig, reporterConfig)
}

// setupCluster provisions a kind cluster with Flux unless E2E_KUBECONFIG is set.
func setupCluster() {
	if useExistingCluster {
		return
	}
	Expect(env.Provision(ctx)).To(Succeed())
	// Flux pinned to the NKP 2.19 release line (see catalog.DefaultFluxVersion / FLUX_VERSION).
	Expect(catalog.InstallFlux(ctx, env, catalog.FluxVersion())).To(Succeed())

	scheme := flux.NewScheme()
	_ = fluxhelmv2.AddToScheme(scheme)
	var err error
	k8sClient, err = genericClient.New(env.K8sClient.Config(), genericClient.Options{Scheme: scheme})
	Expect(err).ShouldNot(HaveOccurred())
}

func teardownCluster() {
	if useExistingCluster || os.Getenv("SKIP_CLUSTER_TEARDOWN") != "" {
		return
	}
	Expect(env.Destroy(ctx)).To(Succeed())
}
