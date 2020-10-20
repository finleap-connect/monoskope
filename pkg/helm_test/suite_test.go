package helm_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	"github.com/kubism/testutil/pkg/helm"
	"github.com/kubism/testutil/pkg/kind"
)

var (
	cluster          *kind.Cluster
	kubeconfig       string
	helmClient       *helm.Client
	stableRepository = &helm.RepositoryEntry{
		Name: "stable",
		URL:  "https://kubernetes-charts.storage.googleapis.com",
	}
	kubismRepository = &helm.RepositoryEntry{
		Name: "kubism.io",
		URL:  "https://kubism.github.io/charts",
	}
)

func TestHelm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t,
		"Helm Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	var err error

	By("bootstrapping test environment")

	cluster, err = kind.NewCluster(kind.ClusterWithWaitForReady(5 * time.Minute))
	Expect(err).ToNot(HaveOccurred())

	kubeconfig, err = cluster.GetKubeConfig()
	Expect(err).ToNot(HaveOccurred())

	helmClient, err = helm.NewClient(kubeconfig, helm.ClientWithDriver("secret"))
	Expect(err).ToNot(HaveOccurred())
	err = helmClient.AddRepository(stableRepository)
	Expect(err).ToNot(HaveOccurred())
	err = helmClient.AddRepository(kubismRepository)
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	if cluster != nil {
		cluster.Close()
	}
	if helmClient != nil {
		helmClient.Close()
	}
})
