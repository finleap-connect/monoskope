package helm

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"

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
		kind.ClusterWithConfig(&v1alpha4.Cluster{
			KubeadmConfigPatchesJSON6902: []v1alpha4.PatchJSON6902{
				{
					Group:   "kubeadm.k8s.io",
					Version: "v1beta2",
					Kind:    "ClusterConfiguration",
					Patch:   "- op: add\r\n  path: /apiServer/certSANs/-\r\n  value: docker",
				},
			},
			KubeadmConfigPatches: []string{
				"kind: InitConfiguration\nnodeRegistration:\n  kubeletExtraArgs:\n    cgroup-root: \"kind\"\n",
			},
		}),
	}
	if KindCluster != "" {
		clusterOptions = append(clusterOptions,
			kind.ClusterWithName(KindCluster),
			kind.ClusterUseExisting(),
			kind.ClusterDoNotDelete(),
		)
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
