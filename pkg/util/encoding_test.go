package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	teststring = "this is a test"
	testbase64 = "dGhpcyBpcyBhIHRlc3Q="
)

var _ = Describe("util.encoding", func() {
	It("can encode bytes to Base64", func() {
		base64 := BytesToBase64([]byte(teststring))
		Expect(testbase64).To(Equal(base64))
	})
	It("can decode bytes to Base64", func() {
		teststringbytes, err := Base64ToBytes(testbase64)
		Expect(err).NotTo(HaveOccurred())
		Expect(teststringbytes).To(Equal([]byte(teststring)))
	})
})
