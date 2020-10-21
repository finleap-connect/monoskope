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
	kubeConfig       string
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

	By("setup kind cluster")
	clusterOptions := []kind.ClusterOption{
		kind.ClusterWithWaitForReady(3 * time.Minute),
	}
	cluster, err = kind.NewCluster(clusterOptions...)
	Expect(err).To(Succeed())

	By("setup kubeconfig")
	kubeConfig, err = cluster.GetKubeConfig()
	Expect(err).To(Succeed())

	helmClient, err = helm.NewClient(kubeConfig)
	Expect(err).To(Succeed())
	err = helmClient.AddRepository(stableRepository)
	Expect(err).To(Succeed())
	err = helmClient.AddRepository(kubismRepository)
	Expect(err).To(Succeed())

	close(done)
}, 240)

var _ = AfterSuite(func() {
	By("tearing down kind cluster")
	if cluster != nil {
		cluster.Close()
	}
	if helmClient != nil {
		helmClient.Close()
	}
})
