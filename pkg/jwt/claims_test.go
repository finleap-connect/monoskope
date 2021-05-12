package jwt

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("jwt/claims", func() {
	It("validate claims", func() {
		t := NewClusterBootstrapToken("me")
		Expect(t.IsExpired()).To(BeFalse())
		Expect(t.IsValid()).To(BeTrue())
	})
})
