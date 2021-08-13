package jwt

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("jwt/claims", func() {
	expectedIssuer := "https://localhost"
	expectedValidity := time.Hour * 1
	It("validate cluster bootstrap token", func() {
		t := NewClusterBootstrapToken(&StandardClaims{}, expectedIssuer, "me")
		Expect(t.Validate(expectedIssuer, AudienceM8Operator, AudienceMonoctl)).ToNot(HaveOccurred())
	})
	It("validate auth token", func() {
		t := NewAuthToken(&StandardClaims{}, expectedIssuer, "me", expectedValidity)
		Expect(t.Validate(expectedIssuer, AudienceMonoctl, AudienceM8Operator)).ToNot(HaveOccurred())
	})
	It("validate auth token", func() {
		t := NewAuthToken(&StandardClaims{}, expectedIssuer, "me", expectedValidity)
		Expect(t.Validate(expectedIssuer, AudienceK8sAuth)).To(HaveOccurred())
	})
	It("fail validate auth token", func() {
		t := NewAuthToken(&StandardClaims{}, expectedIssuer, "me", expectedValidity)
		t.Expiry = jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * -12))
		Expect(t.Validate(expectedIssuer, AudienceMonoctl)).To(HaveOccurred())
	})
})
