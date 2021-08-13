package kind

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"helm.sh/helm/v3/pkg/release"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	"github.com/kubism/testutil/pkg/helm"
	"github.com/kubism/testutil/pkg/kind"
)

var (
	log             logger.Logger
	cluster         *kind.Cluster
	kubeConfig      string
	helmClient      *helm.Client
	helmReleaseName string
	releases        []*release.Release
	repositories    []*helm.RepositoryEntry = []*helm.RepositoryEntry{
		{
			Name: "jetstack",
			URL:  "https://charts.jetstack.io",
		},
		{
			Name: "bitnami",
			URL:  "https://charts.bitnami.com/bitnami",
		},
		{
			Name: "ambassador",
			URL:  "https://getambassador.io",
		},
	}
)

func TestKind(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t,
		"Kind Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error
	log = logger.WithName("kind-test-suite")

	By("setup kind cluster")

	workspaceFolder := os.Getenv("WORKSPACE_FOLDER")
	Expect(workspaceFolder).To(Not(BeEmpty()))

	kindClusterName := os.Getenv("KIND_CLUSTER_NAME")
	Expect(kindClusterName).To(Not(BeEmpty()))

	helmChartPath := os.Getenv("HELM_CHART_PATH")
	Expect(helmChartPath).To(Not(BeEmpty()))

	helmChartValues := os.Getenv("HELM_CHART_VALUES")
	Expect(helmChartValues).To(Not(BeEmpty()))

	helmReleaseName = os.Getenv("HELM_RELEASE_NAME")
	Expect(helmChartValues).To(Not(BeEmpty()))

	log.Info("Starting kind cluster...", "Name", kindClusterName)
	cluster, err = kind.NewCluster(
		kind.ClusterWithWaitForReady(10*time.Minute),
		kind.ClusterWithName(kindClusterName),
		kind.ClusterUseExisting(),
		kind.ClusterDoNotDelete(),
	)
	Expect(err).To(Succeed())

	By("setup kubeconfig")
	log.Info("Setting up kubeconfig...", "Name", kindClusterName)
	kubeConfig, err = cluster.GetKubeConfig()
	Expect(err).To(Succeed())

	log.Info("Setting up helm...")
	helmClient, err = helm.NewClient(kubeConfig)
	Expect(err).To(Succeed())

	for _, repo := range repositories {
		Expect(helmClient.AddRepository(repo)).To(Succeed())
	}

	certManagerValuesFile := filepath.Join(workspaceFolder, "test", "kind", "cert-manager-values.yaml")
	log.Info("Setting up certmanager...", "Values", certManagerValuesFile)
	rls, err := helmClient.Install("jetstack/cert-manager", "v1.1.0", helm.ValuesOptions{
		ValueFiles: []string{certManagerValuesFile}},
		helm.InstallWithReleaseName("cert-manager"),
	)
	Expect(err).ToNot(HaveOccurred())
	Expect(rls).ToNot(BeNil())
	releases = append(releases, rls)

	log.Info("Installing chart...", "ChartPath", helmChartPath, "Values", helmChartValues)
	rls, err = helmClient.Install(helmChartPath, "", helm.ValuesOptions{ValueFiles: []string{helmChartValues}},
		helm.InstallWithReleaseName(helmReleaseName),
	)
	Expect(err).ToNot(HaveOccurred())
	Expect(rls).ToNot(BeNil())
	releases = append(releases, rls)
}, 240)

var _ = AfterSuite(func() {
	By("tearing down kind cluster")

	for _, release := range releases {
		log.Info("Uninstalling helm release...", "Release", release.Name)
		if err := helmClient.Uninstall(release.Name); err != nil {
			log.Error(err, "Uninstalling helm release failed.", "Release", release.Name)
		}
	}

	log.Info("Shutting down cluster...")
	if cluster != nil {
		if err := cluster.Close(); err != nil {
			log.Error(err, "Shutting down cluster failed.")
		}
	}

	log.Info("Freeing helm resources...")
	if helmClient != nil {
		if err := helmClient.Close(); err != nil {
			log.Error(err, "Freeing helm resources failed.")
		}
	}
})
