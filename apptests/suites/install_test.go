package suites

import (
	"fmt"
	"time"

	fluxhelmv2 "github.com/fluxcd/helm-controller/api/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Fen0l/nutanix-nkp-lab-catalog/apptests/catalog"
)

// installTimeout bounds the wait for a HelmRelease to become Ready. Charts
// with many subcharts (Gitea with postgresql-ha, Argo CD) need the upper end.
const installTimeout = 20 * time.Minute

var _ = Describe("Catalog application", Ordered, Label("install"), func() {
	var app *catalog.App

	BeforeAll(func() {
		var err error
		app, err = catalog.New(*appName, *appVersion)
		Expect(err).ToNot(HaveOccurred())
		GinkgoWriter.Printf("Testing %s @ %s\n", app.Name, app.Version)
		setupCluster()
	})

	AfterAll(func() {
		teardownCluster()
	})

	It("installs its prerequisites", func() {
		for _, p := range catalog.Prereqs(app.Name) {
			By(fmt.Sprintf("installing prerequisite %s", p.Name))
			Expect(p.Install(ctx, env)).To(Succeed())
			if p.IsHelm() {
				waitHelmReleaseReady(p.Name, installTimeout)
			}
		}
	})

	It("installs with the catalog defaults and reaches Ready", func() {
		Expect(app.Install(ctx, env)).To(Succeed())
		// The HelmRelease is named after ${releaseName} for multi-instance apps
		// and after the app for singletons; in tests both equal the app name.
		waitHelmReleaseReady(app.Name, installTimeout)
	})
})

// waitHelmReleaseReady polls the HelmRelease in the test namespace until its
// Ready condition is True, printing the latest status on every poll.
func waitHelmReleaseReady(name string, timeout time.Duration) {
	hr := &fluxhelmv2.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: catalog.DefaultNamespace},
	}
	Eventually(func() error {
		if err := k8sClient.Get(ctx, ctrlClient.ObjectKeyFromObject(hr), hr); err != nil {
			return err
		}
		for _, c := range hr.Status.Conditions {
			if c.Type == "Ready" {
				GinkgoWriter.Printf("HelmRelease %s: Ready=%s reason=%s %s\n", hr.Name, c.Status, c.Reason, c.Message)
				if c.Status == metav1.ConditionTrue {
					return nil
				}
				return fmt.Errorf("not ready: %s", c.Message)
			}
		}
		return fmt.Errorf("HelmRelease %s has no Ready condition yet", hr.Name)
	}).WithTimeout(timeout).WithPolling(catalog.PollInterval).Should(Succeed())
}
