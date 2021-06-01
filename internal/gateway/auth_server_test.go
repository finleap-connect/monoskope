package gateway

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gateway Auth Server", func() {
	It("can retrieve openid conf", func() {
		res, err := httpClient.Get(fmt.Sprintf("http://%s/.well-known/openid-configuration", localAddrAuthServer))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(doc.Text()).NotTo(BeEmpty())
	})
	It("can retrieve jwks", func() {
		res, err := httpClient.Get(fmt.Sprintf("http://%s/keys", localAddrAuthServer))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(doc.Text()).NotTo(BeEmpty())
	})
})

var _ = Describe("HealthCheck", func() {
	It("can do readiness checks", func() {
		res, err := httpClient.Get(fmt.Sprintf("http://%s/readyz", localAddrAuthServer))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})
})
