package jwt

import (
	"crypto/rsa"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("jwt/key", func() {
	It("can load public key from file", func() {
		bytes, err := ioutil.ReadFile(testEnv.privateKeyFile)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).ToNot(BeNil())

		privKey, err := LoadPrivateKey(bytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(privKey).ToNot(BeNil())
		Expect(privKey).To(Equal(testEnv.privateKey))
	})
	It("can load private key from file", func() {
		bytes, err := ioutil.ReadFile(testEnv.publicKeyFile)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).ToNot(BeNil())

		pubKey, err := LoadPublicKey(bytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(pubKey).ToNot(BeNil())

		rsaPublicKey, ok := pubKey.(*rsa.PublicKey)
		Expect(ok).To(BeTrue())
		Expect(*rsaPublicKey).To(Equal(testEnv.privateKey.PublicKey))
	})
})
