package helm

import (
	"github.com/kubism/testutil/pkg/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helm chart", func() {
	It("can be installed", func() {
		if !WithKind {
			return
		}

		rls, err := helmClient.Install(HelmChartPath, "", helm.ValuesOptions{ValueFiles: []string{HelmChartValues}},
			helm.InstallWithReleaseName("monoskope"),
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(rls).ToNot(BeNil())
		Expect(helmClient.Uninstall(rls.Name)).To(Succeed())
	})
})
