package helm

import (
	"github.com/kubism/testutil/pkg/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

var _ = Describe("Helm chart", func() {
	It("can be installed", func() {
		if !test.WithKind {
			return
		}

		rls, err := helmClient.Install(test.HelmChartPath, "", helm.ValuesOptions{ValueFiles: []string{test.HelmChartValues}},
			helm.InstallWithReleaseName("monoskope"),
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(rls).ToNot(BeNil())
		Expect(helmClient.Uninstall(rls.Name)).To(Succeed())
	})
})
