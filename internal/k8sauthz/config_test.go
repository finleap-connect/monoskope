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

			conf, err := NewConfigFromFile(test_config)
			Expect(err).ToNot(HaveOccurred())
			Expect(conf).ToNot(BeNil())
			Expect(conf.Repository).ToNot(BeNil())
			Expect(len(conf.Mappings)).To(BeNumerically("==", 2))

			Expect(conf.AllClusters).To(BeTrue())
		})
	})
})
