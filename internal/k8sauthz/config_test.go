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

package k8sauthz

import (
	_ "embed"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:embed test_config.yaml
var test_config []byte

var _ = Describe("internal/k8sauthz", func() {
	Context("GitRepository", func() {
		It("NewGitRepository() creates a new instance with defaults", func() {
			os.Setenv("test1.basic.username", "test1")
			os.Setenv("test1.basic.password", "testpw")

			os.Setenv("test2.ssh.privateKey", `-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAxgq3A6cgAXA03h8Iv8qlKvvixxKOAMrwpa0Qk0iNrNW7s92x
xfCPsB2iJPg8TK4jhIBrFXR5ck80hJ+poImYZP3g2ayTPJAEOIYUP7r4R3NarFMs
mOqXibADhPaQhVYbHRPL5xqe+0lzXiBNkARVMJ8VAj1FdfFW900mieBItULquGgS
4PZeXO5y9a5oOed5lPg6bNXnT1f1X5qjKHIl51w83raHx3wYzaBZytJz3MV8EAar
SpVt/tTw2rakf5aQIeB6f+EWpYEJA+vEAQViIZLIx//tODdSvQb1KjHsRzxPBCZ7
SUNW1dOp1I7ctKiigeJdWCx+fqVnV+jTF3cl7QIDAQABAoIBAQDADLACCNYOzwNV
y/6uDxQxJOd8ac/vkN3hrH0Ad0F8tCA8gubcrzaIcZRGCbxgGBc1bVtrCPKlGY4y
dqsGlaiFd3XPYX2/I2IUdmG4k9YlZAYwStBa0lR2H4yqKnln9GKCxincSYKW6VnJ
k8XjvW/noL282DI99/aAB/8xgaDAtRZfNE/ViVXfmW7rQu6H2aVcujjLWmDQlRys
Cq83sr5AzoQf1ul/i92MvzyCvqQkxr56wK24DSAfSqM2lPtNh/pQB9FMkYHe3aiR
Bc5j+FwLgX4PvGuCWsLLj58CxZZi0MPSRhXFhqUXMwx0WiE+kYHTMPIl298oHAnv
ZCUkQNq5AoGBAPxl7Oi4LCZoxhvbE12q4z+XRx/LRNdUmw/mt1IQHWRZc8r2kydc
g2YTKaGu+3yJJPw0Filw1PbyrSlwEZA768DO8H1jj2PTYihLFFmT71OmtegEWJQ2
UP6mMq5AC7OhUFa5Xb/NqGQIojGIGUPXDD0yPjnHF78lJZLG+FGRVbh/AoGBAMje
Nk3+YK98q/OLGBj52zlUMTkHVcz/L2a8hZg4nsi+AhNqNV5gt/RgGlenE9GR46iP
kAZlNVGToeh1POXJUD5MGvZ+tq4z/x9vtNEkj88YQet7Z39sjct++jHpI87OD5Vh
+FX+PfIkPUHn8MbzuNUURbk/IDcezjdmtXfpT0uTAoGBANEw8VBZBU64XCpOQxjJ
QBw+l9aENQR7vSFKyfocHSBE+7Hm+EQRa1641zKIwyuOx7a4vP2P5RLgRCEqH4R6
BYYCflSOphPI6XwX7j/oWt3sOKyanYgKFMtamHrP3mL9eobGkCsem0h912BSOw59
OEmoUOa2Wro0YZXrj3ibin6lAoGBALFy3TZ5iqB82ssEDf3CfNBgeDuRJ1M5/H4d
/UdEfcGJWfwcz5jYiKnyXFMi3nc9EolOikWcE0ZFzbq2F15IQgNG6Grp9ihe9vqm
jG5WhlTFcTa+hegj1f/35L9C3sdLcqjdK62MNN089oDKjgK0PIzAo/fKtIL9S9/V
dGrCSbLRAoGBALEns5tKchvtSUgFcD0wkBovwxMewNIdJ1M7AtE4fBzX7G2kKMH5
UJY+A+PVvq05FqxKqFcBRLT+em9FGgUCemVBbX6KTuyLd8tAFzZZomRoSA1Umm9A
qjbpxy/532FZXv4+K3Z7PmJqTUKhgOnM/UDHzapBRXGtyicIoqOEiBmi
-----END RSA PRIVATE KEY-----
`)
			os.Setenv("test2.ssh.password", "a25vd24taG9zdC1rZXlz")
			os.Setenv("test2.ssh.known_hosts", "github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==")

			conf, err := NewConfigFromFile(test_config)
			Expect(err).ToNot(HaveOccurred())
			Expect(conf).ToNot(BeNil())
			Expect(len(conf.Repositories)).To(BeNumerically("==", 2))
			Expect(len(conf.Mappings)).To(BeNumerically("==", 2))

			firstRepo := conf.Repositories[0]
			Expect(firstRepo.AllClusters).To(BeTrue())
			Expect(len(firstRepo.Clusters)).To(BeNumerically("==", 0))

			secondRepo := conf.Repositories[1]
			Expect(secondRepo.AllClusters).To(BeFalse())
			Expect(len(secondRepo.Clusters)).To(BeNumerically("==", 2))
		})
	})
})
