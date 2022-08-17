// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwt

import (
	"crypto/rsa"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("jwt/key", func() {
	It("can load private key from file", func() {
		bytes, err := os.ReadFile(testEnv.privateKeyFile)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).ToNot(BeNil())

		privKey, err := LoadPrivateKey(bytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(privKey).ToNot(BeNil())
		Expect(privKey.Key).To(Equal(testEnv.privateKey))
	})
	It("can load public key from file", func() {
		bytes, err := os.ReadFile(testEnv.publicKeyFile)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).ToNot(BeNil())

		pubKey, err := LoadPublicKey(bytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(pubKey).ToNot(BeNil())

		rsaPublicKey, ok := pubKey.Key.(*rsa.PublicKey)
		Expect(ok).To(BeTrue())
		Expect(*rsaPublicKey).To(Equal(testEnv.privateKey.PublicKey))
	})
	It("can load public/private key from cert", func() {
		pemCert :=
			`-----BEGIN CERTIFICATE-----
MIICnTCCAkSgAwIBAgIQMo7x823NtJ/Xyy1Wl+8+yzAKBggqhkjOPQQDAjAnMSUw
IwYDVQQDExxyb290Lm1vbm9za29wZS5jbHVzdGVyLmxvY2FsMB4XDTIxMDYwMjAy
MTAxNVoXDTIxMDYwNDAyMTAxNVowMDESMBAGA1UEChMJTW9ub3Nrb3BlMRowGAYD
VQQDExFtOC1hdXRoZW50aWNhdGlvbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBANCKZWW0el3OzPw7914TC1Ld2At/xIh/3zoiawQcbS8mrjnVMO2oSomY
mks6sEaWp4p80PwJkzSplpgoJmEOYqps+YXo+1NLp66bFPkAbMEZDsZ4QmrQQ7X3
iv5IaDFW4vSGJFSkTQnUmedlhrWguasOD3vL0Pek89L8kQ09+YlDk/fpBZUXFADU
+ef4GjTkWJzkg32dSOudJDYD4wUPczTFlRO097MBBlaMb4LKYfDfjuUKRCOAL3LD
7kKAatHKeoADuBptUv/lQLExGNzlhRteaLocTHHab2hs+NCFYABv2Px5Tcnbw8g+
/r/97gwKkpFeF5p4WhdVgbDYd2MGUlMCAwEAAaN+MHwwHQYDVR0lBBYwFAYIKwYB
BQUHAwIGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAUI+RyVqj0
J9qH8l3pbY9KUkoTgHQwLAYDVR0RBCUwI4IJbG9jYWxob3N0hwR/AAABhxAAAAAA
AAAAAAAAAAAAAAAAMAoGCCqGSM49BAMCA0cAMEQCIEPbvMo2YvqlYQtdkQwlhJci
mTlsDv6VmO4WfCjrQdwLAiA+N0eeiL/yLPC5ReaPYQ7PeoXbc9+EPR2FBDrkiBbA
8w==
-----END CERTIFICATE-----`

		pemPrivKey :=
			`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0IplZbR6Xc7M/Dv3XhMLUt3YC3/EiH/fOiJrBBxtLyauOdUw
7ahKiZiaSzqwRpaninzQ/AmTNKmWmCgmYQ5iqmz5hej7U0unrpsU+QBswRkOxnhC
atBDtfeK/khoMVbi9IYkVKRNCdSZ52WGtaC5qw4Pe8vQ96Tz0vyRDT35iUOT9+kF
lRcUANT55/gaNORYnOSDfZ1I650kNgPjBQ9zNMWVE7T3swEGVoxvgsph8N+O5QpE
I4AvcsPuQoBq0cp6gAO4Gm1S/+VAsTEY3OWFG15ouhxMcdpvaGz40IVgAG/Y/HlN
ydvDyD7+v/3uDAqSkV4XmnhaF1WBsNh3YwZSUwIDAQABAoIBAQCYyc0ghupgcHOf
GhBSzIEvZXo0cpf7qjRS04S0rl8QfLaJiLkgZny18yiYlZcxII//1xMGlb1UiCvd
rwzvbyq60ry+b8QzcuqX8ueax8TmdQVuRA3lVFFHsOYVB9fOzmnZ3a4glYAcA7f+
4VOhHvDpcpPFj766shAyNPnRSebZuWUvvbqMi1wqFva3GXdacor26awHkKsTEk4G
NBJEsVIq/to5T37CyTHHoFOCXDRIxJw1Cd3Nw9r5lNv4JJJttHk8MhOC2ZOtPLWX
Nh1laUPCo8yHfH0xw2fOJLBIFDVFe8QZaJ4Q5vz6yYLahyEMOK/iactZe/Z84RYJ
ZtpcRyWBAoGBANGOWagciFPjIl7ZbU8gCxaaZLKMtYwHagQnFRUE7c8vt/lbn5IE
i2L+xlpn4IV2pPvkkMB3SGex18sHcTbPx0jsvddXperCXOjj6JF9PdyEfmQByvx1
VGLV/8rVZz10isg5f6pNBwEqfUro81hJA1DpIaX4j1JYUHIOMGyEF3QXAoGBAP7C
bqNzlBQG55NjwsL05FGR9Rg/1HHjcijielF+cuZSP6eIbn9ucD1QV0wcyrCE2DDv
YZ0jYUb09ZOttM8UfKF94u0EvDKMGKnt+l1tcFUDLAePaJNi+Lux1l5zi1bfO22h
jZpPKoDAWxfhkODlx+VMT6Y/g8YDpltqXywJqa0lAoGAc7FZglyuT1H46dC0bpjM
RmBa89CHcpWtTDmfhAlCmb5Indydzmm/4pmyPLtY05ZbI85etEOmr8kZ0Dd9o7s2
1OYPMVJsgZ1o2hLplVlFy/dCKEhtHtBQFHj9Tahf5SfwbvZ/qy/3jAc/QRo3Lyiw
Mf1j3FPMHLQxRabbyS1sHWUCgYEA39gXPqcfRTmL4IWXa5Whx8pngJcVI7ylYicd
Mt3YN2etZpcKAA4ZsMYW7lmd/tu62cR8EIY1wxMZdFj8tbdaissBySCP/Bn80dK4
Wb7/JLNUzI/FYztjMghgQz1jAUHEBeAde6hzwA1D/QfFNNaxfVg/4+OK9UHfuhMM
7LTQ0cECgYAjI+Ig1L9MpAVvN5YFsJU9eez8cP7h2CB+UX0sFDiMoPKVBKpxyibJ
BOcqwoOhP3v4Pq/OT6vB27GnRjqST03cJDt17tcgBohyrI7odUEwvMvgrCQQUMtz
hRkBLbTlzXZsdQJDdBeE+6+YFNmowNqHnrhzqIo54VMzjF1so+OynA==
-----END RSA PRIVATE KEY-----`

		pubKey, err := LoadPublicKey([]byte(pemCert))
		Expect(err).ToNot(HaveOccurred())
		Expect(pubKey).ToNot(BeNil())

		rsaPublicKey, ok := pubKey.Key.(*rsa.PublicKey)
		Expect(ok).To(BeTrue())

		privKey, err := LoadPrivateKey([]byte(pemPrivKey))
		Expect(err).ToNot(HaveOccurred())

		rsaPrivKey, ok := privKey.Key.(*rsa.PrivateKey)
		Expect(ok).To(BeTrue())

		Expect(*rsaPublicKey).To(Equal(rsaPrivKey.PublicKey))
	})
})
