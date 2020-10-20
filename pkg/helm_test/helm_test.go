package helm_test

import (
	"github.com/kubism/testutil/pkg/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helm", func() {
	Describe("Install monoskope helm chart", func() {
		It("should install helm chart", func() {
			release, err := helmClient.Install("monoskope", "local", helm.ValuesOptions{ValueFiles: []string{"examples/00-monoskope-dev-values.yaml"}}, helm.InstallWithReleaseName("local"))
			Expect(err).To(BeNil())
			Expect(release).To(Not(nil))
		})
	})
})
