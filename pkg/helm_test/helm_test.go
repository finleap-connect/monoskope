package helm_test

import (
	"github.com/kubism/testutil/pkg/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helm chart", func() {
	It("can be installed", func() {
		rls, err := helmClient.Install("kubism.io/monoskope", "", helm.ValuesOptions{},
			helm.InstallWithReleaseName("monoskope"),
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(rls).ToNot(BeNil())
		Expect(helmClient.Uninstall(rls.Name)).To(Succeed())
	})
})
